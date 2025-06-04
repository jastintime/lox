package main

import "fmt"

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func newEnvironment(enclosing *Environment) Environment {
	return Environment{make(map[string]any), enclosing}
}

func (e Environment) Get(name Token) any {
	value, ok := e.values[name.Lexeme]
	if ok {
		return value
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	panic(RuntimeError{name, "Undefined variable '" + name.Lexeme + "'."})
	return nil
}

func (e *Environment) Assign(name Token, value any) {
	_, ok := e.values[name.Lexeme]
	if ok {
		e.values[name.Lexeme] = value
		return
	}
	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
		return
	}
	panic(RuntimeError{name, "Undefined variable '" + name.Lexeme + "'."})
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e Environment) ancestor(distance int) Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = *environment.enclosing
	}
	return environment
}

func (e Environment) GetAt(distance int, name string) any {
	return e.ancestor(distance).values[name]
}

func (e *Environment) AssignAt(distance int, name Token, value any) {
	e.ancestor(distance).values[name.Lexeme] = value
}

func (e Environment) String() string {
	result := fmt.Sprintf("%v", e.values)
	if e.enclosing != nil {
		result += " -> " + e.enclosing.String()
	}
	return result

}
