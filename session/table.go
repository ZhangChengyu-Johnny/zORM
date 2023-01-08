/*
封装了表的基础操作
*/

package session

import (
	"fmt"
	"reflect"
	"strings"
	"zORM/log"
	"zORM/schema"
)

/* 通过传入的Go对象解析成数据库中的表存入Session */
func (s *Session) Model(obj interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(obj) != reflect.TypeOf(s.refTable.Model) {
		// 解析操作是比较耗时的，所以每次解析前看一下refTable是否发生变化，如果不变就直接返回
		s.refTable = schema.Parse(obj, s.dialect)
	}
	return s
}

/* 返回Session里的表 */
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

/* 利用Session中保存的表信息创建表，之前必须调用过Session.Model(&obj{}) */
func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	// 示例: CREATE TABLE User (Name PRIMIADY KEY PRIMIADY KEY,Age  );
	// fmt.Printf("CREATE TABLE %s (%s);\n", table.Name, desc)
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

/* 利用Session中保存的表信息删除表，之前必须调用过Session.Model(&obj{}) */
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

/* 验证表是否存在 */
func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, values...).QueryRow()
	var tmp string
	row.Scan(&tmp)
	return tmp == s.RefTable().Name
}
