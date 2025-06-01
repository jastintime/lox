package main

import (
	"fmt"
)

type Interpreter struct {
	environment Environment
}

func (i Interpreter) Interpret(statements []Stmt) {
	for _, statement := range statements {
		if hadRuntimeError {
			return
		}
		i.execute(statement)
	}
}

func (i Interpreter) VisitLiteralExpr(expr LiteralExpr) any {
	return expr.Value
}

func (i Interpreter) VisitLogicalExpr(expr LogicalExpr) any {
	left := i.evaluate(expr.Left)
	if expr.Operator.Type == Or {
		if i.isTruthy(left) {
			return left
		} else {
			if !i.isTruthy(left) {
				return left
			}
		}
	}
	return i.evaluate(expr.Right)
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

func (i Interpreter) execute(stmt Stmt) {
	stmt.Accept(i)
}

func (i Interpreter) executeBlock(statements []Stmt, environment Environment) {
	// NOTE: in java a finally was used, we could simply just
	// do i.environment = previous at the end of this function but what
	// if we panic somewhere?
	previous := i.environment
	defer func() {
		i.environment = previous
	}()

	i.environment = environment
	for _, statement := range statements {
		i.execute(statement)
	}

}

func (i Interpreter) VisitBlockStmt(stmt BlockStmt) any {
	i.executeBlock(stmt.Statements, newEnvironment(&i.environment))
	return nil
}

func (i Interpreter) VisitExprStmt(stmt ExprStmt) any {
	i.evaluate(stmt.Expression)
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
	fmt.Println(value)
	return nil
}

func (i Interpreter) VisitVariableStmt(stmt VariableStmt) any {
	var value any
	value = nil
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.environment.Define(stmt.Name.Lexeme, value)
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
	i.environment.Assign(expr.Name, value)
	return value
}

func (i Interpreter) VisitVariableExpr(expr VariableExpr) any {
	return i.environment.Get(expr.Name)
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
