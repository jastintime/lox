package main

type Parser struct {
	Tokens  []*Token
	current int
}

type ParseError struct {
	Token  Token
	Messge string
}

func newParser(tokens []*Token) *Parser {
	return &Parser{tokens, 0}
}

func (p *Parser) Parse() []Stmt {
	var statements []Stmt
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements

}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) declaration() Stmt {
	//NOTE: this is golang try catching :)
	defer func() {
		recovered := recover()
		parseError, ok := recovered.(ParseError)
		if ok {
			emitTokenError(parseError.Token, parseError.Messge)
			p.synchronize()
		}
		//STILL WANT TO PANIC OTHERWISE
		if !ok {
			if recovered != nil {
				panic(recovered)
			}
		}
	}()

	if p.match(Class) {
		return p.classDeclaration()
	}

	if p.match(Fun) {
		return p.function("function")
	}

	if p.match(Var) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) classDeclaration() (result Stmt) {
	name := p.consume(Identifier, "Expect class name.")
	var superclass *VariableExpr = nil
	if p.match(Less) {
		p.consume(Identifier, "Expect superclass name.")
		superclass = &VariableExpr{p.previous()}
	}

	p.consume(LeftBrace, "Expect '{' before class body.")

	var methods []FunctionStmt
	for !p.check(RightBrace) && !p.isAtEnd() {
		methods = append(methods, p.function("method"))
	}
	p.consume(RightBrace, "Expect '}' after class body.")
	return ClassStmt{name, superclass, methods}

}

func (p *Parser) statement() Stmt {
	if p.match(For) {
		return p.forStatement()
	}
	if p.match(If) {
		return p.ifStatement()
	}
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
	return p.expressionStatement()
}

//	func (p *Parser) forStatement() Stmt {
//		p.consume(LeftParen, "Expect '(' after 'for'.")
//
//		var initializer Stmt
//		if p.match(Semicolon) {
//			initializer = nil
//		} else if p.match(Var) {
//			initializer = p.varDeclaration()
//		} else {
//			initializer = p.expressionStatement()
//		}
//
//		var condition Expr = nil
//		if !p.check(Semicolon) {
//			condition = p.expression()
//		}
//		p.consume(Semicolon, "Expect ';' after loop condition.")
//
//		var increment Expr = nil
//		if !p.check(RightParen) {
//			increment = p.expression()
//		}
//		p.consume(RightParen, "Expect ')' after for clauses.")
//		body := p.statement()
//		if increment != nil {
//			body = BlockStmt{[]Stmt{body, ExprStmt{increment}}}
//		}
//
//		if condition == nil {
//			condition = LiteralExpr{true}
//		}
//		body = WhileStmt{condition, body}
//		if initializer != nil {
//			body = BlockStmt{[]Stmt{initializer, body}}
//		}
//
//		return body
//	}
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
	p.consume(Semicolon, "Expect ';' after loop condition.")

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
			if len(parameters) >= 255 {
				panic(ParseError{p.peek(), "Can't have more than 255 parameters."})
			}
			parameters = append(parameters, p.consume(Identifier, "Expect parameter name."))
		}
	}
	p.consume(RightParen, "Expect ')' after parameters.")
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
		varE, ok := expr.(VariableExpr)
		if ok {
			name := varE.Name
			return AssignExpr{name, value}
		}
		get, ok := expr.(GetExpr)
		if ok {
			return SetExpr{get.Object, get.Name, value}
		}
		panic(ParseError{equals, "Invalid assignment target."})
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
			if len(arguments) >= 255 {
				panic(ParseError{p.peek(), "Can't have more than 255 arguments."})
			}
			arguments = append(arguments, p.expression())
			//NOTE: in the java version we do 255 but here we aren't doing a do while so its 254
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
		} else if p.match(Dot) {
			name := p.consume(Identifier, "Expect property name after '.'.")
			expr = GetExpr{expr, name}
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
	if p.match(Super) {
		keyword := p.previous()
		p.consume(Dot, "Expect '.' after 'super'.")
		method := p.consume(Identifier, "Expect superclass method name.")
		return SuperExpr{keyword, method}
	}
	if p.match(This) {
		return ThisExpr{p.previous()}
	}
	if p.match(Identifier) {
		return VariableExpr{p.previous()}
	}
	if p.match(LeftParen) {
		expr := p.expression()
		p.consume(RightParen, "Expect ')' after expression.")
		return GroupingExpr{expr}
	}
	panic(ParseError{p.peek(), "Expect expression."})
	return nil

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
	panic(ParseError{p.peek(), message})
	return Token{}

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
