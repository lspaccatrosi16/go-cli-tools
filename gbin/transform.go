package gbin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

type transformer struct {
}

func new_transformer() *transformer {
	return &transformer{}
}

const MAX_PAYLOAD_LEN = 0xfffffff

/*
GENERAL SPECIFICATION
======================================================================================

CONTROL CODE: 1 BYTE
PAYLOAD LENGTH: 7 BYTES

PAYLOAD
*/

func (t *transformer) encode(v reflect.Value) ([]byte, error) {
	switch v.Kind() {
	case reflect.Map:
		return t.encode_map(v.MapRange())
	case reflect.Struct:
		return t.encode_struct(v)
	case reflect.Pointer:
		return t.encode_ptr(v)
	case reflect.Slice:
		return t.encode_slice(v)
	case reflect.String:
		return t.encode_string(v.String())
	case reflect.Bool:
		return t.encode_bool(v.Bool())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		ui := v.Uint()
		return t.encode_int(int64(ui))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return t.encode_int(v.Int())
	case reflect.Float32, reflect.Float64:
		return t.encode_float(v.Float())
	default:
		return []byte{}, fmt.Errorf("type: %s is not currently supported for serialization", v.Kind())
	}
}

type EncodedType byte

const (
	FLOAT  EncodedType = 1
	INT                = 2
	BOOL               = 4
	STRING             = 8
	SLICE              = 16
	PTR                = 32
	STRUCT             = 64
	MAP                = 128
)

//PAYLOAD: ENCODED, ENCODED
func (t *transformer) encode_map(m *reflect.MapIter) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	for {
		if !m.Next() {
			break
		}

		k := m.Key()
		v := m.Value()
		kEnc, err := t.encode(k)
		if err != nil {
			return []byte{}, err
		}
		vEnc, err := t.encode(v)
		if err != nil {
			return []byte{}, err
		}
		buf.Write(kEnc)
		buf.Write(vEnc)
	}
	return t.format_encode(MAP, buf.Bytes())

}

//PAYLOAD: I64 field index, ENCODED VALUE
func (t *transformer) encode_struct(value reflect.Value) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	n := value.NumField()
	ty := value.Type()
	for i := 0; i < n; i++ {
		val := value.Field(i)
		field := ty.Field(i)
		fIdx := field.Index
		if len(fIdx) != 1 {
			return []byte{}, fmt.Errorf("embedded structs with multi layer access are not supported")
		}
		binary.Write(buf, binary.BigEndian, int64(fIdx[0]))
		fieldVal, err := t.encode(val)
		if err != nil {
			return []byte{}, err
		}
		buf.Write(fieldVal)
	}
	return t.format_encode(STRUCT, buf.Bytes())
}

// PAYLOAD: ENCODED VALUE POINTED AT
func (t *transformer) encode_ptr(value reflect.Value) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	pointedAt := value.Elem()
	encoded, err := t.encode(pointedAt)
	if err != nil {
		return []byte{}, err
	}
	buf.Write(encoded)
	return t.format_encode(PTR, buf.Bytes())
}

//PAYLOAD: SERIES OF BYTES
func (t *transformer) encode_slice(value reflect.Value) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	n := value.Len()
	for i := 0; i < n; i++ {
		el := value.Index(i)
		encoded, err := t.encode(el)
		if err != nil {
			return []byte{}, err
		}
		buf.Write(encoded)
	}
	return t.format_encode(SLICE, buf.Bytes())
}

//PAYLOAD: BINARY ENCODED STRING AS BYTE ARRAY
func (t *transformer) encode_string(s string) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, []byte(s))
	return t.format_encode(STRING, buf.Bytes())
}

//PAYLOAD: BINARY ENCODED BOOL
func (t *transformer) encode_bool(b bool) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, b)
	return t.format_encode(BOOL, buf.Bytes())
}

//PAYLOAD: BINARY ENCODED INT64
func (t *transformer) encode_int(i int64) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, i)
	return t.format_encode(INT, buf.Bytes())
}

//PAYLOAD: BINARY ENCODED FLOAT64
func (t *transformer) encode_float(f float64) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, f)
	return t.format_encode(FLOAT, buf.Bytes())
}

func (t *transformer) format_encode(control EncodedType, payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	payloadLen := len(payload)
	if payloadLen+1 > MAX_PAYLOAD_LEN {
		return []byte{}, fmt.Errorf("payload too big")
	}

	binary.Write(buf, binary.BigEndian, int64(payloadLen))
	buf.Write(payload)

	s := buf.Bytes()
	s[0] = byte(control)
	return s, nil
}
