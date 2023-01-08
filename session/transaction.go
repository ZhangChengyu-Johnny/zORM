/*
把database/sql中三个事务的基础方法封装到Session中
*/

package session

import "zORM/log"

/* 调用Session内的db创建tx，再赋值给Session的tx */
func (s *Session) Begin() (err error) {
	log.Info("transaction begin!")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return
	}
	return
}

func (s *Session) Commit() (err error) {
	log.Info("transaction commit")
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Rollback() (err error) {
	log.Info("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
	}
	return
}
