// A simple example with nested and embedded structs
// There are no arrays or slices, only pointers

package gsts

import (
	"fmt"
	"testing"

	_ "modernc.org/sqlite"
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
	Field string
}

type User struct {
	ID    int
	Name  string
	Score float64
}

func (u *User) String() string {
	return fmt.Sprintf("<ID=%d %s has score %f>", u.ID, u.Name, u.Score)
}

func TestMain(t *testing.T) {
	Init()
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
		Field: "string",
	}
	query = GenerateCreationQuery(alice)
	fmt.Print("\n", query)
}

func TestCreateTable(t *testing.T) {
	Init()
	CreateTable(&Product{})
}

func TestInsert(t *testing.T) {
	var phone = &Product{
		Name:  "CellPhone",
		Price: 42.50,
	}
	Init()
	CreateTable(&Product{})
	Insert(phone)
}

func TestRead(t *testing.T) {
	Init()
	CreateTable(&User{})
	var alice = User{ID: 42, Name: "Alice", Score: 12.34}
	Insert(&alice)

	var result User = Read[User]()
	fmt.Println(result)
}
