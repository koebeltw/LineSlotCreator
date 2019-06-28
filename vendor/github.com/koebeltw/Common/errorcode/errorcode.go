package errorcode

type ErrorCode uint16

const (
	SusseceError ErrorCode = 0
	IndexError             = 1 //索引不對
	SearchError            = 2 //搜尋不到
	DataError              = 3 //資料有誤
	AddError               = 4 //新增失敗
	DelError               = 5 //刪除失敗
	FullError              = 6 //容器已滿
	ZeroError              = 7 //容器全空
	ExistError             = 8 //存在錯誤
)
