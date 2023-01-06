package clause

import (
	"strings"
)

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

type Clause struct {
	sql     map[Type]string        // 分句
	sqlVars map[Type][]interface{} // 分句中的变量
}

/* 根据传入的Type组织出对应分句和参数 */
func (c *Clause) Set(name Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
		c.sqlVars = make(map[Type][]interface{})
	}
	// 根据传入操作Type和参数，调用对应方法构造分句和参数列表
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

/* 根据分句顺序，组织出最终SQL语句和参数列表 */
func (c *Clause) Build(typeOrders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range typeOrders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
