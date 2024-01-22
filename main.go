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

// main is the entry point when running this
// program in standalone mode
func main() {
	reader := bufio.NewReader(os.Stdin)
	kv := pkg.NewStore[string, int]()

	// Perpetual loop till EXIT
	// Print the prompt and then wait for
	// user command
	for {
		fmt.Print(prompt)

		processedInput := ""
		var (
			incomplete bool = true
			bytes      []byte
			err        error
		)

		// reader.ReadLine may not read full line if it exceeds
		// buffer, thus we create a loop to read it fully
		// (this is possibly an overkill for our use-case)
		for incomplete {
			bytes, incomplete, err = reader.ReadLine()
			if err != nil {
				fmt.Println("ERROR: unable to parse input line, try again.")
				break
			}
			processedInput += string(bytes)
		}

		args := strings.Split(processedInput, " ")

		// The switch case is not setup to check args[0]
		// to later incorporate listening for os.INTERRUPT
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
