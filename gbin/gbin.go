package gbin

import (
	"bytes"
	"io"
	"reflect"
)

type Encoder[T any] struct {
}

func New_Encoder[T any]() *Encoder[T] {
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
	tf := new_transformer()
	value := reflect.ValueOf(*data)
	encoded, err := tf.encode(value)
	buf := bytes.NewBuffer(encoded)
	return buf, err
}

type Decoder[T any] struct {
}

func New_Decoder[T any]() *Decoder[T] {
	return &Decoder[T]{}
}

func (d *Decoder[T]) Decode(data []byte) *T {
	buf := bytes.NewBuffer(data)
	return d.DecodeStream(buf)
}

func (d *Decoder[T]) DecodeStream(data io.Reader) *T {
	return new(T)
}
