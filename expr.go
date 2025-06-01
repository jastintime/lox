package main

type ExprVisitor interface {
	VisitLiteralExpr(expr LiteralExpr) any
	VisitAssignExpr(expr AssignExpr) any
	VisitUnaryExpr(expr UnaryExpr) any
	VisitBinaryExpr(expr BinaryExpr) any
	VisitGroupingExpr(expr GroupingExpr) any
	VisitVariableExpr(expr VariableExpr) any
	VisitLogicalExpr(expr LogicalExpr) any
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
