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
	Base      Node
	ExprIndex Node
}

func (ie *IndexExpr) acceptVisitor(*Visitor) {

}

type MemberExpr struct {
	Obj   Node
	Field *Identifier
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

type CallExpr struct {
	Base      Node
	Arguments Node
}

func (t *CallExpr) acceptVisitor(*Visitor) {

}

type Function struct {
	Parameters ArgList
	Body       []Node
}

func (t *Function) acceptVisitor(*Visitor) {

}

type NamedFunction struct {
	Name       Node
	Parameters ArgList
	Body       []Node
}

func (t *NamedFunction) acceptVisitor(*Visitor) {

}

type LocalFunction struct {
	NamedFunction
}

func (t *LocalFunction) acceptVisitor(*Visitor) {

}

type AssignmentExpr struct {
	Vars  []Node
	Exprs []Node
}

func (f *AssignmentExpr) acceptVisitor(*Visitor) {

}

type LocalAssignmentExpr struct {
	Vars  []Node
	Exprs []Node
}

func (f *LocalAssignmentExpr) acceptVisitor(*Visitor) {

}

type DoStmnt struct {
	Block []Node
}

func (s *DoStmnt) acceptVisitor(*Visitor) {

}

type WhileStmnt struct {
	Condition Node
	Block     []Node
}

func (s *WhileStmnt) acceptVisitor(*Visitor) {

}

type RepeatStmnt struct {
	Condition Node
	Block     []Node
}

func (s *RepeatStmnt) acceptVisitor(*Visitor) {

}

type IfStmnt struct {
	Clauses Node
}

func (s *IfStmnt) acceptVisitor(*Visitor) {

}

type IfClause struct {
	Condition Node
	Block     []Node
}

func (s *IfClause) acceptVisitor(*Visitor) {

}

type ElseIfClause struct {
	Condition Node
	Block     []Node
}

func (s *ElseIfClause) acceptVisitor(*Visitor) {

}

type ElseClause struct {
	Block []Node
}

func (s *ElseClause) acceptVisitor(*Visitor) {

}

type ForStmnt struct {
	Start     Node
	Condition Node
	Step      Node
	Block     []Node
}

func (s *ForStmnt) acceptVisitor(*Visitor) {

}
