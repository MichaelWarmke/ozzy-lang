package parser

import (
	"fmt"
	"ozzy/ast"
	"ozzy/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
    let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram returned nil program")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements has %d statements, want 3", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements has %d statements, want 1", len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 993322;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram returned nil program")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements has %d statements, want 3", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement, got %T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral got %q, want \"return\"", returnStmt.TokenLiteral())
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has %d statements, want 1", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.Identifier, got %T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value got %q, want \"foobar\"", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral got %q, want \"foobar\"", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has %d statements, want 1", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	testIntegerLiteral(t, stmt.Expression, int64(5))
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has %d statements, want 1", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
	}

	testBoolean(t, stmt.Expression, true)
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	stmt := getFirstExpression(t, input)

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral, got %T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters got %d, want 2", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function literal statements got %d, want 1", len(function.Body.Statements))
	}

	bodystmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function literal statement is not ast.ExpressionStatement, got %T", function.Body.Statements[0])
	}

	testInfixExpression(t, bodystmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{input: "fn() {}", expected: []string{}},
		{input: "fn(x) {}", expected: []string{"x"}},
		{input: "fn(x, y) {}", expected: []string{"x", "y"}},
	}

	for _, tt := range tests {
		stmt := getFirstExpression(t, tt.input)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expected) {
			t.Errorf("function literal parameters wrong. want %d, got=%d\n", len(tt.expected), len(function.Parameters))
		}

		for i, ident := range tt.expected {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 +5);"

	stmt := getFirstExpression(t, input)
	callExp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression, got %T", stmt.Expression)
	}

	if len(callExp.Arguments) != 3 {
		t.Fatalf("callExp.Arguments has %d arguments, want 3", len(callExp.Arguments))
	}

	testLiteralExpression(t, callExp.Function, "add")

	testLiteralExpression(t, callExp.Arguments[0], 1)
	testInfixExpression(t, callExp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callExp.Arguments[2], 4, "+", 5)
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements has %d statements, want 1", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.PrefixExpression, got %T", program.Statements[0])
		}

		if exp.Operator != tt.operator {
			t.Errorf("exp.Operator got %q, want %q", exp.Operator, tt.operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func testIntegerLiteral(t *testing.T, right ast.Expression, value int64) bool {
	integerLiteral, ok := right.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("right is not ast.IntegerLiteral, got %v", right)
	}

	if integerLiteral.Value != value {
		t.Errorf("right value got %d, want %d", integerLiteral.Value, value)
		return false
	}

	if integerLiteral.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integer literal got %q, want %d", integerLiteral.TokenLiteral(), value)
		return false
	}
	return true
}

func testBoolean(t *testing.T, exp ast.Expression, value bool) bool {
	boolean, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp is not ast.IntegerLiteral, got %v", exp)
	}

	if boolean.Value != value {
		t.Errorf("exp value got %t, want %t", boolean.Value, value)
		return false
	}

	if boolean.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("integer literal got %q, want %t", boolean.TokenLiteral(), value)
		return false
	}
	return true
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements has %d statements, want 1", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b",
			"((-a) * b)",
		},
		{"!-a",
			"(!(-a))",
		},
		{"a + b + c",
			"((a + b) + c)",
		},
		{"a + b - c",
			"((a + b) - c)",
		},
		{"a * b * c",
			"((a * b) * c)",
		},
		{"a * b / c",
			"((a * b) / c)",
		},
		{"a + b / c",
			"(a + (b / c))",
		},
		{"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{"true",
			"true",
		},
		{"false",
			"false",
		},
		{"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{"-(5 + 5)",
			"(-(5 + 5))",
		},
		{"!(true == true)",
			"(!(true == true))",
		},
		{"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("program.String() got=%q want=%q", actual, tt.expected)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, input string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier, got %T", exp)
		return false
	}

	if ident.Value != input {
		t.Errorf("ident.Value got %q, want %q", ident.Value, input)
		return false
	}

	if ident.TokenLiteral() != input {
		t.Errorf("ident.TokenLiteral got %q, want %q", ident.TokenLiteral(), input)
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	}

	t.Errorf("type of exp isn't of type %T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression, got %v", exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("opExp.Operator got %q, want %q", opExp.Operator, operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func TestIfExpressionParsing(t *testing.T) {
	input := `if (x < y) { x }`
	stmt := getFirstExpression(t, input)

	testIfExpression(t, stmt.Expression, "")
}

func TestIfElseExpressionParsing(t *testing.T) {
	input := `if (x < y) { x } else { y }`
	stmt := getFirstExpression(t, input)

	testIfExpression(t, stmt.Expression, "y")
}

func testIfExpression(t *testing.T, stmt ast.Expression, altValue string) {
	exp, ok := stmt.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression, got %v", stmt)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got %T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if altValue == "" && exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%T", exp.Alternative)
	}

	if altValue != "" {
		alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got %T", exp.Alternative.Statements[0])
		}

		if !testLiteralExpression(t, alternative.Expression, altValue) {
			t.Errorf("exp.Alternative.String() got %q, want %q", exp.Alternative.String(), altValue)
		}
	}
}

func getFirstExpression(t *testing.T, input string) *ast.ExpressionStatement {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %v", program.Statements[0])
	}

	return stmt
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, e := range errors {
		t.Errorf("parser error: %s", e)
	}
	t.FailNow()
}

func testLetStatement(t *testing.T, stmt ast.Statement, name string) bool {
	if stmt.TokenLiteral() != "let" {
		t.Errorf("stmt.TokenLiteral got %s, want let", stmt.TokenLiteral())
		return false
	}

	letStmt, ok := stmt.(*ast.LetStatement)
	if !ok {
		t.Errorf("stmt not ast.LetStatement, got %T", stmt)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value got %s, want %s", letStmt.Name.Value, name)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral got %s, want %s", letStmt.Name.TokenLiteral(), name)
		return false
	}

	return true
}
