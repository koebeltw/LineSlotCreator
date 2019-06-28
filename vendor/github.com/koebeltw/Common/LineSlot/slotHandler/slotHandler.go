package slotHandler

import (
	"github.com/koebeltw/LineSlotCreator/ls/loadData"
	"crypto/rand"
	"github.com/koebeltw/Common/type"
	"github.com/koebeltw/Common/util"
	"math"
	"math/big"
)

type slotCh struct {
	noFreeGameCh   chan *Type.Temp
	intoFreeGameCh chan *Type.Temp
	//normalCh       chan *Type.Temp
	//freeGameCh     chan *Type.Temp
}

var SlotCh [256]slotCh

func init() {
	//rand.Seed(time.Now().UnixNano())
	for i := 0; i < len(loadData.DataArray); i++ {
		if loadData.DataArray[i].Kind != "" {
			SlotCh[i].noFreeGameCh = createNoFreeGameCh(loadData.DataArray[i])
			SlotCh[i].intoFreeGameCh = createIntoFreeGameCh(loadData.DataArray[i])
			//SlotCh[i].normalCh = createNormalCh(loadData.DataArray[i])
			//SlotCh[i].freeGameCh = createFreeGameCh(loadData.DataArray[i])
		}
	}
}

func calScatterPrize(wheel []uint16, prize []loadData.OddType) (p Type.Prize) {
	specialNum := 0
	for key, value := range wheel {
		if value == 200 {
			p.LinePicPosition = append(p.LinePicPosition, uint16(key+1))
			p.LinePicNo = append(p.LinePicNo, uint16(value))
			specialNum++
		}
	}

	p.WinPic = 200
	p.PicCount = uint16(specialNum)
	p.Odds = float64(prize[p.WinPic][p.PicCount])
	return p
}

func calOneLinePrize(wheel []uint16, line []uint16, show loadData.ShowType, prize []loadData.OddType) (p Type.Prize) {
	normal := make(map[uint16]uint16, 5)
	wild := make(map[uint16]uint16, 5)
	var picNo uint16
	var normalCount uint16
	var wildCount uint16
	var picCount uint16

	for _, value := range line {
		picNo = wheel[show[value-1]-1]
		if picNo != 255 {
			normalCount++
			normal[picNo] = normal[picNo] + 1
		} else {
			wildCount++
			wild[picNo] = wild[picNo] + 1
		}

		if picNo == 200 {
			normalCount--
			delete(normal, picNo)
			break
		}

		if len(normal) == 2 {
			normalCount--
			delete(normal, picNo)
			break
		}
	}

	picCount = normalCount + wildCount
	if picCount <= 0 {
		normal = nil
		wild = nil
		return
	}

	var winPic uint16
	for key, _ := range normal {
		winPic = key
	}

	var wildNum uint16
	for _, value := range line {
		picNo = wheel[show[value-1]-1]
		if picNo == 255 {
			wildNum++
		} else {
			break
		}
	}

	if winPic > 0 {
		if float64(prize[winPic][picCount]) > float64(prize[255][wildNum]) {
			p.WinPic = winPic
			p.PicCount = picCount
			p.Odds = float64(prize[winPic][picCount])
		} else if float64(prize[winPic][picCount]) <= float64(prize[255][wildNum]) {
			p.WinPic = 255
			p.PicCount = wildNum
			p.Odds = float64(prize[255][wildNum])
		}
	} else {
		p.WinPic = 255
		p.PicCount = wildNum
		p.Odds = float64(prize[255][wildNum])
	}

	if wildCount > 0 {
		p.Odds *= 2
		p.WildCount = wildCount
	}

	if p.Odds > 0 {
		for _, value := range line {
			p.LinePicNo = append(p.LinePicNo, uint16(wheel[show[value-1]-1]))
			p.LinePicPosition = append(p.LinePicPosition, uint16(show[value-1]))
		}
	}

	return p
}

func getRandomResult(slotShow loadData.SlotShow, slotData loadData.SlotData) (temp *Type.Temp) {
	temp = &Type.Temp{}
	p, b := getRandomWheel(slotShow, slotData)
	temp = getResult(p, b, slotShow, slotData)
	return temp
}

func getRandomResultNoFreeGame(slotShow loadData.SlotShow, slotData loadData.SlotData) (temp *Type.Temp) {
	temp = &Type.Temp{}
	for {
		p, b := getRandomWheel(slotShow, slotData)
		temp = getResult(p, b, slotShow, slotData)

		if temp.ScatterCount < 3 {
			break
		}
	}

	return temp
}

func getRandomWheel(slotShow loadData.SlotShow, slotData loadData.SlotData) (r []uint16, b []uint16) {
	//rand.Seed(time.Now().UnixNano())
	r = make([]uint16, len(slotData.Wheel))

	for key, _ := range slotData.Wheel {
		randomNum, _ := rand.Int(rand.Reader, big.NewInt(int64(len(slotData.Wheel[key]))))
		r[key] = uint16(randomNum.Uint64())
	}

	for i := 0; i < int(slotShow.ColCount); i++ {
		for j := 0; j < int(slotShow.RowCount); j++ {
			index := (uint16(i) + r[j]) % uint16(len(slotData.Wheel[j]))
			b = append(b, slotData.Wheel[j][index])
		}

	}
	return
}

func getResult(p []uint16, b []uint16, slotShow loadData.SlotShow, slotData loadData.SlotData) (temp *Type.Temp) {
	temp = &Type.Temp{}
	ScatterPrize := calScatterPrize(b, slotData.Prize)
	if ScatterPrize.Odds > 0 {
		temp.TotalOdds = temp.TotalOdds + ScatterPrize.Odds
		temp.Prize = append(temp.Prize, ScatterPrize)
		temp.ScatterCount = ScatterPrize.PicCount
	}

	for key, value := range slotData.Line {
		OneLinePrize := calOneLinePrize(b, value, slotShow.Show, slotData.Prize)
		OneLinePrize.LineNo = uint16(key + 1)

		if OneLinePrize.Odds > 0 {
			temp.TotalOdds = temp.TotalOdds + OneLinePrize.Odds
			temp.Prize = append(temp.Prize, OneLinePrize)
		}
	}

	for _, Value := range b {
		temp.WheelUint16 = append(temp.WheelUint16, uint16(Value))
	}

	temp.TotalOdds = float64(uint64(math.Round(temp.TotalOdds * 1000))) / 1000
	temp.Positions = p
	return temp
}

func getRandomWheelIntoFreeGame(Positions loadData.PostionType, slotShow loadData.SlotShow, slotData loadData.SlotData, scatterCount uint8) (r []uint16, b []uint16) {
	//rand.Seed(time.Now().UnixNano())
	r = make([]uint16, slotShow.RowCount)

	indexArray := make([]uint16, slotShow.ColCount)
	for i := 0; i < len(indexArray); i++ {
		indexArray[i] = uint16(i)
	}

	for i := 0; i < len(indexArray); i++ {
		randomNum, _ := rand.Int(rand.Reader, big.NewInt(int64(len(indexArray))))
		tempKey := uint16(randomNum.Uint64())
		tempValue := indexArray[i]
		indexArray[i] = indexArray[tempKey]
		indexArray[tempKey] = tempValue
	}

	for i := 0; i < len(indexArray); i++ {
		index := indexArray[i]
		if len(Positions[index]) > 0 && scatterCount > 0 {
			randomNum, _ := rand.Int(rand.Reader, big.NewInt(int64(len(Positions[index]))))
			r[index] = Positions[index][uint16(randomNum.Uint64())]

			scatterCount--
		} else {
			randomNum, _ := rand.Int(rand.Reader, big.NewInt(int64(len(slotData.Wheel[index]))))
			r[index] = uint16(randomNum.Uint64())
		}
	}

	for i := 0; i < int(slotShow.ColCount); i++ {
		for j := 0; j < int(slotShow.RowCount); j++ {
			index := (uint16(i) + r[j]) % uint16(len(slotData.Wheel[j]))
			b = append(b, slotData.Wheel[j][index])
		}
	}

	return
}

func createNormalCh(d *loadData.SlotGameData) (ch chan *Type.Temp) {
	ch = make(chan *Type.Temp, 1000)
	go func(ch chan *Type.Temp) {
		for {
			ch <- getRandomResult(d.SlotShow, d.Normal)
		}
	}(ch)

	return ch
}

func createFreeGameCh(d *loadData.SlotGameData) (ch chan *Type.Temp) {
	ch = make(chan *Type.Temp, 1000)
	go func(ch chan *Type.Temp) {
		for {
			ch <- getRandomResult(d.SlotShow, d.FreeGame)
		}
	}(ch)

	return ch
}

func createIntoFreeGameCh(d *loadData.SlotGameData) (ch chan *Type.Temp) {
	ch = make(chan *Type.Temp, 100)
	go func(ch chan *Type.Temp) {
		for {
			p, b := getRandomWheelIntoFreeGame(d.NormalScatterPosition, d.SlotShow, d.Normal, d.SlotShow.ColCount)
			ch <- getResult(p, b, d.SlotShow, d.Normal)
		}
	}(ch)

	return ch
}

func createNoFreeGameCh(d *loadData.SlotGameData) (ch chan *Type.Temp) {
	ch = make(chan *Type.Temp, 100)
	go func(ch chan *Type.Temp) {
		for {
			ch <- getRandomResultNoFreeGame(d.SlotShow, d.Normal)
		}
	}(ch)

	return ch
}

func PlayOnce(gameKind uint8, freeGameCount uint16, Positions []uint16, intoFreeGameRate int32) (temp *Type.Temp) {
	//rand.Seed(time.Now().UnixNano())

	if freeGameCount <= 0 {
		randomNum, _ := rand.Int(rand.Reader, big.NewInt(10000000))
		if randomNum.Uint64() < uint64(math.Abs(float64(intoFreeGameRate))) {
			if intoFreeGameRate < 0 {
				temp = <-SlotCh[gameKind].noFreeGameCh
			} else if intoFreeGameRate > 0 {
				temp = <-SlotCh[gameKind].intoFreeGameCh
			}
		} else {
			//temp = <-SlotCh[gameKind].normalCh
			temp = getRandomResult(loadData.DataArray[gameKind].SlotShow, loadData.DataArray[gameKind].Normal)
		}
	} else {
		//temp = <-SlotCh[gameKind].freeGameCh
		temp = getRandomResult(loadData.DataArray[gameKind].SlotShow, loadData.DataArray[gameKind].FreeGame)
		temp.FreeGameCount = freeGameCount - 1
	}

	for _, value := range temp.Prize {
		if value.LineNo == 0 {
			if value.PicCount >= 3 {
				temp.FreeGameCount += 15
				break
			}
		}
	}

	return temp
}

func RandomRunWheel(gameKind uint8, count uint32, intoFreeGameRate int32) (recordMapArrays []map[string]interface{}) {
	//rand.Seed(time.Now().UnixNano())

	var v uint32
	if count == 1 {
		v = util.ConvUint32(1)
	} else {
		v = util.ConvUint32(((count - 1) / 10000) + 1)
	}

	recordMapArrays = make([]map[string]interface{}, v)
	for i := 0; i < len(recordMapArrays); i++ {
		recordMapArrays[i] = make(map[string]interface{})
		recordMapArrays[i]["累積總投入"] = float64(0)
		recordMapArrays[i]["累積總吐出"] = float64(0)
		recordMapArrays[i]["累積回吐率"] = float64(0)
		recordMapArrays[i]["累積總輸贏"] = float64(0)

		recordMapArrays[i]["階段總投入"] = float64(0)
		recordMapArrays[i]["階段總吐出"] = float64(0)
		recordMapArrays[i]["階段回吐率"] = float64(0)
		recordMapArrays[i]["階段總輸贏"] = float64(0)

		recordMapArrays[i]["normal次數"] = uint32(0)
		recordMapArrays[i]["normal賠率"] = float64(0)
		recordMapArrays[i]["freeGame次數"] = uint32(0)
		recordMapArrays[i]["freeGame賠率"] = float64(0)

		recordMapArrays[i]["normal獎項"+util.ConvString(0)+"次數"] = uint32(0)
		recordMapArrays[i]["normal獎項"+util.ConvString(0)+"賠率"] = float64(0)
		for key, _ := range loadData.DataArray[gameKind].Normal.Line {
			recordMapArrays[i]["normal獎項"+util.ConvString(key+1)+"次數"] = uint32(0)
			recordMapArrays[i]["normal獎項"+util.ConvString(key+1)+"賠率"] = float64(0)
		}

		recordMapArrays[i]["freeGame獎項"+util.ConvString(0)+"次數"] = uint32(0)
		recordMapArrays[i]["freeGame獎項"+util.ConvString(0)+"賠率"] = float64(0)
		for key, _ := range loadData.DataArray[gameKind].FreeGame.Line {
			recordMapArrays[i]["freeGame獎項"+util.ConvString(key+1)+"次數"] = uint32(0)
			recordMapArrays[i]["freeGame獎項"+util.ConvString(key+1)+"賠率"] = float64(0)
		}
	}

	var totalMoneyIn float64
	var totalMoneyOut float64
	var freeGameCount uint16
	var temp *Type.Temp
	for i := 0; i < int(count); i++ {
		v := util.ConvUint32(i / 10000)
		if freeGameCount <= 0 {
			//temp = getRandomResult(loadData.DataArray[gameKind].Show, loadData.DataArray[gameKind].NormalWheel, loadData.DataArray[gameKind].NormalLine, loadData.DataArray[gameKind].NormalPrize)
			temp = PlayOnce(gameKind, freeGameCount, []uint16{}, intoFreeGameRate)

			totalMoneyIn = totalMoneyIn + 1
			totalMoneyOut = totalMoneyOut + temp.TotalOdds
			recordMapArrays[v]["累積總投入"] = totalMoneyIn
			recordMapArrays[v]["累積總吐出"] = totalMoneyOut
			recordMapArrays[v]["累積回吐率"] = totalMoneyOut / totalMoneyIn * 100
			recordMapArrays[v]["累積總輸贏"] = totalMoneyOut - totalMoneyIn

			recordMapArrays[v]["階段總投入"] = recordMapArrays[v]["階段總投入"].(float64) + 1
			recordMapArrays[v]["階段總吐出"] = recordMapArrays[v]["階段總吐出"].(float64) + temp.TotalOdds
			recordMapArrays[v]["階段回吐率"] = recordMapArrays[v]["階段總吐出"].(float64) / recordMapArrays[v]["階段總投入"].(float64) * 100
			recordMapArrays[v]["階段總輸贏"] = recordMapArrays[v]["階段總吐出"].(float64) - recordMapArrays[v]["階段總投入"].(float64)

			recordMapArrays[v]["normal次數"] = recordMapArrays[v]["normal次數"].(uint32) + 1
			recordMapArrays[v]["normal賠率"] = recordMapArrays[v]["normal賠率"].(float64) + temp.TotalOdds

			for _, value := range temp.Prize {
				recordMapArrays[v]["normal獎項"+util.ConvString(value.LineNo)+"次數"] = recordMapArrays[v]["normal獎項"+util.ConvString(value.LineNo)+"次數"].(uint32) + 1
				recordMapArrays[v]["normal獎項"+util.ConvString(value.LineNo)+"賠率"] = recordMapArrays[v]["normal獎項"+util.ConvString(value.LineNo)+"賠率"].(float64) + value.Odds
			}
		} else {
			//temp = getRandomResult(loadData.DataArray[gameKind].Show, loadData.DataArray[gameKind].FreeGameWheel, loadData.DataArray[gameKind].FreeGameLine, loadData.DataArray[gameKind].FreeGamePrize)
			temp = PlayOnce(gameKind, freeGameCount, []uint16{}, intoFreeGameRate)

			totalMoneyOut = totalMoneyOut + temp.TotalOdds
			recordMapArrays[v]["累積總投入"] = totalMoneyIn
			recordMapArrays[v]["累積總吐出"] = totalMoneyOut
			recordMapArrays[v]["累積回吐率"] = totalMoneyOut / totalMoneyIn * 100
			recordMapArrays[v]["累積總輸贏"] = totalMoneyOut - totalMoneyIn

			recordMapArrays[v]["階段總吐出"] = recordMapArrays[v]["階段總吐出"].(float64) + temp.TotalOdds
			recordMapArrays[v]["階段回吐率"] = recordMapArrays[v]["階段總吐出"].(float64) / recordMapArrays[v]["階段總投入"].(float64) * 100
			recordMapArrays[v]["階段總輸贏"] = recordMapArrays[v]["階段總吐出"].(float64) - recordMapArrays[v]["階段總投入"].(float64)

			recordMapArrays[v]["freeGame次數"] = recordMapArrays[v]["freeGame次數"].(uint32) + 1
			recordMapArrays[v]["freeGame賠率"] = recordMapArrays[v]["freeGame賠率"].(float64) + temp.TotalOdds

			for _, value := range temp.Prize {
				recordMapArrays[v]["freeGame獎項"+util.ConvString(value.LineNo)+"次數"] = recordMapArrays[v]["freeGame獎項"+util.ConvString(value.LineNo)+"次數"].(uint32) + 1
				recordMapArrays[v]["freeGame獎項"+util.ConvString(value.LineNo)+"賠率"] = recordMapArrays[v]["freeGame獎項"+util.ConvString(value.LineNo)+"賠率"].(float64) + value.Odds
			}
			//freeGameCount--
		}

		//for _, value := range temp.Prize {
		//	if value.LineNo == 0 {
		//		if value.PicCount >= 3 {
		//			freeGameCount += 15
		//			break
		//		}
		//	}
		//}
		freeGameCount = temp.FreeGameCount
	}

	return recordMapArrays
}
