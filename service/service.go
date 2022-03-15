package service

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"icapeg/dtos"
	ic "icapeg/icap-client"
	"icapeg/logger"
	"icapeg/service/clamav"
	"icapeg/service/glasswall"
	"icapeg/service/metadefender"
	"icapeg/service/virustotal"
	"icapeg/service/vmray"
)

// The service names
const (
	VENVirusTotal   = "virustotal"
	VENMetaDefender = "metadefender"
	VENVmray        = "vmray"
	SVCClamav       = "clamav"
	VENGlasswall    = "glasswall"
)

type (
	// Service holds the info to distinguish a service
	Service interface {
		SubmitFile(*bytes.Buffer, string) (*dtos.SubmitResponse, error)
		SendFileApi(*bytes.Buffer, string, string) (*http.Response, int, bool, string, error)
		GetSubmissionStatus(string) (*dtos.SubmissionStatusResponse, error)
		GetSampleFileInfo(string, ...dtos.FileMetaInfo) (*dtos.SampleInfo, error)
		GetSampleURLInfo(string, ...dtos.FileMetaInfo) (*dtos.SampleInfo, error)
		SubmitURL(string, string) (*dtos.SubmitResponse, error)
		GetStatusCheckInterval() time.Duration
		GetStatusCheckTimeout() time.Duration
		GetBadFileStatus() []string
		GetOkFileStatus() []string
		StatusEndpointExists() bool
		RespSupported() bool
		ReqSupported() bool
	}

	// LocalService holds the blueprint of a local service
	LocalService interface {
		ScanFileStream(io.Reader, dtos.FileMetaInfo) (*dtos.SampleInfo, error)
		GetBadFileStatus() []string
		GetOkFileStatus() []string
		RespSupported() bool
		ReqSupported() bool
	}

	// ICAPService holds the blueprint of a Remote ICAP service
	ICAPService interface {
		DoReqmod() (*ic.Response, error)
		DoRespmod() (*ic.Response, error)
		DoOptions() (*ic.Response, error)
		GetURL() string
		GetRespmodEndpoint() string
		GetReqmodEndpoint() string
		GetOptionsEndpoint() string
		GetTimeout() time.Duration
		SetHTTPRequest(*http.Request)
		SetHTTPResponse(*http.Response)
		SetHeader(map[string][]string)
		ChangeOptionsEndpoint(string)
	}
)

// IsServiceLocal determines if a service is local or not
func IsServiceLocal(vendor string, serviceName string, logger *logger.ZLogger) bool {
	svc := GetService(vendor, serviceName, serviceName, logger)

	if svc != nil {
		return false
	}

	lsvc := GetLocalService(vendor, serviceName)

	if lsvc != nil {
		return true
	}

	return false
}

// GetService returns a service based on the service name
// change name to vendor and add parameter service name
func GetService(vendorName string, serviceName string, methodName string, logger *logger.ZLogger) Service {
	switch vendorName {
	case VENVirusTotal:
		return virustotal.NewVirusTotalService(serviceName, logger)
	case VENMetaDefender:
		return metadefender.NewMetaDefenderService(serviceName, logger)
	case VENVmray:
		return vmray.NewVmrayService(serviceName, logger)
	case VENGlasswall:
		return glasswall.NewGlasswallService(serviceName, methodName, logger)
	}
	return nil
}

// GetLocalService returns a local service based on the name
func GetLocalService(vendor string, serviceName string) LocalService {
	switch vendor {
	case SVCClamav:
		return clamav.NewClamavService(serviceName)
	}

	return nil
}

// GetICAPService returns a remote ICAP service based on the name
func GetICAPService(name string) ICAPService {
	return NewRemoteICAPService(name)
}
