package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
}

func genBindVars(num int) string {
	// return: "?, ?, ?..."
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _insert(values ...interface{}) (string, []interface{}) {
	// return: "INSERT INTO 表名 (字段1, 字段2, 字段3...)", 空列表
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

func _values(values ...interface{}) (string, []interface{}) {
	// return: "VALUES (?, ?...), (?, ?...), ...", 参数列表
	var bindStr string
	var sql strings.Builder
	var vars []interface{}

	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func _select(values ...interface{}) (string, []interface{}) {
	// return: "SELECT 字段1, 字段2 FROM 表名", 空列表
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	// return: "LIMIT ?", 参数列表
	return "LIMIT ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	// return: "WHERE 表达式", 参数列表
	return fmt.Sprintf("WHERE %s", values[0]), values[1:]
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	// return: "ORDER BY 字段", 空列表
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}
