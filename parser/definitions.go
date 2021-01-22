package parser

import (
	"../lexer"
)

// Visitor represents interface for visiting nodes in the AST
type Visitor interface {
	// to do add visit methods
}

// Node is the interface type in the AST
type Node interface {
	acceptVisitor(*Visitor)
}

// SimpleExpr ..
type SimpleExpr struct {
	Type lexer.TokenType
	Val  string
}

func (se *SimpleExpr) acceptVisitor(*Visitor) {

}

// UnaryExpr ..
type UnaryExpr struct {
	op      lexer.TokenType
	operand Node
}

func (ue *UnaryExpr) acceptVisitor(*Visitor) {

}

// BinExpr ..
type BinExpr struct {
	op    lexer.TokenType
	left  Node
	right Node
}

func (be *BinExpr) acceptVisitor(*Visitor) {

}

// Identifier ..
type Identifier struct {
	Name string
}

func (id *Identifier) acceptVisitor(*Visitor) {

}

type ConstructorExpr struct {
	FieldList []Node
}

func (c *ConstructorExpr) acceptVisitor(*Visitor) {

}

type IndexExpr struct {
	Arr       Node
	ExprIndex Node
}

func (ie *IndexExpr) acceptVisitor(*Visitor) {

}

type MemberExpr struct {
	Obj   Node
	Field Identifier
}

func (m *MemberExpr) acceptVisitor(*Visitor) {

}

type KeyExpr struct {
	LeftExpr  Node
	RightExpr Node
}

func (t *KeyExpr) acceptVisitor(*Visitor) {

}

type ArgList []Node

func (t ArgList) acceptVisitor(*Visitor) {

}

type NameArgList struct {
	Name string
	Arg  Node
}

func (t *NameArgList) acceptVisitor(*Visitor) {

}
