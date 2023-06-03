package evaluator

import (
	"testing"

	"github.com/thedonutfactory/donut/lexer"
	"github.com/thedonutfactory/donut/object"
	"github.com/thedonutfactory/donut/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"10 % 3", 1},
		{"(10 % 3) + 8 % 3", 3},
		{"let five = 5; five++; five", 6},
		{"let five = 5; five--; five", 4},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"\"monkey\" == \"monkey\"", true},
		{"\"lang\" != \"langs\"", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"(2 >= 2) == true", true},
		{"(2 >= 2) == false", false},
		{"(2 <= 2) == true", true},
		{"(2 <= 2) == false", false},
		{"true && true", true},
		{"true && false", false},
		{"false && false", false},
		{"true || true", true},
		{"true || false", true},
		{"false || true", true},
		{"!false && true", true},
		{"true && \"\"", false},
		{"true && \"monkey\"", true},
		{"true && 0", false},
		{"true && 1", true},
		{"true && []", false},
		{"true && [1, 2]", true},
		{"true && {}", false},
		{"true && { \"key\": \"value\" }", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`
	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not a String. Got: %T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. Got: %q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`
	evaluated := testEval(input)

	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. Got: %T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. Got: %q", str.Value)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return 10;
  }

  return 1;
}
`,
			10,
		},
		{
			`
		let f = func(x) {
		  return x;
		  x + 10;
		};
		f(10);`,
			10,
		},
		{
			`
		let f = func(x) {
		   let result = x + 10;
		   return result;
		   return 10;
		};
		f(10);`,
			20,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"Line 0: Type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"Line 0: Type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"Line 0: Unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"Line 0: Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"Line 0: Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"Line 0: Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
if (10 > 1) {
  if (10 > 1) {
    return true + false;
  }

  return 1;
}
`,
			"Line 3: Unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"Line 0: Identifier not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"Line 0: Unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[func(x) { x }];`,
			"Line 0: Unusable as a hash key: FUNCTION",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("No error object returned. Got: %T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("Wrong error message. Expected: %q, Got: %q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestConstStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"const a = 5; a;", 5},
		{"const a = 5 * 5; a", 25},
		{"const a = 5; const b = a; b;", 5},
		{"const a = 5; const b = a; let c = a + b + 5; c", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "func(x) { x + 2; };"
	evaluated := testEval(input)

	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("Object is not Function. Got: %T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("Function has wrong amount of parameters. Got: %+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("Parameter is not 'x'. Got: %q", fn.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fn.Body.String() != expectedBody {
		t.Fatalf("Body is not %q. Got: %q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = func(x) { x; }; identity(5);", 5},
		{"let identity = func(x) { return x; }; identity(5);", 5},
		{"let double = func(x) { x * 2; }; double(5);", 10},
		{"let add = func(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"func(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
let newAdder = func(x) {
	func(y) { x + y };
};

let addTwo = newAdder(2);
addTwo(2);
`

	testIntegerObject(t, testEval(input), 4)
}

func TestEnclosingEnvironments(t *testing.T) {
	input := `
let first = 10;
let second = 10;
let third = 10;

let ourFunction = func(first) {
  let second = 20;

  first + second + third;
};

ourFunction(20) + first + second;`

	testIntegerObject(t, testEval(input), 70)
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "Argument to `len` not supported. Got: INTEGER"},
		{`len("one", "two")`, "Wrong number of arguments. Got: 2, Expected: 1"},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`print("hello", "world!")`, nil},
		{`first([1, 2, 3])`, 1},
		{`first([])`, nil},
		{`first(1)`, "Argument to `first` must be an Array. Got: INTEGER"},
		{`last([1, 2, 3])`, 3},
		{`last([])`, nil},
		{`last(1)`, "Argument to `last` must be an Array. Got: INTEGER"},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([])`, nil},
		{`rest([])`, nil},
		{`push([], 1)`, []int{1}},
		{`push(1, 1)`, "Argument to `push` must be an Array. Got: INTEGER"},
		{`pop([7, 8, 9])`, []int{7, 8}},
		{`pop(["monkey", "lang"])`, object.String{Value: "lang"}},
		{`pop([])`, Null},
		{`pop(1)`, "Argument to `pop` must be an Array. Got: INTEGER"},
		{`pop([1, 2], [1, 2])`, "Wrong number of arguments. Got: 2, Expected: 1"},
		{`split("My name is brad")`, "Wrong number of arguments. Got: 1, Expected: 2"},
		{`split("My name is brad", " ")`, []string{"My", "name", "is", "brad"}},
		{`split("", " ")`, []string{}},
		{`join([])`, "Wrong number of arguments. Got: 1, Expected: 2"},
		{`join([], "")`, object.String{Value: ""}},
		{`join(["My", "name", "is", "brad"], " ")`, object.String{Value: "My name is brad"}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case nil:
			testNullObject(t, evaluated)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("Object is not Error. Got: %T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("Wrong error message. Expected: %q, Got: %q",
					expected, errObj.Message)
			}
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d",
					len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testIntegerObject(t, array.Elements[i], int64(expectedElem))
			}
		}
	}
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)

	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("Object is not an Array. Got: %T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("Array has wrong num of elements. Got:%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestHashLiterals(t *testing.T) {
	input := `let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`

	evaluated := testEval(input)

	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return a Hash. Got: %T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		True.HashKey():                             5,
		False.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. Got: %d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("No pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not an Integer. Got: %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. Expected: %d, Got: %d", result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != Null {
		t.Errorf("object is not Null. Got: %T (%+v)", obj, obj)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not a Boolean. Got: %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. Expected: %t, Got: %t", result.Value, expected)
		return false
	}

	return true
}
