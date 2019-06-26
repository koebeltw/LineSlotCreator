package loadData

import (
	"github.com/koebeltw/Common/packet"
	"github.com/koebeltw/Common/util"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/cstockton/go-conv"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/pquerna/ffjson/ffjson"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ShowType []uint16
type LineType [][]uint16
type WheelType [][]uint16
type OddType []float64
type PostionType [][]uint16

type SlotData struct {
	Line  LineType
	Wheel WheelType
	Prize []OddType
}

type SlotShow struct {
	ColCount uint8
	RowCount uint8
	Show     ShowType
}

type SlotGameData struct {
	Kind     string
	SlotShow SlotShow
	Normal   SlotData
	FreeGame SlotData

	NormalScatterPosition PostionType
	FreeGameScatterPosition PostionType
}

var DataArrayCount = 0
var DataArrayRW = sync.RWMutex{}
var DataArray [256]*SlotGameData

func init() {
	for i := 0; i < len(DataArray); i++ {
		DataArray[i] = &SlotGameData{}
	}
	walkDir("./Data")

	AnalysisScatterPosition()
}

func walkDir(dir string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {

	}

	for _, file := range files {
		filename := file.Name()
		filepath.Walk(dir+"/"+filename, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return filepath.SkipDir
			}

			if !fi.IsDir() {
				return filepath.SkipDir
			}

			walkData(dir, filename)
			return nil
		})
	}
}

func walkData(dir, filename string) {
	gameKind, _ := conv.Uint8(filename)
	dir = dir + "/" + filename
	files, err := ioutil.ReadDir(dir)
	if err != nil {

	}

	DataArrayCount++
	for _, file := range files {
		filename := file.Name()
		filepath.Walk(dir+"/"+filename, func(path string, fi os.FileInfo, err error) error {
			if err != nil {
				return filepath.SkipDir
			}

			if (filepath.Ext(strings.ToLower(path)) != ".txt") &&  (filepath.Ext(strings.ToLower(path)) != ".xlsx") {
				return filepath.SkipDir
			}

			if (strings.Index(strings.ToLower(path), strings.ToLower("Show_"+util.ConvString(gameKind)+"_S.Txt")) > 0) {
				DataArray[gameKind].SlotShow.ColCount, DataArray[gameKind].SlotShow.RowCount, DataArray[gameKind].SlotShow.Show = LoadShow(path)
			} else if (strings.Index(strings.ToLower(path), strings.ToLower("NormalLine_"+util.ConvString(gameKind)+"_S.Txt")) > 0) {
				DataArray[gameKind].Normal.Line = LoadLine(path)
			} else if (strings.Index(strings.ToLower(path), strings.ToLower("NormalWheel_"+util.ConvString(gameKind)+"_S.Txt")) > 0) {
				DataArray[gameKind].Normal.Wheel = LoadWheel(path)
			} else if (strings.Index(strings.ToLower(path), strings.ToLower("NormalPrize_"+util.ConvString(gameKind)+"_S.Txt")) > 0) {
				DataArray[gameKind].Normal.Prize = LoadPrize(path)
			} else if (strings.Index(strings.ToLower(path), strings.ToLower("FreeGameLine_"+util.ConvString(gameKind)+"_S.Txt")) > 0) {
				DataArray[gameKind].FreeGame.Line = LoadLine(path)
			} else if (strings.Index(strings.ToLower(path), strings.ToLower("FreeGameWheel_"+util.ConvString(gameKind)+"_S.Txt")) > 0) {
				DataArray[gameKind].FreeGame.Wheel = LoadWheel(path)
			} else if (strings.Index(strings.ToLower(path), strings.ToLower("FreeGamePrize_"+util.ConvString(gameKind)+"_S.Txt")) > 0) {
				DataArray[gameKind].FreeGame.Prize = LoadPrize(path)
			} else if (strings.Index(strings.ToLower(path), strings.ToLower("Kind.xlsx")) > 0) {
				DataArray[gameKind].Kind = LoadKind(path)
			}

			//fmt.Println(path)
			return nil
		})
	}

}

func LoadKind(fileName string) (kind string) {
	Data, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	rows := Data.GetRows("工作表1")
	for _, row := range rows {
		for _, colCell := range row {
			kind = colCell
		}
	}

	return kind
}

func LoadShow(fileName string) (ColCount uint8, RowCount uint8, show ShowType) {
	file, _ := ioutil.ReadFile(fileName)
	m := make([]map[string]interface{}, 0)
	ffjson.Unmarshal(file, &m)

	s := make([][]uint16, len(m))
	var sh ShowType
	for key, value := range m {
		for i := 1; i < len(value); i++ {
			s[key] = append(s[key], util.ConvUint16(value["Wheel"+util.ConvString(i)]))
			sh = append(sh, util.ConvUint16(value["Wheel"+util.ConvString(i)]))
		}
	}

	show = make([]uint16, len(sh))
	for key, value := range sh {
		show[value-1] = uint16(key + 1)
	}

	return uint8(len(s)), uint8(len(s[0])), show
}

func LoadLine(fileName string) (line LineType) {

	file, _ := ioutil.ReadFile(fileName)
	m := make([]map[string]interface{}, 0)
	ffjson.Unmarshal(file, &m)

	line = make([][]uint16, len(m))
	for key, value := range m {
		for i := 1; i < len(value); i++ {
			line[key] = append(line[key], util.ConvUint16(value["Wheel"+util.ConvString(i)]))
		}
	}

	return
}

func LoadPrize(fileName string) (prizes []OddType) {
	file, _ := ioutil.ReadFile(fileName)
	m := make([]map[string]interface{}, 0)
	ffjson.Unmarshal(file, &m)

	prizes = make([]OddType, 256)
	for _, value := range m {
		picNo := util.ConvUint8(value["Symbol"])
		prizes[picNo] = append(prizes[picNo], 0)
		for i := 1; i <= 5; i++ {
			prizes[picNo] = append(prizes[picNo], util.ConvFloat64(value["Odds"+util.ConvString(i)]))
		}
	}

	return
}

func LoadWheel(fileName string) (wheel WheelType) {
	file, _ := ioutil.ReadFile(fileName)
	m := make([]map[string]interface{}, 0)
	ffjson.Unmarshal(file, &m)

	for _, value := range m {
		for i := 0; i < len(value)-1; i++ {
			if util.ConvUint16(value["Wheel"+util.ConvString(i+1)]) <= 0 {
				continue
			}

			if i >= len(wheel) {
				wheel = append(wheel, []uint16{})
			}

			wheel[i] = append(wheel[i], util.ConvUint16(value["Wheel"+util.ConvString(i+1)]))
		}
	}

	return
}

func AnalysisScatterPosition() {
	for i := 0; i < len(DataArray); i++ {
		DataArray[i].NormalScatterPosition = getScatterPosition(DataArray[i].SlotShow.Show, DataArray[i].Normal.Wheel)
		DataArray[i].FreeGameScatterPosition = getScatterPosition(DataArray[i].SlotShow.Show, DataArray[i].FreeGame.Wheel)
	}
}

func getScatterPosition(show ShowType, wheel WheelType) (Positions PostionType) {
	if show == nil {
		return
	}

	if wheel == nil {
		return
	}

	var basePosition PostionType
	for i := 0; i < len(wheel); i++ {
		basePosition = append(basePosition, []uint16{})
		for j := 0; j < len(wheel[i]); j++ {
			if wheel[i][j] == 200 {
				basePosition[i] = append(basePosition[i], uint16(j))
			}
		}
	}

	for i := 0; i < len(basePosition); i++ {
		var index int
		Positions = append(Positions, []uint16{})
		set := treeset.NewWithIntComparator()
		for j := 0; j < len(basePosition[i]); j++ {
			for k := 0; k < len(show); k++ {
				index = int(basePosition[i][j]) - k
				if index > len(wheel[i]) {
					index = index - len(wheel[i])
				}
				set.Add(index)

				index = int(basePosition[i][j]) - k
				if index < 0 {
					index = index + len(wheel[i])
				}
				set.Add(index)
			}
		}

		for _, value := range set.Values() {
			Positions[i] = append(Positions[i], uint16(value.(int)))
		}
	}

	return
}

func GetDataArrayBytes(callback func([]byte)) {
	pa := packet.NewPacket()
	defer packet.PutPacket(pa)

	for i := 0; i < len(DataArray); i++ {
		if DataArray[i].Kind != "" {
			pa.Reset()
			pa.WriteUint8(uint8(i))

			pa.WriteUint8(util.ConvUint8(len(DataArray[i].Normal.Wheel)))
			for j := 0; j < len(DataArray[i].Normal.Wheel); j++ {
				pa.WriteUint8(util.ConvUint8(len(DataArray[i].Normal.Wheel[j])))

			}

			for j := 0; j < len(DataArray[i].Normal.Wheel); j++ {
				for k := 0; k < len(DataArray[i].Normal.Wheel[j]); k++ {
					pa.WriteUint8(util.ConvUint8(DataArray[i].Normal.Wheel[j][k]))
				}

			}

			pa.WriteUint8(util.ConvUint8(len(DataArray[i].FreeGame.Wheel)))
			for j := 0; j < len(DataArray[i].Normal.Wheel); j++ {
				pa.WriteUint8(util.ConvUint8(len(DataArray[i].FreeGame.Wheel[j])))
			}

			for j := 0; j < len(DataArray[i].Normal.Wheel); j++ {
				for k := 0; k < len(DataArray[i].FreeGame.Wheel[j]); k++ {
					pa.WriteUint8(util.ConvUint8(DataArray[i].FreeGame.Wheel[j][k]))
				}

			}

			if callback != nil {
				callback(pa.CopyBytes())
			}
		}
	}
}
