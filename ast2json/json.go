package ast2json

import (
	"fmt"
	"io"

	"../lexer"
	"../parser"
)

var tokenOp map[lexer.TokenType]string = map[lexer.TokenType]string{
	lexer.STRING:   "string",
	lexer.NUMBER:   "number",
	lexer.NIL:      "nil",
	lexer.ASSIGN:   "=",
	lexer.PLUS:     "+",
	lexer.MINUS:    "-",
	lexer.MULT:     "*",
	lexer.DIV:      "/",
	lexer.POW:      "^",
	lexer.MOD:      "%",
	lexer.CONCAT:   "..",
	lexer.LESSER:   "<",
	lexer.LESSERQ:  "<=",
	lexer.GREATER:  ">",
	lexer.GREATERQ: ">=",
	lexer.EQ:       "==",
	lexer.AND:      "and",
	lexer.OR:       "or",
	lexer.UMINUS:   "-",
	lexer.NOT:      "not",
	lexer.HTAG:     "#"}

type VisitorJSON struct {
	indent int
	writer io.Writer
}

func NewJSONVisitor(writer io.Writer) *VisitorJSON {
	return &VisitorJSON{0, writer}
}

func (v *VisitorJSON) checkAndAccept(node parser.Node) {
	if node != nil {
		node.AcceptVisitor(v)
	} else {
		v.writer.Write([]byte("null"))
	}
}

func (v *VisitorJSON) VisitSimpleExpr(expr *parser.SimpleExpr) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"SimpleExpression\",")
	io.WriteString(v.writer, fmt.Sprintf("\"ValueType\": \"%s\",", tokenOp[expr.Type]))
	io.WriteString(v.writer, fmt.Sprintf("\"Value\": \"%s\"", expr.Val))
	io.WriteString(v.writer, "}")

}

func (v *VisitorJSON) VisitUnaryExpr(expr *parser.UnaryExpr) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"UnaryExpression\",")
	io.WriteString(v.writer, fmt.Sprintf("\"Operator\": \"%s\",", tokenOp[expr.Op]))
	io.WriteString(v.writer, "\"Operand\": ")
	v.checkAndAccept(expr.Operand)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitBinExpr(expr *parser.BinExpr) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"BinaryExpression\",")
	io.WriteString(v.writer, fmt.Sprintf("\"Operator\": \"%s\",", tokenOp[expr.Op]))
	io.WriteString(v.writer, "\"LeftOperand\": ")
	v.checkAndAccept(expr.Left)
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"RightOperand\": ")
	v.checkAndAccept(expr.Right)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitIdentifier(id *parser.Identifier) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"Identifier\",")
	io.WriteString(v.writer, fmt.Sprintf("\"Name\": \"%s\"", id.Name))
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitConstructorExpr(expr *parser.ConstructorExpr) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"ConstructorExpression\",")
	io.WriteString(v.writer, "\"FieldList\": [")
	for i := range expr.FieldList {
		v.checkAndAccept(expr.FieldList[i])
		if i != len(expr.FieldList)-1 {
			io.WriteString(v.writer, ", ")
		}

	}
	io.WriteString(v.writer, "]")
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitIndexExpr(expr *parser.IndexExpr) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"IndexExpression\",")
	io.WriteString(v.writer, "\"BaseExpression\": ")
	v.checkAndAccept(expr.Base)
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"Index\": ")
	v.checkAndAccept(expr.ExprIndex)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitMemberExpr(expr *parser.MemberExpr) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"MemberExpression\",")
	io.WriteString(v.writer, "\"Object\": ")
	v.checkAndAccept(expr.Obj)
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"Field\": ")
	v.checkAndAccept(expr.Field)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitKeyExpr(expr *parser.KeyExpr) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"KeyExpression\",")
	io.WriteString(v.writer, "\"Key\": ")
	v.checkAndAccept(expr.LeftExpr)
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"Value\": ")
	v.checkAndAccept(expr.RightExpr)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitProgram(program parser.Program) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"Program\",")
	io.WriteString(v.writer, "\"Statements\": [")
	for i := range program {
		v.checkAndAccept(program[i])
		if i != len(program)-1 {
			io.WriteString(v.writer, ", ")
		}
	}
	io.WriteString(v.writer, "]")
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitArgList(l parser.ArgList) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"ArgumentList\",")
	io.WriteString(v.writer, "\"Arguments\": [")
	for i := range l {
		v.checkAndAccept(l[i])
		if i != len(l)-1 {
			io.WriteString(v.writer, ", ")
		}
	}
	io.WriteString(v.writer, "]")
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitReturnList(l parser.ReturnList) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"ReturnList\",")
	io.WriteString(v.writer, "\"ReturnValues\": [")
	for i := range l {
		v.checkAndAccept(l[i])
		if i != len(l)-1 {
			io.WriteString(v.writer, ", ")
		}
	}
	io.WriteString(v.writer, "]")
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitCallExpr(expr *parser.CallExpr) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"CallExpression\",")
	io.WriteString(v.writer, "\"Base\": ")
	v.checkAndAccept(expr.Base)
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"Argument\": ")
	v.checkAndAccept(expr.Arguments)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitFunction(f *parser.Function) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"UnnamedFunction\",")
	io.WriteString(v.writer, "\"Parameters\": ")
	v.checkAndAccept(f.Parameters)
	io.WriteString(v.writer, ",")
	v.body2JSON(f.Body)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitNamedFunction(f *parser.NamedFunction) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"Function\",")
	io.WriteString(v.writer, "\"Name\": ")
	v.checkAndAccept(f.FunctionName)
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"Parameters\": ")
	v.checkAndAccept(f.Parameters)
	io.WriteString(v.writer, ",")
	v.body2JSON(f.Body)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitLocalFunction(f *parser.LocalFunction) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"LocalFunction\",")
	io.WriteString(v.writer, "\"Name\": ")
	v.checkAndAccept(f.FunctionName)
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"Parameters\": ")
	v.checkAndAccept(f.Parameters)
	io.WriteString(v.writer, ",")
	v.body2JSON(f.Body)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitAssignmentExpr(expr *parser.AssignmentExpr) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"AssignmentExpression\",")
	io.WriteString(v.writer, "\"Variables\": [")
	for i := range expr.Vars {
		v.checkAndAccept(expr.Vars[i])
		if i != len(expr.Vars)-1 {
			io.WriteString(v.writer, ", ")
		}
	}
	io.WriteString(v.writer, "],")
	io.WriteString(v.writer, "\"Expressions\": [")
	for i := range expr.Exprs {
		expr.Exprs[i].AcceptVisitor(v)
		if i != len(expr.Exprs)-1 {
			io.WriteString(v.writer, ", ")
		}
	}
	io.WriteString(v.writer, "]")
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitLocalAssignmentExpr(expr *parser.LocalAssignmentExpr) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"LocalAssignmentExpression\",")
	io.WriteString(v.writer, "\"Variables\": [")
	for i := range expr.Vars {
		v.checkAndAccept(expr.Vars[i])
		if i != len(expr.Vars)-1 {
			io.WriteString(v.writer, ", ")
		}
	}
	io.WriteString(v.writer, "],")
	io.WriteString(v.writer, "\"Expressions\": [")
	for i := range expr.Exprs {
		v.checkAndAccept(expr.Exprs[i])
		if i != len(expr.Exprs)-1 {
			io.WriteString(v.writer, ", ")
		}
	}
	io.WriteString(v.writer, "]")
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitDoStmnt(st *parser.DoStmnt) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"DoStatement\",")
	v.body2JSON(st.Block)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitWhileStmnt(st *parser.WhileStmnt) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"WhileStatement\",")
	io.WriteString(v.writer, "\"Condition\": ")
	v.checkAndAccept(st.Condition)
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"Body\": [")
	for i := range st.Block {
		v.checkAndAccept(st.Block[i])
		if i != len(st.Block)-1 {
			io.WriteString(v.writer, ", ")
		}
	}
	io.WriteString(v.writer, "]")
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitRepeatStmnt(*parser.RepeatStmnt) {
	// TO DO
}

func (v *VisitorJSON) VisitIfStmnt(st *parser.IfStmnt) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"IfStatement\",")
	io.WriteString(v.writer, "\"Clauses\": ")
	v.checkAndAccept(st.Clauses)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitIfClause(st *parser.IfClause) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"IfClauseStatement\",")
	io.WriteString(v.writer, "\"Condition\": ")
	v.checkAndAccept(st.Condition)
	io.WriteString(v.writer, ",")
	v.body2JSON(st.Block)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitElseIfClause(st *parser.ElseIfClause) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"ElseIfClauseStatement\",")
	io.WriteString(v.writer, "\"Condition\": ")
	v.checkAndAccept(st.Condition)
	io.WriteString(v.writer, ",")
	v.body2JSON(st.Block)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitElseClause(st *parser.ElseClause) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"ElseClauseStatement\",")
	v.body2JSON(st.Block)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitForStmnt(st *parser.ForStmnt) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"ForStatement\",")
	io.WriteString(v.writer, "\"Initialization\": ")
	v.checkAndAccept(st.Start)
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"Condition\": ")
	v.checkAndAccept(st.Condition)
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"Iteration\": ")
	v.checkAndAccept(st.Step)
	io.WriteString(v.writer, ",")
	v.body2JSON(st.Block)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) body2JSON(block []parser.Node) {
	io.WriteString(v.writer, "\"Body\": [")
	for i := range block {
		v.checkAndAccept(block[i])
		if i != len(block)-1 {
			io.WriteString(v.writer, ", ")
		}
	}
	io.WriteString(v.writer, "]")
}
