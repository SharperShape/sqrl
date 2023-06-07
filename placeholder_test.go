package sqrl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuestion(t *testing.T) {
	sql := "x = ? AND y = ?"
	s, _ := Question.ReplacePlaceholders(sql)
	assert.Equal(t, sql, s)
}

func TestDollar(t *testing.T) {
	sql := "x = ? AND y = ?"
	s, _ := Dollar.ReplacePlaceholders(sql)
	assert.Equal(t, "x = $1 AND y = $2", s)
}

func TestDollarMixed(t *testing.T) {
	sql := "x = ? AND y = $20 AND z = $1 AND w = ? AND m = ($20, $1) AND a = ?"
	args := []interface{}{"x", "y", "z", "w", "y", "z", "a"}
	s, args, _ := Dollar.ReplacePlaceholdersMixed(sql, args)
	assert.Equal(t, "x = $1 AND y = $2 AND z = $3 AND w = $4 AND m = ($2, $3) AND a = $5", s)
	assert.Equal(t, []interface{}{"x", "y", "z", "w", "a"}, args)
}

func TestPlaceholders(t *testing.T) {
	assert.Equal(t, Placeholders(2), "?,?")
}

func TestEscape(t *testing.T) {
	sql := "SELECT uuid, \"data\" #> '{tags}' AS tags FROM nodes WHERE  \"data\" -> 'tags' ??| array['?'] AND enabled = ?"
	s, _ := Dollar.ReplacePlaceholders(sql)
	assert.Equal(t, "SELECT uuid, \"data\" #> '{tags}' AS tags FROM nodes WHERE  \"data\" -> 'tags' ?| array['$1'] AND enabled = $2", s)

	sql = "SELECT uuid, \"data\" #> '{tags}' AS tags FROM nodes WHERE  \"data\" -> 'tags' $$| array['$4'] AND enabled = $1"
	s, _, _ = Dollar.ReplacePlaceholdersMixed(sql, []interface{}{1, 2})
	assert.Equal(t, "SELECT uuid, \"data\" #> '{tags}' AS tags FROM nodes WHERE  \"data\" -> 'tags' $| array['$1'] AND enabled = $2", s)
}

func BenchmarkPlaceholdersArray(b *testing.B) {
	var count = b.N
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "?"
	}
	var _ = strings.Join(placeholders, ",")
}

func BenchmarkPlaceholdersStrings(b *testing.B) {
	Placeholders(b.N)
}
