package serde

import (
	"encoding/binary"
	"math"
)

const (
	UINT16_SIZE  = 2
	UINT32_SIZE  = 4
	UINT64_SIZE  = 8
	UINT256_SIZE = 32
)

type Uint256 [UINT256_SIZE]byte

var UINT256_EMPTY = Uint256{}

type SourceBuffer struct {
	s   []byte
	off uint64 // current reading index
}

func (sb *SourceBuffer) Len() uint64 {
	length := uint64(len(sb.s))
	if sb.off >= length {
		return 0
	}
	return length - sb.off
}

func (sb *SourceBuffer) Bytes() []byte {
	return sb.s
}

func (sb *SourceBuffer) OffBytes() []byte {
	return sb.s[sb.off:]
}

func (sb *SourceBuffer) Pos() uint64 {
	return sb.off
}

func (sb *SourceBuffer) Size() uint64 { return uint64(len(sb.s)) }

func (sb *SourceBuffer) NextBytes(n uint32) (data []byte, eof bool) {
	m := uint64(len(sb.s))
	end, overflow := SafeAdd(sb.off, uint64(n))
	if overflow || end > m {
		end = m
		eof = true
	}
	data = sb.s[sb.off:end]
	sb.off = end

	return
}

func (sb *SourceBuffer) Skip(n uint64) (eof bool) {
	m := uint64(len(sb.s))
	end, overflow := SafeAdd(sb.off, n)
	if overflow || end > m {
		end = m
		eof = true
	}
	sb.off = end

	return
}

func (sb *SourceBuffer) NextByte() (data byte, eof bool) {
	if sb.off >= uint64(len(sb.s)) {
		return 0, true
	}

	b := sb.s[sb.off]
	sb.off++
	return b, false
}

func (sb *SourceBuffer) NextUint8() (data uint8, eof bool) {
	var val byte
	val, eof = sb.NextByte()
	return uint8(val), eof
}

func (sb *SourceBuffer) NextBool() (data bool, eof bool) {
	val, eof := sb.NextByte()
	if val == 0 {
		data = false
	} else if val == 1 {
		data = true
	} else {
		eof = true
	}
	return
}

func (sb *SourceBuffer) BackUp(n uint64) {
	sb.off -= n
}

func (sb *SourceBuffer) NextUint16() (data uint16, eof bool) {
	var buf []byte
	buf, eof = sb.NextBytes(UINT16_SIZE)
	if eof {
		return
	}

	return binary.BigEndian.Uint16(buf), eof
}

func (sb *SourceBuffer) NextUint32() (data uint32, eof bool) {
	var buf []byte
	buf, eof = sb.NextBytes(UINT32_SIZE)
	if eof {
		return
	}

	return binary.BigEndian.Uint32(buf), eof
}

func (sb *SourceBuffer) NextUint64() (data uint64, eof bool) {
	var buf []byte
	buf, eof = sb.NextBytes(UINT64_SIZE)
	if eof {
		return
	}

	return binary.BigEndian.Uint64(buf), eof
}

func (sb *SourceBuffer) NextInt32() (data int32, eof bool) {
	var val uint32
	val, eof = sb.NextUint32()
	return int32(val), eof
}

func (sb *SourceBuffer) NextInt64() (data int64, eof bool) {
	var val uint64
	val, eof = sb.NextUint64()
	return int64(val), eof
}

func (sb *SourceBuffer) NextInt16() (data int16, eof bool) {
	var val uint16
	val, eof = sb.NextUint16()
	return int16(val), eof
}

func (sb *SourceBuffer) NextVarBytes() (data []byte, eof bool) {
	count, eof := sb.NextUint32()
	if eof {
		return
	}
	data, eof = sb.NextBytes(count)
	return
}

func (sb *SourceBuffer) NextHash() (data Uint256, eof bool) {
	var buf []byte
	buf, eof = sb.NextBytes(UINT256_SIZE)
	if eof {
		return
	}
	copy(data[:], buf)

	return
}

func (sb *SourceBuffer) NextString() (data string, eof bool) {
	var val []byte
	val, eof = sb.NextVarBytes()
	data = string(val)
	return
}

func (sb *SourceBuffer) NextVarUint() (data uint64, eof bool) {
	var fb byte
	fb, eof = sb.NextByte()
	if eof {
		return
	}

	switch fb {
	case 0xFD:
		val, e := sb.NextUint16()
		if e {
			eof = e
			return
		}
		data = uint64(val)
	case 0xFE:
		val, e := sb.NextUint32()
		if e {
			eof = e
			return
		}
		data = uint64(val)
	case 0xFF:
		val, e := sb.NextUint64()
		if e {
			eof = e
			return
		}
		data = uint64(val)
	default:
		data = uint64(fb)
	}
	return
}

// NewReader returns a new SourceBuffer reading from b.
func NewSourceBuffer(b []byte) *SourceBuffer {
	return &SourceBuffer{b, 0}
}

const (
	MAX_UINT64 = math.MaxUint64
)

func SafeAdd(x, y uint64) (uint64, bool) {
	return x + y, y > MAX_UINT64-x
}
