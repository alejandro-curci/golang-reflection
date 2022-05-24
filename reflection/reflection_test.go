package reflection

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite
}

func TestReflection(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (t testSuite) Test_Concepts() {
	p := Person{
		Name: "Bob",
		Age:  35,
	}

	pType := reflect.TypeOf(p)
	pName := pType.Name()
	pKind := pType.Kind()
	pValue := reflect.ValueOf(p)

	t.Equal("reflection.Person", fmt.Sprintf("%v", pType))
	t.Equal("Person", fmt.Sprintf("%v", pName))
	t.Equal("struct", fmt.Sprintf("%v", pKind))
	t.Equal("{Bob 35}", fmt.Sprintf("%v", pValue))
}

func (t testSuite) Test_UsefulMethods() {
	p := &Person{
		Name: "Bob",
		Age:  35,
	}

	pType := reflect.TypeOf(p)
	pKind := pType.Kind()
	pValue := reflect.ValueOf(p)

	t.Equal("*reflection.Person", fmt.Sprintf("%v", pType))
	t.Equal("ptr", fmt.Sprintf("%v", pKind))
	t.Equal("&{Bob 35}", fmt.Sprintf("%v", pValue))

	pElem := pValue.Elem()
	t.Equal("{Bob 35}", fmt.Sprintf("%v", pElem))

	pFields := pElem.NumField()
	pAgeField := pElem.Type().Field(1).Name
	pAgeType := pElem.Type().Field(1).Type
	pAgeValue := pElem.Field(1)
	t.Equal("2", fmt.Sprintf("%v", pFields))
	t.Equal("Age", fmt.Sprintf("%v", pAgeField))
	t.Equal("int", fmt.Sprintf("%v", pAgeType))
	t.Equal("35", fmt.Sprintf("%v", pAgeValue))
}

func (t testSuite) Test_SetNewValue() {
	modify := func(data interface{}) {
		v := reflect.ValueOf(data)
		elem := v.Elem() // will panic if data is not a pointer
		elem.Set(reflect.ValueOf(Person{"Will", 54}))
	}

	safeModify := func(data interface{}) {
		v := reflect.ValueOf(data)
		if reflect.TypeOf(data).Kind() == reflect.Ptr {
			elem := v.Elem()
			elem.Set(reflect.ValueOf(Person{"Will", 54}))
		}
	}

	p := Person{"Bob", 35}
	t.Equal("{Bob 35}", fmt.Sprintf("%v", p))
	t.Panics(func() {
		modify(p)
	})
	safeModify(&p)
	t.Equal("{Will 54}", fmt.Sprintf("%v", p))
}

func (t testSuite) Test_MakeNewInstance() {
	create := func(pointer interface{}) interface{} {
		dataElem := reflect.ValueOf(pointer).Elem()
		newPointer := reflect.New(dataElem.Type())
		newPointer.Elem().Set(reflect.ValueOf(Person{"Will", 54}))
		return newPointer.Interface()
	}
	p := &Person{"Bob", 35}
	newP := create(p)
	t.Equal("&{Bob 35}", fmt.Sprintf("%v", p))
	t.Equal("&{Will 54}", fmt.Sprintf("%v", newP))
}

func (t testSuite) Test_CreateQuery() {
	e := Employee{
		ID:       12,
		Name:     "Tom",
		Position: "Technical Leader",
		Country:  "South Africa",
		Salary:   23401910,
	}
	q, err := CreateQuery(e)
	t.Equal("insert into Employee values(12, \"Tom\", \"Technical Leader\", \"South Africa\", 23401910)", q)
	t.Nil(err)

	p := Person{
		Name: "Samantha",
		Age:  29,
	}
	q, err = CreateQuery(p)
	t.Equal("insert into Person values(\"Samantha\", 29)", q)
	t.Nil(err)
}

// BENCHMARKING

func (t testSuite) Test_SummingTwo_Timing() {
	var num int64 = 5

	start1 := time.Now()
	SumTwo(&num)
	ms1 := time.Since(start1).Nanoseconds()
	fmt.Println("first sum: ", ms1)
	t.Equal(int64(7), num)

	start2 := time.Now()
	SumTwoWithReflection(&num)
	ms2 := time.Since(start2).Nanoseconds()
	fmt.Println("second sum: ", ms2)
	t.Equal(int64(9), num)
}

func BenchmarkSumTwo(b *testing.B) {
	var num int64 = 10
	for i := 0; i < b.N; i++ {
		SumTwo(&num)
	}
}

func BenchmarkSumTwoWithReflection(b *testing.B) {
	var num int64 = 10
	for i := 0; i < b.N; i++ {
		SumTwoWithReflection(&num)
	}
}
