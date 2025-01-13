package pkg

import "fmt"

func constCheckFunc() {
	testString1 := "test" // want "could be replaced by a constant 'test', repeated cont 5"
	testString2 := "test"
	testString3 := "test"
	testString4 := "test"
	testString5 := "test"

	testInt1 := 123 // want "could be replaced by a constant '123', repeated cont 5"
	testInt2 := 123
	testInt3 := 123
	testInt4 := 123
	testInt5 := 123

	fmt.Println(testString1, testString2, testString3, testString4, testString5)
	fmt.Println(testInt1, testInt2, testInt3, testInt4, testInt5)
}
