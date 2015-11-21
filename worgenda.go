package main

import (
	"fmt"
)

func main() {

	notes, err := ParseFile("./test.org")

	if err != nil {
		fmt.Println(err)
	}

	// for _, n := range notes {
	// 	fmt.Println(n.String())
	// }

	a := NewAgenda()
	a.InsertNotes(notes)
	a.Build()
	fmt.Println(a.String())
}
