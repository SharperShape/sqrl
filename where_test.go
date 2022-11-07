package sqrl

import (
	"testing"

	"bytes"

	"github.com/stretchr/testify/assert"
)

func TestWherePartsAppendToSql(t *testing.T) {
	parts := []Sqlizer{
		NewWherePart("x = ?", 1),
		NewWherePart(nil),
		NewWherePart(Eq{"y": 2}),
	}
	sql := &bytes.Buffer{}
	args, _ := appendToSql(parts, sql, " AND ", []interface{}{})
	assert.Equal(t, "x = ? AND y = ?", sql.String())
	assert.Equal(t, []interface{}{1, 2}, args)
}

func TestWherePartsAppendToSqlErr(t *testing.T) {
	parts := []Sqlizer{NewWherePart(1)}
	_, err := appendToSql(parts, &bytes.Buffer{}, "", []interface{}{})
	assert.Error(t, err)
}

func TestWherePartNil(t *testing.T) {
	sql, _, _ := NewWherePart(nil).ToSql()
	assert.Equal(t, "", sql)
}

func TestWherePartErr(t *testing.T) {
	_, _, err := NewWherePart(1).ToSql()
	assert.Error(t, err)
}

func TestWherePartString(t *testing.T) {
	sql, args, _ := NewWherePart("x = ?", 1).ToSql()
	assert.Equal(t, "x = ?", sql)
	assert.Equal(t, []interface{}{1}, args)
}

func TestWherePartMap(t *testing.T) {
	test := func(pred interface{}) {
		sql, _, _ := NewWherePart(pred).ToSql()
		expect := []string{"x = ? AND y = ?", "y = ? AND x = ?"}
		if sql != expect[0] && sql != expect[1] {
			t.Errorf("expected one of %#v, got %#v", expect, sql)
		}
	}
	m := map[string]interface{}{"x": 1, "y": 2}
	test(m)
	test(Eq(m))
}
