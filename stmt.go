package main

type StmtVisitor interface {
	VisitExprStmt(stmt ExprStmt) any
	VisitPrintStmt(stmt PrintStmt) any
	VisitVariableStmt(stmt VariableStmt) any
	VisitBlockStmt(stmt BlockStmt) any
}

type Stmt interface {
	Accept(visitor StmtVisitor) any
}

type BlockStmt struct {
	Statements []Stmt
}

type ExprStmt struct {
	Expression Expr
}

type PrintStmt struct {
	Expression Expr
}

type VariableStmt struct {
	Name        Token
	Initializer Expr
}

func (b ExprStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitExprStmt(b)
}
func (b PrintStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitPrintStmt(b)
}
func (b VariableStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitVariableStmt(b)
}
func (b BlockStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitBlockStmt(b)
}
