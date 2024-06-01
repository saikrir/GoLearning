package main

import (
	"fmt"

	"github.com/saikrir/multi-init/util"
)

func init() {
	fmt.Println("Inside Main init")
}

func main() {
	fmt.Println("Say Hello ", util.SayHello())
}
