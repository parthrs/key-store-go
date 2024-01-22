package main

import (
	"bufio"
	"fmt"
	"key-store-go/pkg"
	"os"
	"strconv"
	"strings"
)

type KeyStore[K comparable, V any] interface {
	Set(K, V)
	Get(K) V
	Delete(K)
	Count() int
	Begin()
	End()
	Rollback()
	Commit()
}

var prompt string = "> "

func main() {
	reader := bufio.NewReader(os.Stdin)
	kv := pkg.NewStore[string, int]()

	for {
		fmt.Print(prompt)

		processedInput := ""
		incomplete := true
		bytes := []byte{}
		var err error

		for incomplete {
			bytes, incomplete, err = reader.ReadLine()
			if err != nil {
				fmt.Println("ERROR: unable to parse input line, try again.")
				break
			}
			processedInput += string(bytes)
		}
		args := strings.Split(processedInput, " ")
		switch {
		case args[0] == "SET":
			intInput, err := strconv.Atoi(args[2])
			if err != nil {
				fmt.Printf("ERROR: Unable to parse %s as int, try again\n", args[2])
			} else {
				kv.Set(args[1], intInput)
			}
		case args[0] == "GET":
			if len(args) < 2 {
				fmt.Printf("ERROR: Please provide an argument for GET (the key)\n")
			} else {
				val, found := kv.Get(args[1])
				if found {
					fmt.Printf("%d\n", val)
				}
			}
		case args[0] == "BEGIN":
			kv.Begin()
		case args[0] == "DELETE":
			if len(args) < 2 {
				fmt.Printf("ERROR: Please provide an argument for DELETE (the key)\n")
			} else {
				kv.Delete(args[1])
			}
		case args[0] == "COMMIT":
			kv.Commit()
		case args[0] == "COUNT":
			fmt.Printf("%d\n", kv.Count())
		case args[0] == "END":
			kv.End()
		case args[0] == "ROLLBACK":
			kv.Rollback()
		case args[0] == "EXIT":
			os.Exit(0)
		default:
			continue
		}
	}
}
