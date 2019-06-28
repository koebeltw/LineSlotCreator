package orderArrayManager

import (
	"github.com/koebeltw/Common/errorcode"
	"github.com/koebeltw/Common/type"
	"errors"
	"strconv"
	"sync"
)

// OAM OrderArrayManager blabla
type OAM struct {
	rw    sync.RWMutex `bson:"-"`
	datas []interface{}
	max   uint64
	count uint64
}

// NewOAM blabla
func NewOAM(max uint64) *OAM {
	if max == 0 {
		return nil
	}

	return &OAM{max: max, datas: make([]interface{}, max)}
}

func (OAM *OAM) isValid(index uint64) bool {
	return index < OAM.max
}

func (OAM *OAM) isMax(index uint64) bool {
	return index >= OAM.max
}

func (OAM *OAM) isZero() bool {
	return OAM.count == 0
}

func (OAM *OAM) isFull() bool {
	return OAM.count >= OAM.max
}

func (OAM *OAM) IsNil(index uint64) bool {
	return OAM.datas[index] == nil
}

func (OAM *OAM) findEmpty() (index uint64) {
	var count uint64
	for key, vaule := range OAM.datas {
		if OAM.isMax(count) == true {
			return
		}

		if vaule == nil {
			index = uint64(key)
			return
		}

		count = count + 1
	}

	return
}

// Add blabla
func (OAM *OAM) Add(data interface{}) (index uint64, err error) {
	defer OAM.rw.Unlock()
	OAM.rw.Lock()

	if data == nil {
		err = errors.New(strconv.Itoa(errorcode.DataError))
		return
	}

	if OAM.isFull() == true {
		err = errors.New(strconv.Itoa(errorcode.FullError))
	}

	index = OAM.findEmpty()
	OAM.datas[index] = data
	OAM.count = OAM.count + 1
	return
}

// Del blabla
func (OAM *OAM) Del(index uint64) (err error) {
	defer OAM.rw.Unlock()
	OAM.rw.Lock()

	if OAM.isZero() == true {
		err = errors.New(strconv.Itoa(errorcode.ZeroError))
		return
	}

	if OAM.isValid(index) == false {
		err = errors.New(strconv.Itoa(errorcode.IndexError))
		return
	}

	if OAM.datas[index] == nil {
		err = errors.New(strconv.Itoa(errorcode.DataError))
		return
	}

	OAM.count = OAM.count - 1
	OAM.datas[index] = nil
	return
}

// Clear blabla
func (OAM *OAM) Clear() {
	defer OAM.rw.Unlock()
	OAM.rw.Lock()

	if OAM.isZero() == true {
		return
	}

	for key := range OAM.datas {
		OAM.datas[key] = nil
	}

	OAM.count = 0
	return
}

//Put 存儲操作
func (OAM *OAM) Put(index uint64, datas interface{}) (err error) {
	defer OAM.rw.Unlock()
	OAM.rw.Lock()

	if OAM.isValid(index) == false {
		err = errors.New(strconv.Itoa(errorcode.IndexError))
		return
	}

	if OAM.IsNil(index) {
		OAM.count = OAM.count + 1
	}

	OAM.datas[index] = datas
	return
}

//Get 獲取操作
func (OAM *OAM) Get(index uint64) (result interface{}) {
	defer OAM.rw.RUnlock()
	OAM.rw.RLock()

	if OAM.isValid(index) == false {
		return
	}

	return OAM.datas[index]
}

//Find blabla
func (OAM *OAM) Find(checKfunc Type.Checkfunc) (result interface{}) {
	defer OAM.rw.RUnlock()
	OAM.rw.RLock()

	for key, vaule := range OAM.datas {
		if OAM.isMax(uint64(key)) == true {
			return
		}

		if vaule == nil {
			continue
		}

		if ok := checKfunc(key, vaule); ok == true {
			result = vaule
			return
		}
	}

	return
}

// GetForEach blabla
func (OAM *OAM) GetForEach(callBackfunc Type.CallBackfunc) func() {
	return func() { OAM.ForEach(callBackfunc) }
}

// ForEach blabla
func (OAM *OAM) ForEach(callBackfunc Type.CallBackfunc) {
	defer OAM.rw.RUnlock()
	OAM.rw.RLock()

	for key, vaule := range OAM.datas {
		if OAM.isMax(uint64(key)) == true {
			return
		}

		if vaule != nil {
			callBackfunc(key, vaule)
		}
	}
}

// Filter blabla
func (OAM *OAM) Filter(checKfunc Type.Checkfunc) (result []interface{}) {
	defer OAM.rw.RUnlock()
	OAM.rw.RLock()

	for key, vaule := range OAM.datas {
		if OAM.isMax(uint64(key)) == true {
			return
		}

		if vaule != nil {
			if checKfunc(key, vaule) == true {
				result = append(result, vaule)
			}
		}
	}

	return
}

// Every blabla
func (OAM *OAM) Every(checKfunc Type.Checkfunc) (is bool) {
	defer OAM.rw.RUnlock()
	OAM.rw.RLock()

	for key, vaule := range OAM.datas {
		if OAM.isMax(uint64(key)) == true {
			is = true
			return
		}

		if vaule != nil {
			if checKfunc(key, vaule) == false {
				return
			}
		}
	}

	is = true
	return
}

// Any blabla
func (OAM *OAM) Any(checKfunc Type.Checkfunc) (is bool) {
	defer OAM.rw.RUnlock()
	OAM.rw.RLock()

	for key, vaule := range OAM.datas {
		if OAM.isMax(uint64(key)) == true {
			return
		}

		if vaule != nil {
			if checKfunc(key, vaule) == true {
				is = true
				return
			}
		}
	}

	return
}
