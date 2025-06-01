package main

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func newEnvironment(enclosing *Environment) Environment {
	return Environment{make(map[string]any), enclosing}
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e Environment) Get(name Token) any {
	value, ok := e.values[name.Lexeme]
	if ok {
		return value
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	emitRuntimeError(name, "Undefined variable'"+name.Lexeme+"'.")
	return nil
}

func (e *Environment) Assign(name Token, value any) {
	value, ok := e.values[name.Lexeme]
	if ok {
		e.values[name.Lexeme] = value
		return
	}
	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
		return
	}
	emitRuntimeError(name, "Undefined variable '"+name.Lexeme+"'.")
}
