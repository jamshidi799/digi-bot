package main

import "fmt"

func main() {
	x := []string{"a", "b", "c"}
	var y []string
	for _, value := range x {
		y = append(y, value)
	}
	fmt.Print(y)
}
