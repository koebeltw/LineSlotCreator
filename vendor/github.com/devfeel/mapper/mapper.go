package mapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	ZeroValue                reflect.Value
	DefaultTimeWrapper       *TimeWrapper
	fieldNameMap             sync.Map
	registerMap              sync.Map
	enabledTypeChecking      bool
	enabledMapperStructField bool
	enabledAutoTypeConvert   bool
	timeType                 = reflect.TypeOf(time.Now())
	jsonTimeType             = reflect.TypeOf(JSONTime(time.Now()))
	typeWrappers             []TypeWrapper
)

const (
	packageVersion = "0.6"
	mapperTagKey   = "mapper"
	jsonTagKey     = "json"
	IgnoreTagValue = "-"
	nameConnector  = "_"
	formatTime     = "15:04:05"
	formatDate     = "2006-01-02"
	formatDateTime = "2006-01-02 15:04:05"
)

func init() {
	ZeroValue = reflect.Value{}
	DefaultTimeWrapper = NewTimeWrapper()
	typeWrappers = []TypeWrapper{}
	UseWrapper(DefaultTimeWrapper)

	enabledTypeChecking = false
	enabledMapperStructField = true
	enabledAutoTypeConvert = true
}

func PackageVersion() string {
	return packageVersion
}

// UseWrapper register a type wrapper
func UseWrapper(w TypeWrapper) {
	if len(typeWrappers) > 0 {
		typeWrappers[len(typeWrappers)-1].SetNext(w)
	}
	typeWrappers = append(typeWrappers, w)
}

// CheckIsTypeWrapper check value is in type wrappers
func CheckIsTypeWrapper(value reflect.Value) bool {
	for _, w := range typeWrappers {
		if w.IsType(value) {
			return true
		}
	}
	return false
}

// SetEnabledTypeChecking set enabled flag for TypeChecking
// if set true, the field type will be checked for consistency during mapping
// default is false
func SetEnabledTypeChecking(isEnabled bool) {
	enabledTypeChecking = isEnabled
}

// SetEnabledAutoTypeConvert set enabled flag for auto type convert
// if set true, field will auto convert in Time and Unix
// default is true
func SetEnabledAutoTypeConvert(isEnabled bool) {
	enabledAutoTypeConvert = isEnabled
}

// SetEnabledMapperStructField set enabled flag for MapperStructField
// if set true, the reflect.Struct field will auto mapper
// must follow premises:
// 1. fromField and toField type must be reflect.Struct and not time.Time
// 2. fromField and toField must be not same type
// default is enabled
func SetEnabledMapperStructField(isEnabled bool) {
	enabledMapperStructField = isEnabled
}

// Register register struct to init Map
func Register(obj interface{}) error {
	objValue := reflect.ValueOf(obj)
	if objValue == ZeroValue {
		return errors.New("no exists this value")
	}

	return registerValue(objValue.Elem())
}

// registerValue register Value to init Map
func registerValue(objValue reflect.Value) error {
	regValue := objValue
	if objValue == ZeroValue {
		return errors.New("no exists this value")
	}

	if regValue.Type().Kind() == reflect.Ptr {
		regValue = regValue.Elem()
	}

	typeName := regValue.Type().String()
	for i := 0; i < regValue.NumField(); i++ {
		mapFieldName := typeName + nameConnector + GetFieldName(regValue, i)
		realFieldName := regValue.Type().Field(i).Name
		fieldNameMap.Store(mapFieldName, realFieldName)
	}

	//store register flag
	registerMap.Store(typeName, nil)
	return nil
}

// GetTypeName get type name
func GetTypeName(obj interface{}) string {
	object := reflect.ValueOf(obj)
	return object.String()
}

// CheckExistsField check field is exists by name
func CheckExistsField(elem reflect.Value, fieldName string) (realFieldName string, exists bool) {
	typeName := elem.Type().String()
	fileKey := typeName + nameConnector + fieldName
	realName, isOk := fieldNameMap.Load(fileKey)

	if !isOk {
		return "", isOk
	} else {
		return realName.(string), isOk
	}

}

// GetFieldName get fieldName with ElemValue and index
// if config tag string, return tag value
func GetFieldName(objElem reflect.Value, index int) string {
	fieldName := ""
	field := objElem.Type().Field(index)
	tag := getStructTag(field)
	if tag != "" {
		fieldName = tag
	} else {
		fieldName = field.Name
	}
	return fieldName
}

// MapperMap mapper and set value from map to object
// support auto register struct
// now support field type:
// 1.reflect.Bool
// 2.reflect.String
// 3.reflect.Int8\16\32\64
// 4.reflect.Uint8\16\32\64
// 5.reflect.Float32\64
// 6.time.Time
func MapperMap(fromMap map[string]interface{}, toObj interface{}) error {
	toElem := reflect.ValueOf(toObj).Elem()
	if toElem == ZeroValue {
		return errors.New("to obj is not legal value")
	}
	//check register flag
	//if not register, register it
	if !checkIsRegister(toElem) {
		Register(toObj)
	}
	for k, v := range fromMap {
		fieldName := k
		//check field is exists
		realFieldName, exists := CheckExistsField(toElem, fieldName)
		if !exists {
			continue
		}
		fieldInfo, exists := toElem.Type().FieldByName(realFieldName)
		if !exists {
			continue
		}
		fieldKind := fieldInfo.Type.Kind()
		fieldValue := toElem.FieldByName(realFieldName)
		setFieldValue(fieldValue, fieldKind, v)
	}
	return nil
}

// MapperMapSlice mapper from map[string]map[string]interface{} to a slice of any type's ptr
// toSlice must be a slice of any type's ptr.
func MapperMapSlice(fromMaps map[string]map[string]interface{}, toSlice interface{}) error {
	var err error
	toValue := reflect.ValueOf(toSlice)
	if toValue.Kind() != reflect.Ptr {
		return errors.New("toSlice must pointer of slice")
	}
	if toValue.IsNil() {
		return errors.New("toSlice must not nil pointer")
	}

	toElemType := reflect.TypeOf(toSlice).Elem().Elem()
	if toElemType.Kind() != reflect.Ptr {
		return errors.New("slice elem must ptr ")
	}

	direct := reflect.Indirect(toValue)
	//3 elem parse: 1.[]*type 2.*type 3.type
	toElemType = toElemType.Elem()
	for _, v := range fromMaps {
		elem := reflect.New(toElemType)
		err = MapperMap(v, elem.Interface())
		if err == nil {
			direct.Set(reflect.Append(direct, elem))
		}
	}
	return err
}

// MapToJson mapper from map[string]interface{} to json []byte
func MapToJson(fromMap map[string]interface{}) ([]byte, error) {
	json, err := json.Marshal(fromMap)
	if err != nil {
		return nil, err
	}
	return json, nil
}

// JsonToMap mapper from json []byte to map[string]interface{}
func JsonToMap(body []byte, toMap *map[string]interface{}) error {
	err := json.Unmarshal(body, toMap)
	return err
}

// Mapper mapper and set value from struct fromObj to toObj
// not support auto register struct
func Mapper(fromObj, toObj interface{}) error {
	fromElem := reflect.ValueOf(fromObj).Elem()
	toElem := reflect.ValueOf(toObj).Elem()
	if fromElem == ZeroValue {
		return errors.New("from obj is not legal value")
	}
	if toElem == ZeroValue {
		return errors.New("to obj is not legal value")
	}
	return elemMapper(fromElem, toElem)
}

// MapperSlice mapper from slice of struct to a slice of any type
// fromSlice and toSlice must be a slice of any type.
func MapperSlice(fromSlice, toSlice interface{}) error {
	var err error
	toValue := reflect.ValueOf(toSlice)
	if toValue.Kind() != reflect.Ptr {
		return errors.New("toSlice must pointer of slice")
	}
	if toValue.IsNil() {
		return errors.New("toSlice must not nil pointer")
	}

	elemType := reflect.TypeOf(toSlice).Elem().Elem()
	if elemType.Kind() != reflect.Ptr {
		return errors.New("slice elem must ptr ")
	}

	direct := reflect.Indirect(toValue)
	//3 elem parse: 1.[]*type 2.*type 3.type
	elemType = elemType.Elem()

	fromElems := convertToSlice(fromSlice)
	for _, v := range fromElems {
		elem := reflect.New(elemType)
		err = elemMapper(reflect.ValueOf(v).Elem(), elem.Elem())
		if err == nil {
			direct.Set(reflect.Append(direct, elem))
		}
	}
	return err
}

// Mapper mapper and set value from struct fromObj to toObj
// support auto register struct
func AutoMapper(fromObj, toObj interface{}) error {
	return Mapper(fromObj, toObj)
}

func elemMapper(fromElem, toElem reflect.Value) error {
	//check register flag
	//if not register, register it
	if !checkIsRegister(fromElem) {
		registerValue(fromElem)
	}
	if !checkIsRegister(toElem) {
		registerValue(toElem)
	}

	for i := 0; i < fromElem.NumField(); i++ {
		fromFieldInfo := fromElem.Field(i)
		fieldName := GetFieldName(fromElem, i)
		//check field is exists
		realFieldName, exists := CheckExistsField(toElem, fieldName)
		if !exists {
			continue
		}

		toFieldInfo := toElem.FieldByName(realFieldName)
		//check field is same type
		if enabledTypeChecking {
			if fromFieldInfo.Kind() != toFieldInfo.Kind() {
				continue
			}
		}

		if enabledMapperStructField &&
			toFieldInfo.Kind() == reflect.Struct && fromFieldInfo.Kind() == reflect.Struct &&
			toFieldInfo.Type() != fromFieldInfo.Type() &&
			!CheckIsTypeWrapper(toFieldInfo) && !CheckIsTypeWrapper(fromFieldInfo) {
			x := reflect.New(toFieldInfo.Type()).Elem()
			err := elemMapper(fromFieldInfo, x)
			if err != nil {
				fmt.Println("auto mapper field", fromFieldInfo, "=>", toFieldInfo, "error", err.Error())
			} else {
				toFieldInfo.Set(x)
			}
		} else {
			isSet := false
			if enabledAutoTypeConvert {
				if DefaultTimeWrapper.IsType(fromFieldInfo) && toFieldInfo.Kind() == reflect.Int64 {
					fromTime := fromFieldInfo.Interface().(time.Time)
					toFieldInfo.Set(reflect.ValueOf(TimeToUnix(fromTime)))
					isSet = true
				} else if DefaultTimeWrapper.IsType(toFieldInfo) && fromFieldInfo.Kind() == reflect.Int64 {
					fromTime := fromFieldInfo.Interface().(int64)
					toFieldInfo.Set(reflect.ValueOf(UnixToTime(fromTime)))
					isSet = true
				}
			}
			if !isSet {
				toFieldInfo.Set(fromFieldInfo)
			}
		}

	}
	return nil
}

func setFieldValue(fieldValue reflect.Value, fieldKind reflect.Kind, value interface{}) error {
	switch fieldKind {
	case reflect.Bool:
		if value == nil {
			fieldValue.SetBool(false)
		} else if v, ok := value.(bool); ok {
			fieldValue.SetBool(v)
		} else {
			v, _ := Convert(ToString(value)).Bool()
			fieldValue.SetBool(v)
		}

	case reflect.String:
		if value == nil {
			fieldValue.SetString("")
		} else {
			fieldValue.SetString(ToString(value))
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == nil {
			fieldValue.SetInt(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldValue.SetInt(val.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldValue.SetInt(int64(val.Uint()))
			default:
				v, _ := Convert(ToString(value)).Int64()
				fieldValue.SetInt(v)
			}
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == nil {
			fieldValue.SetUint(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fieldValue.SetUint(uint64(val.Int()))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fieldValue.SetUint(val.Uint())
			default:
				v, _ := Convert(ToString(value)).Uint64()
				fieldValue.SetUint(v)
			}
		}
	case reflect.Float64, reflect.Float32:
		if value == nil {
			fieldValue.SetFloat(0)
		} else {
			val := reflect.ValueOf(value)
			switch val.Kind() {
			case reflect.Float64:
				fieldValue.SetFloat(val.Float())
			default:
				v, _ := Convert(ToString(value)).Float64()
				fieldValue.SetFloat(v)
			}
		}
	case reflect.Struct:
		if value == nil {
			fieldValue.Set(reflect.Zero(fieldValue.Type()))
		} else if DefaultTimeWrapper.IsType(fieldValue) {
			var timeString string
			if fieldValue.Type() == timeType {
				timeString = ""
				fieldValue.Set(reflect.ValueOf(value))
			}
			if fieldValue.Type() == jsonTimeType {
				timeString = ""
				fieldValue.Set(reflect.ValueOf(JSONTime(value.(time.Time))))
			}
			switch d := value.(type) {
			case []byte:
				timeString = string(d)
			case string:
				timeString = d
			case int64:
				if enabledAutoTypeConvert {
					//try to transform Unix time to local Time
					t, err := UnixToTimeLocation(value.(int64), time.UTC.String())
					if err != nil {
						return err
					}
					fieldValue.Set(reflect.ValueOf(t))
				}
			}
			if timeString != "" {
				if len(timeString) >= 19 {
					//满足yyyy-MM-dd HH:mm:ss格式
					timeString = timeString[:19]
					t, err := time.ParseInLocation(formatDateTime, timeString, time.UTC)
					if err == nil {
						t = t.In(time.UTC)
						fieldValue.Set(reflect.ValueOf(t))
					}
				} else if len(timeString) >= 10 {
					//满足yyyy-MM-dd格式
					timeString = timeString[:10]
					t, err := time.ParseInLocation(formatDate, timeString, time.UTC)
					if err == nil {
						fieldValue.Set(reflect.ValueOf(t))
					}
				}
			}
		}
	default:
		if reflect.ValueOf(value).Type() == fieldValue.Type() {
			fieldValue.Set(reflect.ValueOf(value))
		}
	}

	return nil
}

func getStructTag(field reflect.StructField) string {
	tagValue := ""
	//1.check mapperTagKey
	tagValue = field.Tag.Get(mapperTagKey)
	if checkTagValidity(tagValue) {
		return tagValue
	}

	//2.check jsonTagKey
	tagValue = field.Tag.Get(jsonTagKey)
	if checkTagValidity(tagValue) {
		// support more tag property, as json tag omitempty 2018-07-13
		return strings.Split(tagValue, ",")[0]
	}

	return ""
}

func checkTagValidity(tagValue string) bool {
	if tagValue != "" && tagValue != IgnoreTagValue {
		return true
	}
	return false
}

func checkIsRegister(objElem reflect.Value) bool {
	typeName := objElem.Type().String()
	_, isOk := registerMap.Load(typeName)
	return isOk
}

//convert slice interface{} to []interface{}
func convertToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		panic("toslice arr not slice")
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}
