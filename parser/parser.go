package parser

import (
	"errors"
	"fmt"

	"../lexer"
)

// Parser is the AST parser
type Parser struct {
	tokens        []lexer.Token
	topstatements Program
	i             int
}

// NewParser constructs a Parser
func NewParser(tokens []lexer.Token) Parser {
	return Parser{tokens, nil, 0}
}

func unOp(tt lexer.TokenType) bool {
	return tt == lexer.UMINUS || tt == lexer.NOT || tt == lexer.HTAG
}

func prodOp(tt lexer.TokenType) bool {
	return tt == lexer.MULT || tt == lexer.DIV || tt == lexer.POW
}

func sumOp(tt lexer.TokenType) bool {
	return tt == lexer.PLUS || tt == lexer.MINUS || tt == lexer.CONCAT || tt == lexer.MOD
}

func relOp(tt lexer.TokenType) bool {
	return tt == lexer.LESSER || tt == lexer.LESSERQ || tt == lexer.GREATER || tt == lexer.GREATERQ || tt == lexer.EQ
}

func termExpr(tt lexer.TokenType) bool {
	return tt == lexer.NIL || tt == lexer.FALSE || tt == lexer.TRUE || tt == lexer.NUMBER || tt == lexer.STRING
}

func (p *Parser) hasTokens() bool {
	return len(p.tokens) > 0
}

func (p *Parser) current() (lexer.Token, error) {
	if p.i >= len(p.tokens) {
		return lexer.Token{}, errors.New("No more tokens")
	}
	return p.tokens[p.i], nil
}

func (p *Parser) next() {
	p.i++
}

func (p *Parser) varList() []Node {
	crr, err := p.current()
	if err != nil {
		return nil
	}

	list := make([]Node, 0, 5)
	variable := p.parseVar()
	if variable == nil {
		return nil
	}
	list = append(list, variable)
	crr, _ = p.current()
	for crr.Type == lexer.COMMA {
		p.next()
		variable = p.parseVar()
		list = append(list, variable)
		crr, _ = p.current()
	}
	return list
}

func (p *Parser) exprList() []Node {
	crr, err := p.current()
	if err != nil {
		return nil
	}

	list := make([]Node, 0, 5)
	expr := p.parseExpression()
	if expr != nil {
		list = append(list, expr)
		crr, _ = p.current()
		for crr.Type == lexer.COMMA {
			p.next()
			expr = p.parseExpression()
			list = append(list, expr)
			crr, _ = p.current()
		}
	}
	return list
}

func (p *Parser) parseField() Node {
	crr, err := p.current()
	if err != nil {
		return nil
	}

	var key Node

	if crr.Type == lexer.LBRACE {
		p.next()
		key = p.parseExpression()
		p.next() // ']'

		crr, _ = p.current()
		if crr.Type != lexer.ASSIGN {
			return nil
		}
		p.next()

	} else if crr.Type == lexer.IDENTIFIER {
		p.next()
		key = &Identifier{crr.Val}
		crr, _ = p.current()
		if crr.Type != lexer.ASSIGN {
			p.i--
			return &KeyExpr{nil, p.parseExpression()}
		}
		p.next()
	}

	expr := p.parseExpression()
	return &KeyExpr{key, expr}
}

func (p *Parser) parseFieldList() []Node {

	fieldList := make([]Node, 0, 10)

	field := p.parseField()

	crr, err := p.current()
	if err != nil {
		return fieldList
	}

	for field != nil {

		fieldList = append(fieldList, field)

		if crr.Type != lexer.SEMICOLON && crr.Type != lexer.COMMA {
			return fieldList
		}
		p.next()
		field = p.parseField()
		crr, err = p.current()
	}
	return fieldList
}

func (p *Parser) parseTableConstructor() Node {
	crr, err := p.current()
	if err != nil {
		return nil
	}

	if crr.Type != lexer.LCBRACE {
		return nil
	}
	p.next()

	fieldList := p.parseFieldList()

	p.next() // '}'
	return &ConstructorExpr{fieldList}

}

func (p *Parser) parseNameAndArgs() Node {
	crr, err := p.current()
	if err != nil {
		return nil
	}

	if crr.Type == lexer.LPAR {
		p.next()
		exprList := p.exprList()
		p.next() // ')'
		return ArgList(exprList)
	}

	constructor := p.parseTableConstructor()
	if constructor != nil {
		return constructor
	}

	if crr.Type != lexer.STRING {
		return nil
	}

	p.next()
	return &SimpleExpr{crr.Type, crr.Val}
}

func (p *Parser) assignmentStatement() Node {
	crrI := p.i
	vars := p.varList()

	crr, _ := p.current()

	if vars == nil || crr.Type != lexer.ASSIGN {
		p.i = crrI
		return nil
	}
	p.next()
	exprs := p.exprList()

	return &AssignmentExpr{vars, exprs}
}
func (p *Parser) statement() Node {

	crr, err := p.current()
	if err != nil {
		return nil
	}

	assignment := p.assignmentStatement()
	if assignment != nil {
		return assignment
	}

	funcCall := p.functionCall()
	if funcCall != nil {
		return funcCall
	}

	switch crr.Type {
	case lexer.DO:
		return p.doStatement()
	case lexer.WHILE:
		return p.whileStatement()
	case lexer.REPEAT:
		return p.repeatStatement()
	case lexer.IF:
		return p.ifStatement()
	case lexer.FOR:
		return p.forStatement()
	case lexer.FUNCTION:
		return p.functionStatement()
	case lexer.LOCAL:
		return p.localStatement()
	case lexer.RETURN:
		return p.returnStatement()
	case lexer.BREAK:
		p.next()
		return &SimpleExpr{lexer.BREAK, "break"}
	}

	return nil
}

func (p *Parser) returnStatement() Node {
	p.next()
	return ReturnList(p.exprList())
}

func (p *Parser) localStatement() Node {
	p.next()

	crr, err := p.current()
	if err != nil {
		return nil
	}

	if crr.Type == lexer.FUNCTION {
		namedFunction, assert := p.functionStatement().(*NamedFunction)
		if assert && namedFunction != nil {
			return &LocalFunction{namedFunction}
		}
	}

	assignment, assert := p.assignmentStatement().(*AssignmentExpr)
	if assert && assignment != nil {
		return &LocalAssignmentExpr{assignment}
	}

	return nil
}

func (p *Parser) functionStatement() Node {
	p.next()

	// function name
	crr, _ := p.current()
	if crr.Type != lexer.IDENTIFIER {
		return nil
	}
	var id Node
	id = &Identifier{crr.Val}
	p.next()

	crr, _ = p.current()
	for crr.Type == lexer.DOT {
		p.next()
		id = &MemberExpr{id, &Identifier{crr.Val}}
		p.next()
		crr, _ = p.current()
	}

	p.next()
	var args []Node
	crr, _ = p.current()

	if crr.Type == lexer.IDENTIFIER {
		for crr.Type == lexer.IDENTIFIER {
			args = append(args, &Identifier{crr.Val})
			p.next()
			p.next() // skip ',' or ')'
			crr, _ = p.current()
		}
	}

	if crr.Type == lexer.VARAGS {
		p.next()
		args = append(args, &Identifier{crr.Val})
	}
	crr, _ = p.current()
	if crr.Type == lexer.RPAR {
		p.next() // skip ')'
	}

	//parse body
	block := p.block()
	crr, _ = p.current()
	if crr.Type != lexer.END {
		return nil
	}
	p.next()

	return &NamedFunction{id, args, block}
}

// only regular for
func (p *Parser) forStatement() Node {
	p.next() // 'for'
	init := p.assignmentStatement()
	p.next() // 'do'
	block := p.block()
	p.next() // 'end'

	expr, ok := init.(*AssignmentExpr)
	if ok && len(expr.Exprs) > 2 {
		return &ForStmnt{expr.Exprs[0], expr.Exprs[1], expr.Exprs[2], block}
	}
	return &ForStmnt{expr.Exprs[0], expr.Exprs[1], nil, block}

}
func (p *Parser) ifStatement() Node {
	clauses := make([]Node, 0, 3)
	p.next() // 'if'
	expr := p.parseExpression()
	p.next() // 'then'
	block := p.block()
	clauses = append(clauses, &IfClause{expr, block})

	crr, err := p.current()
	if err != nil {
		return nil
	}
	for crr.Type == lexer.ELSEIF {
		p.next() // 'elseif'
		expr := p.parseExpression()
		p.next() // 'then'
		block := p.block()
		clauses = append(clauses, &ElseIfClause{expr, block})
		crr, err = p.current()
	}

	if crr.Type == lexer.ELSE {
		p.next() // 'else'
		block := p.block()
		p.next() // 'end'
		clauses = append(clauses, &ElseClause{block})
	}

	crr, err = p.current()
	if crr.Type == lexer.END {
		p.next()
	}

	return &IfStmnt{ArgList(clauses)}
}

func (p *Parser) repeatStatement() Node {
	p.next()
	block := p.block()
	p.next() // 'until'
	cond := p.parseExpression()
	return &RepeatStmnt{cond, block}
}

func (p *Parser) whileStatement() Node {
	p.next()
	cond := p.parseExpression()
	p.next() // 'do'
	block := p.block()
	p.next() // 'end'
	return &WhileStmnt{cond, block}
}

func (p *Parser) doStatement() Node {
	p.next()
	block := p.block()
	p.next() // 'end'
	return &DoStmnt{block}
}

func (p *Parser) block() []Node {
	statements := make([]Node, 0, 10)
	statement := p.statement()
	for statement != nil {
		statements = append(statements, statement)
		statement = p.statement()
	}

	return statements
}

func (p *Parser) functionExpr() Node {
	crr, err := p.current()
	if err != nil {
		return nil
	}

	if crr.Type != lexer.FUNCTION {
		return nil
	}
	p.next()
	crr, err = p.current()
	if crr.Type != lexer.LPAR {
		return nil
	}

	p.next()
	crr, _ = p.current()
	var args []Node

	if crr.Type == lexer.IDENTIFIER {
		for crr.Type == lexer.IDENTIFIER {
			args = append(args, &Identifier{crr.Val})
			p.next()
			p.next() // skip ',' or ')'
			crr, _ = p.current()
		}
	}

	if crr.Type == lexer.VARAGS {
		p.next()
		args = append(args, &Identifier{crr.Val})
	}

	crr, _ = p.current()
	if crr.Type == lexer.RPAR {
		p.next() // skip ')'
	}

	//parse body
	block := p.block()
	crr, _ = p.current()
	if crr.Type != lexer.END {
		return nil
	}
	p.next()

	return &Function{ArgList(args), block}
}

// and prefix expr at the same time
func (p *Parser) functionCall() Node {
	_, err := p.current()
	if err != nil {
		return nil
	}

	callee := p.varOrExpr()
	if callee == nil {
		return nil
	}
	args := p.parseNameAndArgs()
	for args != nil {
		callee = &CallExpr{callee, args}
		args = p.parseNameAndArgs()
	}
	return callee
}

func (p *Parser) varSuffix(pref Node) Node {
	crr, err := p.current()
	if err != nil {
		return nil
	}

	arg := p.parseNameAndArgs()
	callExpr := arg

	if arg != nil {
		arg = p.parseNameAndArgs()
		for arg != nil {
			callExpr = &CallExpr{callExpr, arg}
			arg = p.parseNameAndArgs()
		}
	}

	crr, err = p.current()
	if err == nil && crr.Type == lexer.LBRACE {
		p.next()
		expr := p.parseExpression()
		p.next() // ']'
		if callExpr == nil {
			return &IndexExpr{pref, expr}
		}
		return &IndexExpr{callExpr, expr}
	} else if err == nil && crr.Type == lexer.DOT {
		p.next()
		crr, err = p.current()
		if crr.Type != lexer.IDENTIFIER {
			panic(fmt.Errorf("Expected identifier, but received %o", crr.Type))
		}
		p.next()
		if callExpr == nil {
			return &MemberExpr{pref, &Identifier{crr.Val}}
		}
		return &MemberExpr{callExpr, &Identifier{crr.Val}}
	}

	return nil
}

func (p *Parser) parseVar() Node {
	crr, err := p.current()
	if err != nil {
		return nil
	}

	var varExpr Node

	if crr.Type == lexer.IDENTIFIER {
		p.next()
		varExpr = &Identifier{crr.Val}
		crrI := p.i
		suff := p.varSuffix(varExpr)
		if suff == nil {
			p.i = crrI
			return varExpr
		}
		varExpr = suff
		crrI = p.i
		suff = p.varSuffix(varExpr)
		if suff != nil {
			for suff != nil {
				varExpr = &CallExpr{varExpr, suff}
				suff = p.varSuffix(varExpr)
			}
			return varExpr
		}
		p.i = crrI
		return varExpr

	} else if crr.Type == lexer.LPAR {
		p.next()
		varExpr = p.parseExpression()
		p.next() // ')'
		crrI := p.i
		suff := p.varSuffix(varExpr)

		if suff == nil {
			p.i = crrI
			return varExpr
		}

		varExpr = suff
		crrI = p.i
		suff = p.varSuffix(varExpr)
		if suff != nil {
			for suff != nil {
				varExpr = &CallExpr{varExpr, suff}
				suff = p.varSuffix(varExpr)
			}
			return varExpr
		}
		p.i = crrI
		return varExpr
	}

	return nil
}

func (p *Parser) varOrExpr() Node {
	crr, err := p.current()

	if err != nil {
		return nil
	}

	if termExpr(crr.Type) {
		p.next()
		return &SimpleExpr{crr.Type, crr.Val}
	}

	expr := p.parseVar()
	if expr != nil {
		return expr
	}

	return expr
}

func (p *Parser) parseProdExpr() Node {
	crr, err := p.current()

	if err == nil && unOp(crr.Type) {
		p.next()
		expr := p.varOrExpr()
		return &UnaryExpr{crr.Type, expr}
	}
	return p.varOrExpr()
}

func (p *Parser) parseSumExpr() Node {

	left := p.parseProdExpr()

	crr, err := p.current()
	for err == nil && prodOp(crr.Type) {
		p.next()
		expr := p.parseProdExpr()
		left = &BinExpr{crr.Type, left, expr}
		crr, err = p.current()
	}

	return left
}

func (p *Parser) parseRelExpr() Node {

	left := p.parseSumExpr()

	crr, err := p.current()
	for err == nil && sumOp(crr.Type) {
		p.next()
		expr := p.parseSumExpr()
		left = &BinExpr{crr.Type, left, expr}
		crr, err = p.current()
	}

	return left
}

func (p *Parser) parseAndExpr() Node {

	left := p.parseRelExpr()

	crr, err := p.current()
	for err == nil && relOp(crr.Type) {
		p.next()
		expr := p.parseRelExpr()
		left = &BinExpr{crr.Type, left, expr}
		crr, err = p.current()
	}

	return left
}

func (p *Parser) parseOrExpr() Node {

	left := p.parseAndExpr()

	crr, err := p.current()
	for err == nil && crr.Type == lexer.AND {
		p.next()
		expr := p.parseAndExpr()
		left = &BinExpr{crr.Type, left, expr}
		crr, err = p.current()
	}

	return left
}

func (p *Parser) parseBinExpr() Node {

	left := p.parseOrExpr()
	crr, err := p.current()
	for err == nil && crr.Type == lexer.OR {
		p.next()
		expr := p.parseOrExpr()
		left = &BinExpr{crr.Type, left, expr}
		crr, err = p.current()
	}

	return left
}

func (p *Parser) parseExpression() Node {

	expr := p.functionExpr()
	if expr != nil {
		return expr
	}

	crrI := p.i
	expr = p.parseBinExpr()
	if expr != nil {
		crr, _ := p.current()
		if crr.Type != lexer.LPAR {
			return expr
		}
		p.i = crrI
	}

	expr = p.parseTableConstructor()
	if expr != nil {
		return expr
	}

	return p.functionCall()
}

// Run builds the AST
func (p *Parser) Run() Program {
	statement := p.statement()
	statements := make([]Node, 0, 10)
	for statement != nil {
		statements = append(statements, statement)
		if p.i >= len(p.tokens) {
			break
		}
		statement = p.statement()
	}

	p.topstatements = statements
	return statements
}
