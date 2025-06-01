package main

type StmtVisitor interface {
	VisitExprStmt(stmt ExprStmt) any
	VisitPrintStmt(stmt PrintStmt) any
	VisitVariableStmt(stmt VariableStmt) any
	VisitBlockStmt(stmt BlockStmt) any
	VisitIfStmt(stmt IfStmt) any
	VisitWhileStmt(stmt WhileStmt) any
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

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type PrintStmt struct {
	Expression Expr
}

type VariableStmt struct {
	Name        Token
	Initializer Expr
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
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
func (b IfStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitIfStmt(b)
}

func (b WhileStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitWhileStmt(b)
}
