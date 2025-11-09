package exprparser

import "fmt"

// TokenType represents different token types
type TokenType int

const (
	TokenNumber TokenType = iota
	TokenIdent
	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash
	TokenLParen
	TokenRParen
	TokenEOF
)

// Token represents a lexical token
type Token struct {
	Type  TokenType
	Value string
}

// Expr is the base interface for all expressions
type Expr interface {
	Eval(ctx Context) (Value, error)
	String() string
}

// Context provides variable bindings for evaluation
type Context map[string]Value

// Value represents a runtime value
type Value struct {
	// TODO: Support int, float, string, bool
}

// BinaryExpr represents binary operations
type BinaryExpr struct {
	Left  Expr
	Op    string
	Right Expr
}

func (e *BinaryExpr) Eval(ctx Context) (Value, error) {
	// TODO: Implement evaluation
	return Value{}, nil
}

func (e *BinaryExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left, e.Op, e.Right)
}

// Lexer tokenizes input
type Lexer struct {
	input string
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) NextToken() Token {
	// TODO: Implement lexer
	return Token{Type: TokenEOF}
}

// Parser parses expressions
type Parser struct {
	lexer   *Lexer
	current Token
}

func NewParser(input string) *Parser {
	return &Parser{
		lexer: NewLexer(input),
	}
}

func (p *Parser) ParseExpression() (Expr, error) {
	// TODO: Implement recursive descent parser
	return nil, nil
}
