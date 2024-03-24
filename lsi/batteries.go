package lsi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mykolq/zabbix_agent2_plugins/commonfunctions"
)

func GetLsiBatteryInfo(CliDirPath, CliFullPath, CliName, method string, timeout int, UseSudo bool) (Batteries BatteryInfo, errStr error) {

	Batteries.CV = make([]CVInfo, 0)
	Batteries.BBU = make([]BBUInfo, 0)

	BatteryInfoCmdParamMap := map[string][]string{
		"bbu": {"/call/bbu", "show", "all", "j"},
		"cv":  {"/call/cv", "show", "all", "j"}}

	if method == "bbu" {
		BatteriesInfoCmdResult, err := commonfunctions.ExecWithContextTimeout(CliDirPath, CliFullPath, timeout, BatteryInfoCmdParamMap[method], UseSudo)
		errStr = err

		var BatteriesAllInfo LsiBbuResponse
		json.Unmarshal(BatteriesInfoCmdResult, &BatteriesAllInfo)

		Batteries.BBU = GetBBUInfo(BatteriesAllInfo)
	}

	if method == "cv" {
		BatteriesInfoCmdResult, err := commonfunctions.ExecWithContextTimeout(CliDirPath, CliFullPath, timeout, BatteryInfoCmdParamMap[method], UseSudo)
		errStr = err

		var BatteriesAllInfo LsiCvResponse
		json.Unmarshal(BatteriesInfoCmdResult, &BatteriesAllInfo)

		Batteries.CV = GetCacheVaultInfo(BatteriesAllInfo)
	}

	if method == "combined" {
		for k, _ := range BatteryInfoCmdParamMap {
			timeout = timeout / 2
			BatteriesInfoCmdResult, err := commonfunctions.ExecWithContextTimeout(CliDirPath, CliFullPath, timeout, BatteryInfoCmdParamMap[k], UseSudo)
			if err != nil {
				errStr = fmt.Errorf("\n%v", err)
			}

			if k == "bbu" {
				var BatteriesAllInfo LsiBbuResponse
				json.Unmarshal(BatteriesInfoCmdResult, &BatteriesAllInfo)
				Batteries.BBU = GetBBUInfo(BatteriesAllInfo)
			}
			if k == "cv" {
				var BatteriesAllInfo LsiCvResponse
				json.Unmarshal(BatteriesInfoCmdResult, &BatteriesAllInfo)
				Batteries.CV = GetCacheVaultInfo(BatteriesAllInfo)
			}
		}
	}

	if method == "stub" {
		Batteries.CV = nil
		Batteries.BBU = nil
		errStr = nil
	}

	return
}

func GetBBUInfo(BatteriesAllInfo LsiBbuResponse) (BbuBatteries []BBUInfo) {

	var BbuBatteryInfo BBUInfo
	BbuBatteries = make([]BBUInfo, 0)

	for _, CtlBBuInfo := range BatteriesAllInfo.Controllers {
		CommandStatus := CtlBBuInfo.CommandStatus.Status
		if CommandStatus != "Failure" && CommandStatus != "" && CtlBBuInfo.CommandStatus.Description != "No Controller found" {

			BbuBatteryInfo.ControllerId = CtlBBuInfo.CtlResponseData.CtlBasics.Basics.Num
			for _, property := range CtlBBuInfo.CtlResponseData.BBUInfo.BBUInfo {
				if property.Property == "Type" {
					BbuBatteryInfo.Subtype = strings.Replace(string(property.Value), "\"", "", -1)
				}
				if property.Property == "Voltage" {
					BbuBatteryInfo.Voltage = commonfunctions.StrToInt(commonfunctions.ReplaceUnits(string(property.Value)))
				}
				if property.Property == "Current" {
					BbuBatteryInfo.Amperage = commonfunctions.StrToInt(commonfunctions.ReplaceUnits(string(property.Value)))
				}
				if property.Property == "Temperature" {
					BbuBatteryInfo.Temperature = commonfunctions.StrToInt(commonfunctions.ReplaceUnits(string(property.Value)))
				}
				if property.Property == "Battery State" {
					BbuBatteryInfo.State = strings.Replace(string(property.Value), "\"", "", -1)
				}
			}
			for _, property := range CtlBBuInfo.CtlResponseData.BBUInfo.BBUStatus {
				if property.Property == "Charging Status" {
					BbuBatteryInfo.ChargingStatus = strings.Replace(string(property.Value), "\"", "", -1)
				}
				if property.Property == "Replacement required" {
					BbuBatteryInfo.NeedReplace = strings.Replace(string(property.Value), "\"", "", -1)
				}
				if property.Property == "Remaining Capacity Low" {
					BbuBatteryInfo.RemainingCapacityLow = strings.Replace(string(property.Value), "\"", "", -1)
				}
				if property.Property == "Voltage" {
					BbuBatteryInfo.VoltageState = strings.Replace(string(property.Value), "\"", "", -1)
				}
			}
			for _, property := range CtlBBuInfo.CtlResponseData.BBUInfo.BBUCapacityInfo {
				if property.Property == "Relative State of Charge" {
					BbuBatteryInfo.ChargingStatus = strings.Replace(string(property.Value), "\"", "", -1)
				}
			}

			BbuBatteries = append(BbuBatteries, BbuBatteryInfo)
		}
	}
	return
}

func GetCacheVaultInfo(BatteriesAllInfo LsiCvResponse) (CVBatteries []CVInfo) {

	var CvBatteryInfo CVInfo
	CVBatteries = make([]CVInfo, 0)

	for _, CtlCvInfo := range BatteriesAllInfo.Controllers {
		CommandStatus := CtlCvInfo.CommandStatus.Status
		if CommandStatus != "Failure" && CommandStatus != "" && CtlCvInfo.CommandStatus.Description != "No Controller found" {

			CvBatteryInfo.ControllerId = CtlCvInfo.CommandStatus.Controller

			for _, property := range CtlCvInfo.CtlResponseData.CacheVaultInfo {
				if property.Property == "Type" {
					CvBatteryInfo.Subtype = strings.Replace(string(property.Value), "\"", "", -1)
				}
				if property.Property == "State" {
					CvBatteryInfo.State = strings.Replace(string(property.Value), "\"", "", -1)
				}
				if property.Property == "Temperature" {
					tempStr := strings.Replace(string(property.Value), "\"", "", -1)
					CvBatteryInfo.Temperature = commonfunctions.StrToInt(commonfunctions.ReplaceUnits(tempStr))
				}
			}
			for _, property := range CtlCvInfo.CtlResponseData.FirmwareStatus {
				if property.Property == "Replacement required" {
					CvBatteryInfo.NeedReplace = strings.Replace(string(property.Value), "\"", "", -1)
				}
			}
			for _, property := range CtlCvInfo.CtlResponseData.DesignInfo {
				if property.Property == "Serial Number" {
					CvBatteryInfo.SerialNumber = strings.Replace(string(property.Value), "\"", "", -1)
				}
			}

			CVBatteries = append(CVBatteries, CvBatteryInfo)

		}
	}
	return
}
