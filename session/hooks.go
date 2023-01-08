/*
CallMethod方法:
根据传入的钩子函数的名称，调用钩子函数
*/
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
	/*
		如果传入空的obj，就构造一个空的Session.Model实例，例如 &User{}，再根据名字找到这个实例中的钩子函数
		如果传入非空obj，例如 &User{"Jack", 18}，直接从这个实例中根据名字找到钩子函数
	*/

	fm := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if obj != nil {
		fm = reflect.ValueOf(obj).MethodByName(method)
	}
	// 组织参数列表，钩子函数约定好只有一个参数*Session，如:func (session *Session) error
	param := []reflect.Value{reflect.ValueOf(s)}
	if fm.IsValid() {
		// 检查方法是否有效
		if v := fm.Call(param); len(v) > 0 {
			// 调用钩子函数拿返回值，钩子函数约定只有一个返回值 error
			if err, ok := v[0].Interface().(error); ok {
				// 拿返回值转error类型，如果成功表示返回了一个error实例；如果失败表示返回的是nil
				log.Error(err)
			}
		}
	}
}
