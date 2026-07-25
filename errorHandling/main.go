package main

import (
	"log"
	"os"
)

func main() {
	f, err := os.OpenFile("test.txt", os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.Default().Println("File opened successfully")

	n, err := f.WriteString("Hello, Daniel!!\n")
	if err != nil {
		log.Fatal(err)
	}
	log.Default().Printf("Wrote %d bytes to file", n)

	// defer f.Close()
}
