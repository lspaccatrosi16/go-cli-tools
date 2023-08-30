package gbin

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"runtime"

	"github.com/lspaccatrosi16/go-cli-tools/pkgError"
	"github.com/lspaccatrosi16/go-cli-tools/stack"
)

var wrapEncode = pkgError.WrapErrorFactory("gbin/encode")
var wrapDecode = pkgError.WrapErrorFactory("gbin/decode")

func addStack(err error, trace string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s at %s", err.Error(), trace)
}

type Encoder[T any] struct {
}

func NewEncoder[T any]() *Encoder[T] {
	if runtime.GOARCH != "amd64" {
		panic("only supports 64-bit architectures currently")
	}
	if !reflect.ValueOf(*new(T)).IsValid() {
		panic("type parameter must not be an interface{}")
	}

	return &Encoder[T]{}
}

func (e *Encoder[T]) Encode(data *T) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	st, err := e.EncodeStream(data)
	if err != nil {
		return []byte{}, err
	}
	io.Copy(buf, st)
	return buf.Bytes(), nil
}

func (e *Encoder[T]) EncodeStream(data *T) (io.Reader, error) {
	panicked := true
	tf := newEncodeTransformer()
	defer func() {
		if panicked {
			fmt.Println("ENCODE TRACE:")
			fmt.Println(tf.trace())
		}
	}()
	value := reflect.ValueOf(*data)
	encoded, err := tf.encode(value)
	buf := bytes.NewBuffer(encoded)
	panicked = false
	return buf, wrapEncode(addStack(err, tf.trace()))
}

type Decoder[T any] struct {
}

func NewDecoder[T any]() *Decoder[T] {
	if runtime.GOARCH != "amd64" {
		panic("only supports 64-bit architectures currently")
	}
	if !reflect.ValueOf(*new(T)).IsValid() {
		panic("type parameter must not be an interface{}")
	}

	return &Decoder[T]{}
}

func (d *Decoder[T]) Decode(data []byte) (*T, error) {
	buf := bytes.NewBuffer(data)
	return d.DecodeStream(buf)
}

func (d *Decoder[T]) DecodeStream(data io.Reader) (*T, error) {
	panicked := true
	buf := bytes.NewBuffer([]byte{})
	io.Copy(buf, data)
	emptyStack := stack.NewStack[string]()
	tf := newDecodeTransformer(*buf, emptyStack)
	defer func() {
		if panicked {
			fmt.Println("DECODE TRACE:")
			fmt.Println(tf.trace())
		}
	}()
	val, err := tf.decode()
	if err != nil {
		return nil, wrapDecode(addStack(err, tf.trace()))
	}
	as := newAssigner[T]()
	checked, err := as.assign(val)
	if err != nil {
		return nil, wrapDecode(addStack(err, tf.trace()))
	}
	panicked = false
	return checked, nil
}
