package parser

import (
	"errors"

	"../lexer"
)

// Parser builds the AST
type Parser struct {
	tokens []lexer.Token
	Start  Node
	i      int
}

// NewParser constructor
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
	return tt == lexer.NIL || tt == lexer.FALSE || tt == lexer.TRUE || tt == lexer.NUMBER || tt == lexer.STRING || tt == lexer.VARAGS || tt == lexer.IDENTIFIER
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

func (p *Parser) parseField() Node {
	crr, err := p.current()
	if err != nil {
		panic(err)
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
			return nil
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
		panic(err)
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
		panic(err)
	}

	name := ""

	if crr.Type == lexer.COLON {
		p.next()
		crr, _ = p.current()
		name = crr.Val
		p.next()
	}

	// args nonterminal
	if crr.Type == lexer.LPAR {
		argList := make([]Node, 0, 10)
		p.next()
		expr := p.parseExpression()
		if expr != nil {
			argList = append(argList, expr)
			crr, _ = p.current()
			for crr.Type == lexer.COMMA {
				p.next()
				expr = p.parseExpression()
				argList = append(argList, expr)
			}
			p.next() // ')'
		}
		return &NameArgList{name, ArgList(argList)}
	}

	constructor := p.parseTableConstructor()
	if constructor != nil {
		return &NameArgList{name, constructor}
	}

	if crr.Type != lexer.STRING {
		return nil
	}

	return &NameArgList{name, &SimpleExpr{crr.Type, crr.Val}}
}

/*//!!!! TO DO
func (p *Parser) parseVar() (Node, error) {
	crr, err := p.current()
	if err != nil {
		panic(err)
	}

	var node Node

	if crr.Type == lexer.IDENTIFIER {
		p.next()
		node = &Identifier{crr.Val}
	} else if crr.Type != lexer.LPAR {
		return nil, fmt.Errorf("Not a var")
	} else {
		p.next()
		node, _ = p.parseExpression()
		p.next() // ')'
		p.varSuffix()
	}

	for {

	}

	return node, err
}*/

// change name of func
func (p *Parser) varOrExpr() Node {
	crr, err := p.current()
	// identifier also is simpleExpr, change later

	if err != nil {
		panic(err)
	}

	if err == nil && termExpr(crr.Type) {
		p.next()
		return &SimpleExpr{crr.Type, crr.Val}
	}

	// is function
	// is prefixexpr

	if crr.Type != lexer.LPAR {
		return nil
	}

	p.next()
	expr := p.parseExpression()
	p.next()
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
		//p.next() bomb?
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

	// case function
	// case tableConstructor

	expr := p.parseTableConstructor()
	if expr != nil {
		return expr
	}

	return p.parseBinExpr()
}

func (p *Parser) Run() []Node {
	node := p.parseExpression()
	statements := make([]Node, 0, 10)
	for node != nil {
		statements = append(statements, node)
		if p.i >= len(p.tokens) {
			break
		}
		node = p.parseExpression()
	}

	p.Start = statements[0]
	return statements
}
