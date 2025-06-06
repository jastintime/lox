package main

import (
	"fmt"
	"time"
)

type Interpreter struct {
	Environment Environment
	Globals     Environment
	Locals      map[Expr]int
}
type ReturnValue struct {
	Value any
}

func (r ReturnValue) Unbox() any {
	inside := r.Value
	switch v := inside.(type) {
	case ReturnValue:
		return v.Unbox()
	}
	return inside
}

func newInterpreter() Interpreter {
	globals := newEnvironment(nil)
	globals.Define("clock", Clock{})
	environment := globals
	return Interpreter{environment, globals, make(map[Expr]int)}
}

type Clock struct{}

func (c Clock) Arity() int {
	return 0
}
func (c Clock) Call(interpreter Interpreter, arguments []any) any {
	return float64(time.Now().UnixMilli()) / 1000.0
}

func (c Clock) String() string {
	return "<native fn>"
}

func (i Interpreter) Interpret(statements []Stmt) {
	defer func() {
		panicked := recover()
		runtimeError, ok := panicked.(RuntimeError)
		if ok {
			emitRuntimeError(runtimeError)
			return
		}
		if panicked != nil {
			panic(panicked)
		}
	}()
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i Interpreter) evaluate(expr Expr) any {
	return expr.Accept(i)
}

func (i Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
}

func (i Interpreter) Resolve(expr Expr, depth int) {
	i.Locals[expr] = depth
}

func (i Interpreter) executeBlock(statements []Stmt, environment Environment) {
	// NOTE: in java a finally was used, we could simply just
	// do i.Environment = previous at the end of this function but what
	// if we panic somewhere?
	previous := i.Environment
	defer func() {
		i.Environment = previous
	}()

	i.Environment = environment
	for _, statement := range statements {
		i.execute(statement)
	}

}

func (i Interpreter) VisitBlockStmt(stmt BlockStmt) any {
	i.executeBlock(stmt.Statements, newEnvironment(&i.Environment))
	return nil
}

func (i Interpreter) VisitClassStmt(stmt ClassStmt) any {
	var superclass any
	if stmt.Superclass != nil {
		superclass = i.evaluate(stmt.Superclass)
		_, ok := superclass.(LoxClass)
		if !ok {
			panic(RuntimeError{stmt.Superclass.Name, "Superclass must be a class."})
		}
	}
	i.Environment.Define(stmt.Name.Lexeme, nil)

	if stmt.Superclass != nil {
		old := i.Environment
		i.Environment = newEnvironment(&old)
		i.Environment.Define("super", superclass)
	}

	methods := make(map[string]LoxFunction)
	for _, method := range stmt.Methods {
		function := newLoxFunction(method, i.Environment, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = function
	}

	var class LoxClass
	s, ok := superclass.(LoxClass)
	if ok {
		class = LoxClass{stmt.Name.Lexeme, &s, methods}
	} else {
		class = LoxClass{stmt.Name.Lexeme, nil, methods}
	}

	if ok {
		i.Environment = *i.Environment.enclosing
	}
	i.Environment.Assign(stmt.Name, class)
	return nil
}

func (i Interpreter) VisitExprStmt(stmt ExprStmt) any {
	i.evaluate(stmt.Expression)
	return nil
}

func (i Interpreter) VisitFunctionStmt(stmt FunctionStmt) any {
	function := LoxFunction{stmt, i.Environment, false}
	i.Environment.Define(stmt.Name.Lexeme, function)
	return nil
}

func (i Interpreter) VisitIfStmt(stmt IfStmt) any {
	if i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i Interpreter) VisitPrintStmt(stmt PrintStmt) any {
	value := i.evaluate(stmt.Expression)
	v, ok := value.(ReturnValue)
	if ok {
		value = v.Unbox()
	}
	if value == nil {
		fmt.Println("nil")
	} else {
		fmt.Println(value)
	}
	return nil
}

func (i Interpreter) VisitReturnStmt(stmt ReturnStmt) any {
	var value any = nil
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}
	panic(ReturnValue{value})
}

func (i Interpreter) VisitVariableStmt(stmt VariableStmt) any {
	var value any
	value = nil
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.Environment.Define(stmt.Name.Lexeme, value)
	return nil
}

func (i Interpreter) VisitWhileStmt(stmt WhileStmt) any {
	for i.isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}
	return nil
}

func (i Interpreter) VisitAssignExpr(expr AssignExpr) any {
	value := i.evaluate(expr.Value)
	distance, ok := i.Locals[expr]
	//	fmt.Println("DISTANCE = ",distance)
	if ok {
		i.Environment.AssignAt(distance, expr.Name, value)
	} else {
		i.Globals.Assign(expr.Name, value)
	}
	return value
}

func (i Interpreter) VisitBinaryExpr(expr BinaryExpr) any {
	// NOTE: Scuffed, but anything that panics here is handled
	// just don't want to write a bunch of stuff to push it up.

	left := i.evaluate(expr.Left)
	//	fmt.Println("left is ", left)
	//	fmt.Println(expr.Operator.Type)
	right := i.evaluate(expr.Right)
	//	fmt.Println("right is ", right)
	//fmt.Println("Type of left is ", reflect.TypeOf(expr.Left))
	switch expr.Operator.Type {
	case Greater:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case GreaterEqual:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case Less:
		//fmt.Println("HERE")
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case LessEqual:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case Minus:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case BangEqual:
		return !isEqual(left, right)
	case EqualEqual:
		return isEqual(left, right)
	case Plus:
		leftValue, isLeftDouble := left.(float64)
		rightValue, isRightDouble := right.(float64)
		if isLeftDouble && isRightDouble {
			return leftValue + rightValue
		}
		leftString, isLeftString := left.(string)
		rightString, isRightString := right.(string)
		if isLeftString && isRightString {
			return leftString + rightString
		}
		panic(RuntimeError{expr.Operator, "Operands must be two numbers or two strings."})
	case Slash:
		i.checkNumberOperands(expr.Operator, left, right)
		if right.(float64) == 0 {
			// one of the challenges
			panic(RuntimeError{expr.Operator, "division by zero! panic"})
		}
		return left.(float64) / right.(float64)
	case Star:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)
	}

	// Unreachable
	return nil
}

func (i Interpreter) VisitCallExpr(expr CallExpr) any {
	callee := i.evaluate(expr.Callee)
	var arguments []any
	for _, argument := range expr.Arguments {
		arguments = append(arguments, i.evaluate(argument))
	}
	function, ok := callee.(LoxCallable)
	if !ok {
		panic(RuntimeError{expr.Paren, "Can only call functions and classes."})
		return nil
	}
	if len(arguments) != function.Arity() {
		panic(RuntimeError{expr.Paren, fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments))})
		return nil
	}
	return function.Call(i, arguments)
}

func (i Interpreter) VisitGetExpr(expr GetExpr) any {
	object := i.evaluate(expr.Object)
	value, ok := object.(LoxInstance)
	if ok {
		return value.Get(expr.Name)
	}
	panic(RuntimeError{expr.Name, "Only instances have properties."})
	return nil
}

func (i Interpreter) VisitGroupingExpr(expr GroupingExpr) any {
	return i.evaluate(expr.Expression)
}

func (i Interpreter) VisitLiteralExpr(expr LiteralExpr) any {
	return expr.Value
}

func (i Interpreter) VisitLogicalExpr(expr LogicalExpr) any {
	left := i.evaluate(expr.Left)
	if expr.Operator.Type == Or {
		if i.isTruthy(left) {
			return left
		}
	}
	if expr.Operator.Type == And {
		if !i.isTruthy(left) {
			return left
		}
	}
	return i.evaluate(expr.Right)
}

func (i Interpreter) VisitSetExpr(expr SetExpr) any {
	object := i.evaluate(expr.Object)
	o, ok := object.(LoxInstance)
	if !ok {
		panic(RuntimeError{expr.Name, "Only instances have fields."})
	}
	value := i.evaluate(expr.Value)
	o.Set(expr.Name, value)
	return value
}

func (i Interpreter) VisitSuperExpr(expr SuperExpr) any {
	distance := i.Locals[expr]
	superclass, _ := i.Environment.GetAt(distance, "super").(LoxClass)
	object, _ := i.Environment.GetAt(distance-1, "this").(LoxInstance)
	method, exist := superclass.FindMethod(expr.Method.Lexeme)
	if !exist {
		panic(RuntimeError{expr.Method, "Undefined property '" + expr.Method.Lexeme + "'."})
	}
	return method.Bind(object)
}

func (i Interpreter) VisitThisExpr(expr ThisExpr) any {
	return i.lookupVariable(expr.Keyword, expr)
}

func (i Interpreter) VisitUnaryExpr(expr UnaryExpr) any {
	right := i.evaluate(expr.Right)
	switch expr.Operator.Type {
	case Bang:
		return !i.isTruthy(right)
	case Minus:
		i.checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	}
	//unreachable
	return nil
}

func (i Interpreter) VisitVariableExpr(expr VariableExpr) any {
	//fmt.Println("im looking up", expr.Name)
	//fmt.Println("I GOT " ,i.lookupVariable(expr.Name, expr))
	return i.lookupVariable(expr.Name, expr)
}

func (i Interpreter) lookupVariable(name Token, expr Expr) any {
	distance, ok := i.Locals[expr]
	//fmt.Println("distance i got is ",i.Locals[expr])
	if ok {
		//	fmt.Println("I am okay")
		//	fmt.Println("got from the environment",i.Environment.GetAt(distance, name.Lexeme))
		return i.Environment.GetAt(distance, name.Lexeme)
	}
	return i.Globals.Get(name)

}

func (i Interpreter) checkNumberOperand(operator Token, operand any) {
	_, ok := operand.(float64)
	if ok {
		return
	}
	panic(RuntimeError{operator, "Operand must be a number."})
}

func (i Interpreter) checkNumberOperands(operator Token, left any, right any) {
	_, okleft := left.(float64)
	_, okright := right.(float64)
	if okleft && okright {
		return
	}
	//fmt.Println("operator where we crash is ", operator.Lexeme)
	//fmt.Println("type of left is ",reflect.TypeOf(left))
	//fmt.Println("left ", left,"right ", right)
	panic(RuntimeError{operator, "Operands must be numbers."})
	return
}

func (i Interpreter) isTruthy(object any) bool {
	if object == nil {
		return false
	}
	value, ok := object.(bool)
	if ok {
		return value
	}
	return true
}

func isEqual(a any, b any) (ret bool) {
	//NOTE: This function is neccesary because we can't compare maps in go by default
	// we need to use the maps.Equal function or else we panic
	defer func() {
		panicked := recover()
		if panicked != nil {
			class1, ok := a.(LoxClass)
			class2, ok2 := b.(LoxClass)
			//
			if ok && ok2 {
				ret = class1.Equals(class2)
				return
			}
			fun1, ok3 := a.(LoxFunction)
			fun2, ok4 := b.(LoxFunction)
			if ok3 && ok4 {
				ret = fun1.Equals(fun2)
				return
			}
			instance1, ok5 := a.(LoxInstance)
			instance2, ok6 := a.(LoxInstance)
			if ok5 && ok6 {
				ret = instance1.Equals(instance2)
				return
			}

			panic(panicked)
		}
	}()

	return a == b

}
