package test

import (
	clause2 "GolangORM/clause"
	"reflect"
	"testing"
)

func testSelect(t *testing.T) {
	var clause clause2.Clause
	clause.Set(clause2.Limit, 3)
	clause.Set(clause2.Select, "User", []string{"*"})
	clause.Set(clause2.Where, "Name = ?", "Tom")
	clause.Set(clause2.OrderBy, "Age ASC")
	sql, vars := clause.Build(clause2.Select, clause2.Where, clause2.OrderBy, clause2.Limit)
	t.Log(sql, vars)
	if sql != "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testSelect(t)
	})
}
