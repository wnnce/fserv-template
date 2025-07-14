package sqlbuild

import (
	"reflect"
	"testing"
)

func TestPostgresDeleteBuilder_Basic(t *testing.T) {
	builder := NewDeleteBuilder("users").Where("id").Eq(1).BuildAsDelete()
	sql := builder.SQL()
	expected := "DELETE FROM users WHERE id = $1"
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	args := builder.Args()
	if !reflect.DeepEqual(args, []any{1}) {
		t.Errorf("expected args [1], got %v", args)
	}
}

func TestPostgresDeleteBuilder_NoWhere(t *testing.T) {
	builder := NewDeleteBuilder("users")
	sql := builder.SQL()
	expected := "DELETE FROM users"
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
}

func TestPostgresDeleteBuilder_MultiWhere(t *testing.T) {
	builder := NewDeleteBuilder("users").Where("id").Eq(1).And("name").Eq("Tom").BuildAsDelete()
	sql := builder.SQL()
	expected := "DELETE FROM users WHERE id = $1 AND name = $2"
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	args := builder.Args()
	if !reflect.DeepEqual(args, []any{1, "Tom"}) {
		t.Errorf("expected args [1,Tom], got %v", args)
	}
}
