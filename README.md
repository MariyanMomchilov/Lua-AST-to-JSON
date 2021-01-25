# Frontend for Lua

The purpose of this project is to build Abstract Syntax Tree(AST) from Lua code and then serialize the AST in JSON format. This is done is 3 steps:
1. Splitting the Lua code into a series of tokens
2. Creating the AST from the tokens
3. Traverse the tree and serialize each node 

Currently still working on the parser