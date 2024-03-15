package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generatorMap map[Type]generator

func init() {
	generatorMap = make(map[Type]generator)
	generatorMap[Select] = _select
	generatorMap[Insert] = _insert
	generatorMap[Update] = _update
	generatorMap[Delete] = _delete
	generatorMap[Values] = _values
	generatorMap[Where] = _where
	generatorMap[Limit] = _limit
	generatorMap[OrderBy] = _orderBy
	generatorMap[Count] = _count
}

func getVarsPlaceholder(num int) string {
	return strings.Repeat("?, ", num-1) + "?"
}

func _select(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	filedNames := strings.Join(values[1].([]string), ", ")
	return fmt.Sprintf("SELECT %s FROM %s", filedNames, tableName), []interface{}{}
}

func _insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	filedNames := strings.Join(values[1].([]string), ", ")
	return fmt.Sprintf("INSERT INTO %s(%s)", tableName, filedNames), []interface{}{}
}

func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	args := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k, v := range args {
		keys = append(keys, fmt.Sprintf("%s = ?", k))
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}

func _delete(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	return fmt.Sprintf("DELETE FROM %s", tableName), []interface{}{}
}

// 对每一个变量进行遍历，生成一个？的 Values 子句
func _values(values ...interface{}) (string, []interface{}) {
	var sqlBuilder strings.Builder
	var bindStr string
	var vars []interface{}
	sqlBuilder.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = getVarsPlaceholder(len(v))
		}
		sqlBuilder.WriteString(fmt.Sprintf("(%s)", bindStr))
		if i < len(values)-1 {
			sqlBuilder.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sqlBuilder.String(), vars
}

// 传入 condition
func _where(values ...interface{}) (string, []interface{}) {
	desc, args := values[0].(string), values[1:]
	return fmt.Sprintf("WHERE %s", desc), args
}

func _limit(values ...interface{}) (string, []interface{}) {
	return "LIMIT ?", values
}

// 传入 fieldName
func _orderBy(values ...interface{}) (string, []interface{}) {
	fieldName := values[0].(string)
	return fmt.Sprintf("ORDER BY %s", fieldName), []interface{}{}
}

func _count(values ...interface{}) (string, []interface{}) {
	tableName := values[0].(string)
	return fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName), []interface{}{}
}
