package sqlbuild

import (
	"reflect"
	"strings"
	"testing"
)

func TestPostgresUpdateBuilder_Basic(t *testing.T) {
	builder := NewUpdateBuilder("users").Set("name", "Tom").Where("id").Eq(1).BuildAsUpdate()
	sql := builder.SQL()
	expected := "UPDATE users SET name = $1 WHERE id = $2"
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	args := builder.Args()
	if !reflect.DeepEqual(args, []any{"Tom", 1}) {
		t.Errorf("expected args [Tom,1], got %v", args)
	}
}

func TestPostgresUpdateBuilder_SetRaw(t *testing.T) {
	builder := NewUpdateBuilder("users").SetRaw("updated_at", "now()")
	sql := builder.SQL()
	if !strings.Contains(sql, "updated_at = now()") {
		t.Errorf("expected raw set in sql, got %q", sql)
	}
}

func TestPostgresUpdateBuilder_SetBySlice(t *testing.T) {
	builder := NewUpdateBuilder("users").SetBySlice([]string{"name", "age"}, []any{"Tom", 18})
	sql := builder.SQL()
	if !strings.Contains(sql, "name = $1, age = $2") {
		t.Errorf("expected set by slice in sql, got %q", sql)
	}
}

func TestPostgresUpdateBuilder_SetByMap(t *testing.T) {
	builder := NewUpdateBuilder("users").SetByMap(map[string]any{"name": "Tom", "age": 18})
	sql := builder.SQL()
	if !strings.Contains(sql, "name = $1") || !strings.Contains(sql, "age = $2") {
		t.Errorf("expected set by map in sql, got %q", sql)
	}
}

func TestPostgresUpdateBuilder_Returning(t *testing.T) {
	builder := NewUpdateBuilder("users").Set("name", "Tom").Returning("id")
	sql := builder.SQL()
	if !strings.Contains(sql, "RETURNING id") {
		t.Errorf("expected returning in sql, got %q", sql)
	}
}

func TestPostgresUpdateBuilder_NoSet(t *testing.T) {
	builder := NewUpdateBuilder("users")
	sql := builder.SQL()
	expected := "UPDATE users"
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
}

func TestPostgresUpdateBuilder_ConditionSet(t *testing.T) {
	builder := NewUpdateBuilder("users").SetByCondition(false, "name", "Tom").SetByCondition(true, "age", 18)
	sql := builder.SQL()
	if !strings.Contains(sql, "age = $1") || strings.Contains(sql, "name = $2") {
		t.Errorf("expected only age set, got %q", sql)
	}
}
