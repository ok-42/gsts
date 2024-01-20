// Generate SQLite query for a Go struct

package gsts

import (
	"fmt"
	"reflect"
	"strings"
)

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
func GetAttrs(
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
			GetAttrs(
				reflect.New(fieldType).Elem().Interface(),
				field.Name+"_",
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

// Generate a CREATE TABLE query
func GenerateCreationQuery(
	data map[string]string,
	name string,
	columnNames []string,
	columnTypes []string,
) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("create table %s (\n", name))
	n := len(columnNames)
	for i := 0; i < n; i++ {
		name := columnNames[i]
		goType := columnTypes[i]
		sqlType := data[goType]
		sb.WriteString(fmt.Sprintf("    %s %s,\n", name, sqlType))
	}
	sb.WriteString(")\n")
	return sb.String()
}

func main() {

}
