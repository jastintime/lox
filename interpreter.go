package main

import (
	"fmt"
	"strconv"
)

type Interpreter struct{}

func (i Interpreter) Interpret(expression Expr) {
	value := i.evaluate(expression)
	if hadRuntimeError {
		return
	}
	v, ok := value.(float64)
	if ok {
		// remove trailing .0
		value = strconv.FormatFloat(v, 'f', -1, 64)
	}
	fmt.Println(value)
}

func (i Interpreter) VisitLiteralExpr(expr LiteralExpr) any {
	return expr.Value
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

func (i Interpreter) checkNumberOperand(operator Token, operand any) {
	_, ok := operand.(float64)
	if ok {
		return
	}
	emitRuntimeError(operator, "Operand must be a number.")
}

func (i Interpreter) checkNumberOperands(operator Token, left any, right any) bool {
	_, okleft := left.(float64)
	_, okright := right.(float64)
	if okleft && okright {
		return true
	}
	emitRuntimeError(operator, "Operands must be numbers.")
	return false
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

func (i Interpreter) VisitGroupingExpr(expr GroupingExpr) any {
	return i.evaluate(expr.Expression)
}

func (i Interpreter) evaluate(expr Expr) any {
	return expr.Accept(i)
}

func (i Interpreter) VisitBinaryExpr(expr BinaryExpr) any {
	// NOTE: Scuffed, but anything that panics here is handled
	// just don't want to write a bunch of stuff to push it up.
	defer func() {
		if recover() != nil {
		}
	}()

	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)
	switch expr.Operator.Type {
	case Greater:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case GreaterEqual:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case Less:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case LessEqual:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case Minus:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)
	case BangEqual:
		// NOTE: no isEqual, golang handles this correctly
		return !(left == right)
	case EqualEqual:
		return left == right
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
		emitRuntimeError(expr.Operator, "Operands must be two numbers or two strings.")
	case Slash:
		i.checkNumberOperands(expr.Operator, left, right)
		if right.(float64) == 0 {
			// one of the challenges
			emitRuntimeError(expr.Operator, "division by zero! panic")
		}
		return left.(float64) / right.(float64)
	case Star:
		i.checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)
	}
	// Unreachable
	return nil
}
