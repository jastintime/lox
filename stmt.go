package main

import "slices"

type StmtVisitor interface {
	VisitExprStmt(stmt ExprStmt) any
	VisitPrintStmt(stmt PrintStmt) any
	VisitVariableStmt(stmt VariableStmt) any
	VisitBlockStmt(stmt BlockStmt) any
	VisitIfStmt(stmt IfStmt) any
	VisitWhileStmt(stmt WhileStmt) any
	VisitFunctionStmt(stmt FunctionStmt) any
	VisitReturnStmt(stmt ReturnStmt) any
	VisitClassStmt(stmt ClassStmt) any
}

type Stmt interface {
	Accept(visitor StmtVisitor) any
}

type BlockStmt struct {
	Statements []Stmt
}

type ClassStmt struct {
	Name       Token
	Superclass *VariableExpr
	Methods    []FunctionStmt
}

type ExprStmt struct {
	Expression Expr
}

type FunctionStmt struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

func (f FunctionStmt) Equals(other FunctionStmt) bool {
	if f.Name != other.Name {
		return false
	}
	if !slices.Equal(f.Params, other.Params) {
		return false
	}
	if !slices.Equal(f.Body, other.Body) {
		return false
	}
	return true
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type PrintStmt struct {
	Expression Expr
}

type ReturnStmt struct {
	Keyword Token
	Value   Expr
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

func (b FunctionStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitFunctionStmt(b)
}

func (b ReturnStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitReturnStmt(b)
}
func (b ClassStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitClassStmt(b)
}
