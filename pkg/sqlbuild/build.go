package sqlbuild

import (
	"strings"
	"sync"
)

// SQLBuilder SQL构造器接口
type SQLBuilder interface {
	table(tableName string) SQLBuilder
	SQL() string
	StatementParameterBuffer
}

func NewSelectBuilder(table string) SelectBuilder {
	selectBuilder := &PostgresSelectBuilder{
		tableName: table,
		limit:     -1,
		offset:    -1,
	}
	selectBuilder.builder = selectBuilder
	return selectBuilder
}

func NewUpdateBuilder(table string) UpdateBuilder {
	updateBuilder := &PostgresUpdateBuilder{
		tableName: table,
	}
	updateBuilder.builder = updateBuilder
	return updateBuilder
}

func NewDeleteBuilder(table string) DeleteBuilder {
	deleteBuilder := &PostgresDeleteBuilder{
		tableName: table,
	}
	deleteBuilder.builder = deleteBuilder
	return deleteBuilder
}

func NewInsertBuilder(table string) InsertBuilder {
	return &PostgresInsertBuilder{
		tableName: table,
	}
}

var defaultPool *sqlBuilderPool

func init() {
	defaultPool = newSQLBuilderPool()
}

// 字段条件和stringBuilder对象池
type sqlBuilderPool struct {
	fieldPool  *sync.Pool
	bufferPool *sync.Pool
}

func newSQLBuilderPool() *sqlBuilderPool {
	return &sqlBuilderPool{
		fieldPool: &sync.Pool{
			New: func() any {
				return &Field{}
			},
		},
		bufferPool: &sync.Pool{
			New: func() any {
				return &strings.Builder{}
			},
		},
	}
}

func (self *sqlBuilderPool) GetField() *Field {
	return self.fieldPool.Get().(*Field)
}
func (self *sqlBuilderPool) GetStringBuilder() *strings.Builder {
	return self.bufferPool.Get().(*strings.Builder)
}

// RecycleField 回收Field
func (self *sqlBuilderPool) RecycleField(field *Field) {
	field.Recycle()
	self.fieldPool.Put(field)
}

// RecycleStringBuilder 回收stringBuilder
func (self *sqlBuilderPool) RecycleStringBuilder(builder *strings.Builder) {
	builder.Reset()
	self.bufferPool.Put(builder)
}
