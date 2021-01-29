package lexer

// TokenType enum
type TokenType int

const (
	// keywords
	END TokenType = iota
	IN
	REPEAT
	BREAK
	FALSE
	LOCAL
	RETURN
	DO
	FOR
	NIL
	THEN
	ELSE
	FUNCTION
	TRUE
	ELSEIF
	IF
	UNTIL
	WHILE

	IDENTIFIER
	STRING
	NUMBER
	COMMENT
	EOF
	INVALID

	DOT
	COMMA
	SEMICOLON
	COLON
	LPAR
	RPAR
	LBRACE  // [
	RBRACE  // ]
	LCBRACE // {
	RCBRACE // }
	VARAGS  // ...

	// bin op's
	ASSIGN
	PLUS
	MINUS
	MULT
	DIV
	POW
	MOD
	CONCAT
	LESSER   // <
	LESSERQ  // <=
	GREATER  // >
	GREATERQ // >=
	EQ       // ==
	AND
	OR

	//unary op's
	UMINUS
	NOT
	HTAG // #
)
