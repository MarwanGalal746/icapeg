package service

import (
	"icapeg/service/services/clamav"
	"icapeg/service/services/echo"
	"icapeg/service/services/virustotal"
	"icapeg/utils"
)

//Vendors names
const (
	VendorEcho       = "echo"
	VendorClamav     = "clamav"
	VendorVirustotal = "virustotal"
)

type (
	// Service holds the info to distinguish a service
	Service interface {
		Processing(bool) (int, interface{}, map[string]string)
		ISTagValue() string
	}
)

// GetService returns a service based on the service name
// change name to vendor and add parameter service name
func GetService(vendor, serviceName, methodName string, httpMsg *utils.HttpMsg) Service {
	switch vendor {
	case VendorEcho:
		return echo.NewEchoService(serviceName, methodName, httpMsg)
	case VendorVirustotal:
		return virustotal.NewVirustotalService(serviceName, methodName, httpMsg)
	case VendorClamav:
		return clamav.NewClamavService(serviceName, methodName, httpMsg)

	}
	return nil
}

func InitServiceConfig(vendor, serviceName string) {
	switch vendor {
	case VendorEcho:
		echo.InitEchoConfig(serviceName)
	case VendorClamav:
		clamav.InitClamavConfig(serviceName)
	case VendorVirustotal:
		virustotal.InitVirustotalConfig(serviceName)

	}
}
