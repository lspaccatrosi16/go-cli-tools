package parser

import (
	"errors"
	"testing"
)

func TestParser(t *testing.T) {
	dummyData := Object{
		"foo": "bar",
		"one": Object{
			"two": Object{
				"three": 3,
			},
		},
	}

	dummyString := "Foo is always equal to $foo, and not baz. One + Two = $one.two.three"
	expectedString := "Foo is always equal to bar, and not baz. One + Two = 3"

	runParserTest(t, ParserTest{Input: dummyString, Expected: expectedString, Data: dummyData})

}

func TestParserError(t *testing.T) {
	dummyData := Object{
		"fox": "bar",
	}

	_, parseErrors := ParseString("Hi $foo", dummyData)

	if len(parseErrors) == 0 {
		t.Fatal("expected errors from incorrect input")
	}
}

type ParserTest struct {
	Input    string
	Expected string
	Data     Object
}

func runParserTest(t *testing.T, test ParserTest) {
	parseResult, parseErrors := ParseString(test.Input, test.Data)

	if len(parseErrors) != 0 {
		e := errors.Join(parseErrors...)
		t.Fatal(e.Error())
	}

	if parseResult != test.Expected {
		t.Fatalf("expected '%s' but got '%s'", test.Expected, parseResult)
	}
}
