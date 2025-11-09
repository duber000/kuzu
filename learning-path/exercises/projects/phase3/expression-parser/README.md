# Project 3.1: Expression Parser

## Overview
Build a recursive descent parser for SQL-like expressions with proper error messages and AST construction.

**Duration:** 10-12 hours
**Difficulty:** Medium

## Learning Objectives
- Implement recursive descent parsing
- Build Abstract Syntax Trees (AST)
- Handle operator precedence
- Generate helpful error messages
- Support expression evaluation

## Features to Implement

### Expression Types
- Arithmetic: `+, -, *, /, %`
- Comparison: `=, !=, <, >, <=, >=`
- Logical: `AND, OR, NOT`
- Functions: `SUM(), COUNT(), AVG()`
- Literals: integers, floats, strings, booleans

### API Design
```go
type Expr interface {
	Eval(ctx Context) (Value, error)
}

type Parser struct {
	lexer *Lexer
}

func (p *Parser) ParseExpression() (Expr, error)
```

## Test Cases
- Operator precedence: `1 + 2 * 3`
- Parentheses: `(1 + 2) * 3`
- Error handling: `1 + + 2`
- Complex expressions: `(a > 5 AND b < 10) OR c = 20`

## Time Estimate
Core: 6-8 hours, Testing: 2-3 hours, Extensions: 2-3 hours
