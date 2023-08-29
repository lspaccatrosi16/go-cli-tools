package gbin_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/lspaccatrosi16/go-cli-tools/gbin"
)

func TestString(t *testing.T) {
	data := "abcd"
	if pass := runTest(data); !pass {
		t.Fail()
	}
}

func TestInt(t *testing.T) {
	data := int(622711)
	if pass := runTest(data); !pass {
		t.Fail()
	}
}

func TestFloat(t *testing.T) {
	data := 1541523.21231
	if pass := runTest(data); !pass {
		t.Fail()
	}
}

func TestBool(t *testing.T) {
	data := true
	if pass := runTest(data); !pass {
		t.Fail()
	}
}

func TestPtr(t *testing.T) {
	data := "abcdefg"
	if pass := runTest(&data); !pass {
		t.Fail()
	}
}

func TestSlice(t *testing.T) {
	data := []string{"a", "b", "c", "d", "E", "f"}
	if pass := runTest(data); !pass {
		t.Fail()
	}
}

func TestMap(t *testing.T) {
	testMap := map[string]int{
		"a": 2,
		"b": 3,
		"c": 4,
	}
	if pass := runTest(testMap); !pass {
		t.Fail()
	}
}

func TestStruct(t *testing.T) {
	testStruct := struct {
		A string
		B map[string]int
		c bool
	}{
		A: "Hi there",
		B: map[string]int{
			"1":   2,
			"2":   3,
			"2+2": 5,
		},
		c: false,
	}
	if pass := runTest(testStruct); !pass {
		t.Fail()
	}
}

func TestError(t *testing.T) {
	testData := map[string]struct {
		A string
		B interface{}
	}{
		"a": {A: "ma", B: 67},
		"b": {A: "foo", B: "bar"},
	}
	if pass := runTest(testData); pass {
		t.Fail()
	}
}

func runTest[T any](data T) bool {
	encoder := gbin.NewEncoder[T]()
	decoder := gbin.NewDecoder[T]()
	encoded, err := encoder.Encode(&data)
	if err != nil {
		fmt.Println("ENCODE ERROR:")
		fmt.Println(err)
		return false
	}

	decoded, err := decoder.Decode(encoded)
	if err != nil {
		fmt.Printf("% x\n", encoded)
		fmt.Println("DECODE ERROR:")
		fmt.Println(err)
		return false
	}

	if reflect.DeepEqual(data, *decoded) {
		return true
	} else {
		fmt.Println("ORIGINAL")
		fmt.Printf("%#v\n", data)
		fmt.Println("DECODED")
		fmt.Printf("%#v\n", *decoded)
		return false
	}
}
