package session

import (
	"errors"
	"reflect"
	"zORM/clause"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	// 插入前调用钩子
	for _, v := range values {
		s.CallMethod(BeforeInsert, v)
	}
	// values: &user1{"Tim", 18}, &user2{"Tom", 25}...
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		// 构造子句
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		// table.RecordValues(&user1{"Time", 18}) = [Tom, 18]
		recordValues = append(recordValues, table.RecordValues(value))
	}
	// 用构造好的子句构造最终SQL语句
	// recordValues = [[Tom, 18], [Sam, 25]]
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	// 执行SQL语句，返回结果
	result, err := s.Raw(sql, vars...).Exec()
	// 插入后调用钩子
	for _, v := range values {
		s.CallMethod(AfterInsert, v)
	}
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Find(obj interface{}) error {
	// 查询前调用钩子
	s.CallMethod(BeforeQuery, nil)
	objValue := reflect.Indirect(reflect.ValueOf(obj))
	objType := objValue.Type().Elem()
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
		// 查询后调用钩子，操作每一行
		s.CallMethod(AfterQuery, dest.Addr().Interface())
		objValue.Set(reflect.Append(objValue, dest))
	}
	return rows.Close()
}

func (s *Session) Update(kv ...interface{}) (int64, error) {
	// 更新前调用钩子
	s.CallMethod(BeforeUpdate, nil)
	// 参数可以接收一个字典{字段1:值1, 字段2:值2}，
	// 也可以接收一个列表[字段1，值1，字段2，值2]
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	// 构造SQL语句
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	// 执行并返回结果
	result, err := s.Raw(sql, vars...).Exec()
	// 更新后调用钩子
	s.CallMethod(AfterUpdate, nil)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	// 删除前调用钩子
	s.CallMethod(BeforeDelete, nil)
	// 构造SQL语句
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	// 删除后调用钩子
	s.CallMethod(AfterDelete, nil)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	// desc: "name = ?", args: "Tom"
	var vars []interface{}
	vars = append(vars, desc)
	vars = append(vars, args...)
	// vars: ["name = ?", "Tom"]
	s.clause.Set(clause.WHERE, vars...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

func (s *Session) First(value interface{}) error {
	// 查询前调用钩子
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
