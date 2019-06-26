package loadConfig

import (
	"fmt"
	"github.com/zieckey/goini"
	"io/ioutil"
	"path/filepath"
)

type ConfigStruct struct {
	GSIP string `json:"-"`
	LSPort string `json:"-"`
	//SQLAccount        string `json:"SQL帳號"`
	//SQLPassword       string `json:"SQL密碼"`
	//SQLIP             string `json:"SQLIP"`
	//SQLPort           string `json:"SQLPort"`
	//CalLineNum        uint32 `json:"同時計算連線的最大數"`
	//IsPrint           bool   `json:"是否顯示訊息"`
	//IsReWriteData     bool   `json:"是否覆蓋資料"`
	//IsSaveResult      bool   `json:"是否儲存Result"`
	//OnceSaveResultNum uint32 `json:"Result一次寫入SQL的最大數"`
	//IsSaveTemp        bool   `json:"是否儲存Temp"`
	//OnceSaveTempNum   uint32 `json:"Temp一次寫入SQL的最大數"`
}

func init() {
	loadConfig()
}

var Config ConfigStruct

func loadConfig() {
	contents, _ := ioutil.ReadFile("./config.ini")

	filename := filepath.Join("./config.ini")
	ini := goini.New()
	if err := ini.ParseFile(filename); err != nil {
		if err := ini.Parse(contents[3:], goini.DefaultLineSeparator, goini.DefaultKeyValueSeparator); err != nil {
			fmt.Printf("parse INI file %v failed : %v\n", filename, err.Error())
			return
		}
	}

	Config.GSIP, _ = ini.Get("GSIP")
	Config.LSPort, _ = ini.Get("LSPort")

	//Config.SQLAccount, _ = ini.Get("SQLAccount")
	//Config.SQLPassword, _ = ini.Get("SQLPassword")
	//Config.SQLIP, _ = ini.Get("SQLIP")
	//Config.SQLPort, _ = ini.Get("SQLPort")
	//Config.SQLIP, _ = ini.Get("SQLIP")
	//
	//if v, ok := ini.GetInt("CalLineNum"); ok {
	//	Config.CalLineNum, _ = conv.Uint32(v)
	//}
	//
	//if v, ok := ini.GetInt("IsSaveResult"); ok {
	//	Config.IsSaveResult, _ = conv.Bool(v)
	//}
	//
	//if v, ok := ini.GetInt("IsSaveTemp"); ok {
	//	Config.IsSaveTemp, _ = conv.Bool(v)
	//}
	//
	//if v, ok := ini.GetInt("OnceSaveResultNum"); ok {
	//	Config.OnceSaveResultNum, _ = conv.Uint32(v)
	//}
	//
	//if v, ok := ini.GetInt("OnceSaveTempNum"); ok {
	//	Config.OnceSaveTempNum, _ = conv.Uint32(v)
	//}
	//
	//if v, ok := ini.GetInt("IsPrint"); ok {
	//	Config.IsPrint, _ = conv.Bool(v)
	//}
	//
	//if v, ok := ini.GetInt("IsReWriteData"); ok {
	//	Config.IsReWriteData, _ = conv.Bool(v)
	//}
}
