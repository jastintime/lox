package main

type ExprVisitor interface {
	VisitLiteralExpr(expr LiteralExpr) any
	VisitUnaryExpr(expr UnaryExpr) any
	VisitBinaryExpr(expr BinaryExpr) any
	VisitGroupingExpr(expr GroupingExpr) any
}

type Expr interface {
	Accept(visitor ExprVisitor) any
}

type LiteralExpr struct {
	Value any
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

type GroupingExpr struct {
	Expression Expr
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
