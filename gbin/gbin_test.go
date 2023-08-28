package gbin_test

import (
	"fmt"
	"testing"

	"github.com/lspaccatrosi16/go-cli-tools/gbin"
)

type StructTestData struct {
	A string
}

func TestString(t *testing.T) {
	data := "abcd"
	if pass := runTest(data); !pass {
		t.Fail()
	}
}

func TestInt(t *testing.T) {
	data := int64(622711)
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

func TestStruct(t *testing.T) {
	testStruct := StructTestData{
		A: "Hi there",
	}
	if pass := runTest(testStruct); !pass {
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

	fmt.Printf("% x\n", encoded)

	decoded, err := decoder.Decode(encoded)
	if err != nil {
		fmt.Println("DECODE ERROR:")
		fmt.Println(err)
		return false
	}

	fmt.Println("ORIGINAL")
	fmt.Printf("%#v\n", data)
	fmt.Println("DECODED")
	fmt.Printf("%#v\n", *decoded)

	return true
}
