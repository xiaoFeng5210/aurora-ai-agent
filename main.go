package main

import (
	"fmt"
	"reflect"
)

func main() {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	x := User{
		Name: "John",
		Age:  30,
	}

	f, _ := reflect.TypeOf(x).FieldByName("Name")
	fmt.Println(f.Tag.Get("json"))
}
