package main

import (
	"github.com/jastintime/lox/ClassType"
	"github.com/jastintime/lox/FunctionType"
)

type Resolver struct {
	interpreter     Interpreter
	scopes          Stack[map[string]bool]
	currentFunction functionType.FunctionType
	currentClass    classType.ClassType
}

func newResolver(interpreter Interpreter) Resolver {
	var stack = Stack[map[string]bool]{}
	return Resolver{interpreter, stack, functionType.None, classType.None}
}

func (r Resolver) resolve(a any) {
	switch arg := a.(type) {
	case []Stmt:
		for _, statement := range arg {
			r.resolve(statement)
		}
	case Stmt:
		arg.Accept(r)
	case Expr:
		arg.Accept(r)
	default:
		panic("unexpected resolve")
	}
}

func (r Resolver) VisitBlockStmt(stmt BlockStmt) any {
	r.beginScope()
	r.resolve(stmt.Statements)
	r.endScope()
	return nil
}

func (r Resolver) VisitClassStmt(stmt ClassStmt) any {
	enclosingClass := r.currentClass
	r.currentClass = classType.Class

	r.declare(stmt.Name)
	r.define(stmt.Name)

	if stmt.Superclass != nil && stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme {
		emitTokenError(stmt.Superclass.Name, "A class can't inherit from itself.")
	}

	if stmt.Superclass != nil {
		r.currentClass = classType.Subclass
		r.resolve(*stmt.Superclass)
	}

	if stmt.Superclass != nil {
		r.beginScope()
		r.scopes.Peek()["super"] = true
	}

	r.beginScope()
	r.scopes.Peek()["this"] = true

	for _, method := range stmt.Methods {
		declaration := functionType.Method
		if method.Name.Lexeme == "init" {
			declaration = functionType.Initializer
		}

		r.resolveFunction(method, declaration)
	}

	r.endScope()
	if stmt.Superclass != nil {
		r.endScope()
	}
	r.currentClass = enclosingClass
	return nil
}

func (r Resolver) VisitExprStmt(stmt ExprStmt) any {
	r.resolve(stmt.Expression)
	return nil
}

func (r Resolver) VisitFunctionStmt(stmt FunctionStmt) any {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	r.resolveFunction(stmt, functionType.Function)
	return nil
}

func (r Resolver) VisitIfStmt(stmt IfStmt) any {
	r.resolve(stmt.Condition)
	r.resolve(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolve(stmt.ElseBranch)
	}
	return nil
}

func (r Resolver) VisitPrintStmt(stmt PrintStmt) any {
	r.resolve(stmt.Expression)
	return nil
}

func (r Resolver) VisitReturnStmt(stmt ReturnStmt) any {
	if r.currentFunction == functionType.None {
		emitTokenError(stmt.Keyword, "Can't return from top-level code.")
	}
	if stmt.Value != nil {
		if r.currentFunction == functionType.Initializer {
			emitTokenError(stmt.Keyword, "Can't return a value from an initializer.")
		}

		r.resolve(stmt.Value)
	}
	return nil
}

func (r Resolver) VisitVariableStmt(stmt VariableStmt) any {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		r.resolve(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r Resolver) VisitWhileStmt(stmt WhileStmt) any {
	r.resolve(stmt.Condition)
	r.resolve(stmt.Body)
	return nil
}

func (r Resolver) VisitAssignExpr(expr AssignExpr) any {
	r.resolve(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r Resolver) VisitBinaryExpr(expr BinaryExpr) any {
	r.resolve(expr.Left)
	r.resolve(expr.Right)
	return nil
}

func (r Resolver) VisitCallExpr(expr CallExpr) any {
	r.resolve(expr.Callee)
	for _, argument := range expr.Arguments {
		r.resolve(argument)
	}
	return nil
}

func (r Resolver) VisitGetExpr(expr GetExpr) any {
	r.resolve(expr.Object)
	return nil
}

func (r Resolver) VisitGroupingExpr(expr GroupingExpr) any {
	r.resolve(expr.Expression)
	return nil
}

func (r Resolver) VisitLiteralExpr(expr LiteralExpr) any {
	return nil
}

func (r Resolver) VisitLogicalExpr(expr LogicalExpr) any {
	r.resolve(expr.Left)
	r.resolve(expr.Right)
	return nil
}

func (r Resolver) VisitSetExpr(expr SetExpr) any {
	r.resolve(expr.Value)
	r.resolve(expr.Object)
	return nil
}

func (r Resolver) VisitSuperExpr(expr SuperExpr) any {
	if r.currentClass == classType.None {
		emitTokenError(expr.Keyword, "Can't use 'super' outside of a class.")
	} else if r.currentClass != classType.Subclass {
		emitTokenError(expr.Keyword, "Can't use 'super' in a class with no superclass.")
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil
}

func (r Resolver) VisitThisExpr(expr ThisExpr) any {
	if r.currentClass == classType.None {
		emitTokenError(expr.Keyword, "Can't use 'this' outside of a class.")
		return nil
	}
	r.resolveLocal(expr, expr.Keyword)
	return nil
}

func (r Resolver) VisitUnaryExpr(expr UnaryExpr) any {
	r.resolve(expr.Right)
	return nil
}

func (r Resolver) VisitVariableExpr(expr VariableExpr) any {
	if !r.scopes.IsEmpty() {
		variable, inScope := r.scopes.Peek()[expr.Name.Lexeme]
		if inScope && variable == false {
			emitTokenError(expr.Name, "Can't read local variable in its own initializer.")
		}
	}
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r Resolver) resolveFunction(function FunctionStmt, t functionType.FunctionType) any {
	enclosingFunction := r.currentFunction
	r.currentFunction = t
	r.beginScope()
	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}
	r.resolve(function.Body)
	r.endScope()
	r.currentFunction = enclosingFunction
	return nil
}

func (r *Resolver) beginScope() {
	r.scopes.Push(make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(name Token) {
	if r.scopes.IsEmpty() {
		return
	}
	scope := r.scopes.Peek()
	_, dup := scope[name.Lexeme]
	if dup {
		emitTokenError(name, "Already a variable with this name in this scope.")
	}

	scope[name.Lexeme] = false
}

func (r *Resolver) define(name Token) {
	if r.scopes.IsEmpty() {
		return
	}
	r.scopes.Peek()[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := r.scopes.Size() - 1; i >= 0; i-- {
		_, ok := r.scopes.Get(i)[name.Lexeme]
		if ok {
			r.interpreter.Resolve(expr, r.scopes.Size()-1-i)
			return
		}
	}
}
