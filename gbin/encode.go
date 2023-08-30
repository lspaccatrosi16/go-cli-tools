package gbin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/bits"
	"reflect"

	"github.com/lspaccatrosi16/go-cli-tools/stack"
)

type encodeTransformer struct {
	stack stack.Stack[string]
}

func newEncodeTransformer() *encodeTransformer {
	return &encodeTransformer{}
}

func (t *encodeTransformer) trace() string {
	buf := bytes.NewBufferString("/")
	t.stack.Reverse()
	for {
		val, ok := t.stack.Pop()
		if !ok {
			break
		}
		formatted := fmt.Sprintf("%s/", val)
		buf.WriteString(formatted)
	}

	return buf.String()
}

func (t *encodeTransformer) encode(v reflect.Value) ([]byte, error) {
	switch v.Kind() {
	case reflect.Interface:
		return t.encode_interface(v)
	case reflect.Map:
		return t.encode_map(v)
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
	case reflect.Int:
		return t.encode_int(int(v.Int()))
	case reflect.Int64:
		return t.encode_int64(v.Int())
	case reflect.Uint:
		return t.encode_uint(uint(v.Uint()))
	case reflect.Uint64:
		return t.encode_uint64(v.Uint())
	case reflect.Uint8:
		return t.encode_uint8(uint8(v.Uint()))
	case reflect.Float64:
		return t.encode_float64(v.Float())
	default:
		return []byte{}, fmt.Errorf("type %s is not currently supported for serialization", v.Kind())
	}
}

func (t *encodeTransformer) encode_interface(i reflect.Value) ([]byte, error) {
	it := i.Type()
	stackEntry := fmt.Sprintf("interface(%s)", it.Name())
	buf := bytes.NewBuffer([]byte{})
	t.stack.Push(stackEntry)
	innerVal := i.Elem()
	encoded, err := t.encode(innerVal)
	if err != nil {
		return []byte{}, err
	}
	t.stack.Pop()
	buf.Write(encoded)
	return t.format_encode(INTERFACE, buf.Bytes())
}

//PAYLOAD: ENCODED, ENCODED
func (t *encodeTransformer) encode_map(m reflect.Value) ([]byte, error) {
	mt := m.Type()
	stackEntry := fmt.Sprintf("map[%s]%s(%s)", mt.Key().Kind(), mt.Elem().Kind(), mt.Name())
	t.stack.Push(stackEntry)
	mi := m.MapRange()
	buf := bytes.NewBuffer([]byte{})
	for {
		if !mi.Next() {
			break
		}
		k := mi.Key()
		v := mi.Value()
		kEntry := fmt.Sprintf("key[%v]", k.Interface())
		t.stack.Push(kEntry)
		kEnc, err := t.encode(k)
		if err != nil {
			return []byte{}, err
		}
		t.stack.Pop()
		vEntry := fmt.Sprintf("val[%v]", k.Interface())
		t.stack.Push(vEntry)
		vEnc, err := t.encode(v)
		if err != nil {
			return []byte{}, err
		}
		t.stack.Pop()
		buf.Write(kEnc)
		buf.Write(vEnc)
	}
	t.stack.Pop()
	return t.format_encode(MAP, buf.Bytes())
}

//PAYLOAD: STRING FIELD NAME, ENCODED VALUE
func (t *encodeTransformer) encode_struct(value reflect.Value) ([]byte, error) {
	st := value.Type()
	stackEntry := fmt.Sprintf("struct(%s)", st.Name())
	t.stack.Push(stackEntry)
	buf := bytes.NewBuffer([]byte{})
	n := value.NumField()
	ty := value.Type()
	for i := 0; i < n; i++ {
		val := value.Field(i)
		field := ty.Field(i)
		if !field.IsExported() {
			continue
		}
		fEntry := fmt.Sprintf("field[%s]", field.Name)
		t.stack.Push(fEntry)
		t.stack.Push("key")
		fName := field.Name
		encodedName, err := t.encode(reflect.ValueOf(fName))
		if err != nil {
			return []byte{}, err
		}
		t.stack.Pop()
		t.stack.Push("value")
		fieldVal, err := t.encode(val)
		if err != nil {
			return []byte{}, err
		}
		t.stack.Pop()
		t.stack.Pop()
		buf.Write(encodedName)
		buf.Write(fieldVal)
	}
	t.stack.Pop()
	return t.format_encode(STRUCT, buf.Bytes())
}

// PAYLOAD: ENCODED VALUE POINTED AT
func (t *encodeTransformer) encode_ptr(value reflect.Value) ([]byte, error) {
	t.stack.Push("ptr")
	buf := bytes.NewBuffer([]byte{})
	pointedAt := value.Elem()
	encoded, err := t.encode(pointedAt)
	if err != nil {
		return []byte{}, err
	}
	t.stack.Pop()
	buf.Write(encoded)
	return t.format_encode(PTR, buf.Bytes())
}

//PAYLOAD: SERIES OF BYTES
func (t *encodeTransformer) encode_slice(value reflect.Value) ([]byte, error) {
	st := value.Type()
	stackEntry := fmt.Sprintf("slice[%s]", st.Elem().Kind())
	t.stack.Push(stackEntry)
	buf := bytes.NewBuffer([]byte{})
	n := value.Len()
	for i := 0; i < n; i++ {
		elEntry := fmt.Sprintf("el%d", i)
		t.stack.Push(elEntry)
		el := value.Index(i)
		encoded, err := t.encode(el)
		if err != nil {
			return []byte{}, err
		}
		t.stack.Pop()
		buf.Write(encoded)
	}
	t.stack.Pop()
	return t.format_encode(SLICE, buf.Bytes())
}

//PAYLOAD: BINARY ENCODED STRING AS BYTE ARRAY
func (t *encodeTransformer) encode_string(s string) ([]byte, error) {
	t.stack.Push("string")
	buf := bytes.NewBuffer([]byte{})
	for _, c := range s {
		buf.WriteRune(c)
	}
	t.stack.Pop()
	return t.format_encode(STRING, buf.Bytes())
}

//PAYLOAD: BINARY ENCODED BOOL
func (t *encodeTransformer) encode_bool(b bool) ([]byte, error) {
	t.stack.Push("bool")
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, BYTE_ORDER, b)
	if err != nil {
		return []byte{}, err
	}
	t.stack.Pop()
	return t.format_encode(BOOL, buf.Bytes())
}

//PAYLOAD: BINARY ENCODED INT64
func (t *encodeTransformer) encode_int(i int) ([]byte, error) {
	t.stack.Push("int")
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, BYTE_ORDER, int64(i))
	if err != nil {
		return []byte{}, err
	}
	t.stack.Pop()
	return t.format_encode(INT, buf.Bytes())
}

func (t *encodeTransformer) encode_int64(i int64) ([]byte, error) {
	t.stack.Push("int64")
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, BYTE_ORDER, i)
	if err != nil {
		return []byte{}, err
	}
	t.stack.Pop()
	return t.format_encode(INT64, buf.Bytes())
}

func (t *encodeTransformer) encode_uint(i uint) ([]byte, error) {
	t.stack.Push("uint")
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, BYTE_ORDER, uint64(i))
	if err != nil {
		return []byte{}, err
	}
	t.stack.Pop()
	return t.format_encode(UINT, buf.Bytes())
}

func (t *encodeTransformer) encode_uint64(i uint64) ([]byte, error) {
	t.stack.Push("uint64")
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, BYTE_ORDER, i)
	if err != nil {
		return []byte{}, err
	}
	t.stack.Pop()
	return t.format_encode(UINT64, buf.Bytes())
}

func (t *encodeTransformer) encode_uint8(i uint8) ([]byte, error) {
	t.stack.Push("uint8")
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, BYTE_ORDER, i)
	if err != nil {
		return []byte{}, err
	}
	t.stack.Pop()
	return t.format_encode(UINT8, buf.Bytes())
}

//PAYLOAD: BINARY ENCODED FLOAT64
func (t *encodeTransformer) encode_float64(f float64) ([]byte, error) {
	t.stack.Push("float64")
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, BYTE_ORDER, f)
	if err != nil {
		return []byte{}, err
	}
	t.stack.Pop()
	return t.format_encode(FLOAT64, buf.Bytes())
}

func (t *encodeTransformer) format_encode(objectType EncodedType, payload []byte) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	payloadLen := uint64(len(payload))
	if payloadLen+1 > MAX_PAYLOAD_LEN {
		return []byte{}, fmt.Errorf("payload too big")
	}
	lenLen := (bits.Len64(payloadLen) / 8) + 1
	if lenLen > 3 {
		return []byte{}, fmt.Errorf("payload len does not fit in control byte")
	}
	controlByte := (byte(objectType) << 3) | byte(lenLen)
	buf.WriteByte(controlByte)
	tmpBuf := bytes.NewBuffer([]byte{})
	payloadLenShifted := uint64(payloadLen << (64 - (lenLen * 8)))
	binary.Write(tmpBuf, BYTE_ORDER, payloadLenShifted)
	for i := 0; i < lenLen; i++ {
		buf.WriteByte(tmpBuf.Bytes()[i])
	}
	buf.Write(payload)
	return buf.Bytes(), nil
}
