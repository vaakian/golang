package main

import (
	"flag"
	"fmt"
)

func main() {
	var name = flag.String("name", "Bee", "name for project")
	var age = flag.Int("age", 21, "age for owner")

	flag.Parse()

	fmt.Printf("name=%s\nage=%d\n", *name, *age)
}
