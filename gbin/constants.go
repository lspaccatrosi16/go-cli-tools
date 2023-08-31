package gbin

import (
	"encoding/binary"
	"reflect"
)

const MAX_PAYLOAD_LEN = 0xfffffff

var BYTE_ORDER = binary.BigEndian

type EncodedType byte

const (
	INTERFACE EncodedType = iota
	MAP
	STRUCT
	PTR
	SLICE
	STRING
	BOOL
	INT
	INT64
	UINT
	UINT64
	UINT8
	FLOAT64
	INVALID
)

var kindControl map[reflect.Kind]EncodedType = map[reflect.Kind]EncodedType{
	reflect.Invalid:   INVALID,
	reflect.Interface: INTERFACE,
	reflect.Map:       MAP,
	reflect.Struct:    STRUCT,
	reflect.Pointer:   PTR,
	reflect.Slice:     SLICE,
	reflect.String:    STRING,
	reflect.Bool:      BOOL,
	reflect.Int:       INT,
	reflect.Int64:     INT64,
	reflect.Uint:      UINT,
	reflect.Uint64:    UINT64,
	reflect.Uint8:     UINT8,
	reflect.Float64:   FLOAT64,
}

var controlKind map[EncodedType]reflect.Kind = map[EncodedType]reflect.Kind{
	INTERFACE: reflect.Interface,
	MAP:       reflect.Map,
	STRUCT:    reflect.Struct,
	PTR:       reflect.Pointer,
	SLICE:     reflect.Slice,
	STRING:    reflect.String,
	BOOL:      reflect.Bool,
	INT:       reflect.Int,
	INT64:     reflect.Int64,
	UINT:      reflect.Uint,
	UINT64:    reflect.Uint64,
	UINT8:     reflect.Uint8,
	FLOAT64:   reflect.Float64,
	INVALID:   reflect.Invalid,
}

var kindComparableType map[reflect.Kind]reflect.Type = map[reflect.Kind]reflect.Type{
	reflect.String:  reflect.TypeOf(""),
	reflect.Bool:    reflect.TypeOf(true),
	reflect.Int:     reflect.TypeOf(int(0)),
	reflect.Int64:   reflect.TypeOf(int64(0)),
	reflect.Uint:    reflect.TypeOf(uint(0)),
	reflect.Uint64:  reflect.TypeOf(uint64(0)),
	reflect.Uint8:   reflect.TypeOf(uint8(0)),
	reflect.Float64: reflect.TypeOf(float64(0)),
}
