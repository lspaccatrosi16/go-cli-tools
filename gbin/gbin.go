package gbin

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
)

type Encoder[T any] struct {
}

func NewEncoder[T any]() *Encoder[T] {
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

	checked, ok := val.Interface().(T)
	if !ok {
		eBuff := bytes.NewBufferString("Could not convert reflected value to value\n")
		fmt.Fprintf(eBuff, "NAME: \"%s\"\n", val.Type().Name())
		fmt.Fprintln(eBuff, "VALUE:")
		fmt.Fprintf(eBuff, "%#v", val.Interface())
		return nil, errors.New(eBuff.String())
	}

	return &checked, nil
}
