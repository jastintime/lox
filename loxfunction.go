package main

type LoxFunction struct {
	Declaration   FunctionStmt
	Closure       Environment
	isInitializer bool
}

func (l *LoxFunction) Bind(instance LoxInstance) LoxFunction {
	environment := newEnvironment(&l.Closure)
	environment.Define("this", instance)
	return newLoxFunction(l.Declaration, environment, l.isInitializer)
}

func newLoxFunction(declaration FunctionStmt, closure Environment, isInitializer bool) LoxFunction {
	return LoxFunction{declaration, closure, isInitializer}
}

// BEAUTY
func (l LoxFunction) Call(interpreter Interpreter, arguments []any) (result any) {
	defer func() {
		result = recover()
		if l.isInitializer {
			result = l.Closure.GetAt(0, "this")
		}
	}()

	environment := newEnvironment(&l.Closure)
	for i := 0; i < len(l.Declaration.Params); i++ {
		environment.Define(l.Declaration.Params[i].Lexeme, arguments[i])
	}
	interpreter.executeBlock(l.Declaration.Body, environment)
	if l.isInitializer {
		l.Closure.GetAt(0, "this")
	}
	return nil
}

func (l LoxFunction) Arity() int {
	return len(l.Declaration.Params)
}

func (l LoxFunction) String() string {
	return "<fn " + l.Declaration.Name.Lexeme + ">"
}
