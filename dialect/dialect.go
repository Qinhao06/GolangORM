package dialect

import "reflect"

type Dialect interface {
	DataTypeOf(value reflect.Value) string
	TableExistSql(tableName string) (string, []interface{})
}

var dialects = map[string]Dialect{}

func Register(name string, d Dialect) {
	dialects[name] = d
}

func GetDialect(name string) (Dialect, bool) {
	dialect, ok := dialects[name]
	return dialect, ok
}
