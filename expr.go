package main

type ExprVisitor interface {
	VisitLiteralExpr(expr LiteralExpr) any
	VisitAssignExpr(expr AssignExpr) any
	VisitUnaryExpr(expr UnaryExpr) any
	VisitBinaryExpr(expr BinaryExpr) any
	VisitGroupingExpr(expr GroupingExpr) any
	VisitVariableExpr(expr VariableExpr) any
	VisitLogicalExpr(expr LogicalExpr) any
	VisitCallExpr(expr CallExpr) any
	VisitGetExpr(expr GetExpr) any
	VisitSetExpr(Expr SetExpr) any
	VisitThisExpr(Expr ThisExpr) any
	VisitSuperExpr(Expr SuperExpr) any
}

type Expr interface {
	Accept(visitor ExprVisitor) any
}

type LiteralExpr struct {
	Value any
}

type LogicalExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

type SetExpr struct {
	Object Expr
	Name   Token
	Value  Expr
}

type SuperExpr struct {
	Keyword Token
	Method  Token
}

type ThisExpr struct {
	Keyword Token
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

type AssignExpr struct {
	Name  Token
	Value Expr
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

type CallExpr struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

type GetExpr struct {
	Object Expr
	Name   Token
}

type GroupingExpr struct {
	Expression Expr
}

type VariableExpr struct {
	Name Token
}

func (b LiteralExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteralExpr(b)
}
func (b UnaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnaryExpr(b)
}

func (b BinaryExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinaryExpr(b)
}

func (b GroupingExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitGroupingExpr(b)
}
func (b VariableExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitVariableExpr(b)
}
func (b AssignExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitAssignExpr(b)
}
func (b LogicalExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLogicalExpr(b)
}
func (b CallExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitCallExpr(b)
}
func (b GetExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitGetExpr(b)
}
func (b SetExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitSetExpr(b)
}
func (b ThisExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitThisExpr(b)
}
func (b SuperExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitSuperExpr(b)
}
