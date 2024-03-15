package dialect

import (
	"GolangORM/log"
	"fmt"
	"reflect"
	"time"
)

type Sqlite struct {
}

func (s *Sqlite) DataTypeOf(value reflect.Value) string {
	switch value.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := value.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	log.Error(fmt.Sprintf("invalid sql type %s (%s)", value.Type().Name(), value.Kind()))
	panic(fmt.Sprintf("invalid sql type %s (%s)", value.Type().Name(), value.Kind()))
}

func (s *Sqlite) TableExistSql(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' AND name = ?", args
}

var _ Dialect = (*Sqlite)(nil)
