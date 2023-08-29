package gbin

import (
	"bytes"
	"io"
	"reflect"
	"runtime"
)

type Encoder[T any] struct {
}

func NewEncoder[T any]() *Encoder[T] {
	if runtime.GOARCH != "amd64" {
		panic("only supports 64-bit architectures currently")
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
	value := reflect.ValueOf(*data)
	encoded, err := tf.encode(value)
	buf := bytes.NewBuffer(encoded)
	return buf, err
}

type Decoder[T any] struct {
}

func NewDecoder[T any]() *Decoder[T] {
	if runtime.GOARCH != "amd64" {
		panic("only supports 64-bit architectures currently")
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
	val, err := tf.decode()
	if err != nil {
		return nil, err
	}

	as := newAssigner[T]()

	checked, err := as.assign(val)

	if err != nil {
		return nil, err
	}

	return checked, nil
}
