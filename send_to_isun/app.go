package main

import (
	"fmt"
	"os"
)

func main() {
	// Чтение значения переменной среды
	dbUser := os.Getenv("DATE")
	fmt.Println("DATE:", dbUser)

	// Чтение значения аргумента командной строки
	args := os.Args
	if len(args) > 1 {
		configFile := args[1]
		fmt.Println("Config file:", configFile)
	}
}
