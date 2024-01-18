// A simple example with nested and embedded structs
// There are no arrays or slices, only pointers

package gsts

import (
	"fmt"
	"testing"
)

// Only basic data types
type Product struct {
	Name  string
	Price float64
}

// Nested struct
type Order struct {
	ID          int
	ProductAttr Product
}

// Embedded struct, pointer
type Customer struct {
	Name  string
	Email string
	*Order
}

func TestMain(t *testing.T) {

	var phone = Product{
		Name:  "CellPhone",
		Price: 42.50,
	}
	var order = &Order{
		ID:          4815162342,
		ProductAttr: phone,
	}
	var alice = &Customer{
		Name:  "Alice",
		Email: "alice@example.com",
		Order: order,
	}

	fieldNames, fieldTypes, n := GetAttrs(alice)
	fmt.Println("Attribute            Type")
	fmt.Println("----------------------------------")
	for i := 0; i < n; i++ {
		fmt.Printf("%-25v %-10s\n", fieldNames[i], fieldTypes[i])
	}

	query := GenerateCreationQuery(MakeDataTypes(), "Customer", fieldNames, fieldTypes)
	fmt.Print("\n", query)
}
