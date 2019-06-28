package util

import (
	"fmt"
	"github.com/cstockton/go-conv"
	"net"
)

func ConvString(form interface{}) string {
	v, _ := conv.String(form)
	return v
}

func ConvUint8(form interface{}) uint8 {
	s, _ := conv.Uint8(form)
	return s
}

func ConvUint16(form interface{}) uint16 {
	v, _ := conv.Uint16(form)
	return v
}

func ConvUint32(form interface{}) uint32 {
	v, _ := conv.Uint32(form)
	return v
}

func ConvUint64(form interface{}) uint64 {
	s, _ := conv.Uint64(form)
	return s
}

func ConvInt8(form interface{}) int8 {
	v, _ := conv.Int8(form)
	return v
}

func ConvInt16(form interface{}) int16 {
	v, _ := conv.Int16(form)
	return v
}

func ConvInt32(form interface{}) int32 {
	v, _ := conv.Int32(form)
	return v
}

func ConvInt64(form interface{}) int64 {
	v, _ := conv.Int64(form)
	return v
}

func ConvFloat32(form interface{}) float32 {
	v, _ := conv.Float32(form)
	return v
}

func ConvFloat64(form interface{}) float64 {
	v, _ := conv.Float64(form)
	return v
}

func GetIPs() (ips []string) {

	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interface addrs: %v", err)
		return ips
	}

	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

