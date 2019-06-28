package tcp

import (
	"math"
	"reflect"
	"strings"
	"fmt"
	"github.com/cstockton/go-conv"
)

// Eventfunc blabla
type Eventfunc func(s Session, b []byte)

// eventHandler blabla
type EventHandler [math.MaxUint8 + 1][math.MaxUint8 + 1]func(session Session, b []byte)

// GetEventHandler blabla
func GetEventHandler(event interface{}, method func(s Session, b []byte)) (e *EventHandler) {
	e = &EventHandler{}

	if method != nil {
		e.SetAll(method)
	}

	rValue := reflect.ValueOf(event)
	rType := reflect.TypeOf(event)

	if rType.Kind() != reflect.Struct {
		return
	}

	args := make([]reflect.Value, 0)

	for i := 0; i < rValue.NumMethod(); i++ {
		//fmt.Println(runtime.FuncForPC(rValue.Method(i).Pointer()).Name())
		result := rValue.Method(i).Call(args)

		name := rType.Method(i).Name
		if strings.Index(name, `Event`) < 0 {
			continue
		}

		name = strings.Replace(rType.Method(i).Name, "Event", "", -1)
		strs := strings.Split(name, "Dash")
		var err error
		var msgNo, subNo uint8
		msgNo, err = conv.Uint8(strs[0])
		if err != nil {
			fmt.Println("msgNo Error:", err)
			continue
		}
		subNo, err = conv.Uint8(strs[1])
		if err != nil {
			fmt.Println("subNo Error:", err)
			continue
		}

		// log.Println(result[0].Interface().(byte), result[1].Interface().(byte))
		e[msgNo][subNo] = result[0].Interface().(Eventfunc)
	}

	return
}

// Set blabla
func (e EventHandler) Set(msgNo byte, subNo byte, method Eventfunc) {
	e[msgNo][subNo] = method
}

// SetAll blabla
func (e EventHandler) SetAll(method Eventfunc) {
	for i := 0; i < len(e); i++ {
		for j := 0; j < len(e[i]); j++ {
			e[i][j] = method
		}
	}

}
