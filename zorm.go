package zORM

import (
	"database/sql"
	"zORM/log"
	"zORM/session"
)

type Engine struct {
	db *sql.DB
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}

	// ping一下测试
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	e = &Engine{db: db}
	log.Info("Connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to closed database")
	}
	log.Info("Close database success")
}

func (enging *Engine) NewSession() *session.Session {
	return session.New(enging.db)
}
