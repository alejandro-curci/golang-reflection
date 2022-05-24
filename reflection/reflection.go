package reflection

import (
	"errors"
	"fmt"
	"reflect"
)

// ***** WHAT IS REFLECTION ? *****
// It is the ability of a program to inspect its variables and values at run time.
// Go implements reflection through the "reflect" package in the standard library.

// ***** WHY DO YOU NEED IT ? *****
// In golang, each and every variable in our program is defined by us,
// so we know its type at compile time itself. Sadly, this is not always true.
// Sometimes you want to work with variables at runtime using information
// that did not exist when the program was written.
// Reflection gives you the ability to examine types at runtime.

// ***** COMMON USE CASES *****
// Marshall and Unmarshall from the "json" package
// Generic SQL query creator

// ***** TYPES IN REFLECT PACKAGE *****
// It is built around 3 concepts -> Types, Kinds and Values

type Person struct {
	Name string
	Age  int
}

// Type -> Person (actually, it is reflection.Person, the package is also part of a type)
// Kind -> struct
// Value -> {"Bob", 38}

// ***** USEFUL METHODS *****

// Name() -> returns the type's name within its package
// Elem() -> dereferences a pointer
// Type() -> returns the type of a value
// NumField() -> returns the number of fields in a struct
// Field(i int) -> returns the reflect.Value of the ith field
// Int(), String(), Bool() -> returns the underlying value as an int64, string, bool, respectively

// ***** HOW TO MODIFY A VALUE ? *****

// You can set a new value following these steps:
// 1- Get the reflected value of the variable with reflect.ValueOf()
// 2- Dereference the value with Elem()
// 3- Set a new reflected value with Set()
// 4- You can return an interface converting to reflect type with Interface()
// This will only work when passing pointers as interface{}
// Passing a copy will not work, you can read it but not modify it.

// ***** HOW TO MAKE A NEW INSTANCE ? *****

// You need to create a new pointer from an existing pointer, so you follow these steps:
// 1- Dereference the pointer
// 2- Create a new pointer with reflect.New(), passing the element's type
// 3- Set a reflect value to the new pointer
// 4- Return the new pointer as interface{} with Interface()

// ***** DOWNSIDES *****
// 1- Performance penalty -> check benchmark examples
// 2- More complicated code -> "Clear is better than clever. Reflection is never clear.", Rob Pike
// 3- Unsafe methods -> everything panics when it is not properly used

// EXAMPLES FOR BENCHMARKING

func SumTwo(num *int64) {
	*num += 2
}

func SumTwoWithReflection(num interface{}) {
	n := reflect.ValueOf(num).Elem()
	nInt := n.Int()
	n.Set(reflect.ValueOf(nInt + 2))
}

// EXAMPLE OF QUERY GENERATOR

type Employee struct {
	ID       int
	Name     string
	Position string
	Country  string
	Salary   int
}

func CreateQuery(q interface{}) (string, error) {
	if reflect.ValueOf(q).Kind() == reflect.Struct {
		t := reflect.TypeOf(q).Name()
		query := fmt.Sprintf("insert into %s values(", t)
		v := reflect.ValueOf(q)
		for i := 0; i < v.NumField(); i++ {
			switch v.Field(i).Kind() {
			case reflect.Int:
				if i == 0 {
					query = fmt.Sprintf("%s%d", query, v.Field(i).Int())
				} else {
					query = fmt.Sprintf("%s, %d", query, v.Field(i).Int())
				}
			case reflect.String:
				if i == 0 {
					query = fmt.Sprintf("%s\"%s\"", query, v.Field(i).String())
				} else {
					query = fmt.Sprintf("%s, \"%s\"", query, v.Field(i).String())
				}
			default:
				return "", errors.New("unsupported type")
			}
		}
		query = fmt.Sprintf("%s)", query)
		return query, nil

	}
	return "", errors.New("unsupported type")
}
