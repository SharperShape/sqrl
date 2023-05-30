package sqrl

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// PlaceholderFormat is the interface that wraps the ReplacePlaceholders method.
//
// ReplacePlaceholders takes a SQL statement and replaces each question mark
// placeholder with a (possibly different) SQL placeholder.
type PlaceholderFormat interface {
	ReplacePlaceholders(sql string) (string, error)
}

var (
	// Question is a PlaceholderFormat instance that leaves placeholders as
	// question marks.
	Question = questionFormat{}

	// Dollar is a PlaceholderFormat instance that replaces placeholders with
	// dollar-prefixed positional placeholders (e.g. $1, $2, $3).
	Dollar = dollarFormat{}
)

type questionFormat struct{}

func (_ questionFormat) ReplacePlaceholders(sql string) (string, error) {
	return sql, nil
}

type dollarFormat struct{}

func (_ dollarFormat) ReplacePlaceholders(sql string) (string, error) {
	return replacePlaceholders(sql, func(buf *bytes.Buffer, i int) error {
		fmt.Fprintf(buf, "$%d", i)
		return nil
	})
}

func (_ dollarFormat) ReplacePlaceholdersMixed(sql string, args []interface{}) (string, []interface{}, error) {
	return replacePlaceholdersMixed(sql, args, func(buf *bytes.Buffer, i int) error {
		fmt.Fprintf(buf, "$%d", i)
		return nil
	})
}

// Placeholders returns a string with count ? placeholders joined with commas.
func Placeholders(count int) string {
	if count < 1 {
		return ""
	}

	return strings.Repeat(",?", count)[1:]
}

func replacePlaceholders(sql string, replace func(buf *bytes.Buffer, i int) error) (string, error) {
	buf := &bytes.Buffer{}
	i := 0
	for {
		p := strings.Index(sql, "?")
		if p == -1 {
			break
		}

		if len(sql[p:]) > 1 && sql[p:p+2] == "??" { // escape ?? => ?
			buf.WriteString(sql[:p])
			buf.WriteString("?")
			if len(sql[p:]) == 1 {
				break
			}
			sql = sql[p+2:]
		} else {
			i++
			buf.WriteString(sql[:p])
			if err := replace(buf, i); err != nil {
				return "", err
			}
			sql = sql[p+1:]
		}
	}

	buf.WriteString(sql)
	return buf.String(), nil
}

func replacePlaceholdersMixed(
	sql string, args []interface{}, replace func(buf *bytes.Buffer, i int) error,
) (string, []interface{}, error) {
	buf := &bytes.Buffer{}
	i := 0
	newArgs := make([]interface{}, 0, len(args))
	arg := 0
	var renumbered map[int]int
	for {
		p := strings.IndexAny(sql, "?$")
		if p == -1 {
			break
		}

		switch {
		case len(sql[p:]) > 1 && sql[p:p+2] == "??": // escape ?? => ?
			buf.WriteString(sql[:p])
			buf.WriteString("?")
			if len(sql[p:]) == 1 {
				break
			}
			sql = sql[p+2:]

		case sql[p:p+1] == "?":
			i++
			newArgs = append(newArgs, args[arg])
			arg++
			buf.WriteString(sql[:p])
			if err := replace(buf, i); err != nil {
				return "", nil, err
			}
			sql = sql[p+1:]

		case len(sql[p:]) > 1 && sql[p:p+2] == "$$": // escape $$ => $
			buf.WriteString(sql[:p])
			buf.WriteString("$")
			if len(sql[p:]) == 1 {
				break
			}
			sql = sql[p+2:]

		// If there are already some dollar placeholders, we renumber them,
		// but make sure that if a single argument is used in multiple places we preserve that.
		default: // $ placeholder
			end := p + 1
			for ; end < len(sql) && sql[end] >= '0' && sql[end] <= '9'; end++ {
			}
			if end == p+1 { // just a $, but not a placeholder
				buf.WriteString(sql[:p+1])
				sql = sql[p+1:]
			} else {
				num, err := strconv.Atoi(sql[p+1 : end])
				if err != nil {
					panic(err) // should never happen
				}
				j, ok := renumbered[num]
				if !ok {
					if renumbered == nil {
						renumbered = make(map[int]int)
					}
					i++
					renumbered[num] = i
					j = i
					newArgs = append(newArgs, args[arg])
				}
				arg++
				buf.WriteString(sql[:p])
				if err := replace(buf, j); err != nil {
					return "", nil, err
				}
				sql = sql[end:]
			}
		}
	}

	buf.WriteString(sql)
	return buf.String(), newArgs, nil
}
