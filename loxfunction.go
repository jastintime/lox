package main

type LoxFunction struct {
	Declaration FunctionStmt
	Closure     Environment
}

func newLoxFunction(declaration FunctionStmt, closure Environment) LoxFunction {
	return LoxFunction{declaration, closure}
}

// BEAUTY
func (l LoxFunction) Call(interpreter Interpreter, arguments []any) (result any) {
	defer func() {
		result = recover()
	}()

	environment := newEnvironment(&l.Closure)
	for i := 0; i < len(l.Declaration.Params); i++ {
		environment.Define(l.Declaration.Params[i].Lexeme, arguments[i])
	}
	interpreter.executeBlock(l.Declaration.Body, environment)
	return nil
}

func (l LoxFunction) Arity() int {
	return len(l.Declaration.Params)
}

func (l LoxFunction) String() string {
	return "<fn " + l.Declaration.Name.Lexeme + ">"
}
