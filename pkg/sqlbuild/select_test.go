package sqlbuild

import (
	"reflect"
	"strings"
	"testing"
)

func TestPostgresSelectBuilder_Basic(t *testing.T) {
	builder := NewSelectBuilder("users").Select("id", "name").Where("id").Eq(1).BuildAsSelect().GroupBy("name").OrderBy("id").Limit(10).Offset(5)
	sql := builder.SQL()
	expected := "SELECT id, name FROM users WHERE id = $1 GROUP BY name ORDER BY id LIMIT 10 OFFSET 5"
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
	args := builder.Args()
	if !reflect.DeepEqual(args, []any{1}) {
		t.Errorf("expected args [1], got %v", args)
	}
}

func TestPostgresSelectBuilder_Empty(t *testing.T) {
	builder := NewSelectBuilder("users")
	sql := builder.SQL()
	expected := "SELECT * FROM users"
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
}

func TestPostgresSelectBuilder_Join(t *testing.T) {
	builder := NewSelectBuilder("users u").Select("u.id", "p.phone").LeftJoin("phones p").On("p.user_id").EqRaw("u.id").BuildAsSelect()
	sql := builder.SQL()
	if !strings.Contains(sql, "LEFT JOIN phones p ON p.user_id = u.id") {
		t.Errorf("expected join in sql, got %q", sql)
	}
}

func TestPostgresSelectBuilder_CountSql(t *testing.T) {
	builder := NewSelectBuilder("users").Where("id").Eq(1).BuildAsSelect()
	sql := builder.CountSQL()
	expected := "SELECT COUNT(*) as total FROM users WHERE id = $1"
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
}
