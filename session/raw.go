/*
封装了Session的基础方法

New方法:
新建一个Session

Clear方法:
清除Session中保存的SQL语句

DB方法:
返回事务句柄或DB句柄

RAW方法:
将SQL语句和对应的参数保存进Session

Exec方法:
封装了database/sql的具体执行

QueryRow方法：
封装了database/sql的单条查询

QueryRows方法:
封装了database/sql的多条查询
*/
package session

import (
	"database/sql"
	"strings"
	"zORM/clause"
	"zORM/dialect"
	"zORM/log"
	"zORM/schema"
)

/* 把tx和db封装成一个接口便于将判断方法抽出来 */
type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

type Session struct {
	db       *sql.DB         // database中sql.Open()返回的句柄
	tx       *sql.Tx         // database中执行事务的句柄
	dialect  dialect.Dialect // SQL解释器
	refTable *schema.Schema  // 通过dialect解析出来的数据库中的表
	sql      strings.Builder // 用户传入的带占位符的SQL语句
	sqlVars  []interface{}   // SQL语句中对应的值
	clause   clause.Clause   // 分句
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()              // 清空SQL
	s.sqlVars = nil            // 清空SQL参数列表
	s.clause = clause.Clause{} // 清空分句
}

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	} else {
		return s.db
	}
}

/* 传入SQL语句和参数，写入Session保存起来 */
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

/* 执行Session中的SQL语句 */
func (s *Session) Exec() (result sql.Result, err error) {
	// 每次执行完清空Session的变量，使Session可以复用，开启一次会话执行多次SQL
	defer s.Clear()
	// 打印[info ]等级日志
	log.Info(s.sql.String(), s.sqlVars)
	// 将用户编译好存在Session里的SQL语句交给dabase/sql包去具体执行
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		// 如果执行出错打印[error]等级日志
		log.Error(err)
	}
	return
}

/* 查询单条数据 */
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	// 将用户编译的SQL语句交给dabase/sql包的QueryRow()去具体执行
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

/* 查询多条数据 */
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	// 将用户编译的SQL语句交给dabase/sql包的Query()去具体执行
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
