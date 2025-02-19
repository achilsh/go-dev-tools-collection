package main

import "fmt"

func foo() *int {
	return nil
}
func bar() {
	fmt.Println(*foo()) // nilness reports NO error here, but NilAway does.
}
func main() {
	var xy *int
	
	fmt.Println(*xy)
	//
	bar()
}
