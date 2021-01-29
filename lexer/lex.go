package lexer

import (
	"errors"
)

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isLetter(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

func isValidStartIdentifier(c byte) bool {
	return c == '_' || isLetter(c)
}

func isValidIdentifier(c byte) bool {
	return c == '_' || isDigit(c) || isLetter(c)
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n'
}

func isHex(c byte) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f')
}

// Token represents a single token in the Lexer
type Token struct {
	Type  TokenType
	Val   string
	atRow int
	atCol int
}

// Lexer represents the unit which will parse the source file
type Lexer struct {
	src      string
	tokens   []Token
	crrRow   int
	crrCol   int
	keywords map[string]TokenType
	i        int
}

func (lex *Lexer) prev() {
	lex.i--
}

// New constructs new lexer
func (lex *Lexer) New(src string) Lexer {

	kwrds := map[string]TokenType{
		"and":      AND,
		"end":      END,
		"in":       IN,
		"repeat":   REPEAT,
		"break":    BREAK,
		"false":    FALSE,
		"local":    LOCAL,
		"return":   RETURN,
		"do":       DO,
		"for":      FOR,
		"nil":      NIL,
		"then":     THEN,
		"else":     ELSE,
		"function": FUNCTION,
		"not":      NOT,
		"true":     TRUE,
		"elseif":   ELSEIF,
		"if":       IF,
		"or":       OR,
		"until":    UNTIL,
		"while":    WHILE}

	return Lexer{src: src, tokens: nil, crrRow: 0, crrCol: 0, keywords: kwrds, i: 0}
}

func (lex *Lexer) current() (byte, error) {
	if lex.i >= len(lex.src) {
		return '0', errors.New("EOF")
	}
	return lex.src[lex.i], nil
}

func (lex *Lexer) reslice() error {
	if lex.i >= len(lex.src) {
		lex.src = ""
		return errors.New("EOF")
	}
	lex.src = lex.src[lex.i:]
	lex.i = 0
	return nil
}

func (lex *Lexer) next() {
	if lex.src[0] == '\n' {
		lex.crrRow++
		lex.crrCol = 0
	} else {
		lex.crrCol++
	}
	lex.i++
}

func (lex *Lexer) matchOne(m string) bool {
	for i := 0; i < len(m) && len(lex.src) > lex.i; i++ {
		if lex.src[lex.i] == m[i] {
			lex.i++
			return true
		}
	}
	return false
}

func (lex *Lexer) parseNumber() (Token, error) {

	char, err := lex.current()
	if err == nil && !isDigit(char) {
		return Token{}, errors.New("not a number")
	}

	lex.next()

	if lex.src[lex.i-1] == '0' && lex.matchOne("x") {
		lex.crrCol++
		char, err = lex.current()
		for err == nil && isHex(char) {
			lex.next()
			char, err = lex.current()
		}
		num := lex.src[:lex.i]
		lex.reslice()
		return Token{Type: NUMBER, Val: num, atRow: lex.crrRow, atCol: lex.crrCol}, nil
	}

	dot := false
	char, err = lex.current()
	for err == nil && (isDigit(char) || (char == '.' && !dot)) {
		if char == '.' {
			dot = true
		}
		lex.next()
		char, err = lex.current()
	}
	num := lex.src[:lex.i]
	lex.reslice()
	return Token{Type: NUMBER, Val: num, atRow: lex.crrRow, atCol: lex.crrCol}, nil
}

func (lex *Lexer) parseString() (Token, error) {
	if !lex.matchOne("\"'[") {
		return Token{Type: NUMBER, Val: "", atRow: lex.crrRow, atCol: lex.crrCol}, errors.New("not a string")
	}

	charM := lex.src[lex.i-1]
	prev := charM

	// case [[ ]] - multiline string
	if charM == '[' {
		if !lex.matchOne("[") {
			lex.reslice()
			return Token{Type: LBRACE, Val: "[", atRow: lex.crrRow, atCol: lex.crrCol}, nil
		}

		crr, err := lex.current()
		for err == nil && !(crr == ']' && crr == prev) {
			prev = crr
			lex.next()
			crr, err = lex.current()
		}
		if err != nil {
			panic(err)
		}

		str := lex.src[2 : lex.i-1]
		lex.next()
		lex.reslice()

		return Token{Type: STRING, Val: str, atRow: lex.crrRow, atCol: lex.crrCol}, nil
	}

	lex.next()
	crr, err := lex.current()
	for err == nil && crr != '\n' && (crr != charM || (crr == charM && prev == '\\')) {
		prev = crr
		lex.next()
		crr, err = lex.current()
	}

	if err != nil {
		panic(err)
	}

	if crr == '\n' {
		return Token{Type: NUMBER, Val: "", atRow: lex.crrRow, atCol: lex.crrCol}, errors.New("Expected \" to end string")
	}
	str := lex.src[1:lex.i]
	lex.next()
	lex.reslice()

	return Token{Type: STRING, Val: str, atRow: lex.crrRow, atCol: lex.crrCol}, nil
}

func (lex *Lexer) parseComment() (Token, error) {
	crr, err := lex.current()
	if err != nil {
		panic(err)
	}

	lex.next()
	if crr == '-' && lex.matchOne("-") {
		lex.reslice()
		crr, err = lex.current()
		var comment string
		lex.next()
		if crr == '[' && lex.matchOne("[") {
			lex.reslice()
			crr, err = lex.current()
			prev := crr
			for err == nil && !(crr == ']' && prev == ']') {
				lex.next()
				prev = crr
				crr, err = lex.current()
			}

			if err != nil {
				comment = lex.src[:len(lex.src)]
				lex.reslice()
				if len(lex.src) == 0 {
					return Token{Type: COMMENT, Val: comment, atRow: lex.crrRow, atCol: lex.crrCol}, nil
				}
				return Token{Type: COMMENT, Val: comment, atRow: lex.crrRow, atCol: lex.crrCol}, errors.New("EOF")
			}
			comment = lex.src[:lex.i-1]
			lex.next()
			lex.reslice()
			return Token{Type: COMMENT, Val: comment, atRow: lex.crrRow, atCol: lex.crrCol}, nil
		}

		for err == nil && crr != '\n' {
			lex.next()
			crr, err = lex.current()
		}

		if err != nil {
			comment = lex.src[:len(lex.src)]
			lex.reslice()
			if len(lex.src) == 0 {
				return Token{Type: COMMENT, Val: comment, atRow: lex.crrRow, atCol: lex.crrCol}, nil
			}
			return Token{Type: COMMENT, Val: comment, atRow: lex.crrRow, atCol: lex.crrCol}, errors.New("EOF")
		}
		comment = lex.src[:lex.i]
		lex.next()
		lex.reslice()
		return Token{Type: COMMENT, Val: comment, atRow: lex.crrRow, atCol: lex.crrCol}, nil
	}

	lex.prev()
	return Token{Type: INVALID, Val: "", atRow: lex.crrRow, atCol: lex.crrCol}, errors.New("not a comment")
}

func (lex *Lexer) parseIdentifier() (Token, error) {

	crr, err := lex.current()
	if err != nil {
		panic(err)
	}

	if isValidStartIdentifier(crr) {
		lex.next()
		crr, err = lex.current()
		for err == nil && isValidIdentifier(crr) {
			lex.next()
			crr, err = lex.current()
		}
		str := lex.src[:lex.i]
		lex.reslice()
		val, hasKey := lex.keywords[str]

		// keyword case
		if hasKey {
			return Token{Type: val, Val: str, atRow: lex.crrRow, atCol: lex.crrCol}, nil
		}
		return Token{Type: IDENTIFIER, Val: str, atRow: lex.crrRow, atCol: lex.crrCol}, nil
	}

	return Token{Type: INVALID, Val: "", atRow: lex.crrRow, atCol: lex.crrCol}, errors.New("not a identifier")
}

func (lex *Lexer) smallerToken() (Token, error) {
	crr, err := lex.current()
	if err != nil {
		panic(err)
	}

	var token = Token{Type: INVALID, Val: "", atRow: lex.crrRow, atCol: lex.crrCol}
	err = nil

	switch crr {

	case '+':
		token = Token{Type: PLUS, Val: "+", atRow: lex.crrRow, atCol: lex.crrCol}
	case '-':
		token = Token{Type: MINUS, Val: "-", atRow: lex.crrRow, atCol: lex.crrCol}
	case '*':
		token = Token{Type: MULT, Val: "*", atRow: lex.crrRow, atCol: lex.crrCol}
	case '/':
		token = Token{Type: DIV, Val: "/", atRow: lex.crrRow, atCol: lex.crrCol}
	case '%':
		token = Token{Type: MOD, Val: "%", atRow: lex.crrRow, atCol: lex.crrCol}
	case '^':
		token = Token{Type: POW, Val: "^", atRow: lex.crrRow, atCol: lex.crrCol}
	case '#':
		token = Token{Type: HTAG, Val: "#", atRow: lex.crrRow, atCol: lex.crrCol}
	case '(':
		token = Token{Type: LPAR, Val: "(", atRow: lex.crrRow, atCol: lex.crrCol}
	case ')':
		token = Token{Type: RPAR, Val: ")", atRow: lex.crrRow, atCol: lex.crrCol}
	case '[':
		token = Token{Type: LBRACE, Val: "[", atRow: lex.crrRow, atCol: lex.crrCol}
	case ']':
		token = Token{Type: RBRACE, Val: "]", atRow: lex.crrRow, atCol: lex.crrCol}
	case '{':
		token = Token{Type: LCBRACE, Val: "{", atRow: lex.crrRow, atCol: lex.crrCol}
	case '}':
		token = Token{Type: RCBRACE, Val: "}", atRow: lex.crrRow, atCol: lex.crrCol}
	case ';':
		token = Token{Type: SEMICOLON, Val: ";", atRow: lex.crrRow, atCol: lex.crrCol}
	case ':':
		token = Token{Type: COLON, Val: ":", atRow: lex.crrRow, atCol: lex.crrCol}
	case ',':
		token = Token{Type: COMMA, Val: ",", atRow: lex.crrRow, atCol: lex.crrCol}
	case '=':
		lex.next()
		crr, err = lex.current()
		if err == nil && crr == '=' {
			lex.next()
			lex.reslice()
			return Token{Type: EQ, Val: "==", atRow: lex.crrRow, atCol: lex.crrCol}, err
		}
		lex.reslice()
		return Token{Type: ASSIGN, Val: "=", atRow: lex.crrRow, atCol: lex.crrCol}, err
	case '<':
		lex.next()
		crr, err = lex.current()
		if err == nil && crr == '=' {
			lex.next()
			lex.reslice()
			return Token{Type: LESSERQ, Val: "<=", atRow: lex.crrRow, atCol: lex.crrCol}, err
		}
		lex.reslice()
		return Token{Type: LESSER, Val: "<", atRow: lex.crrRow, atCol: lex.crrCol}, err
	case '>':
		lex.next()
		crr, err = lex.current()
		if err == nil && crr == '=' {
			lex.next()
			lex.reslice()
			return Token{Type: GREATERQ, Val: ">=", atRow: lex.crrRow, atCol: lex.crrCol}, err
		}
		lex.reslice()
		return Token{Type: GREATER, Val: ">", atRow: lex.crrRow, atCol: lex.crrCol}, err
	case '.':
		lex.next()
		crr, err = lex.current()
		if err == nil && crr == '.' {
			lex.next()
			crr, err = lex.current()
			if err == nil && crr == '.' {
				lex.next()
				lex.reslice()
				return Token{Type: VARAGS, Val: "...", atRow: lex.crrRow, atCol: lex.crrCol}, err
			}
			lex.reslice()
			return Token{Type: CONCAT, Val: "..", atRow: lex.crrRow, atCol: lex.crrCol}, err
		}
		lex.reslice()
		return Token{Type: DOT, Val: ".", atRow: lex.crrRow, atCol: lex.crrCol}, err
	}

	if token.Type == INVALID {
		return token, errors.New("continue parsing")
	}

	lex.next()
	lex.reslice()
	return token, nil
}

func (lex *Lexer) nextToken() (Token, error) {

	char, err := lex.current()
	for err == nil && len(lex.src) > 0 && isWhitespace(char) {
		lex.next()
		char, err = lex.current()
	}

	if err := lex.reslice(); err != nil {
		return Token{Type: EOF, Val: "", atRow: lex.crrRow, atCol: lex.crrCol}, err
	}

	token, err := lex.parseComment()
	if err == nil || len(lex.src) == 0 {
		return token, nil
	}

	token, err = lex.parseString()
	if err == nil {
		return token, nil
	}

	token, err = lex.smallerToken()
	if err == nil {
		return token, nil
	}

	token, err = lex.parseNumber()
	if err == nil {
		return token, nil
	}

	token, err = lex.parseIdentifier()
	if err == nil {
		return token, nil
	}

	return Token{Type: INVALID, Val: "", atRow: lex.crrRow, atCol: lex.crrCol}, nil
}

// Run produces a list of tokens from the source
func (lex *Lexer) Run() ([]Token, error) {
	for len(lex.src) > 0 {
		token, err := lex.nextToken()
		if err != nil {
			return lex.tokens, err
		}
		if token.Type != COMMENT {
			lex.tokens = append(lex.tokens, token)
		}
	}
	return lex.tokens, nil
}
