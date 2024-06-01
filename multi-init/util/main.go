package util

import "fmt"

func init() {
	fmt.Println("Iniside util init1")
}

func init() {
	fmt.Println("Iniside util init2")
}

func SayHello() string {
	return "Hello World"
}
