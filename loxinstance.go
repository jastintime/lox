package main

import "maps"

type LoxInstance struct {
	Class  LoxClass
	fields map[string]any
}

func newLoxInstance(class LoxClass) LoxInstance {
	return LoxInstance{class, make(map[string]any)}
}

func (l LoxInstance) String() string {
	return l.Class.Name + " instance"
}

func (l LoxInstance) Get(name Token) any {
	value, ok := l.fields[name.Lexeme]
	if ok {
		return value
	}
	method, exist := l.Class.FindMethod(name.Lexeme)
	if exist {
		return method.Bind(l)
	}
	panic(RuntimeError{name, "Undefined property '" + name.Lexeme + "'."})
	return nil
}

func (l LoxInstance) Set(name Token, value any) {
	l.fields[name.Lexeme] = value
}

func (l LoxInstance) Equals(other LoxInstance) bool {
	return l.Class.Equals(other.Class) && maps.EqualFunc(l.fields, other.fields, isEqual)
}
