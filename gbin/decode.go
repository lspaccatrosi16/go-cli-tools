package gbin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"

	"github.com/lspaccatrosi16/go-cli-tools/stack"
)

type decodeTransformer struct {
	data  *bytes.Buffer
	stack *stack.Stack[string]
}

func newDecodeTransformer(buf bytes.Buffer, stack *stack.Stack[string]) *decodeTransformer {
	return &decodeTransformer{
		data:  &buf,
		stack: stack,
	}
}

func (t *decodeTransformer) trace() string {
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

func (t *decodeTransformer) decode() (*reflect.Value, error) {
	if t.data.Len() < 1 {
		return nil, fmt.Errorf("no header found")
	}
	control, _ := t.data.ReadByte()
	objectType := control >> 3
	lenLen := uint64(control & 0b00000111)
	payloadLenBuff, _ := t.readN(lenLen)
	payloadLenBuffBytes := payloadLenBuff.Bytes()
	payloadLenArr := make([]byte, 8-lenLen)
	payloadLenArr = append(payloadLenArr, payloadLenBuffBytes...)
	payloadLen := binary.BigEndian.Uint64(payloadLenArr)
	switch EncodedType(objectType) {
	case INTERFACE:
		return t.decode_interface(payloadLen)
	case MAP:
		return t.decode_map(payloadLen)
	case STRUCT:
		return t.decode_struct(payloadLen)
	case PTR:
		return t.decode_ptr(payloadLen)
	case SLICE:
		return t.decode_slice(payloadLen)
	case STRING:
		return t.decode_string(payloadLen)
	case BOOL:
		return t.decode_bool(payloadLen)
	case INT:
		return t.decode_int(payloadLen)
	case INT64:
		return t.decode_int64(payloadLen)
	case UINT:
		return t.decode_uint(payloadLen)
	case UINT64:
		return t.decode_uint64(payloadLen)
	case UINT8:
		return t.decode_uint8(payloadLen)
	case FLOAT64:
		return t.decode_float(payloadLen)
	case INVALID:
		return t.decode_nil(payloadLen)
	default:
		return nil, fmt.Errorf("encoding error: unknown type code 0x%x", control)
	}
}

//why is this necessary to trick it into making an interface for me?
func iface() []interface{} {
	return []interface{}{struct{}{}, "a"}
}

func (t *decodeTransformer) decode_interface(stop uint64) (*reflect.Value, error) {
	t.stack.Push("interface")
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	dec := newDecodeTransformer(*buf, t.stack)
	inner, err := dec.decode()
	if err != nil {
		return nil, err
	}
	interfaceType := reflect.TypeOf(iface()).Elem()
	outer := *inner
	if inner.Kind() != reflect.Invalid {
		outer = inner.Convert(interfaceType)
	} else {
		outer = reflect.New(interfaceType).Elem()
	}
	t.stack.Pop()
	return &outer, nil
}

func (t *decodeTransformer) decode_map(stop uint64) (*reflect.Value, error) {
	t.stack.Push("map")
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	dec := newDecodeTransformer(*buf, t.stack)
	zeroKey, err := dec.decode()
	if err != nil {
		return nil, err
	}
	zeroVal, err := dec.decode()
	if err != nil {
		return nil, err
	}
	kType := zeroKey.Type()
	vType := zeroVal.Type()
	kKind := kType.Kind()
	vKind := vType.Kind()
	keyArr := []*reflect.Value{}
	valArr := []*reflect.Value{}
	count := 0
	for {
		if dec.data.Len() == 0 {
			break
		}
		t.stack.Push(fmt.Sprintf("key%d", count))
		k, err := dec.decode()
		if err != nil {
			return nil, err
		}
		t.stack.Pop()
		t.stack.Push(fmt.Sprintf("%v", k.Interface()))
		v, err := dec.decode()
		if err != nil {
			return nil, err
		}
		t.stack.Pop()

		if k.Kind() != kKind {
			return nil, fmt.Errorf("map key types must be consistent (found %s but expected %s)", k.Kind(), kKind)
		} else if v.Kind() != vKind {
			return nil, fmt.Errorf("maps to interfaces are not supported")
		}

		keyArr = append(keyArr, k)
		valArr = append(valArr, v)
		count++
	}
	kType, fk := kindComparableType[kKind]
	if !fk {
		return nil, fmt.Errorf("found illegal key type for map: %s", kKind)
	}
	mapType := reflect.MapOf(kType, vType)
	m := reflect.MakeMap(mapType)
	for i, k := range keyArr {
		v := valArr[i]
		m.SetMapIndex(*k, *v)
	}
	t.stack.Pop()
	return &m, nil
}

func (t *decodeTransformer) decode_struct(stop uint64) (*reflect.Value, error) {
	t.stack.Push("struct")
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	fields := []reflect.StructField{}
	kvMap := map[string]*reflect.Value{}
	dec := newDecodeTransformer(*buf, t.stack)
	count := 0
	for {
		if dec.data.Len() == 0 {
			break
		}
		t.stack.Push(fmt.Sprintf("key%d", count))
		key, err := dec.decode()
		if err != nil {
			return nil, err
		}
		if key.Kind() != reflect.String {
			return nil, fmt.Errorf("encoded struct key must be of type string, not %s", key.Kind())
		}
		t.stack.Pop()
		t.stack.Push(fmt.Sprintf("%v", key.String()))
		val, err := dec.decode()
		if err != nil {
			return nil, err
		}
		t.stack.Pop()
		field := reflect.StructField{
			Name: key.String(),
			Type: val.Type(),
		}
		fields = append(fields, field)
		kvMap[key.String()] = val
	}
	count++
	strType := reflect.StructOf(fields)
	str := reflect.New(strType).Elem()
	for k, v := range kvMap {
		str.FieldByName(k).Set(*v)
	}
	t.stack.Pop()
	return &str, nil
}

func (t *decodeTransformer) decode_ptr(stop uint64) (*reflect.Value, error) {
	t.stack.Push("ptr")
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	dec := newDecodeTransformer(*buf, t.stack)
	zeroVal, err := dec.decode()
	if err != nil {
		return nil, err
	}
	inner, err := dec.decode()
	if err != nil {
		return nil, err
	}
	outer := reflect.New(zeroVal.Type())
	if inner.Kind() != reflect.Invalid {
		outer.Elem().Set(*inner)
	} else {
		nilVal := reflect.New(zeroVal.Addr().Type()).Elem()
		// fmt.Printf("%s %T\n", nilVal.Kind(), nilVal.Interface())
		// fmt.Printf("%s %T\n", nilVal.Elem().Kind(), nilVal.Elem())
		outer = nilVal
		// if !outer.IsNil() {
		// 	panic("should be nil")

		// }
		//make sure that its nil not ptr to zero value
	}
	t.stack.Pop()
	return &outer, nil
}

func (t *decodeTransformer) decode_slice(stop uint64) (*reflect.Value, error) {
	t.stack.Push("slice")
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	dec := newDecodeTransformer(*buf, t.stack)
	zeroVal, err := dec.decode()
	if err != nil {
		return nil, err
	}
	count := 0
	vals := []*reflect.Value{}
	sliceType := zeroVal.Type()
	for {
		if dec.data.Len() == 0 {
			break
		}
		t.stack.Push(fmt.Sprintf("el%d", count))
		val, err := dec.decode()
		if err != nil {
			return nil, err
		}
		if val.Kind() != sliceType.Kind() {
			return nil, fmt.Errorf("slice key types must be consistent (found %s but expected %s)", val.Kind(), sliceType.Kind())
		}
		t.stack.Pop()
		vals = append(vals, val)
		count++
	}
	slice := reflect.New(reflect.SliceOf(sliceType)).Elem()
	for _, val := range vals {
		slice = reflect.Append(slice, *val)
	}
	t.stack.Pop()
	return &slice, nil
}

func (t *decodeTransformer) decode_string(stop uint64) (*reflect.Value, error) {
	t.stack.Push("string")
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	str := string(buf.Bytes())

	val := reflect.New(reflect.TypeOf(str)).Elem()
	val.SetString(str)
	t.stack.Pop()
	return &val, nil
}

func (t *decodeTransformer) decode_bool(stop uint64) (*reflect.Value, error) {
	t.stack.Push("bool")
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	var bVal bool
	err = binary.Read(buf, BYTE_ORDER, &bVal)
	if err != nil {
		return nil, err
	}
	val := reflect.New(reflect.TypeOf(bVal)).Elem()
	val.SetBool(bVal)
	t.stack.Pop()
	return &val, nil
}

func (t *decodeTransformer) decode_int(stop uint64) (*reflect.Value, error) {
	t.stack.Push("int")
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	var iVal int64
	err = binary.Read(buf, BYTE_ORDER, &iVal)
	if err != nil {
		return nil, err
	}
	val := reflect.New(reflect.TypeOf(int(iVal))).Elem()
	val.SetInt(iVal)
	t.stack.Pop()
	return &val, nil
}

func (t *decodeTransformer) decode_int64(stop uint64) (*reflect.Value, error) {
	t.stack.Push("int64")
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	var iVal int64
	err = binary.Read(buf, BYTE_ORDER, &iVal)
	if err != nil {
		return nil, err
	}
	val := reflect.New(reflect.TypeOf(iVal)).Elem()
	val.SetInt(iVal)
	t.stack.Pop()
	return &val, nil
}
func (t *decodeTransformer) decode_uint(stop uint64) (*reflect.Value, error) {
	t.stack.Push("uint")
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	var uVal uint64
	err = binary.Read(buf, BYTE_ORDER, &uVal)
	if err != nil {
		return nil, err
	}
	val := reflect.New(reflect.TypeOf(uint(uVal))).Elem()
	val.SetUint(uVal)
	t.stack.Pop()
	return &val, nil
}

func (t *decodeTransformer) decode_uint64(stop uint64) (*reflect.Value, error) {
	t.stack.Push("uint64")
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	var uVal uint64
	err = binary.Read(buf, BYTE_ORDER, &uVal)
	if err != nil {
		return nil, err
	}
	val := reflect.New(reflect.TypeOf(uVal)).Elem()
	val.SetUint(uVal)
	t.stack.Pop()
	return &val, nil
}

func (t *decodeTransformer) decode_uint8(stop uint64) (*reflect.Value, error) {
	t.stack.Push("uint8")
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	var uVal uint8
	err = binary.Read(buf, BYTE_ORDER, &uVal)
	if err != nil {
		return nil, err
	}
	val := reflect.New(reflect.TypeOf(uVal)).Elem()
	val.SetUint(uint64(uVal))
	t.stack.Pop()
	return &val, nil
}

func (t *decodeTransformer) decode_float(stop uint64) (*reflect.Value, error) {
	t.stack.Push("float64")
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	var fVal float64
	err = binary.Read(buf, BYTE_ORDER, &fVal)
	if err != nil {
		return nil, err
	}
	val := reflect.New(reflect.TypeOf(fVal)).Elem()
	val.SetFloat(fVal)
	t.stack.Pop()
	return &val, nil
}

func (t *decodeTransformer) decode_nil(stop uint64) (*reflect.Value, error) {
	t.stack.Push("invalid")
	val := reflect.ValueOf((*interface{})(nil)).Elem()
	t.stack.Pop()
	return &val, nil
}

func (t *decodeTransformer) readN(n uint64) (*bytes.Buffer, error) {
	arr := []byte{}
	for i := uint64(0); i < n; i++ {
		b, e := t.data.ReadByte()
		if e != nil {
			return nil, e
		}

		arr = append(arr, b)

	}
	buf := bytes.NewBuffer(arr)
	return buf, nil
}
