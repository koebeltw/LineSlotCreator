package Type

type Function func()

// ChecKfunc blabla
type Checkfunc func(key interface{}, value interface{}) bool

// CallBackfunc blabla
type CallBackfunc func(key interface{}, value interface{})

type Addr struct {
	ID   int
	Name string
	IP   string
	Port int
}

type IP struct {
	Local  Addr
	Server Addr
	Client []Addr
}

type Prize struct {
	LinePicNo       []uint16 `json:"連線圖案"`
	LinePicPosition []uint16 `json:"連線位置"`
	LineNo          uint16   `json:"獎項編號"`
	Odds            float64  `json:"賠率"`
	WinPic          uint16    `json:"圖案編號"`
	PicCount        uint16    `json:"圖案數量"`
	WildCount       uint16    `json:"Wild數量"`
}

type Result struct {
	ID        uint64 `xorm:"pk autoincr index"`
	TotalOdds float64
	Wheel     []uint16
}

type Temp struct {
	ID           uint64   `json:"-" xorm:"pk autoincr index"`
	FreeGameCount uint16 `json:"FreeGame次數"`
	Positions    []uint16 `json:"位置"`
	TotalOdds    float64  `json:"總賠率"`
	//JumboWin     uint8    `json:"大獎種類"`
	WheelUint16  []uint16 `json:"滾輪"`
	ScatterCount uint16    `json:"Scatter數量"`
	Prize        []Prize  `json:"獎項"`
}

type OddsCount struct {
	Odds  uint32
	Count uint32
}

type LineOddsCount struct {
	LineNo     uint16
	WinPic     uint16
	PicCount   uint16
	TotalCount uint32
	TotalOdds  uint32
}

type Record struct {
	TotalCount uint32
	TotalOdds  uint32
}

