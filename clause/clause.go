package clause

import "strings"

type Type int

const (
	Select Type = iota
	Insert
	Update
	Delete
	Values
	Limit
	OrderBy
	Where
	Count
)

type Clause struct {
	sql     map[Type]string
	sqlArgs map[Type][]interface{}
}

func (c *Clause) Set(name Type, args ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlArgs = make(map[Type][]interface{})
	}
	c.sql[name], c.sqlArgs[name] = generatorMap[name](args...)
}

func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqlSlice []string
	var args []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqlSlice = append(sqlSlice, sql)
			args = append(args, c.sqlArgs[order]...)
		}
	}
	return strings.Join(sqlSlice, " "), args
}
