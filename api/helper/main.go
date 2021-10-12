package main

import "fmt"

func main() {
	fmt.Println("Hello, What's your favorite number?")
	var i int
	fmt.Scanf("%d\n", &i)
	fmt.Println("Ah I like ", i, " too.")
}
