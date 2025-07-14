package db

import (
	"context"
	"math"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wnnce/fserv-template/pkg/sqlbuild"
)

type PageData[T any] struct {
	Current int   `json:"current,omitempty"`
	Size    int   `json:"size,omitempty"`
	Total   int64 `json:"total,omitempty"`
	Pages   int   `json:"pages,omitempty"`
	Records []*T  `json:"records,omitempty"`
}

func SelectPage[T any](
	ctx context.Context,
	builder sqlbuild.SelectBuilder,
	page, size int,
	safe bool,
	db *pgxpool.Pool,
) (*PageData[T], error) {
	var total int64
	row := db.QueryRow(ctx, builder.CountSQL(), builder.Args()...)
	if err := row.Scan(&total); err != nil {
		return nil, err
	}
	if total == 0 {
		return &PageData[T]{
			Current: page,
			Size:    size,
			Total:   total,
			Pages:   0,
			Records: make([]*T, 0),
		}, nil
	}
	offset := ComputeOffset(total, page, size, safe)
	builder.Limit(int64(size)).Offset(offset)
	rows, err := db.Query(ctx, builder.SQL(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	records, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*T, error) {
		return pgx.RowToAddrOfStructByNameLax[T](row)
	})
	if err != nil {
		return nil, err
	}
	pages := int(math.Ceil(float64(total) / float64(size)))
	if page > pages && safe {
		page = pages
	}
	return &PageData[T]{
		Current: page,
		Size:    size,
		Total:   total,
		Pages:   pages,
		Records: records,
	}, nil
}

func ComputeOffset(total int64, page, size int, safe bool) int64 {
	if page < 1 {
		return 0
	}
	offset := int64((page - 1) * size)
	if !safe || offset < total {
		return offset
	}
	int64Size := int64(size)
	pages := int64(math.Ceil(float64(total) / float64(int64Size)))
	if pages > 0 {
		return (pages - 1) * int64Size
	}
	return 0
}
