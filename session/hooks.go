package session

import (
	"reflect"
	"zORM/log"
)

const (
	BeforeQuery = "BeforeQuery"
	AfterQuery  = "AfterQuery"

	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"

	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"

	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
)

func (s *Session) CallMethod(method string, obj interface{}) {
	// 如果传入空的obj，就构造一个空的Session.Model对象传递给钩子
	// 如果传入非空obj，那么将obj直接传递给钩子函数
	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if obj != nil {
		fm = reflect.ValueOf(obj).MethodByName(method)
	}
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		// 传入Session对象，把钩子函数调用起来(钩子函数: func (session *Session) error )
		if v := fm.Call(param); len(v) > 0 {
			// 正常情况钩子函数返回 <nil>，ok为false
			if err, ok := v[0].Interface().(error); ok {
				// 查看调用结果看是否出错
				log.Error(err)
			}
		}
	}
}
