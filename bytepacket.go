package main

import (
	"strings"
)

type BytePacketBuffer struct {
	Buf [512]byte
	Pos uintptr
}

func newBytePacketBuffer() BytePacketBuffer {
	return BytePacketBuffer{
		Buf: [512]byte{0},
		Pos: 0,
	}
}

func (b *BytePacketBuffer) pos() uintptr {
	return b.Pos
}

func (b *BytePacketBuffer) step(steps uintptr) {
	b.Pos += steps
}

func (b *BytePacketBuffer) seek(pos uintptr) {
	b.Pos = pos
}

func (b *BytePacketBuffer) read() uint8 {
	if b.Pos >= 512 {
		panic("end of buffer")
	}
	res := b.Buf[b.Pos]
	b.Pos += 1
	return res
}

func (b *BytePacketBuffer) get(pos uintptr) uint8 {
	if b.Pos >= 512 {
		panic("end of buffer")
	}
	return b.Buf[pos]
}

func (b *BytePacketBuffer) getRange(start uintptr, len uintptr) []byte {
	if start+len >= 512 {
		panic("end of buffer")
	}
	return b.Buf[start:(start + len)]
}

func (b *BytePacketBuffer) readUint16() uint16 {
	res := uint16(b.read())<<8 | uint16(b.read())
	return res
}

func (b *BytePacketBuffer) readUint32() uint32 {
	res := uint32(b.read())<<24 | uint32(b.read())<<16 | uint32(b.read())<<8 | uint32(b.read())
	return res
}

func (b *BytePacketBuffer) readQname(outstr *string) {
	pos := b.pos()

	jumped := false
	max_jumps := 5
	jumps_performed := 0

	delimiter := ""

	for {
		if jumps_performed > max_jumps {
			panic("max jumps exceeded")
		}

		len := b.get(pos)

		if (len & 0xC0) == 0xC0 {
			if !jumped {
				b.seek(pos + 2)
			}
			b2 := uint16(b.get(pos + 1))
			offset := ((uint16(len) ^ 0xC0) << 8) | b2
			pos = uintptr(offset)

			jumped = true
			jumps_performed += 1

			continue
		} else {
			pos += 1

			if len == 0 {
				break
			}

			*outstr += delimiter
			strBuf := b.getRange(pos, uintptr(len))

			*outstr += strings.ToLower(string(strBuf))

			delimiter = "."
			pos += uintptr(len)
		}
	}

	if !jumped {
		b.seek(pos)
	}
}

func (buffer *BytePacketBuffer) write(val uint8) {
	if buffer.Pos >= 512 {
		panic("end of buffer")
	}
	buffer.Buf[buffer.Pos] = val
	buffer.Pos += 1
}

func (buffer *BytePacketBuffer) writeUint8(val uint8) {
	buffer.write(val)
}

func (buffer *BytePacketBuffer) writeUint16(val uint16) {
	buffer.write(uint8(val >> 8))
	buffer.write(uint8(val & 0xFF))
}

func (buffer *BytePacketBuffer) writeUint32(val uint32) {
	buffer.write(uint8(val>>24) & 0xFF)
	buffer.write(uint8((val >> 16) & 0xFF))
	buffer.write(uint8((val >> 8) & 0xFF))
	buffer.write(uint8((val >> 0) & 0xFF))
}

func (buffer *BytePacketBuffer) writeQname(qname *string) {
	qnameArr := strings.Split(*qname, ".")
	for i := 0; i < len(qnameArr); i++ {
		length := len(qnameArr[i])
		if length > 0x3f {
			panic("label exceeds 63 chars length")
		}
		buffer.writeUint8(uint8(length))
		for j := 0; j < length; j++ {
			buffer.writeUint8(uint8(qnameArr[i][j]))
		}
	}
	buffer.writeUint8(0)
}

func (buffer *BytePacketBuffer) set(pos uintptr, val uint8) {
	buffer.Buf[pos] = val
}

func (buffer *BytePacketBuffer) setUint16(pos uintptr, val uint16) {
	buffer.set(pos, uint8(val>>8))
	buffer.set(pos+1, uint8(val&0xFF))
}
