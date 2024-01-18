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
func GetAttrs(obj interface{}) ([]string, []string, int) {
	var t reflect.Type = reflect.TypeOf(obj)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	var n int = t.NumField()

	var outNames = make([]string, n)
	var outTypes = make([]string, n)
	var innerNames []string
	var innerTypes []string

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
			outNames[i] = "(nested)"
			outTypes[i] = "(nested)"
			resNames, resTypes, nn := GetAttrs(reflect.New(fieldType).Elem().Interface())
			for i := 0; i < nn; i++ {
				resNames[i] = field.Name + "_" + resNames[i]
			}
			innerNames = append(innerNames, resNames...)
			innerTypes = append(innerTypes, resTypes...)
		} else

		// Add field name and type to the result
		{
			outNames[i] = field.Name
			outTypes[i] = typ.Name()
		}
	}

	outNames = append(outNames, innerNames...)
	outTypes = append(outTypes, innerTypes...)
	return outNames, outTypes, len(outNames)
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
