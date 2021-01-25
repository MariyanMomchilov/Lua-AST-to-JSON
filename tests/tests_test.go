package tests_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

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

func TestParser1(t *testing.T) {
	file, _ := os.Open("testFile1.txt")
	src, _ := ioutil.ReadAll(file)
	var lex lexer.Lexer
	lex = lex.New(string(src))
	tokens, _ := lex.Run()
	fmt.Println(tokens)
	parser := parser.NewParser(tokens)
	nodes := parser.Run()
	fmt.Println(nodes)
}

func TestParser2(t *testing.T) {
	var lex lexer.Lexer
	lex = lex.New(`arr3("abc", d, 1 + 3)(opp)`)
	tokens, _ := lex.Run()
	//fmt.Println(tokens)
	parser := parser.NewParser(tokens)
	node := parser.FunctionCall()
	fmt.Println(node)
}

func TestParser5(t *testing.T) {
	var lex lexer.Lexer
	lex = lex.New("arr[i] == i + 1")
	tokens, _ := lex.Run()
	//fmt.Println(tokens)
	parser := parser.NewParser(tokens)
	node := parser.ParseExpression()
	fmt.Println(node)
}

func TestParser3(t *testing.T) {
	file, _ := os.Open("testFile2.txt")
	src, _ := ioutil.ReadAll(file)
	var lex lexer.Lexer
	lex = lex.New(string(src))
	tokens, _ := lex.Run()
	fmt.Println(tokens)
	parser := parser.NewParser(tokens)
	nodes := parser.Run()
	fmt.Println(nodes)
}

func TestParser4(t *testing.T) {
	file, _ := os.Open("testFile3.txt")
	src, _ := ioutil.ReadAll(file)
	var lex lexer.Lexer
	lex = lex.New(string(src))
	tokens, _ := lex.Run()
	fmt.Println(tokens)
	parser := parser.NewParser(tokens)
	nodes := parser.Run()
	fmt.Println(nodes)
}
