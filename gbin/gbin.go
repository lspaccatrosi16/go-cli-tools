package gbin

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"runtime"

	"github.com/lspaccatrosi16/go-cli-tools/pkgError"
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
	tf := newEncodeTransformer()
	defer func() {
		err := recover()
		fmt.Println("TRACE:")
		fmt.Println(tf.trace())
		panic(err)
	}()
	value := reflect.ValueOf(*data)
	encoded, err := tf.encode(value)
	buf := bytes.NewBuffer(encoded)
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
	buf := bytes.NewBuffer([]byte{})
	io.Copy(buf, data)
	tf := newDecodeTransformer(*buf)
	defer func() {
		err := recover()
		fmt.Println("TRACE:")
		fmt.Println(tf.trace())
		panic(err)
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
	return checked, nil
}
