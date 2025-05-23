package main

import (
	"fmt"
)

type AstPrinter struct{}

func (a AstPrinter) print(expr Expr) string {
	return expr.Accept(a).(string)
}

func (a AstPrinter) VisitBinaryExpr(expr BinaryExpr) any {
	return a.Parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (a AstPrinter) VisitGroupingExpr(expr GroupingExpr) any {
	return a.Parenthesize("group", expr.Expression)
}

func (a AstPrinter) VisitLiteralExpr(expr LiteralExpr) any {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprint(expr.Value)
}

func (a AstPrinter) VisitUnaryExpr(expr UnaryExpr) any {
	return a.Parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (a AstPrinter) Parenthesize(name string, exprs ...Expr) string {
	var str string
	str += "(" + name
	for _, expr := range exprs {
		str += " "
		str += a.print(expr)
	}
	str += ")"
	return str
}
