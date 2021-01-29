package parser

import (
	"../lexer"
)

// Visitor represents interface for visiting nodes in the AST
type Visitor interface {
	// to do add visit methods
	VisitSimpleExpr(*SimpleExpr)
	VisitUnaryExpr(*UnaryExpr)
	VisitBinExpr(*BinExpr)
	VisitIdentifier(*Identifier)
	VisitConstructorExpr(*ConstructorExpr)
	VisitIndexExpr(*IndexExpr)
	VisitMemberExpr(*MemberExpr)
	VisitKeyExpr(*KeyExpr)
	VisitProgram(Program)
	VisitArgList(ArgList)
	VisitReturnList(ReturnList)
	VisitCallExpr(*CallExpr)
	VisitFunction(*Function)
	VisitNamedFunction(*NamedFunction)
	VisitLocalFunction(*LocalFunction)
	VisitAssignmentExpr(*AssignmentExpr)
	VisitLocalAssignmentExpr(*LocalAssignmentExpr)
	VisitDoStmnt(*DoStmnt)
	VisitWhileStmnt(*WhileStmnt)
	VisitRepeatStmnt(*RepeatStmnt)
	VisitIfStmnt(*IfStmnt)
	VisitIfClause(*IfClause)
	VisitElseIfClause(*ElseIfClause)
	VisitElseClause(*ElseClause)
	VisitForStmnt(*ForStmnt)
}

// Node is the interface type in the AST
type Node interface {
	AcceptVisitor(Visitor)
}

// SimpleExpr ..
type SimpleExpr struct {
	Type lexer.TokenType
	Val  string
}

func (se *SimpleExpr) AcceptVisitor(v Visitor) {
	v.VisitSimpleExpr(se)
}

// UnaryExpr ..
type UnaryExpr struct {
	Op      lexer.TokenType
	Operand Node
}

func (ue *UnaryExpr) AcceptVisitor(v Visitor) {
	v.VisitUnaryExpr(ue)
}

// BinExpr ..
type BinExpr struct {
	Op    lexer.TokenType
	Left  Node
	Right Node
}

func (be *BinExpr) AcceptVisitor(v Visitor) {
	v.VisitBinExpr(be)
}

// Identifier ..
type Identifier struct {
	Name string
}

func (id *Identifier) AcceptVisitor(v Visitor) {
	v.VisitIdentifier(id)
}

type ConstructorExpr struct {
	FieldList []Node
}

func (c *ConstructorExpr) AcceptVisitor(v Visitor) {
	v.VisitConstructorExpr(c)
}

type IndexExpr struct {
	Base      Node
	ExprIndex Node
}

func (ie *IndexExpr) AcceptVisitor(v Visitor) {
	v.VisitIndexExpr(ie)
}

type MemberExpr struct {
	Obj   Node
	Field *Identifier
}

func (m *MemberExpr) AcceptVisitor(v Visitor) {
	v.VisitMemberExpr(m)
}

type KeyExpr struct {
	LeftExpr  Node
	RightExpr Node
}

func (k *KeyExpr) AcceptVisitor(v Visitor) {
	v.VisitKeyExpr(k)
}

type Program []Node

func (p Program) AcceptVisitor(v Visitor) {
	v.VisitProgram(p)
}

type ArgList []Node

func (l ArgList) AcceptVisitor(v Visitor) {
	v.VisitArgList(l)
}

type ReturnList []Node

func (l ReturnList) AcceptVisitor(v Visitor) {
	v.VisitReturnList(l)
}

type CallExpr struct {
	Base      Node
	Arguments Node
}

func (e *CallExpr) AcceptVisitor(v Visitor) {
	v.VisitCallExpr(e)
}

type Function struct {
	Parameters ArgList
	Body       []Node
}

func (f *Function) AcceptVisitor(v Visitor) {
	v.VisitFunction(f)
}

type NamedFunction struct {
	FunctionName Node
	Parameters   ArgList
	Body         []Node
}

func (f *NamedFunction) AcceptVisitor(v Visitor) {
	v.VisitNamedFunction(f)
}

type LocalFunction struct {
	*NamedFunction
}

func (f *LocalFunction) AcceptVisitor(v Visitor) {
	v.VisitLocalFunction(f)
}

type AssignmentExpr struct {
	Vars  []Node
	Exprs []Node
}

func (e *AssignmentExpr) AcceptVisitor(v Visitor) {
	v.VisitAssignmentExpr(e)
}

type LocalAssignmentExpr struct {
	*AssignmentExpr
}

func (e *LocalAssignmentExpr) AcceptVisitor(v Visitor) {
	v.VisitLocalAssignmentExpr(e)
}

type DoStmnt struct {
	Block []Node
}

func (s *DoStmnt) AcceptVisitor(v Visitor) {
	v.VisitDoStmnt(s)
}

type WhileStmnt struct {
	Condition Node
	Block     []Node
}

func (s *WhileStmnt) AcceptVisitor(v Visitor) {
	v.VisitWhileStmnt(s)
}

type RepeatStmnt struct {
	Condition Node
	Block     []Node
}

func (s *RepeatStmnt) AcceptVisitor(v Visitor) {
	v.VisitRepeatStmnt(s)
}

type IfStmnt struct {
	Clauses Node
}

func (s *IfStmnt) AcceptVisitor(v Visitor) {
	v.VisitIfStmnt(s)
}

type IfClause struct {
	Condition Node
	Block     []Node
}

func (s *IfClause) AcceptVisitor(v Visitor) {
	v.VisitIfClause(s)
}

type ElseIfClause struct {
	Condition Node
	Block     []Node
}

func (s *ElseIfClause) AcceptVisitor(v Visitor) {
	v.VisitElseIfClause(s)
}

type ElseClause struct {
	Block []Node
}

func (s *ElseClause) AcceptVisitor(v Visitor) {
	v.VisitElseClause(s)
}

type ForStmnt struct {
	Start     Node
	Condition Node
	Step      Node
	Block     []Node
}

func (s *ForStmnt) AcceptVisitor(v Visitor) {
	v.VisitForStmnt(s)
}
