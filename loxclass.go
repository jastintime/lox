package main

type LoxClass struct {
	Name       string
	Superclass *LoxClass
	Methods    map[string]LoxFunction
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
	if l.Superclass != nil {
		return l.Superclass.FindMethod(name)
	}
	return LoxFunction{}, false
}

func (l LoxClass) Equals(other LoxClass) bool {
	if (l.Name == other.Name && l.Superclass == other.Superclass) == false {
		return false
	}
	for k, v := range l.Methods {
		ov, ok := other.Methods[k]
		if !ok {
			return false
		}

		if !v.Equals(ov) {
			return false
		}
	}
	return true
}
