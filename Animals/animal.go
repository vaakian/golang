package Animals

import "fmt"

type Animal struct {
	Name string
}

func (ani *Animal) Eat(food string) {
	fmt.Println("i am " + ani.Name + ", eating " + food)
}

type Dog struct {
	*Animal
}

func (dog *Dog) Run() {
	fmt.Println("Dog is running!")
}
