package gbin

import "encoding/binary"

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
)
