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

func TestInt64(t *testing.T) {
	data := int64(1234)
	if pass := runTest(data); !pass {
		t.Fail()
	}
}

func TestUint(t *testing.T) {
	data := uint(1234)
	if pass := runTest(data); !pass {
		t.Fail()
	}
}

func TestUint64(t *testing.T) {
	data := uint64(1234)
	if pass := runTest(data); !pass {
		t.Fail()
	}
}

func TestUint8(t *testing.T) {
	data := uint8(122)
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

func TestInterface(t *testing.T) {
	if pass := runTest(struct{ A interface{} }{A: "aaa"}); !pass {
		t.Fail()
	}
}

func TestSimple(t *testing.T) {
	testData := map[string]interface{}{
		"a": (int)(2),
	}
	if pass := runTest(testData); !pass {
		t.Fail()
	}
}

func TestComplex(t *testing.T) {
	a := (int64)(2)
	complexDataStructure := map[string]interface{}{
		"a": struct {
			A string
			B string
			C int
			D int64
		}{"Foo", "Bar", 60, 70},
		"b": map[string]int{"1": 1, "2": 2},
		"c": struct {
			M *map[string]int
			I *int64
		}{
			M: &map[string]int{
				"1": 1,
				"2": 3,
				"3": 6,
			},
			I: &a,
		},
	}
	if pass := runTest(complexDataStructure); !pass {
		t.Fail()
	}
}

func TestEmptySlice(t *testing.T) {
	data := struct {
		A string
		B []string
	}{"AAA", []string{}}
	if pass := runTest(data); !pass {
		t.Fail()
	}
}

func TestEmptyMap(t *testing.T) {
	data := map[string]int{}
	if pass := runTest(data); !pass {
		t.Fail()
	}
}

func TestNil(t *testing.T) {
	data := (*int)(nil)
	if pass := runTest(data); !pass {
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
