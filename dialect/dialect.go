/*
RegisterDialect方法:
不同数据库在完成自己的实例后通过该方法把数据库解释器注册到全局中

GetDialect方法:
根据数据库名获取对应的数据库解释器
*/
package dialect

import "reflect"

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	// Go类型转SQL类型的方法
	DataTypeOf(typ reflect.Value) string
	// 生成判断表是否存在的SQL语句
	TableExistSQL(tableName string) (string, []interface{})
}

// 用于注册dialect实例。如果新增对某个数据库的支持，那么调用该方法注册到全局
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// 用于获取dialect实例
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
