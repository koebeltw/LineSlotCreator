package SQL

import (
	"diamond/util"
	"diamond/ls/loadConfig"
	"diamond/ls/loadXlsx"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"sync"
	"time"
)

var Engine *xorm.Engine
var saveResultCount uint32
var saveTemp_WinCount uint32
var saveTemp_NoWinCount uint32
var totalSaveResultCount uint32
var totalSaveTempCount uint32

var calLinech chan []uint8
var resultWG = sync.WaitGroup{}
var tempWG = sync.WaitGroup{}
var rw = sync.RWMutex{}

var results []*Type.Result
var resultCh chan *Type.Result

var temps_Win []*Type.Temp
var temps_NoWin []*Type.Temp
var tempCh chan *Type.Temp

func init() {

}

func LoadSQL(config loadConfig.ConfigStruct){
	SQLroot := loadConfig.Config.SQLAccount + ":" + loadConfig.Config.SQLPassword + "@tcp(" + loadConfig.Config.SQLIP + ":" + loadConfig.Config.SQLPort + ")/?charset=utf8"
	x, err := xorm.NewEngine("mysql", SQLroot)
	if err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}

	s := util.ConvString(time.Now().Format("20060102150405"))

	if loadConfig.Config.IsReWriteData {
		s = "slot"
	} else {
		s = "slot" + s
	}

	_, err = x.Query(`Drop DATABASE IF EXISTS ` + s)
	if err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}

	_, err = x.Query(`CREATE DATABASE IF NOT EXISTS ` + s)
	if err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}
	x.Close()

	Database := loadConfig.Config.SQLAccount + `:` + loadConfig.Config.SQLPassword + `@/` + s + `?charset=utf8`
	x, err = xorm.NewEngine("mysql", Database)
	if err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}

	//x.ShowSQL(true)
	//設置連接池的空閒數大小
	x.SetMaxIdleConns(10000)
	//設置最大打開連接數
	x.SetMaxOpenConns(10000)

	if err := x.Sync2(new(Type.Result)); err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}

	if err := x.Table("temp_Win").Sync2(new(Type.Temp)); err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}

	if err := x.Table("temp_NoWin").Sync2(new(Type.Temp)); err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}

	if err := x.Sync2(new(Type.OddsCount)); err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}

	if err := x.Sync2(new(Type.LineOddsCount)); err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}

	Engine = x

	calLinech = make(chan []uint8, loadConfig.Config.CalLineNum)
	results = make([]*Type.Result, loadConfig.Config.OnceSaveResultNum)
	resultCh = make(chan *Type.Result, loadConfig.Config.OnceSaveResultNum*100)
	temps_Win = make([]*Type.Temp, loadConfig.Config.OnceSaveTempNum)
	temps_NoWin = make([]*Type.Temp, loadConfig.Config.OnceSaveTempNum)
	tempCh = make(chan *Type.Temp, loadConfig.Config.OnceSaveTempNum*100)
}

func saveResult(config loadConfig.ConfigStruct, save []*Type.Result) {
	if config.IsSaveResult == false {
		return
	}

	if _, err := Engine.Insert(&save); err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}

	if config.IsPrint {
		log.Printf("ResultCout: " + util.ConvString(totalSaveResultCount))
	}
}

func saveTemp_Win(config loadConfig.ConfigStruct, save []*Type.Temp) {
	if config.IsSaveTemp == false {
		return
	}

	if _, err := Engine.Table("temp_Win").Insert(&save); err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}

	if config.IsPrint {
		log.Printf("TempCout: " + util.ConvString(totalSaveTempCount))
	}
}

func saveTemp_NoWin(config loadConfig.ConfigStruct, save []*Type.Temp) {
	if config.IsSaveTemp == false {
		return
	}

	if _, err := Engine.Table("temp_NoWin").Insert(&save); err != nil {
		log.Fatalf("Fail to sync database: %v\n", err)
	}

	if config.IsPrint {
		log.Printf("TempCout: " + util.ConvString(totalSaveTempCount))
	}
}

func revResult(config loadConfig.ConfigStruct) {
	saveResultCount = 0
	for Value := range resultCh {
		results[saveResultCount] = Value
		saveResultCount++
		totalSaveResultCount++

		if saveResultCount >= config.OnceSaveResultNum {
			saveResult(config, results[:saveResultCount])
			saveResultCount = 0
		}
		resultWG.Done()
	}
}

func revTemp(config loadConfig.ConfigStruct) {
	saveTemp_WinCount = 0
	saveTemp_NoWinCount = 0
	for Value := range tempCh {
		totalSaveTempCount++

		if Value.TotalOdds > 0 {
			temps_Win[saveTemp_WinCount] = Value
			saveTemp_WinCount++
		} else {
			temps_NoWin[saveTemp_NoWinCount] = Value
			saveTemp_NoWinCount++
		}

		if saveTemp_WinCount >= config.OnceSaveTempNum {
			saveTemp_Win(config, temps_Win[:saveTemp_WinCount])
			saveTemp_WinCount = 0
		} else if saveTemp_NoWinCount >= config.OnceSaveTempNum {
			saveTemp_NoWin(config, temps_NoWin[:saveTemp_NoWinCount])
			saveTemp_NoWinCount = 0
		}

		tempWG.Done()
	}
}

func fullRunWheel(show loadXlsx.ShowType, wheel loadXlsx.WheelType) {
	log.Println("數據產出中 請耐心請等待")
	go revResult(loadConfig.Config)
	go revTemp(loadConfig.Config)

	totalSaveResultCount = 0
	totalSaveTempCount = 0
	r := make([]uint8, len(wheel))
	for {
		b := []uint8{}
		for i := 0; i < len(show); i++ {
			for j := 0; j < len(show[i]); j++ {
				index := (uint8(i) + r[j]) % uint8(len(wheel[j]))
				b = append(b, wheel[j][index])
			}
		}

		resultWG.Add(1)
		tempWG.Add(1)
		calLinech <- b
		r[0]++

		for key, Value := range r {
			if int(Value) >= len(wheel[key]) {
				if key == (len(wheel) - 1) {
					return
				} else {
					r[key] = 0
					r[key+1]++
				}
			}
		}
	}

	resultWG.Wait()
	tempWG.Wait()
}
