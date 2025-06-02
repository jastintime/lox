package main

type Parser struct {
	Tokens  []*Token
	current int
}

func newParser(tokens []*Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) declaration() Stmt {
	fmatch := p.match(Fun)
	if hadError {
		p.synchronize()
		return nil
	}
	if fmatch {
		return p.function("function")
	}

	vmatch := p.match(Var)
	if hadError {
		p.synchronize()
		return nil
	}
	if vmatch {
		return p.varDeclaration()
	}
	smatch := p.statement()
	if hadError {
		p.synchronize()
		return nil
	}
	return smatch
}

func (p *Parser) statement() Stmt {
	if p.match(Print) {
		return p.printStatement()
	}
	if p.match(Return) {
		return p.returnStatement()
	}
	if p.match(While) {
		return p.WhileStatement()
	}
	if p.match(LeftBrace) {
		return BlockStmt{p.block()}
	}
	if p.match(For) {
		return p.forStatement()
	}
	if p.match(If) {
		return p.ifStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() Stmt {
	p.consume(LeftParen, "Expect '(' after 'for'.")

	var initializer Stmt
	if p.match(Semicolon) {
		initializer = nil
	} else if p.match(Var) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr = nil
	if !p.check(Semicolon) {
		condition = p.expression()
	}
	p.consume(Semicolon, "Expect ';' after loop condition")

	var increment Expr = nil
	if !p.check(RightParen) {
		increment = p.expression()
	}
	p.consume(RightParen, "Expect ')' after for clauses.")
	body := p.statement()
	if increment != nil {
		body = BlockStmt{[]Stmt{body, ExprStmt{increment}}}
	}

	if condition == nil {
		condition = LiteralExpr{true}
	}
	body = WhileStmt{condition, body}
	if initializer != nil {
		body = BlockStmt{[]Stmt{initializer, body}}
	}

	return body
}

func (p *Parser) ifStatement() Stmt {
	p.consume(LeftParen, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(RightParen, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch Stmt = nil
	if p.match(Else) {
		elseBranch = p.statement()
	}
	return IfStmt{condition, thenBranch, elseBranch}
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(Semicolon, "Expected ';' after value.")
	return PrintStmt{value}
}

func (p *Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr = nil
	if !p.check(Semicolon) {
		value = p.expression()
	}
	p.consume(Semicolon, "Expect ';' after return value.")
	return ReturnStmt{keyword, value}
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(Identifier, "Expect variable name.")
	var initializer Expr
	initializer = nil
	if p.match(Equal) {
		initializer = p.expression()
	}
	p.consume(Semicolon, "Expect ';' after variable declaration.")
	return VariableStmt{name, initializer}
}

func (p *Parser) WhileStatement() Stmt {
	p.consume(LeftParen, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(RightParen, "Expect ')' after 'while'.")
	body := p.statement()
	return WhileStmt{condition, body}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(Semicolon, "Expect ';' after expression.")
	return ExprStmt{expr}
}

func (p *Parser) function(kind string) FunctionStmt {
	name := p.consume(Identifier, "Expect"+kind+" name.")
	p.consume(LeftParen, "Expect '(' after "+kind+"name.")
	var parameters []Token
	if !p.check(RightParen) {
		parameters = append(parameters, p.consume(Identifier, "Expect parameter name."))
		for p.match(Comma) {
			// NOTE: once again because of difference with do while
			if len(parameters) >= 254 {
				emitTokenError(p.peek(), "Can't have more than 255 parameters.")
			}
			parameters = append(parameters, p.consume(Identifier, "Expect parameter name."))
		}
	}
	p.consume(RightParen, "Expect ')' after parameters. ")
	p.consume(LeftBrace, "Expect '{' before "+kind+" body.")
	body := p.block()
	return FunctionStmt{name, parameters, body}

}

func (p *Parser) block() []Stmt {
	var statements []Stmt
	for !p.check(RightBrace) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(RightBrace, "Expected '}' after block.")
	return statements
}

func (p *Parser) assignment() Expr {
	expr := p.or()
	if p.match(Equal) {
		equals := p.previous()
		value := p.assignment()
		expr, ok := expr.(VariableExpr)
		if ok {
			name := expr.Name
			return AssignExpr{name, value}
		}
		emitTokenError(equals, "Invalid assignment target.")
	}
	return expr

}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(Or) {
		operator := p.previous()
		right := p.and()
		expr = LogicalExpr{expr, operator, right}
	}
	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(And) {
		operator := p.previous()
		right := p.equality()
		expr = LogicalExpr{expr, operator, right}
	}
	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()
	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = BinaryExpr{expr, operator, right}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right := p.term()
		expr = BinaryExpr{expr, operator, right}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(Minus, Plus) {
		operator := p.previous()
		right := p.factor()
		expr = BinaryExpr{expr, operator, right}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(Slash, Star) {
		operator := p.previous()
		right := p.unary()
		expr = BinaryExpr{expr, operator, right}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(Bang, Minus) {
		operator := p.previous()
		right := p.unary()
		return UnaryExpr{operator, right}
	}
	return p.call()
}

func (p *Parser) finishCall(callee Expr) Expr {
	var arguments []Expr
	if !p.check(RightParen) {
		p.match(Comma)
		arguments = append(arguments, p.expression())
		for p.match(Comma) {
			arguments = append(arguments, p.expression())
			//NOTE: in the java version we do 255 but here we aren't doing a do while so its 254
			if len(arguments) >= 254 {
				emitTokenError(p.peek(), "Can't have more than 255 arguments.")
			}
		}
	}
	paren := p.consume(RightParen, "Expect ')' after arguments.")
	return CallExpr{callee, paren, arguments}

}

func (p *Parser) call() Expr {
	expr := p.primary()

	for {
		if p.match(LeftParen) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) primary() Expr {
	if p.match(False) {
		return LiteralExpr{false}
	}
	if p.match(True) {
		return LiteralExpr{true}
	}
	if p.match(Nil) {
		return LiteralExpr{nil}
	}
	if p.match(Number, String) {
		return LiteralExpr{p.previous().Literal}
	}
	if p.match(Identifier) {
		return VariableExpr{p.previous()}
	}
	if p.match(LeftParen) {
		expr := p.expression()
		p.consume(RightParen, "Expect ')' after expression.")
		return GroupingExpr{expr}
	}
	// NOTE: should never arrive here as this is used internally and only by unary which only calls us
	// for the matched statements.
	emitTokenError(p.peek(), "Expect expression.")

	return nil

}

func (p *Parser) Parse() []Stmt {
	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements

}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(t TokenType, message string) Token {
	if p.check(t) {
		return p.advance()
	}
	emitTokenError(p.peek(), message)
	return Token{}

}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == Semicolon {
			return
		}
		switch p.peek().Type {
		case Class | Fun | Var | For | If | While | Print | Return:
			return
		}
		p.advance()
	}

}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() Token {
	return *p.Tokens[p.current]
}
func (p *Parser) previous() Token {
	return *p.Tokens[p.current-1]
}
