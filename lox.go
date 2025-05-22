package main

import (
	"bufio"
	"fmt"
	"os"
)

// Global variable gross but better than OOP for a single variable.
var hadError = false

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage : golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		err := runFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
		}
	} else {
		err := runPrompt()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func runFile(path string) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	sourceStr := string(source)
	run(sourceStr)
	if hadError {
		os.Exit(65)
	}
	return nil
}

func runPrompt() error {
	reader := bufio.NewReader(os.Stdin)
	ast := AstPrinter{}
	ast.main()
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		run(line)
		hadError = false
	}
	return nil
}

func run(source string) {
	//for _, token := range source {
	//	fmt.Printf("%c\n", token)
	//}
	scanner := newScanner(source)
	scanner.ScanTokens()

}
func emitError(line int, message string) {
	report(line, "", message)
}
func report(line int, where string, message string) {
	fmt.Println("[line ", line, "] Error", where, ": ", message)
	hadError = true

}
