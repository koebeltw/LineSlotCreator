package webSocketHandler

import (
	"github.com/koebeltw/Common/util"
	"github.com/koebeltw/Common/LineSlot/slotHandler"
	"github.com/koebeltw/LineSlotCreator/ls/loadConfig"
	"fmt"
	"github.com/devfeel/mapper"
	"github.com/pquerna/ffjson/ffjson"
	"github.com/zieckey/goini"
	"io/ioutil"
	"path/filepath"
	"strings"
)

var Account string
var Password string

func init() {

	contents, _ := ioutil.ReadFile("./user.ini")
	filename := filepath.Join("./user.ini")
	ini := goini.New()
	if err := ini.ParseFile(filename); err != nil {
		if err := ini.Parse(contents[3:], goini.DefaultLineSeparator, goini.DefaultKeyValueSeparator); err != nil {
			fmt.Printf("parse INI file %v failed : %v\n", filename, err.Error())
			return
		}
	}

	Account, _ = ini.Get("Account")
	Password, _ = ini.Get("Password")
}

func CanLogin(msg []byte) (r []byte, ok bool) {
	m := make(map[string]interface{}, 0)
	ffjson.Unmarshal(msg, &m)
	r = []byte{}
	switch {
	case m["登入"] != nil:
		v := m["登入"].(map[string]interface{})
		if canLogin(v["account"].(string), v["password"].(string)) {
			r, _ = ffjson.Marshal(map[string]interface{}{"登入": true})
			ok = true
		} else {
			r, _ = ffjson.Marshal(map[string]interface{}{"登入": false})
			ok = false
		}
	default:
		ok = false
	}

	return
}

func GetMsg(msg []byte) (r []byte) {
	m := make(map[string]interface{}, 0)
	ffjson.Unmarshal(msg, &m)
	r = []byte{}
	switch {
	case m["要求設定"] != nil:
		r, _ = ffjson.Marshal(map[string]interface{}{"要求設定": loadConfig.Config})
	case m["改變設定"] != nil:
		mapper.MapperMap(m["改變設定"].(map[string]interface{}), &loadConfig.Config)
	case m["數據"] != nil:
		type Data struct {
			GameKind uint8
			Count    uint32
			IntoFreeGameRate int32
		}
		var data Data
		mapper.MapperMap(m["數據"].(map[string]interface{}), &data)
		r, _ = ffjson.Marshal(map[string]interface{}{"數據": slotHandler.RandomRunWheel(data.GameKind, data.Count, data.IntoFreeGameRate)})
	case m["測試"] != nil:
		v := m["測試"].(map[string]interface{})
		type Test struct {
			GameKind      uint8
			FreeGameCount int32
			Positions      []uint16
		}

		Positions := make([]uint16, 0, 255)
		if v["Positions"] != nil {
			s := util.ConvString(v["Positions"])
			s = strings.Replace(s, "[", "", -1)
			s = strings.Replace(s, "]", "", -1)
			sa := strings.Split(s, " ")
			for _, value := range sa {
				Positions = append(Positions, util.ConvUint16(value))
			}
		}

		r, _ = ffjson.Marshal(map[string]interface{}{"測試": slotHandler.PlayOnce(util.ConvUint8(v["GameKind"]), util.ConvUint16(v["FreeGameCount"]), Positions, 0)})
	default:
	}

	return
}

//type WebSocketConn struct {
//	conn *websocket.Conn
//}

//var WebSocketUsers sync.Map

func canLogin(account string, password string) (bool) {
	return Account == account && Password == password
}
