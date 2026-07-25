package main

import (
	"fmt"
	"log"
	"os"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (p *Person) Greet() {
	fmt.Printf("Hello, my name is %s and I am %d years old.\n", p.Name, p.Age)
}

func (p Person) OpenFile() (bool, error) {
	fmt.Printf("%s is opening a file!\n", p.Name)
	// _, err := os.OpenFile("err.txt", os.O_APPEND, 0666)
	_, err := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		// log.Fatal(err)
		return false, err
	}
	return true, nil
}

func main() {
	p := Person{Name: "Daniel", Age: 26}
	p.Greet()
	f, err := p.OpenFile()
	fmt.Printf("File opened: %v\n", f)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Default().Printf("File opened successfully: %v", f)
	}

	// defer f.Close()
}
