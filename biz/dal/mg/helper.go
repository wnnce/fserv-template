package mg

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// CursorToAddrSlice iterates over a MongoDB cursor and decodes each document into a pointer of type T.
// It returns a slice of pointers to the decoded values.
//
// This function is useful when you need to retain references to each element,
// such as for mutation or when working with pointer-based data structures.
//
// Example:
//
//	var users []*User
//	users, err := CursorToAddrSlice[User](ctx, cursor)
func CursorToAddrSlice[T any](ctx context.Context, cursor *mongo.Cursor) ([]*T, error) {
	result := make([]*T, 0)
	for cursor.Next(ctx) {
		var row T
		if err := cursor.Decode(&row); err != nil {
			return nil, err
		}
		result = append(result, &row)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

// CursorToSlice iterates over a MongoDB cursor and decodes each document into a value of type T.
// It returns a slice of decoded values.
//
// This function is preferred when you don't need pointers and want to avoid extra heap allocations.
//
// Example:
//
//	var users []User
//	users, err := CursorToSlice[User](ctx, cursor)
func CursorToSlice[T any](ctx context.Context, cursor *mongo.Cursor) ([]T, error) {
	result := make([]T, 0)
	for cursor.Next(ctx) {
		var row T
		if err := cursor.Decode(&row); err != nil {
			return nil, err
		}
		result = append(result, row)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
