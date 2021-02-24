package tests_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"../ast2json"
	ast2jsonipl "../ast2jsonIPL"
	"../lexer"
	"../parser"
)

// TestLex test
func TestLex(t *testing.T) {
	src := "123 4.35  123 0x123fa --[[ asd \n dfsd ]] \"asdd dfgdf\"  for int _asf32523 [[ asddfsf assdf\n fsfgsd]] -- asfsdf "
	var lex lexer.Lexer
	lex = lex.New(src)
	tokens, _ := lex.Run()
	fmt.Println(tokens)
}

func TestLex2(t *testing.T) {
	file, _ := os.Open("lexTest.txt")
	src, _ := ioutil.ReadAll(file)
	var lex lexer.Lexer
	lex = lex.New(string(src))
	tokens, _ := lex.Run()
	fmt.Println(tokens)
}

func TestParser(t *testing.T) {
	file, _ := os.Open("parserTest.txt")
	src, _ := ioutil.ReadAll(file)
	var lex lexer.Lexer
	lex = lex.New(string(src))
	tokens, _ := lex.Run()
	parser := parser.NewParser(tokens)
	node := parser.Run()

	jsonfile, _ := os.Create("test.json")
	defer file.Close()
	visitor := ast2json.NewJSONVisitor(jsonfile)
	node.AcceptVisitor(visitor)
}

func TestParserIPL(t *testing.T) {
	file, _ := os.Open("parserTestIPL.txt")
	src, _ := ioutil.ReadAll(file)
	var lex lexer.Lexer
	lex = lex.New(string(src))
	tokens, _ := lex.Run()
	parser := parser.NewParser(tokens)
	node := parser.Run()

	jsonfile, _ := os.Create("testIPL.json")
	defer file.Close()
	visitor := ast2jsonipl.NewJSONVisitor(jsonfile)
	node.AcceptVisitor(visitor)
}

func TestParserIPL2(t *testing.T) {
	file, _ := os.Open("parserTestIPL2.txt")
	src, _ := ioutil.ReadAll(file)
	var lex lexer.Lexer
	lex = lex.New(string(src))
	tokens, _ := lex.Run()
	parser := parser.NewParser(tokens)
	node := parser.Run()

	jsonfile, _ := os.Create("testIPL2.json")
	defer file.Close()
	visitor := ast2jsonipl.NewJSONVisitor(jsonfile)
	node.AcceptVisitor(visitor)
}
