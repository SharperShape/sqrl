package sqrl

import "fmt"

type WherePart part

func NewWherePart(pred interface{}, args ...interface{}) Sqlizer {
	return &WherePart{pred: pred, args: args}
}

func (p WherePart) Pred() interface{} {
	return p.pred
}

func (p WherePart) Args() []interface{} {
	return p.args
}

func (p WherePart) ToSql() (sql string, args []interface{}, err error) {
	switch pred := p.pred.(type) {
	case nil:
		// no-op
	case Sqlizer:
		return pred.ToSql()
	case map[string]interface{}:
		return Eq(pred).ToSql()
	case string:
		sql = pred
		args = p.args
	default:
		err = fmt.Errorf("expected string-keyed map or string, not %T", pred)
	}
	return
}
