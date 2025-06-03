package main

type LoxClass struct {
	Name    string
	Methods map[string]LoxFunction
}

func (l LoxClass) String() string {
	return l.Name
}

func (l LoxClass) Call(interpreter Interpreter, arguments []any) any {
	instance := newLoxInstance(l)
	initializer, exist := l.FindMethod("init")
	if exist {
		initializer.Bind(instance).Call(interpreter, arguments)
	}

	return instance
}

func (l LoxClass) Arity() int {
	initializer, exist := l.FindMethod("init")
	if !exist {
		return 0
	}
	return initializer.Arity()
}

func (l LoxClass) FindMethod(name string) (LoxFunction, bool) {
	value, ok := l.Methods[name]
	if ok {
		return value, true
	}
	return LoxFunction{}, false
}
