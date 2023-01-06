package test

import (
	"fmt"
	"reflect"
	"testing"
)

type TestStruct struct {
	Field1 string `zorm:"column:field1;type:string"`
	Field2 int8
}

func TestTag(t *testing.T) {
	modelType := reflect.Indirect(reflect.ValueOf(&TestStruct{})).Type()
	fmt.Println(modelType)
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		tag, ok := p.Tag.Lookup("zorm")
		if ok {
			fmt.Printf("字段名:%v；字段类型:%v；字段TAG:%v\n", p.Name, p.Type, tag)
		} else {
			fmt.Printf("字段名:%v；字段类型:%v；字段TAG:%v\n", p.Name, p.Type, ok)
		}

	}
}
