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

func TableExists(tableName string) bool {
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name=?", tableName)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		return true
	}
	return false
}

// Accepts pointer to an empty struct
func CreateTable(schema interface{}) {
	query = GenerateCreationQuery(schema)
	res, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
}

// Accepts pointer to a struct
func Insert(obj interface{}) (LastInsertId int64) {

	cNames, cTypes := GetAttrs(obj)
	n := len(cNames)

	// Generate INSERT query

	var (
		tableName  = reflect.TypeOf(obj).Elem().Name()
		columns    = strings.Join(cNames, ",")
		questMarks = strings.TrimSuffix(strings.Repeat("?,", n), ",")
	)
	query = fmt.Sprintf(
		"insert into %s (%s) values (%s)",
		tableName,
		columns,
		questMarks,
	)

	// Prepare INSERT args

	args := make([]interface{}, n)
	for i := 0; i < n; i++ {
		field := reflect.Indirect(reflect.ValueOf(obj)).FieldByName(cNames[i])
		switch cTypes[i] {
		case "int":
			args[i] = field.Int()
		case "string":
			args[i] = field.String()
		case "float64":
			args[i] = field.Float()
		}
	}

	res, err = db.Exec(query, args...)
	if err != nil {
		panic(err)
	}

	LastInsertId, err = res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return
}

func Read[T comparable]() T {

	var zero *T
	var typ reflect.Type = reflect.TypeOf(zero).Elem()
	elem := reflect.New(typ).Elem()
	columnNames, columnTypes := GetAttrs(elem.Interface())
	nColumns := len(columnNames)

	// Execute SELECT
	query = fmt.Sprintf("select * from %s where ID = ?", typ.Name())
	rows, _ := db.Query(query, 42)

	for rows.Next() {

		// Create nColumns pointers to corresponding types
		var pointers = make([]interface{}, nColumns)
		for i := 0; i < nColumns; i++ {
			switch columnTypes[i] {
			case "int":
				pointers[i] = new(int)
			case "string":
				pointers[i] = new(string)
			case "float64":
				pointers[i] = new(float64)
			}
		}
		err = rows.Scan(pointers...)
		if err != nil {
			panic(err)
		}

		for i := 0; i < nColumns; i++ {
			field := elem.FieldByName(columnNames[i])
			var ptr = pointers[i]
			switch columnTypes[i] {
			case "int":
				ptr, _ := ptr.(*int)
				field.SetInt(int64(*ptr))
			case "string":
				ptr, _ := ptr.(*string)
				field.SetString(*ptr)
			case "float64":
				ptr, _ := ptr.(*float64)
				field.SetFloat(*ptr)
			}
		}
	}
	return elem.Interface().(T)
}
