package exprparser

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"1 + 2", "(1 + 2)"},
		{"1 + 2 * 3", "(1 + (2 * 3))"},
		{"(1 + 2) * 3", "((1 + 2) * 3)"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// TODO: Parse and verify AST
			t.Skip("not implemented")
		})
	}
}

func TestEval(t *testing.T) {
	// TODO: Test expression evaluation
	t.Skip("not implemented")
}

func TestErrors(t *testing.T) {
	// TODO: Test error messages
	t.Skip("not implemented")
}
