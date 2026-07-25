package main

import (
	"fmt"
	"runtime"
)

func main() {
	var a string = "Hello, World!"
	fmt.Println(a)

	primes := [6]int{2, 3, 5, 7, 11, 13}
	fmt.Println(primes)

	var s []int = primes[1:4]
	fmt.Println(s)

	m := make(map[string]int)
	m["Answer"] = 42
	fmt.Println("The value:", m["Answer"])

	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}

	fmt.Println("The sum is:", sum)
	fmt.Println("After sum")
	sumWhile := 1
	fmt.Println("after sumWhile")
	for sumWhile < 10 {
		sumWhile += sumWhile
	}
	fmt.Println("after sumWhile loop")
	fmt.Println("The sum is:", sumWhile)

	sumForRange := 0
	for i, value := range primes {
		sumForRange += value
		fmt.Println("Index:", i, "Value:", value)
	}
	fmt.Println("The sum is:", sumForRange)

	// for {
	// 	fmt.Println("This will loop forever")
	// }
	numberIfTest := 7
	if (numberIfTest % 2) == 0 {
		fmt.Println(numberIfTest, "is even")
	} else {
		fmt.Println(numberIfTest, "is odd")
	}

	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		fmt.Printf("%s.", os)
	}

}
