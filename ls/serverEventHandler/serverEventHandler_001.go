package serverEventHandler

import (
	"github.com/koebeltw/LineSlotCreator/ls/slotHandler"
	"github.com/koebeltw/Common/packet"
	"github.com/koebeltw/Common/session"
	"github.com/koebeltw/Common/util"
)

type PlayerData struct {
	GameKind uint8
	PlayerID int32
	Who      int32
	Serial   int32
}

//// Event001Dash001 blabla
//func (h ServerEventHandler) Event001Dash001() (session.Eventfunc) {
//	return func(session session.Session, b []byte) {
//		type readData struct {
//			PlayerData
//			intoFreeGameRate int32
//		}
//
//		go func() {
//			pa := packet.NewPacketByBytes(b)
//			defer packet.PutPacket(pa)
//
//			ReadData := readData{}
//			pa.ReadInterface(&ReadData)
//			Positions := []uint16{}
//
//			temp := slotHandler.PlayOnce(ReadData.GameKind, 0, Positions, util.ConvInt32(ReadData.intoFreeGameRate))
//			pa.Reset()
//			pa.WriteInterface(ReadData)
//			pa.WriteUint16(temp.FreeGameCount)
//			pa.WriteFloat64(temp.TotalOdds)
//			pa.WriteUint8(uint8(len(temp.Positions)))
//			for i := 0; i < len(temp.Positions); i++ {
//				pa.WriteUint8(uint8(temp.Positions[i]))
//			}
//
//			session.SendMsg(001, 001, pa.CopyBytes())
//		}()
//	}
//}

// Event001Dash001 blabla
func (h ServerEventHandler) Event001Dash001() (session.Eventfunc) {
	return func(session session.Session, b []byte) {
		type readData struct {
			PlayerData
			IntoFreeGameRate int32
			FreeGameCount uint16
			PositionCount uint16
		}

		go func() {
			pa := packet.NewPacketByBytes(b)
			defer packet.PutPacket(pa)

			ReadData := readData{}
			pa.ReadInterface(&ReadData)
			Positions := []uint16{}
			if ReadData.PositionCount > 0 {
				for i := 0; i < int(ReadData.PositionCount); i++ {
					v, _ := pa.ReadUint8()
					Positions = append(Positions, uint16(v))
				}
			}

			temp := slotHandler.PlayOnce(ReadData.GameKind, ReadData.FreeGameCount, Positions, util.ConvInt32(ReadData.IntoFreeGameRate))
			pa.Reset()
			pa.WriteInterface(ReadData.PlayerData)
			pa.WriteUint16(temp.FreeGameCount)
			pa.WriteFloat64(temp.TotalOdds)
			pa.WriteUint8(uint8(len(temp.Positions)))
			pa.WriteUint8(uint8(len(temp.WheelUint16)))

			for i := 0; i < len(temp.Positions); i++ {
				pa.WriteUint8(uint8(temp.Positions[i]))
			}

			for i := 0; i < len(temp.WheelUint16); i++ {
				pa.WriteUint8(uint8(temp.WheelUint16[i]))
			}

			//fmt.Println("===========================")
			//bytes, _ := ffjson.Marshal(ReadData.PlayerData)
			//fmt.Println(string(bytes))
			//bytes, _ = ffjson.Marshal(temp)
			//fmt.Println(string(bytes))
			session.SendMsg(001, 001, pa.CopyBytes())
		}()
	}
}
