package test

import (
	"fmt"
	"testing"
	c "zORM/clause"
)

func testSelect() {
	var clause c.Clause
	clause.Set(c.LIMIT, 3)
	clause.Set(c.SELECT, "User", []string{"*"})
	clause.Set(c.WHERE, "Name = ?", "Tom")
	clause.Set(c.ORDERBY, "Age ASC")
	sql, vars := clause.Build(c.SELECT, c.WHERE, c.ORDERBY, c.LIMIT)
	fmt.Println(sql, vars)
}

func TestClauseBuild(t *testing.T) {
	testSelect()
}
