package tcp

import (
	"github.com/koebeltw/Common/packet"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var moneyMsgHeadSize = binary.Size(moneyMsgHead{})
func NewMoneyMsgHead() Coder { return moneyMsgHead{}}

// moneyMsgHead blabla
type moneyMsgHead struct {
	HeadCode uint16
	Size     uint32
	MsgNo    uint8
	SubNo    uint8
}

// Decode blabla
func (m moneyMsgHead) Decode(c Session) (r EventMsg, err error) {
	reader := c.GetConn().(io.Reader)
	buffer := c.GetBuffer()

	size := moneyMsgHeadSize
	_, err = io.ReadFull(reader, buffer[0:size])
	if err != nil {
		fmt.Println(err)
		return EventMsg{}, err
	}

	pa := packet.GetPacket()
	defer packet.PutPacket(pa)
	pa.Write(buffer[0:size])
	//packet := packet.NewPacketByBytes(buffer[0:size])
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

	fmt.Printf("Size:%d\n", size)

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

	// time.Sleep(time.Second * 10)
	return EventMsg{MsgNo: m.MsgNo, SubNo: m.SubNo, Buffer: pa.CopyBytes()}, nil
}

// Encode blabla
func (m moneyMsgHead) Encode(c Session, msgNo byte, subNo byte, buffer []byte) (r []byte, err error) {
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

	if pa.WriteInterface(pa); err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println("packet.Bytes():", pa.Bytes())
	r = pa.CopyBytes()
	return
}
