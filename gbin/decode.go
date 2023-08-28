package gbin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

type decodeTransformer struct {
	data *bytes.Buffer
}

func newDecodeTransformer(buf bytes.Buffer) *decodeTransformer {
	return &decodeTransformer{
		data: &buf,
	}
}

func (t *decodeTransformer) decode() (*reflect.Value, error) {
	if t.data.Len() < 8 {
		return nil, fmt.Errorf("no header found")
	}

	control, _ := t.data.ReadByte()

	payloadLenBuff, _ := t.readN(7)
	payloadLenArr := []byte{0x00}
	payloadLenArr = append(payloadLenArr, payloadLenBuff.Bytes()...)
	payloadLen := binary.BigEndian.Uint64(payloadLenArr)

	switch EncodedType(control) {
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
	case FLOAT:
		return t.decode_float(payloadLen)
	default:
		return nil, fmt.Errorf("encoding error: unknown type code %b", control)
	}
}

func (t *decodeTransformer) decode_map(stop uint64) (*reflect.Value, error) {
	fmt.Println("decode map")
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	dec := newDecodeTransformer(*buf)
	var kType, vType *reflect.Type
	keyArr := []*reflect.Value{}
	valArr := []*reflect.Value{}
	for {
		if dec.data.Len() == 0 {
			break
		}
		k, err := dec.decode()
		if err != nil {
			return nil, err
		}
		v, err := dec.decode()
		if err != nil {
			return nil, err
		}
		if kType == nil && vType == nil {
			kT := k.Type()
			kType = &kT
			vT := v.Type()
			vType = &vT
		} else {
			if k.Type().Kind() != k.Type().Kind() {
				return nil, fmt.Errorf("map key types must be consistent")
			} else if v.Type().Kind() != v.Type().Kind() {
				return nil, fmt.Errorf("maps to interfaces are not supported")
			}
		}
		keyArr = append(keyArr, k)
		valArr = append(valArr, v)
	}
	mapType := reflect.MapOf(*kType, *vType)
	m := reflect.New(mapType).Elem()
	for i, k := range keyArr {
		v := valArr[i]
		m.SetMapIndex(*k, *v)
	}

	return &m, nil
}

func (t *decodeTransformer) decode_struct(stop uint64) (*reflect.Value, error) {
	fmt.Println("decode struct")

	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	fields := []reflect.StructField{}
	kvMap := map[string]*reflect.Value{}
	dec := newDecodeTransformer(*buf)
	for {
		if dec.data.Len() == 0 {
			break
		}
		key, err := dec.decode()
		if err != nil {
			return nil, err
		}
		if key.Kind() != reflect.String {
			return nil, fmt.Errorf("encoded struct key must be of type string, not %s", key.Kind())
		}
		val, err := dec.decode()
		if err != nil {
			return nil, err
		}
		field := reflect.StructField{
			Name: key.String(),
			Type: val.Type(),
		}
		fields = append(fields, field)
		kvMap[key.String()] = val
	}
	strType := reflect.StructOf(fields)
	str := reflect.New(strType).Elem()
	for k, v := range kvMap {
		str.FieldByName(k).Set(*v)
	}
	return &str, nil
}

func (t *decodeTransformer) decode_ptr(stop uint64) (*reflect.Value, error) {
	fmt.Println("decode ptr")

	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	dec := newDecodeTransformer(*buf)
	inner, err := dec.decode()
	if err != nil {
		return nil, err
	}
	outer := inner.Addr()
	return &outer, nil
}

func (t *decodeTransformer) decode_slice(stop uint64) (*reflect.Value, error) {
	fmt.Println("decode slice")

	var nilSlice []any
	slice := reflect.New(reflect.TypeOf(nilSlice)).Elem()
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	dec := newDecodeTransformer(*buf)
	for {
		if dec.data.Len() == 0 {
			break
		}
		val, err := dec.decode()
		if err != nil {
			return nil, err
		}
		reflect.Append(slice, *val)
	}
	return &slice, nil
}

func (t *decodeTransformer) decode_string(stop uint64) (*reflect.Value, error) {
	buf, err := t.readN(stop)
	if err != nil {
		return nil, err
	}
	str := string(buf.Bytes())

	val := reflect.New(reflect.TypeOf(str)).Elem()
	val.SetString(str)
	return &val, nil

}

func (t *decodeTransformer) decode_bool(stop uint64) (*reflect.Value, error) {
	fmt.Println("decode bool")

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
	return &val, nil
}

func (t *decodeTransformer) decode_int(stop uint64) (*reflect.Value, error) {
	fmt.Println("decode int")

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
	return &val, nil
}

func (t *decodeTransformer) decode_float(stop uint64) (*reflect.Value, error) {
	fmt.Println("decode float")

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
