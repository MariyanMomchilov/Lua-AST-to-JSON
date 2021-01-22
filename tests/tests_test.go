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
	file, _ := os.Open("text")
	src, _ := ioutil.ReadAll(file)
	var lex lexer.Lexer
	lex = lex.New(string(src))
	tokens, _ := lex.Run()
	fmt.Println(tokens)
}

func TestParser1(t *testing.T) {
	file, _ := os.Open("p1")
	src, _ := ioutil.ReadAll(file)
	var lex lexer.Lexer
	lex = lex.New(string(src))
	tokens, _ := lex.Run()
	fmt.Println(tokens)
	parser := parser.NewParser(tokens)
	nodes := parser.Run()
	fmt.Println(nodes)
}
