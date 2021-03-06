package ast2jsonipl

import (
	"fmt"
	"io"
	"strconv"

	"../lexer"
	"../parser"
)

var tokenOp map[lexer.TokenType]string = map[lexer.TokenType]string{
	lexer.STRING:   "String",
	lexer.NUMBER:   "Number",
	lexer.NIL:      "Null",
	lexer.ASSIGN:   "Equal",
	lexer.PLUS:     "Plus",
	lexer.MINUS:    "Minus",
	lexer.MULT:     "Star",
	lexer.DIV:      "Division",
	lexer.MOD:      "Modulo",
	lexer.CONCAT:   "Plus",
	lexer.LESSER:   "Less",
	lexer.LESSERQ:  "LessEqual",
	lexer.GREATER:  "Greater",
	lexer.GREATERQ: "GreaterEqual",
	lexer.EQ:       "EqualEqual",
	lexer.AND:      "LogicalAnd",
	lexer.OR:       "LogicalOr",
	lexer.UMINUS:   "Minus",
	lexer.NOT:      "LogicalNot"}

type VisitorJSON struct {
	writer io.Writer
}

func NewJSONVisitor(writer io.Writer) *VisitorJSON {
	return &VisitorJSON{writer}
}

func (v *VisitorJSON) checkAndAccept(node parser.Node) {
	if node != nil {
		node.AcceptVisitor(v)
	} else {
		v.writer.Write([]byte("Null"))
	}
}

func (v *VisitorJSON) VisitSimpleExpr(expr *parser.SimpleExpr) {
	io.WriteString(v.writer, "{")
	if expr.Val == "true" || expr.Val == "false" {
		io.WriteString(v.writer, "\"ExpressionType\": \"LiteralBoolean\",")
		io.WriteString(v.writer, fmt.Sprintf("\"Value\": %s", expr.Val))
		io.WriteString(v.writer, "}")
		return
	}
	toNum, err := strconv.Atoi(expr.Val)
	if err != nil {
		io.WriteString(v.writer, "\"ExpressionType\": \"LiteralString\",")
		io.WriteString(v.writer, fmt.Sprintf("\"Value\": \"%s\"", expr.Val))
		io.WriteString(v.writer, "}")
		return
	}

	io.WriteString(v.writer, "\"ExpressionType\": \"LiteralNumber\",")
	io.WriteString(v.writer, fmt.Sprintf("\"Value\": %d", toNum))
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitUnaryExpr(expr *parser.UnaryExpr) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"UnaryExpression\",")
	io.WriteString(v.writer, "\"Expr\": ")
	v.checkAndAccept(expr.Operand)
	io.WriteString(v.writer, fmt.Sprintf(", \"Operator\": \"%s\"", tokenOp[expr.Op]))
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitBinExpr(expr *parser.BinExpr) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"BinaryExpression\",")
	io.WriteString(v.writer, "\"Left\": ")
	v.checkAndAccept(expr.Left)
	io.WriteString(v.writer, ", \"Right\": ")
	v.checkAndAccept(expr.Right)
	io.WriteString(v.writer, fmt.Sprintf(", \"Operator\": \"%s\"", tokenOp[expr.Op]))
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitIdentifier(id *parser.Identifier) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"IdentifierExpression\",")
	io.WriteString(v.writer, fmt.Sprintf("\"Name\": \"%s\"", id.Name))
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitConstructorExpr(expr *parser.ConstructorExpr) {
}

func (v *VisitorJSON) VisitIndexExpr(expr *parser.IndexExpr) {
}

func (v *VisitorJSON) VisitMemberExpr(expr *parser.MemberExpr) {
}

func (v *VisitorJSON) VisitKeyExpr(expr *parser.KeyExpr) {
}

func (v *VisitorJSON) VisitProgram(program parser.Program) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"TopStatements\",")
	io.WriteString(v.writer, "\"Values\": [")
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
	io.WriteString(v.writer, "\"ExpressionType\": \"ListExpression\",")
	io.WriteString(v.writer, "\"Values\": [")
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
	io.WriteString(v.writer, "\"Identifier\": ")
	v.checkAndAccept(expr.Base)
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"Arguments\": ")
	v.checkAndAccept(expr.Arguments)
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitFunction(f *parser.Function) {
}

func (v *VisitorJSON) VisitNamedFunction(f *parser.NamedFunction) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"FunctionDeclaration\",")
	io.WriteString(v.writer, "\"Name\": ")
	val, ok := f.FunctionName.(*parser.Identifier)
	if ok {
		io.WriteString(v.writer, "\""+val.Name+"\"")
	} else {
		v.checkAndAccept(f.FunctionName)
	}
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"ArgumentsIdentifiers\": ")
	v.writeParamList(f.Parameters)
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"Body\": ")
	v.checkAndAccept(parser.Program(f.Body))
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitLocalFunction(f *parser.LocalFunction) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"FunctionDeclaration\",")
	io.WriteString(v.writer, "\"Name\": ")
	val, ok := f.FunctionName.(*parser.Identifier)
	if ok {
		io.WriteString(v.writer, "\""+val.Name+"\"")
	} else {
		v.checkAndAccept(f.FunctionName)
	}
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"ArgumentsIdentifiers\": ")
	v.writeParamList(f.Parameters)
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"Body\": ")
	v.checkAndAccept(parser.Program(f.Body))
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitAssignmentExpr(expr *parser.AssignmentExpr) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"VariableDefinitionExpression\",")
	io.WriteString(v.writer, "\"Name\": ")
	if len(expr.Vars) > 0 {
		simpExpr, ok := expr.Vars[0].(*parser.Identifier)
		if ok {
			io.WriteString(v.writer, fmt.Sprintf("\"%s\"", simpExpr.Name))
		}
	} else {
		io.WriteString(v.writer, "\"Null\"")
	}
	io.WriteString(v.writer, ", \"Value\": ")
	if len(expr.Exprs) > 0 {
		v.checkAndAccept(expr.Exprs[0])
	} else {
		io.WriteString(v.writer, "\"Null\"")
	}
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitLocalAssignmentExpr(expr *parser.LocalAssignmentExpr) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"VariableDefinitionExpression\",")
	io.WriteString(v.writer, "\"Name\": ")
	if len(expr.Vars) > 0 {
		simpExpr, ok := expr.Vars[0].(*parser.Identifier)
		if ok {
			io.WriteString(v.writer, fmt.Sprintf("\"%s\"", simpExpr.Name))
		}
	} else {
		io.WriteString(v.writer, "\"Null\"")
	}
	io.WriteString(v.writer, ", \"Value\": ")
	if len(expr.Exprs) > 0 {
		v.checkAndAccept(expr.Exprs[0])
	} else {
		io.WriteString(v.writer, "\"Null\"")
	}
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitDoStmnt(st *parser.DoStmnt) {
}

func (v *VisitorJSON) VisitWhileStmnt(st *parser.WhileStmnt) {
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"WhileStatement\",")
	io.WriteString(v.writer, "\"Condition\": ")
	v.checkAndAccept(st.Condition)
	io.WriteString(v.writer, ",")
	io.WriteString(v.writer, "\"Body\": ")
	v.checkAndAccept(parser.Program(st.Block))
	io.WriteString(v.writer, "}")
}

func (v *VisitorJSON) VisitIfStmnt(st *parser.IfStmnt) {
	argList := st.Clauses.(parser.ArgList)
	if len(argList) == 0 {
		//io.WriteString(v.writer, "\"Null\"")
		return
	}
	io.WriteString(v.writer, "{")
	io.WriteString(v.writer, "\"ExpressionType\": \"IfStatement\",")
	switch t := argList[0].(type) {
	case *parser.IfClause:
		io.WriteString(v.writer, "\"Condition\": ")
		v.checkAndAccept(t.Condition)
		io.WriteString(v.writer, ", \"IfStatement\": ")
		v.checkAndAccept(parser.Program(t.Block))
		io.WriteString(v.writer, ", \"ElseStatement\": ")
	case *parser.ElseIfClause:
		io.WriteString(v.writer, "\"Condition\": ")
		v.checkAndAccept(t.Condition)
		io.WriteString(v.writer, ", \"IfStatement\": ")
		v.checkAndAccept(parser.Program(t.Block))
		io.WriteString(v.writer, ", \"ElseStatement\": ")
	case *parser.ElseClause:
		io.WriteString(v.writer, "\"Condition\": \"Null\"")
		io.WriteString(v.writer, ", \"IfStatement\": ")
		v.checkAndAccept(parser.Program(t.Block))
		io.WriteString(v.writer, ", \"ElseStatement\": ")
	}
	v.VisitIfStmnt(&parser.IfStmnt{argList[1:]})

	io.WriteString(v.writer, "}")

}

func (v *VisitorJSON) VisitIfClause(st *parser.IfClause) {
}

func (v *VisitorJSON) VisitElseIfClause(st *parser.ElseIfClause) {
}

func (v *VisitorJSON) VisitElseClause(st *parser.ElseClause) {
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
	io.WriteString(v.writer, ", \"Body\": ")
	v.checkAndAccept(parser.Program(st.Block))
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

func (v *VisitorJSON) writeParamList(list parser.ArgList) {
	io.WriteString(v.writer, "[")
	for i, node := range list {
		if i == len(list)-1 {
			io.WriteString(v.writer, "\""+node.(*parser.Identifier).Name+"\"")
			break
		}
		io.WriteString(v.writer, "\""+node.(*parser.Identifier).Name+"\",")
	}
	io.WriteString(v.writer, "]")
}
