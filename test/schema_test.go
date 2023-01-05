package test

import (
	"fmt"
	"testing"
	"zORM/dialect"
	"zORM/schema"
)

func TestSchema(t *testing.T) {
	sqlite3Translator, _ := dialect.GetDialect("sqlite3")
	schema := schema.Parse(&User{}, sqlite3Translator)
	fmt.Println("Parse TableName:", schema.Name)
	fmt.Println("Parse Field Count:", len(schema.Fields))
	fmt.Println("Parse Name Tag:", schema.GetField("Name").Tag)
}
