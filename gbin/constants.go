package gbin

import "encoding/binary"

const MAX_PAYLOAD_LEN = 0xfffffff

var BYTE_ORDER = binary.BigEndian

type EncodedType byte

const (
	FLOAT     EncodedType = 1
	INT                   = 2
	BOOL                  = 4
	STRING                = 8
	SLICE                 = 16
	PTR                   = 32
	STRUCT                = 64
	MAP                   = 128
	INTERFACE             = 255
)
