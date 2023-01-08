/*
通过反射的包的reflect.ValueOf().Kind()方法获取对象具体值的系统基础类型，再映射成数据库类型
*/
package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type sqlite3 struct{}

var _ Dialect = (*sqlite3)(nil)

// 初始化包时会把sqlite3注册到全局
func init() {
	RegisterDialect("sqlite3", &sqlite3{})
}

/* Go类型转SQL类型的方法 */
func (s *sqlite3) DataTypeOf(typ reflect.Value) string {
	// (reflect.Value).Kind() 和 (reflect.Type).Kind() 相同
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uintptr:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

/* 生成判断表是否存在的SQL语句 */
func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", args
}
