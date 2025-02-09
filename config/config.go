package config

import (
	"fmt"
	"github.com/spf13/viper"
	"icapeg/logging"
	"icapeg/readValues"
	"os"
	"strings"
)

type serviceIcapInfo struct {
	Vendor         string
	ServiceCaption string
	ServiceTag     string
	ReqMode        bool
	RespMode       bool
	ShadowService  bool
	PreviewEnabled bool
	PreviewBytes   string
}

// AppConfig represents the app configuration
type AppConfig struct {
	Port                 int
	LogLevel             string
	LoggingServerURL     string
	LoggingFlushDuration float64
	WriteLogsToConsole   bool
	//RespScannerVendor       string
	//ReqScannerVendor        string
	RespScannerVendorShadow string
	ReqScannerVendorShadow  string
	BypassExtensions        []string
	ProcessExtensions       []string
	PreviewBytes            string
	PreviewEnabled          bool
	PropagateError          bool
	VerifyServerCert        bool
	DebuggingHeaders        bool
	Services                []string
	ServicesInstances       map[string]*serviceIcapInfo
}
type Extension struct {
	Name string
	Exts []string
}

var AppCfg AppConfig

// Init initializes the configuration
func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if readValues.IsSecExists("app") {
		fmt.Println("app section doesn't exist in config file")
	}
	AppCfg = AppConfig{
		Port:                    readValues.ReadValuesInt("app.port"),
		LogLevel:                readValues.ReadValuesString("app.log_level"),
		LoggingServerURL:        readValues.ReadValuesString("app.log_service_url"),
		WriteLogsToConsole:      readValues.ReadValuesBool("app.write_logs_to_console"),
		RespScannerVendorShadow: strings.ToLower(readValues.ReadValuesString("app.resp_scanner_vendor_shadow")),
		ReqScannerVendorShadow:  strings.ToLower(readValues.ReadValuesString("app.req_scanner_vendor_shadow")),
		PreviewBytes:            readValues.ReadValuesString("app.preview_bytes"),
		PreviewEnabled:          readValues.ReadValuesBool("app.preview_enabled"),
		PropagateError:          readValues.ReadValuesBool("app.propagate_error"),
		VerifyServerCert:        readValues.ReadValuesBool("app.verify_server_cert"),
		DebuggingHeaders:        readValues.ReadValuesBool("app.debugging_headers"),
		Services:                readValues.ReadValuesSlice("app.services"),
	}

	//this loop to make sure that all services in the array of services has sections in the config file and from request mode and response mode
	//there is one at least from them are enabled in every service
	AppCfg.ServicesInstances = make(map[string]*serviceIcapInfo)
	for i := 0; i < len(AppCfg.Services); i++ {
		serviceName := AppCfg.Services[i]
		if !readValues.IsSecExists(serviceName) {
			fmt.Println(serviceName + " section doesn't exist")
			os.Exit(1)
		}
		if !readValues.ReadValuesBool(serviceName+".req_mode") && !readValues.ReadValuesBool(serviceName+".resp_mode") {
			fmt.Println("Request mode and response mode are disabled together in " + serviceName + " service")
			os.Exit(1)
		}
		if readValues.ReadValuesInt(serviceName+".max_filesize") < 0 {
			fmt.Println("max_filesize value in config.toml file is not valid")
			os.Exit(1)
		}

		ext := make(map[string]bool)
		asterisks := 0
		//bypass
		bypass := readValues.ReadValuesSlice(serviceName + ".bypass_extensions")
		for i := 0; i < len(bypass); i++ {
			if bypass[i] == "*" && len(bypass) != 1 {
				fmt.Println("bypass_extensions array has one asterisk \"*\" and other extensions but asterisk should be the only element in the array otherwise add extensions as you want")
				os.Exit(1)
			}
			if bypass[i] == "*" {
				asterisks++
			}
			if ext[bypass[i]] == false {
				ext[bypass[i]] = true
			} else {
				fmt.Println("This extension \"" + bypass[i] + "\" was stored in multiple arrays (bypass_extensions or reject_extensions)")
				os.Exit(1)
			}
		}
		//process
		process := readValues.ReadValuesSlice(serviceName + ".process_extensions")
		for i := 0; i < len(process); i++ {
			if process[i] == "*" && len(process) != 1 {
				fmt.Println("process_extensions array has one asterisk \"*\" and other extensions but asterisk should be the only element in the array otherwise add extensions as you want")
				os.Exit(1)
			}
			if process[i] == "*" {
				asterisks++
			}
			if ext[process[i]] == false {
				ext[process[i]] = true
			} else {
				fmt.Println("This extension \"" + process[i] + "\" is stored in multiple arrays")
				os.Exit(1)
			}
		}
		//reject
		reject := readValues.ReadValuesSlice(serviceName + ".reject_extensions")
		for i := 0; i < len(reject); i++ {
			if reject[i] == "*" && len(reject) != 1 {
				fmt.Println("reject_extensions array has one asterisk \"*\" and other extensions but asterisk should be the only element in the array otherwise add extensions as you want")
				os.Exit(1)
			}
			if reject[i] == "*" {
				asterisks++
			}
			if ext[reject[i]] == false {
				ext[reject[i]] = true
			} else {
				fmt.Println("This extension \"" + reject[i] + "\" is stored in multiple arrays")
				os.Exit(1)
			}
		}
		if asterisks != 1 {
			fmt.Println("There is no \"*\" stored in any extension arrays")
			os.Exit(1)
		}

		AppCfg.ServicesInstances[serviceName] = &serviceIcapInfo{
			Vendor:         readValues.ReadValuesString(serviceName + ".vendor"),
			ServiceTag:     readValues.ReadValuesString(serviceName + ".service_tag"),
			ServiceCaption: readValues.ReadValuesString(serviceName + ".service_caption"),
			ReqMode:        readValues.ReadValuesBool(serviceName + ".req_mode"),
			RespMode:       readValues.ReadValuesBool(serviceName + ".resp_mode"),
			ShadowService:  readValues.ReadValuesBool(serviceName + ".shadow_service"),
			PreviewBytes:   readValues.ReadValuesString(serviceName + ".preview_bytes"),
			PreviewEnabled: readValues.ReadValuesBool(serviceName + ".preview_enabled"),
		}
	}
	logging.InitLogger()
}

// InitTestConfig initializes the app with the test config file (for integration test)
func InitTestConfig() {
	AppCfg = AppConfig{
		Port:                 readValues.ReadValuesInt("app.port"),
		LogLevel:             readValues.ReadValuesString("app.log_level"),
		LoggingServerURL:     readValues.ReadValuesString("app.log_service_url"),
		LoggingFlushDuration: float64(readValues.ReadValuesInt("app.log_flush_duration")),
		//RespScannerVendor:       strings.ToLower(readValues.ReadValuesString("app.resp_scanner_vendor")),
		//ReqScannerVendor:        strings.ToLower(readValues.ReadValuesString("app.req_scanner_vendor")),
		RespScannerVendorShadow: strings.ToLower(readValues.ReadValuesString("app.resp_scanner_vendor_shadow")),
		ReqScannerVendorShadow:  strings.ToLower(readValues.ReadValuesString("app.req_scanner_vendor_shadow")),
		PreviewBytes:            readValues.ReadValuesString("app.preview_bytes"),
		PropagateError:          readValues.ReadValuesBool("app.propagate_error"),
		DebuggingHeaders:        readValues.ReadValuesBool("app.debugging_headers"),
	}
}

// App returns the the app configuration instance
func App() *AppConfig {
	return &AppCfg
}
