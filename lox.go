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
	interpreter := newInterpreter()
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
	interpreter := newInterpreter()
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
	resolver := newResolver(interpreter)
	resolver.resolve(statements)
	if hadError {
		return
	}
	interpreter.Interpret(statements)

}
func emitError(line int, message string) {
	report(line, "", message)
}
func emitRuntimeError(err RuntimeError) {
	fmt.Fprintln(os.Stderr, err.Message)
	fmt.Fprintf(os.Stderr, "[line %d]\n", err.Token.Line)
	hadRuntimeError = true
}

func emitTokenError(t Token, message string) {
	if t.Type == EOF {
		report(t.Line, " at end", message)
	} else {
		if t.Type == Identifier || t.Type == Number {
			report(t.Line, " at '"+t.Lexeme+"'", message)
			return
		}
		report(t.Line, " at '"+t.Type.String()+"'", message)
	}
}

func report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)

	hadError = true

}
