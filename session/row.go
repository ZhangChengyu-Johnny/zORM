package session

import (
	"database/sql"
	"strings"
	"zORM/dialect"
	"zORM/log"
	"zORM/schema"
)

type Session struct {
	db       *sql.DB         // database中sql.Open()返回的句柄
	dialect  dialect.Dialect // SQL解释器
	refTable *schema.Schema  // 通过dialect解析出来的数据库中的表
	sql      strings.Builder // 用户传入的带占位符的SQL语句
	sqlVars  []interface{}   // SQL语句中对应的值
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

// 编译SQL语句
func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// 执行方法
func (s *Session) Exec() (result sql.Result, err error) {
	// 每次执行完清空Session的变量，使Session可以复用，开启一次会话执行多次SQL
	defer s.Clear()
	// 打印[info ]等级日志
	log.Info(s.sql.String(), s.sqlVars)
	// 将用户编译的SQL语句交给dabase/sql的db.Exec()去具体执行
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		// 如果执行出错打印[error]等级日志
		log.Error(err)
	}
	return
}

// 查询单条数据
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	// 将用户编译的SQL语句交给dabase/sql的db.QueryRow()去具体执行
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// 查询多条数据
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)

	// 将用户编译的SQL语句交给dabase/sql的db.Query()去具体执行
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}
