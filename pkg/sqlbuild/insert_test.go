package sqlbuild

import (
	"reflect"
	"strings"
	"testing"
)

type testRow struct {
	ID   int
	Name string
}

func TestPostgresInsertBuilder_Basic(t *testing.T) {
	builder := NewInsertBuilder("users").Insert("name", "Tom")
	sql := builder.SQL()
	expected := "INSERT INTO users (name) VALUES ($1)"
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	args := builder.Args()
	if !reflect.DeepEqual(args, []any{"Tom"}) {
		t.Errorf("expected args [Tom], got %v", args)
	}
}

func TestPostgresInsertBuilder_InsertRaw(t *testing.T) {
	builder := NewInsertBuilder("users").InsertRaw("created_at", "now()")
	sql := builder.SQL()
	if !strings.Contains(sql, "(created_at)") || !strings.Contains(sql, "VALUES (now())") {
		t.Errorf("expected raw insert in sql, got %q", sql)
	}
}

func TestPostgresInsertBuilder_InsertBySlice(t *testing.T) {
	builder := NewInsertBuilder("users").InsertBySlice([]string{"name", "age"}, []any{"Tom", 18})
	sql := builder.SQL()
	if !strings.Contains(sql, "name,age") {
		t.Errorf("expected insert by slice in sql, got %q", sql)
	}
}

func TestPostgresInsertBuilder_InsertByMap(t *testing.T) {
	builder := NewInsertBuilder("users").InsertByMap(map[string]any{"name": "Tom", "age": 18})
	sql := builder.SQL()
	if !strings.Contains(sql, "name") || !strings.Contains(sql, "age") {
		t.Errorf("expected insert by map in sql, got %q", sql)
	}
}

func TestPostgresInsertBuilder_FieldsValues(t *testing.T) {
	builder := NewInsertBuilder("users").Fields("name", "age").Values("Tom", 18)
	sql := builder.SQL()
	if !strings.Contains(sql, "name,age") {
		t.Errorf("expected fields/values in sql, got %q", sql)
	}
}

func TestPostgresInsertBuilder_Returning(t *testing.T) {
	builder := NewInsertBuilder("users").Insert("name", "Tom").Returning("id")
	sql := builder.SQL()
	if !strings.Contains(sql, "RETURNING id") {
		t.Errorf("expected returning in sql, got %q", sql)
	}
}

func TestPostgresInsertBuilder_NoParams(t *testing.T) {
	builder := NewInsertBuilder("users")
	sql := builder.SQL()
	expected := "INSERT INTO users () VALUES ()"
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
}

func TestBatchInsertBuilder(t *testing.T) {
	rows := []testRow{{1, "Tom"}, {2, "Jerry"}}
	sql, args := BatchInsertBuilder[testRow]("users", rows, func(r testRow) []any { return []any{r.ID, r.Name} }, "id", "name")
	if !strings.Contains(sql, "INSERT INTO users (id, name) VALUES ($1,$2),($3,$4)") {
		t.Errorf("expected batch insert sql, got %q", sql)
	}
	if !reflect.DeepEqual(args, []any{1, "Tom", 2, "Jerry"}) {
		t.Errorf("expected args [1,Tom,2,Jerry], got %v", args)
	}
}
