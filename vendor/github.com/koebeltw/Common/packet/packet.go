package packet

import (
	"bytes"
	"encoding/binary"
	"sync"
	"math"
)

var packetPool = sync.Pool{New: func() interface{} {return NewPacket()}}

type Packet struct {
	buf *bytes.Buffer
}

func NewPacket() *Packet {
	p := &bytes.Buffer{}
	p.Grow(512)
	return &Packet{
			buf: p,
		}
}

func GetPacket() *Packet {
	return packetPool.Get().(*Packet)
}

func PutPacket(p *Packet) {
	p.Reset()
	packetPool.Put(p)
}

func NewPacketByBytes(buf []byte) *Packet {
	return &Packet{
		buf: bytes.NewBuffer(buf),
	}
}

func (p *Packet) Reset() {
	p.buf.Reset()
}

func (p *Packet) Len() int {
	return p.buf.Len()
}

func (p *Packet) Buffer() *bytes.Buffer{
	return p.buf
}

func (p *Packet) Bytes() []byte {
	return p.buf.Bytes()
}

func (p *Packet) CopyBytes() []byte {
	return append([]byte{}, p.buf.Bytes()...)
}

func (p *Packet) Copy(p2 *Packet) {
	p.Reset()
	p.Write(p2.CopyBytes())
}

func (p *Packet) Next(n int) []byte {
	return p.buf.Next(n)
}

func (p *Packet) Read(buf []byte) (n int, err error) {
	return p.buf.Read(buf)
}

func (p *Packet) ReadBool() (bool, bool) {
	buf := p.buf.Next(1)
	if len(buf) != 1 {
		return false, false
	}
	return buf[0] > 0, true
}

func (p *Packet) ReadUint8() (uint8, bool) {
	buf := p.buf.Next(1)
	if len(buf) != 1 {
		return 0, false
	}
	return buf[0], true
}

func (p *Packet) ReadUint16() (uint16, bool) {
	buf := p.buf.Next(2)
	if len(buf) != 2 {
		return 0, false
	}
	return binary.LittleEndian.Uint16(buf), true
}

func (p *Packet) ReadUint32() (uint32, bool) {
	buf := p.buf.Next(4)
	if len(buf) != 4 {
		return 0, false
	}
	return binary.LittleEndian.Uint32(buf), true
}

func (p *Packet) ReadUint64() (uint64, bool) {
	buf := p.buf.Next(8)
	if len(buf) != 8 {
		return 0, false
	}
	return binary.LittleEndian.Uint64(buf), true
}

func (p *Packet) ReadFloat32() (float32, bool) {
	buf := p.buf.Next(8)
	if len(buf) != 8 {
		return 0, false
	}
	return math.Float32frombits(binary.LittleEndian.Uint32(buf)), true
}

func (p *Packet) ReadFloat64() (float64, bool) {
	buf := p.buf.Next(8)
	if len(buf) != 8 {
		return 0, false
	}
	return math.Float64frombits(binary.LittleEndian.Uint64(buf)), true
}

func (p *Packet) ReadInterface(v interface{}) error {
	return binary.Read(p.buf, binary.LittleEndian, v)
}

func (p *Packet) Write(b []byte) (n int, err error) {
	return p.buf.Write(b)
}

func (p *Packet) WriteString(s string) (n int, err error) {
	return p.buf.WriteString(s)
}

func (p *Packet) WriteBool(v bool) {
	binary.Write(p.buf, binary.LittleEndian, v)
}

func (p *Packet) WriteUint8(v uint8) {
	binary.Write(p.buf, binary.LittleEndian, v)
}

func (p *Packet) WriteUint16(v uint16) {
	binary.Write(p.buf, binary.LittleEndian, v)
}

func (p *Packet) WriteUint32(v uint32) {
	binary.Write(p.buf, binary.LittleEndian, v)
}

func (p *Packet) WriteUint64(v uint64) {
	binary.Write(p.buf, binary.LittleEndian, v)
}

func (p *Packet) WriteFloat32(v float32) {
	binary.Write(p.buf, binary.LittleEndian, math.Float32bits(v))
}

func (p *Packet) WriteFloat64(v float64) {
	binary.Write(p.buf, binary.LittleEndian, math.Float64bits(v))
}

func (p *Packet) WriteInterface(v interface{}) error {
	return binary.Write(p.buf, binary.LittleEndian, v)
}

func Uint64ToBytes(v uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, v)
	return b
}

func Uint32ToBytes(v uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	return b
}

func Uint16ToBytes(v uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, v)
	return b
}

func BytesToUint64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

func BytesToUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func BytesToUint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}
