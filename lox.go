package main

import (
	"bufio"
	"fmt"
	"os"
)

// Global variable gross but better than OOP for a single variable.
var hadError = false
var hadRuntimeError = false

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
	interpreter := Interpreter{newEnvironment(nil)}
	run(sourceStr, interpreter)
	if hadError {
		os.Exit(65)
	}
	if hadRuntimeError {
		os.Exit(70)
	}
	return nil
}

func runPrompt() error {
	reader := bufio.NewReader(os.Stdin)
	interpreter := Interpreter{newEnvironment(nil)}
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		run(line, interpreter)
		hadError = false
		hadRuntimeError = false
	}
	return nil
}

func run(source string, interpreter Interpreter) {
	//	ast := AstPrinter{} // not printing the ast, we are interpreting now!
	scanner := newScanner(source)
	tokens := scanner.ScanTokens()
	parser := Parser{tokens, 0}
	statements := parser.Parse()
	if hadError {
		return
	}
	interpreter.Interpret(statements)

}
func emitError(line int, message string) {
	report(line, "", message)
}
func emitRuntimeError(operator Token, message string) {
	fmt.Println(operator.Type, message, "\n[line ", operator.Line, "]")
	hadRuntimeError = true
}

func emitTokenError(t Token, message string) {
	if t.Type == EOF {
		report(t.Line, " at end", message)
	} else {
		report(t.Line, " at '"+t.Lexeme+"'", message)
	}
}

func report(line int, where string, message string) {
	fmt.Println("[line ", line, "] Error", where, ": ", message)
	hadError = true

}
