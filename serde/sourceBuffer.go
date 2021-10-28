package serde

import (
	"encoding/binary"
	"math"
)

const (
	uint16Size  = 2
	uint32Size  = 4
	uint64Size  = 8
	uint256Size = 32
)

//Uint256 is an alias for [32]byte
type Uint256 [uint256Size]byte

//SourceBuffer encapsulates a bytes array
type SourceBuffer struct {
	s   []byte
	off uint64 // current reading index
}

// NewSourceBuffer returns a new SourceBuffer reading from b.
func NewSourceBuffer(b []byte) *SourceBuffer {
	return &SourceBuffer{b, 0}
}

//Len returns the number of unused bytes
func (sb *SourceBuffer) Len() uint64 {
	length := uint64(len(sb.s))
	if sb.off >= length {
		return 0
	}
	return length - sb.off
}

//Bytes returns the entire array of bytes
func (sb *SourceBuffer) Bytes() []byte {
	return sb.s
}

//OffBytes returns unused bytes
func (sb *SourceBuffer) OffBytes() []byte {
	return sb.s[sb.off:]
}

//Pos returns offset position
func (sb *SourceBuffer) Pos() uint64 {
	return sb.off
}

//Size returns the length of the byte array
func (sb *SourceBuffer) Size() uint64 { return uint64(len(sb.s)) }

//NextBytes returns the next n bytes after offset and increase offset position
func (sb *SourceBuffer) NextBytes(n uint32) ([]byte, bool) {
	eof := false
	m := uint64(len(sb.s))
	end, overflow := safeAdd(sb.off, uint64(n))
	if overflow || end > m {
		end = m
		eof = true
	}
	data := sb.s[sb.off:end]
	sb.off = end

	return data, eof
}

//Skip increase offset position with n bytes
func (sb *SourceBuffer) Skip(n uint64) bool {
	eof := false
	m := uint64(len(sb.s))
	end, overflow := safeAdd(sb.off, n)
	if overflow || end > m {
		end = m
		eof = true
	}
	sb.off = end

	return eof
}

//NextByte returns the next byte after offset and increase offset position
func (sb *SourceBuffer) NextByte() (byte, bool) {
	if sb.off >= uint64(len(sb.s)) {
		return 0, true
	}

	data := sb.s[sb.off]
	sb.off++
	return data, false
}

//NextUint8 returns the next byte after offset as uint8 and increase offset position
func (sb *SourceBuffer) NextUint8() (uint8, bool) {
	val, eof := sb.NextByte()
	return uint8(val), eof
}

//NextBool returns the next byte after offset as bool and increase offset position
func (sb *SourceBuffer) NextBool() (bool, bool) {
	val, eof := sb.NextByte()
	data := false
	if val == 0 {
		data = false
	} else if val == 1 {
		data = true
	} else {
		eof = true
	}
	return data, eof
}

//BackUp decrease offset position with n bytes
func (sb *SourceBuffer) BackUp(n uint64) {
	sb.off -= n
}

//NextUint16 returns the next 2 bytes after offset as uint16 and increase offset position
func (sb *SourceBuffer) NextUint16() (uint16, bool) {
	buf, eof := sb.NextBytes(uint16Size)
	if eof {
		return 0, eof
	}

	return binary.BigEndian.Uint16(buf), eof
}

//NextUint32 returns the next 4 bytes after offset as uint32 and increase offset position
func (sb *SourceBuffer) NextUint32() (uint32, bool) {
	buf, eof := sb.NextBytes(uint32Size)
	if eof {
		return 0, eof
	}

	return binary.BigEndian.Uint32(buf), eof
}

//NextUint64 returns the next 8 bytes after offset as uint64 and increase offset position
func (sb *SourceBuffer) NextUint64() (uint64, bool) {
	buf, eof := sb.NextBytes(uint64Size)
	if eof {
		return 0, eof
	}

	return binary.BigEndian.Uint64(buf), eof
}

//NextInt32 returns the next 4 bytes after offset as int32 and increase offset position
func (sb *SourceBuffer) NextInt32() (int32, bool) {
	val, eof := sb.NextUint32()
	return int32(val), eof
}

//NextInt64 returns the next 8 bytes after offset as int64 and increase offset position
func (sb *SourceBuffer) NextInt64() (int64, bool) {
	val, eof := sb.NextUint64()
	return int64(val), eof
}

//NextInt16 returns the next 2 bytes after offset as int16 and increase offset position
func (sb *SourceBuffer) NextInt16() (int16, bool) {
	val, eof := sb.NextUint16()
	return int16(val), eof
}

//NextVarBytes uses the next 4 bytes to determine the number of bytes after the offset to be returned and returns them
//and increase offset position
func (sb *SourceBuffer) NextVarBytes() ([]byte, bool) {
	count, eof := sb.NextUint32()
	if eof {
		return []byte{}, eof
	}
	data, eof := sb.NextBytes(count)
	return data, eof
}

//NextHash returns the next 32 bytes after offset as Uint256 and increase offset position
func (sb *SourceBuffer) NextHash() (Uint256, bool) {
	buf, eof := sb.NextBytes(uint256Size)
	if eof {
		return Uint256{}, eof
	}
	var data Uint256
	copy(data[:], buf)

	return data, eof
}

//NextString uses the next 4 bytes to determine the number of bytes after the offset to be returned and returns them
//as string and increase offset position
func (sb *SourceBuffer) NextString() (string, bool) {
	val, eof := sb.NextVarBytes()
	data := string(val)
	return data, eof
}

func safeAdd(x, y uint64) (uint64, bool) {
	return x + y, y > math.MaxUint64-x
}
