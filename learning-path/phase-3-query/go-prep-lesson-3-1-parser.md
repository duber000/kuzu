# Phase 3 Lesson 3.1: Go Prep - Parser

**Prerequisites:** Phase 2 complete (Graph Structure)
**Time:** 4-5 hours Go prep + 20-25 hours implementation
**Main Curriculum:** `complete-graph-database-learning-path-go1.25.md` Lesson 3.1

## Overview

Query parsing transforms text into Abstract Syntax Trees (ASTs). Before implementing a Cypher parser, master these Go concepts:
- Recursive descent parsing techniques
- AST design with interfaces and type switches
- Error handling with position tracking
- String parsing and tokenization
- Visitor pattern for AST traversal

**This lesson builds the foundation for understanding query intent!**

## Go Concepts for This Lesson

### 1. Tokenization Basics

**Break input into tokens before parsing!**

```go
package main

import (
    "fmt"
    "strings"
    "unicode"
)

type TokenType int

const (
    TOKEN_EOF TokenType = iota
    TOKEN_IDENT
    TOKEN_NUMBER
    TOKEN_LPAREN
    TOKEN_RPAREN
    TOKEN_PLUS
    TOKEN_MINUS
    TOKEN_STAR
    TOKEN_SLASH
)

type Token struct {
    Type    TokenType
    Value   string
    Line    int
    Column  int
}

type Lexer struct {
    input  string
    pos    int
    line   int
    column int
}

func NewLexer(input string) *Lexer {
    return &Lexer{
        input:  input,
        pos:    0,
        line:   1,
        column: 1,
    }
}

func (l *Lexer) NextToken() Token {
    // Skip whitespace
    for l.pos < len(l.input) && unicode.IsSpace(rune(l.input[l.pos])) {
        if l.input[l.pos] == '\n' {
            l.line++
            l.column = 1
        } else {
            l.column++
        }
        l.pos++
    }

    if l.pos >= len(l.input) {
        return Token{Type: TOKEN_EOF, Line: l.line, Column: l.column}
    }

    start := l.pos
    startCol := l.column

    ch := l.input[l.pos]

    // Single character tokens
    switch ch {
    case '(':
        l.pos++
        l.column++
        return Token{Type: TOKEN_LPAREN, Value: "(", Line: l.line, Column: startCol}
    case ')':
        l.pos++
        l.column++
        return Token{Type: TOKEN_RPAREN, Value: ")", Line: l.line, Column: startCol}
    case '+':
        l.pos++
        l.column++
        return Token{Type: TOKEN_PLUS, Value: "+", Line: l.line, Column: startCol}
    case '-':
        l.pos++
        l.column++
        return Token{Type: TOKEN_MINUS, Value: "-", Line: l.line, Column: startCol}
    case '*':
        l.pos++
        l.column++
        return Token{Type: TOKEN_STAR, Value: "*", Line: l.line, Column: startCol}
    case '/':
        l.pos++
        l.column++
        return Token{Type: TOKEN_SLASH, Value: "/", Line: l.line, Column: startCol}
    }

    // Numbers
    if unicode.IsDigit(rune(ch)) {
        for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
            l.pos++
            l.column++
        }
        return Token{
            Type:   TOKEN_NUMBER,
            Value:  l.input[start:l.pos],
            Line:   l.line,
            Column: startCol,
        }
    }

    // Identifiers
    if unicode.IsLetter(rune(ch)) {
        for l.pos < len(l.input) && (unicode.IsLetter(rune(l.input[l.pos])) || unicode.IsDigit(rune(l.input[l.pos]))) {
            l.pos++
            l.column++
        }
        return Token{
            Type:   TOKEN_IDENT,
            Value:  l.input[start:l.pos],
            Line:   l.line,
            Column: startCol,
        }
    }

    // Unknown character
    l.pos++
    l.column++
    return Token{Type: TOKEN_EOF, Value: string(ch), Line: l.line, Column: startCol}
}

func main() {
    input := "foo + 123 * (bar - 456)"
    lexer := NewLexer(input)

    fmt.Println("Tokens:")
    for {
        tok := lexer.NextToken()
        fmt.Printf("  %d:%d %v %q\n", tok.Line, tok.Column, tok.Type, tok.Value)
        if tok.Type == TOKEN_EOF {
            break
        }
    }
}
```

**Output:**
```
Tokens:
  1:1 1 "foo"
  1:5 4 "+"
  1:7 2 "123"
  1:11 9 "*"
  1:13 3 "("
  1:14 1 "bar"
  1:18 5 "-"
  1:20 2 "456"
  1:23 4 ")"
  1:24 0 ""
```

**Key insight:** Track line and column for error messages!

### 2. AST Design with Interfaces

**Use interfaces for polymorphic AST nodes!**

```go
package main

import (
    "fmt"
    "strings"
)

// All AST nodes implement this interface
type Expr interface {
    String() string
    Accept(Visitor) interface{}
}

// Number literal
type NumberExpr struct {
    Value int
}

func (n *NumberExpr) String() string {
    return fmt.Sprintf("%d", n.Value)
}

func (n *NumberExpr) Accept(v Visitor) interface{} {
    return v.VisitNumber(n)
}

// Variable reference
type IdentExpr struct {
    Name string
}

func (i *IdentExpr) String() string {
    return i.Name
}

func (i *IdentExpr) Accept(v Visitor) interface{} {
    return v.VisitIdent(i)
}

// Binary operation
type BinaryExpr struct {
    Left  Expr
    Op    string
    Right Expr
}

func (b *BinaryExpr) String() string {
    return fmt.Sprintf("(%s %s %s)", b.Left.String(), b.Op, b.Right.String())
}

func (b *BinaryExpr) Accept(v Visitor) interface{} {
    return v.VisitBinary(b)
}

// Visitor pattern for traversal
type Visitor interface {
    VisitNumber(*NumberExpr) interface{}
    VisitIdent(*IdentExpr) interface{}
    VisitBinary(*BinaryExpr) interface{}
}

func main() {
    // Build AST for: (x + 10) * 2
    ast := &BinaryExpr{
        Left: &BinaryExpr{
            Left:  &IdentExpr{Name: "x"},
            Op:    "+",
            Right: &NumberExpr{Value: 10},
        },
        Op:    "*",
        Right: &NumberExpr{Value: 2},
    }

    fmt.Println("AST:", ast.String())
}
```

**Output:**
```
AST: (((x) + (10)) * (2))
```

### 3. Recursive Descent Parser

**Parse expressions using recursive functions!**

```go
package main

import (
    "fmt"
    "strconv"
)

type Parser struct {
    tokens  []Token
    current int
}

func NewParser(tokens []Token) *Parser {
    return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) peek() Token {
    if p.current >= len(p.tokens) {
        return Token{Type: TOKEN_EOF}
    }
    return p.tokens[p.current]
}

func (p *Parser) advance() Token {
    tok := p.peek()
    if p.current < len(p.tokens) {
        p.current++
    }
    return tok
}

func (p *Parser) expect(typ TokenType) (Token, error) {
    tok := p.peek()
    if tok.Type != typ {
        return tok, fmt.Errorf("expected %v but got %v at %d:%d", typ, tok.Type, tok.Line, tok.Column)
    }
    return p.advance(), nil
}

// Grammar:
// expr   -> term (('+' | '-') term)*
// term   -> factor (('*' | '/') factor)*
// factor -> NUMBER | IDENT | '(' expr ')'

func (p *Parser) ParseExpr() (Expr, error) {
    return p.parseTerm()
}

func (p *Parser) parseTerm() (Expr, error) {
    left, err := p.parseFactor()
    if err != nil {
        return nil, err
    }

    for p.peek().Type == TOKEN_PLUS || p.peek().Type == TOKEN_MINUS {
        op := p.advance()
        right, err := p.parseFactor()
        if err != nil {
            return nil, err
        }
        left = &BinaryExpr{
            Left:  left,
            Op:    op.Value,
            Right: right,
        }
    }

    return left, nil
}

func (p *Parser) parseFactor() (Expr, error) {
    left, err := p.parseAtom()
    if err != nil {
        return nil, err
    }

    for p.peek().Type == TOKEN_STAR || p.peek().Type == TOKEN_SLASH {
        op := p.advance()
        right, err := p.parseAtom()
        if err != nil {
            return nil, err
        }
        left = &BinaryExpr{
            Left:  left,
            Op:    op.Value,
            Right: right,
        }
    }

    return left, nil
}

func (p *Parser) parseAtom() (Expr, error) {
    tok := p.peek()

    switch tok.Type {
    case TOKEN_NUMBER:
        p.advance()
        val, _ := strconv.Atoi(tok.Value)
        return &NumberExpr{Value: val}, nil

    case TOKEN_IDENT:
        p.advance()
        return &IdentExpr{Name: tok.Value}, nil

    case TOKEN_LPAREN:
        p.advance()
        expr, err := p.ParseExpr()
        if err != nil {
            return nil, err
        }
        if _, err := p.expect(TOKEN_RPAREN); err != nil {
            return nil, err
        }
        return expr, nil

    default:
        return nil, fmt.Errorf("unexpected token %v at %d:%d", tok.Type, tok.Line, tok.Column)
    }
}

func main() {
    input := "x + 10 * 2"
    lexer := NewLexer(input)

    var tokens []Token
    for {
        tok := lexer.NextToken()
        if tok.Type == TOKEN_EOF {
            break
        }
        tokens = append(tokens, tok)
    }

    parser := NewParser(tokens)
    ast, err := parser.ParseExpr()
    if err != nil {
        fmt.Println("Parse error:", err)
        return
    }

    fmt.Println("Parsed:", ast.String())
}
```

**Output:**
```
Parsed: ((x) + ((10) * (2)))
```

**Key insight:** Operator precedence is encoded in the grammar structure!

### 4. Error Handling with Position Tracking

**Provide helpful error messages!**

```go
package main

import (
    "fmt"
    "strings"
)

type ParseError struct {
    Message string
    Line    int
    Column  int
    Input   string
}

func (e *ParseError) Error() string {
    lines := strings.Split(e.Input, "\n")
    if e.Line > 0 && e.Line <= len(lines) {
        line := lines[e.Line-1]
        pointer := strings.Repeat(" ", e.Column-1) + "^"
        return fmt.Sprintf("%s at line %d, column %d:\n%s\n%s",
            e.Message, e.Line, e.Column, line, pointer)
    }
    return fmt.Sprintf("%s at line %d, column %d", e.Message, e.Line, e.Column)
}

func (p *Parser) parseAtomWithError() (Expr, error) {
    tok := p.peek()

    switch tok.Type {
    case TOKEN_NUMBER:
        p.advance()
        val, _ := strconv.Atoi(tok.Value)
        return &NumberExpr{Value: val}, nil

    case TOKEN_IDENT:
        p.advance()
        return &IdentExpr{Name: tok.Value}, nil

    case TOKEN_LPAREN:
        p.advance()
        expr, err := p.ParseExpr()
        if err != nil {
            return nil, err
        }
        if _, err := p.expect(TOKEN_RPAREN); err != nil {
            return nil, &ParseError{
                Message: "missing closing parenthesis",
                Line:    tok.Line,
                Column:  tok.Column,
                Input:   p.originalInput,
            }
        }
        return expr, nil

    default:
        return nil, &ParseError{
            Message: fmt.Sprintf("unexpected token '%s'", tok.Value),
            Line:    tok.Line,
            Column:  tok.Column,
            Input:   p.originalInput,
        }
    }
}

func main() {
    input := "x + (10 * 2"  // Missing )
    // Parse and show error
    fmt.Println("Input:", input)
    // ... parse ...
    // Output:
    // missing closing parenthesis at line 1, column 5:
    // x + (10 * 2
    //     ^
}
```

### 5. Type Switches for AST Traversal

**Pattern match on node types!**

```go
package main

import (
    "fmt"
)

// Evaluate an expression
func Eval(expr Expr, vars map[string]int) (int, error) {
    switch e := expr.(type) {
    case *NumberExpr:
        return e.Value, nil

    case *IdentExpr:
        val, ok := vars[e.Name]
        if !ok {
            return 0, fmt.Errorf("undefined variable: %s", e.Name)
        }
        return val, nil

    case *BinaryExpr:
        left, err := Eval(e.Left, vars)
        if err != nil {
            return 0, err
        }
        right, err := Eval(e.Right, vars)
        if err != nil {
            return 0, err
        }

        switch e.Op {
        case "+":
            return left + right, nil
        case "-":
            return left - right, nil
        case "*":
            return left * right, nil
        case "/":
            if right == 0 {
                return 0, fmt.Errorf("division by zero")
            }
            return left / right, nil
        default:
            return 0, fmt.Errorf("unknown operator: %s", e.Op)
        }

    default:
        return 0, fmt.Errorf("unknown expression type: %T", expr)
    }
}

func main() {
    // AST for: x + 10 * 2
    ast := &BinaryExpr{
        Left: &IdentExpr{Name: "x"},
        Op:   "+",
        Right: &BinaryExpr{
            Left:  &NumberExpr{Value: 10},
            Op:    "*",
            Right: &NumberExpr{Value: 2},
        },
    }

    vars := map[string]int{"x": 5}
    result, err := Eval(ast, vars)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Printf("Result: %d\n", result)  // 5 + 10*2 = 25
}
```

**Output:**
```
Result: 25
```

### 6. Visitor Pattern for Complex Traversals

**Separate tree traversal from operations!**

```go
package main

import (
    "fmt"
)

// Evaluator visitor
type EvalVisitor struct {
    vars map[string]int
}

func (ev *EvalVisitor) VisitNumber(n *NumberExpr) interface{} {
    return n.Value
}

func (ev *EvalVisitor) VisitIdent(i *IdentExpr) interface{} {
    val, ok := ev.vars[i.Name]
    if !ok {
        panic(fmt.Sprintf("undefined: %s", i.Name))
    }
    return val
}

func (ev *EvalVisitor) VisitBinary(b *BinaryExpr) interface{} {
    left := b.Left.Accept(ev).(int)
    right := b.Right.Accept(ev).(int)

    switch b.Op {
    case "+":
        return left + right
    case "-":
        return left - right
    case "*":
        return left * right
    case "/":
        return left / right
    default:
        panic("unknown op: " + b.Op)
    }
}

// Pretty printer visitor
type PrintVisitor struct{}

func (pv *PrintVisitor) VisitNumber(n *NumberExpr) interface{} {
    return fmt.Sprintf("%d", n.Value)
}

func (pv *PrintVisitor) VisitIdent(i *IdentExpr) interface{} {
    return i.Name
}

func (pv *PrintVisitor) VisitBinary(b *BinaryExpr) interface{} {
    left := b.Left.Accept(pv).(string)
    right := b.Right.Accept(pv).(string)
    return fmt.Sprintf("(%s %s %s)", left, b.Op, right)
}

func main() {
    ast := &BinaryExpr{
        Left:  &IdentExpr{Name: "x"},
        Op:    "+",
        Right: &NumberExpr{Value: 10},
    }

    // Evaluate
    ev := &EvalVisitor{vars: map[string]int{"x": 5}}
    result := ast.Accept(ev)
    fmt.Println("Result:", result)

    // Pretty print
    pv := &PrintVisitor{}
    pretty := ast.Accept(pv)
    fmt.Println("Pretty:", pretty)
}
```

**Output:**
```
Result: 15
Pretty: (x + 10)
```

**Key insight:** Visitor pattern lets you add new operations without modifying AST nodes!

### 7. Keyword Handling

**Reserved words need special treatment!**

```go
package main

import (
    "fmt"
    "strings"
)

var keywords = map[string]TokenType{
    "MATCH":  TOKEN_MATCH,
    "WHERE":  TOKEN_WHERE,
    "RETURN": TOKEN_RETURN,
    "CREATE": TOKEN_CREATE,
    "DELETE": TOKEN_DELETE,
}

const (
    TOKEN_EOF TokenType = iota
    TOKEN_IDENT
    TOKEN_MATCH
    TOKEN_WHERE
    TOKEN_RETURN
    TOKEN_CREATE
    TOKEN_DELETE
)

func (l *Lexer) scanIdentifier() Token {
    start := l.pos
    startCol := l.column

    for l.pos < len(l.input) && isIdentChar(l.input[l.pos]) {
        l.pos++
        l.column++
    }

    value := l.input[start:l.pos]
    upper := strings.ToUpper(value)

    // Check if it's a keyword
    if typ, ok := keywords[upper]; ok {
        return Token{
            Type:   typ,
            Value:  upper,
            Line:   l.line,
            Column: startCol,
        }
    }

    return Token{
        Type:   TOKEN_IDENT,
        Value:  value,
        Line:   l.line,
        Column: startCol,
    }
}

func isIdentChar(ch byte) bool {
    return (ch >= 'a' && ch <= 'z') ||
        (ch >= 'A' && ch <= 'Z') ||
        (ch >= '0' && ch <= '9') ||
        ch == '_'
}
```

## Pre-Implementation Exercises

### Exercise 1: Build a Simple Tokenizer

```go
package main

import (
    "fmt"
)

// TODO: Implement a tokenizer for these tokens:
// - Identifiers (alphanumeric + underscore)
// - Numbers (integers only)
// - Operators: +, -, *, /
// - Parentheses: (, )
// - Comparison: =, <, >, <=, >=, !=

type TokenType int

const (
    TOKEN_EOF TokenType = iota
    // TODO: Add token types
)

type Token struct {
    Type   TokenType
    Value  string
    Line   int
    Column int
}

type Lexer struct {
    // TODO: Add fields
}

func NewLexer(input string) *Lexer {
    // TODO: Implement
    return nil
}

func (l *Lexer) NextToken() Token {
    // TODO: Implement
    return Token{}
}

func main() {
    input := "age >= 18 AND name = 'Alice'"
    lexer := NewLexer(input)

    for {
        tok := lexer.NextToken()
        fmt.Printf("%v %q\n", tok.Type, tok.Value)
        if tok.Type == TOKEN_EOF {
            break
        }
    }
}
```

### Exercise 2: Parse Arithmetic Expressions

```go
package main

// TODO: Implement parser for expressions with:
// - Numbers and identifiers
// - Operators: +, -, *, / (with correct precedence!)
// - Parentheses for grouping

// Grammar:
// expr   -> term (('+' | '-') term)*
// term   -> factor (('*' | '/') factor)*
// factor -> NUMBER | IDENT | '(' expr ')'

type Expr interface {
    String() string
}

type Parser struct {
    // TODO: Add fields
}

func NewParser(tokens []Token) *Parser {
    // TODO: Implement
    return nil
}

func (p *Parser) ParseExpr() (Expr, error) {
    // TODO: Implement
    return nil, nil
}

func main() {
    // TODO: Parse "2 + 3 * 4" and verify AST is correct
    // Should be: (2 + (3 * 4)) not ((2 + 3) * 4)
}
```

### Exercise 3: AST Evaluation

```go
package main

// TODO: Implement an evaluator for your AST from Exercise 2

func Eval(expr Expr, vars map[string]int) (int, error) {
    // TODO: Use type switch to evaluate
    return 0, nil
}

func main() {
    // TODO: Parse "x + y * 2" and evaluate with x=10, y=5
    // Expected result: 10 + 5*2 = 20
}
```

### Exercise 4: Error Messages with Context

```go
package main

// TODO: Enhance your parser to provide helpful error messages

type ParseError struct {
    Message string
    Line    int
    Column  int
    Source  string
}

func (e *ParseError) Error() string {
    // TODO: Format like:
    // Error at line 2, column 5:
    //   x + * 3
    //       ^
    // Unexpected token '*'
    return ""
}

func main() {
    // TODO: Test with invalid inputs:
    // - "2 + + 3" (unexpected +)
    // - "2 + (3" (missing ))
    // - "2 +" (unexpected EOF)
}
```

### Exercise 5: Simple Cypher MATCH Parser

```go
package main

// TODO: Parse simple MATCH statements:
// MATCH (n:Person) WHERE n.age > 18 RETURN n.name

type CypherStmt interface {
    String() string
}

type MatchStmt struct {
    Pattern   *Pattern
    Where     Expr
    Return    []string
}

type Pattern struct {
    NodeVar   string
    NodeLabel string
}

func ParseCypher(input string) (CypherStmt, error) {
    // TODO: Implement
    return nil, nil
}

func main() {
    // TODO: Parse and print AST
}
```

## Performance Benchmarks

### Benchmark 1: Tokenization Speed

```go
func BenchmarkLexer(b *testing.B) {
    input := "MATCH (p:Person)-[:KNOWS]->(f:Person) WHERE p.age > 18 RETURN p.name, f.name"

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        lexer := NewLexer(input)
        for {
            tok := lexer.NextToken()
            if tok.Type == TOKEN_EOF {
                break
            }
        }
    }
}
```

**Target: < 1μs for typical queries**

### Benchmark 2: Parser Performance

```go
func BenchmarkParser(b *testing.B) {
    input := "x + y * (z - 10) / 2"
    lexer := NewLexer(input)

    var tokens []Token
    for {
        tok := lexer.NextToken()
        tokens = append(tokens, tok)
        if tok.Type == TOKEN_EOF {
            break
        }
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        parser := NewParser(tokens)
        _, _ = parser.ParseExpr()
    }
}
```

**Target: < 10μs for typical expressions**

## Common Gotchas to Avoid

### Gotcha 1: Not Handling EOF Properly

```go
// WRONG: Doesn't check EOF
func (p *Parser) peek() Token {
    return p.tokens[p.current]  // Panic if current >= len!
}

// RIGHT: Return EOF token
func (p *Parser) peek() Token {
    if p.current >= len(p.tokens) {
        return Token{Type: TOKEN_EOF}
    }
    return p.tokens[p.current]
}
```

### Gotcha 2: Incorrect Operator Precedence

```go
// WRONG: All operators have same precedence
func (p *Parser) parseExpr() (Expr, error) {
    left, _ := p.parseAtom()
    for isOperator(p.peek()) {
        op := p.advance()
        right, _ := p.parseAtom()
        left = &BinaryExpr{left, op.Value, right}
    }
    return left, nil
}
// Parses "2 + 3 * 4" as "(2 + 3) * 4" - WRONG!

// RIGHT: Separate precedence levels
func (p *Parser) parseExpr() (Expr, error) {
    return p.parseAddSub()  // Lowest precedence
}

func (p *Parser) parseAddSub() (Expr, error) {
    left, _ := p.parseMulDiv()
    for p.peek().Type == TOKEN_PLUS || p.peek().Type == TOKEN_MINUS {
        op := p.advance()
        right, _ := p.parseMulDiv()
        left = &BinaryExpr{left, op.Value, right}
    }
    return left, nil
}

func (p *Parser) parseMulDiv() (Expr, error) {
    left, _ := p.parseAtom()
    for p.peek().Type == TOKEN_STAR || p.peek().Type == TOKEN_SLASH {
        op := p.advance()
        right, _ := p.parseAtom()
        left = &BinaryExpr{left, op.Value, right}
    }
    return left, nil
}
```

### Gotcha 3: Forgetting Position Tracking

```go
// WRONG: No position info
type Token struct {
    Type  TokenType
    Value string
}

// RIGHT: Track for error messages
type Token struct {
    Type   TokenType
    Value  string
    Line   int
    Column int
}
```

### Gotcha 4: Mutating Tokens During Parse

```go
// WRONG: Modifying token slice
func (p *Parser) advance() Token {
    tok := p.tokens[0]
    p.tokens = p.tokens[1:]  // Creates garbage!
    return tok
}

// RIGHT: Use index
func (p *Parser) advance() Token {
    tok := p.peek()
    p.current++
    return tok
}
```

## Checklist Before Starting Lesson 3.1

- [ ] I understand tokenization and lexical analysis
- [ ] I can implement a recursive descent parser
- [ ] I know how to design AST nodes with interfaces
- [ ] I can use type switches for AST traversal
- [ ] I understand the visitor pattern
- [ ] I can handle operator precedence correctly
- [ ] I know how to track source positions for errors
- [ ] I can write helpful error messages
- [ ] I understand the difference between keywords and identifiers
- [ ] I've benchmarked lexer and parser performance

## Next Steps

Once you've completed these exercises and understand the concepts:

**→ Start the main implementation:** `complete-graph-database-learning-path-go1.25.md` Lesson 3.1

You'll implement:
- Cypher tokenizer with keyword handling
- Recursive descent parser for MATCH, WHERE, RETURN
- AST nodes for patterns, predicates, expressions
- Error messages with source context
- AST visitor for query planning
- Test suite with valid and invalid queries

**Time estimate:** 20-25 hours for full implementation

**Parsing is where the query journey begins!** 🚀
