package tcp

import (
	"github.com/koebeltw/Common/packet"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var MsgHeadSize = binary.Size(MsgHead{})

func NewMsgHead() Coder { return MsgHead{}}

// MsgHead blabla
type MsgHead struct {
	HeadCode     uint16
	Size         uint32
	MsgNo        uint8
	SubNo        uint8
	PackCompress bool
	// GateID       uint8
	// ServerID     uint8
	SessionID      int32
	ProtocolSerial uint8
}


// Decode blabla
func (m MsgHead) Decode(c Session) (r EventMsg, err error) {
	reader := c.GetConn().(io.Reader)
	buffer := c.GetBuffer()

	size := MsgHeadSize
	_, err = io.ReadFull(reader, buffer[0:size])
	if err != nil {
		fmt.Println(err)
		return EventMsg{}, err
	}

	pa := packet.GetPacket()
	defer packet.PutPacket(pa)
	pa.Write(buffer[0:size])
	if pa.ReadInterface(&m); err != nil {
		fmt.Println(err)
		return
	}

	if m.HeadCode != 29112 {
		return EventMsg{}, errors.New("HeadCode Eeeor")
	}

	size = int(m.Size) - binary.Size(m)
	if size < 0 {
		return EventMsg{}, errors.New("Size Eeeor")
	}

	//fmt.Printf("Size:%d\n", size)

	for size > 0 {
		if size > len(buffer) {
			if _, err := io.ReadFull(reader, buffer[:]); err != nil {
				return EventMsg{}, err
			}

			pa.Write(buffer[:])
		} else {
			if _, err := io.ReadFull(reader, buffer[0:size]); err != nil {
				return EventMsg{}, err
			}

			pa.Write(buffer[0:size])
		}

		size = size - len(buffer)
	}

	// log.Println("Buffer:", string(packet.Bytes()))
	// log.Println("Buffer:", string(c.recvMsg[0:size]))

	// time.Sleep(time.Second * 10)
	return EventMsg{MsgNo: m.MsgNo, SubNo: m.SubNo, Buffer: pa.CopyBytes()}, nil
}

// Encode blabla
func (m MsgHead) Encode(c Session, msgNo byte, subNo byte, buffer []byte) (r []byte, err error) {
	m.HeadCode = 29112
	m.Size = uint32(len(buffer))
	m.MsgNo = msgNo
	m.SubNo = subNo

	m.Size = m.Size + uint32(binary.Size(m))

	pa := packet.GetPacket()
	defer packet.PutPacket(pa)
	//packet := packet.NewPacket()
	if pa.WriteInterface(m); err != nil {
		fmt.Println(err)
		return
	}

	if pa.WriteInterface(buffer); err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println("packet.CopyBytes():", pa.CopyBytes())
	r = pa.CopyBytes()
	return
}
