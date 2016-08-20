package main

/*
 */
import "C"
import "unsafe"

import "fmt"

//
//export get_module
func get_module() unsafe.Pointer {
	fmt.Println("hhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhhh")
	return unsafe.Pointer(uintptr(5))
}

func main() {}
