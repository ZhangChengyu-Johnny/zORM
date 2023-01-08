/*
将一个数据库的表抽象成Schema对象，将表中的列抽象成Field对象

GetField方法:
获取表中的列(Schema里的Field)

Parse方法:
根据一个Go的实例化对象(指针类型、非指针类型都可)，获取他的类型
*/

package schema

import (
	"go/ast"
	"reflect"
	"zORM/dialect"
)

// 数据库中的列
type Field struct {
	Name string
	Type string
	Tag  string
}

// 数据库中的表
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

/* 把任何Go对象解析成Schema */
func Parse(obj interface{}, translator dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(obj)).Type()
	schema := &Schema{
		Model:    obj,
		Name:     modelType.Name(), // obj原结构体名
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			// 验证字段是否公有
			field := &Field{
				Name: p.Name,
				Type: translator.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("zorm"); ok {
				// zorm:"PRIMIARY KEY", Lookup会提取 v = "PRIMIARY KEY"
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

/* 在Insert前需要做一个字段转换，根据数据库中列的顺序，从对象里找到对应的值 */
func (schema *Schema) RecordValues(obj interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(obj))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
