package gbin_test

import (
	"fmt"
	"testing"

	"github.com/lspaccatrosi16/go-cli-tools/gbin"
)

type TestData struct {
	Bar
	a string
	b int
	c *map[string]int
}

type Bar struct {
	a string
	b string
	d uintptr
}

func TestEncode(t *testing.T) {
	testStruct := TestData{
		a: "foooooooobar",
		b: 77821124512,
		c: &map[string]int{
			"asdasda": 212431412,
			"bbbbbb":  644387,
			"ccccccc": 7677,
		},
	}

	encoder := gbin.New_Encoder[TestData]()
	out, err := encoder.Encode(&testStruct)
	if err != nil {
		fmt.Println(err.Error())
		t.Fail()
	} else {

		fmt.Printf("% x \n", out)
		fmt.Printf("%s \n", []byte(out))
	}

}
