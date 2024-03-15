package schema

import (
	"GolangORM/dialect"
	"go/ast"
	"reflect"
)

type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema 用于解析制定的go 中的结构体，而 table.go用于生成数据库表指令
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:      dest,
		Name:       modelType.Name(),
		FieldNames: make([]string, 0),
		fieldMap:   make(map[string]*Field),
		Fields:     make([]*Field, 0),
	}
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if field.Anonymous || !ast.IsExported(field.Name) {
			continue
		}
		newField := &Field{
			Name: field.Name,
			Type: d.DataTypeOf(reflect.Indirect(reflect.New(field.Type))),
			Tag:  field.Tag.Get("orm"),
		}
		schema.Fields = append(schema.Fields, newField)
		schema.FieldNames = append(schema.FieldNames, field.Name)
		schema.fieldMap[field.Name] = newField
	}
	return schema
}

func (s *Schema) RecordFields(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var values []interface{}
	for _, field := range s.Fields {
		values = append(values, destValue.FieldByName(field.Name).Interface())
	}
	return values
}
