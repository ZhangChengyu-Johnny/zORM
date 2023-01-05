package session

import (
	"reflect"
	"zORM/clause"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		// 构造子句
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	// 用构造好的子句构造最终SQL语句
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	// 执行SQL语句，返回结果
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Find(obj interface{}) error {
	objValue := reflect.Indirect(reflect.ValueOf(obj))
	objType := objValue.Type().Elem()
	// 1. 根据obj找到obj原结构体，构造出一个空的对象
	// 2. 传入Model方法解析出本次查找的数据库表
	table := s.Model(reflect.New(objType).Elem().Interface()).RefTable()
	// 构造SELECT的SQL语句
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	// 执行SQL语句
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		// 遍历SQL语句的结果
		// 根据查询对象obj的原结构体构建一个新对象
		dest := reflect.New(objType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			// 把等下用来解析的空指针以Session.refTable.FieldNames的顺序推进values
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		// 传入空指针列表values进行单条记录的解析
		if err := rows.Scan(values...); err != nil {
			return err
		}
		objValue.Set(reflect.Append(objValue, dest))
	}

	return rows.Close()
}
