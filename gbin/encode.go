package gbin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

type encodeTransformer struct {
}

func newEncodeTransformer() *encodeTransformer {
	return &encodeTransformer{}
}

/*
GENERAL SPECIFICATION
======================================================================================

HEADER:
CONTROL CODE: 1 BYTE

(treat as UINT64)
PAYLOAD LENGTH: 7 BYTES

PAYLOAD
*/

func (t *encodeTransformer) encode(v reflect.Value) ([]byte, error) {
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
	case reflect.Int, reflect.Int64:
		return t.encode_int(v.Int())
	case reflect.Float64:
		return t.encode_float(v.Float())
	default:
		return []byte{}, fmt.Errorf("type: %s is not currently supported for serialization", v.Kind())
	}
}

//PAYLOAD: ENCODED, ENCODED
func (t *encodeTransformer) encode_map(m *reflect.MapIter) ([]byte, error) {
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

//PAYLOAD: STRING FIELD NAME, ENCODED VALUE
func (t *encodeTransformer) encode_struct(value reflect.Value) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	n := value.NumField()
	ty := value.Type()
	for i := 0; i < n; i++ {
		val := value.Field(i)
		field := ty.Field(i)
		if !field.IsExported() {
			continue
		}
		fName := field.Name
		encodedName, err := t.encode(reflect.ValueOf(fName))
		if err != nil {
			return []byte{}, err
		}
		fieldVal, err := t.encode(val)
		if err != nil {
			return []byte{}, err
		}
		buf.Write(encodedName)
		buf.Write(fieldVal)
	}
	return t.format_encode(STRUCT, buf.Bytes())
}

// PAYLOAD: ENCODED VALUE POINTED AT
func (t *encodeTransformer) encode_ptr(value reflect.Value) ([]byte, error) {
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
func (t *encodeTransformer) encode_slice(value reflect.Value) ([]byte, error) {
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
func (t *encodeTransformer) encode_string(s string) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	for _, c := range s {
		buf.WriteRune(c)
	}
	return t.format_encode(STRING, buf.Bytes())
}

//PAYLOAD: BINARY ENCODED BOOL
func (t *encodeTransformer) encode_bool(b bool) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, BYTE_ORDER, b)
	return t.format_encode(BOOL, buf.Bytes())
}

//PAYLOAD: BINARY ENCODED INT64
func (t *encodeTransformer) encode_int(i int64) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, BYTE_ORDER, i)
	return t.format_encode(INT, buf.Bytes())
}

//PAYLOAD: BINARY ENCODED FLOAT64
func (t *encodeTransformer) encode_float(f float64) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, BYTE_ORDER, f)
	return t.format_encode(FLOAT, buf.Bytes())
}

func (t *encodeTransformer) format_encode(control EncodedType, payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	payloadLen := len(payload)
	if payloadLen+1 > MAX_PAYLOAD_LEN {
		return []byte{}, fmt.Errorf("payload too big")
	}

	binary.Write(buf, BYTE_ORDER, uint64(payloadLen))
	buf.Write(payload)

	s := buf.Bytes()
	s[0] = byte(control)
	return s, nil
}
