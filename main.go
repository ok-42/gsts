// Generate SQLite query for a Go struct

package gsts

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

var (
	db    *sql.DB
	err   error
	query string
	res   sql.Result
)

var dataTypes = new(map[string]string)

func Init() {
	db, err = sql.Open("sqlite", ":memory:")
	if err != nil {
		panic(err)
	}
	*dataTypes = MakeDataTypes()
}

// Mapping of Go and SQLite data types
func MakeDataTypes() map[string]string {
	data := map[string]string{
		"string":  "TEXT",
		"int":     "INTEGER",
		"float64": "REAL",
	}
	return data
}

// Returns struct attributes and their types respectively
func getAttrs(
	obj interface{},
	prefix string,
	index *int,
	outNames []string,
	outTypes []string,
) {
	var t reflect.Type = reflect.TypeOf(obj)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	var n int = t.NumField()

	for i := 0; i < n; i++ {

		// Get i-th field of the struct
		var field reflect.StructField = t.Field(i)

		// Set type and kind
		var typ reflect.Type = field.Type
		var kind = typ.Kind()

		fieldType := field.Type

		// If pointer, dereference it and update type and kind
		if kind == reflect.Pointer {
			fieldType = typ.Elem()
			kind = fieldType.Kind()
		}

		// Extract fields from an embedded or nested struct
		if kind == reflect.Struct {
			var newPrefix string
			if prefix == "" {
				newPrefix = ""
			} else {
				newPrefix = prefix
			}
			getAttrs(
				reflect.New(fieldType).Elem().Interface(),
				newPrefix+field.Name+"_",
				index,
				outNames,
				outTypes,
			)
		} else

		// Add field name and type to the result
		{
			outNames[*index] = prefix + field.Name
			outTypes[*index] = typ.Name()
			*index++
		}
	}
}

func GetAttrs(obj interface{}) ([]string, []string) {
	var n int = 20
	fieldNames := make([]string, n)
	fieldTypes := make([]string, n)
	getAttrs(obj, "", new(int), fieldNames, fieldTypes)
	for i, v := range fieldNames {
		if v == "" {
			n = i
			break
		}
	}
	fieldNames = fieldNames[:n]
	fieldTypes = fieldTypes[:n]
	return fieldNames, fieldTypes
}

// Generate a CREATE TABLE query
func GenerateCreationQuery(obj interface{}) string {
	fieldNames, fieldTypes := GetAttrs(obj)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("create table %s (\n", reflect.TypeOf(obj).Elem().Name()))
	for i, name := range fieldNames {
		goType := fieldTypes[i]
		sqlType := (*dataTypes)[goType]
		sb.WriteString(fmt.Sprintf("    %s %s,\n", name, sqlType))
	}
	return strings.TrimSuffix(sb.String(), ",\n") + "\n)\n"
}
