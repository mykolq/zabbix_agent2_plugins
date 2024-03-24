package lsi

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"git.skbkontur.ru/n.kulikov/zabbix_agent2_plugins/commonfunctions"
)

func fixLsiLdstate(inpstate string) (state string) {
	switch inpstate {
	case "Optl":
		return "Optimal"
	case "OfLn":
		return "Offline"
	case "Pdgd":
		return "Partially degraded"
	case "Dgrd":
		return "Degraded"

	default:
		return inpstate
	}
}

func GetLsiInfo(CliDirPath, CliFullPath, CliName string, timeout int, UseSudo bool) (lsiInfo LsiInfo, err error) {

	var StorageControllers = make([]StorageControllerInfo, 0)
	var LogicalDisks = make([]LogicalDiskInfo, 0)
	var PhysicalDisks = make([]PhysicalDiskInfo, 0)
	var Batteries = make([]LsiCtlBatteryInfo, 0)
	var SmartMap = make(map[string]string)

	LdDgReplaceRe := regexp.MustCompile(`\d{1,}\/`)

	//по умолчанию считаем что
	lsiInfo.Flags.PdInfoFromCtl = false
	lsiInfo.Flags.BatteryExists = false
	lsiInfo.Flags.PdInfoArgs = "stub"
	lsiInfo.Flags.BatteryInfoArgs = "stub"

	CliCtlsAllInfoParam := []string{"/call", "show", "all", "j"}
	CliCtlsAllInfoResult, err := commonfunctions.ExecWithContextTimeout(CliDirPath, CliFullPath, timeout, CliCtlsAllInfoParam, UseSudo)
	if err != nil {
		return lsiInfo, err
	}

	var CtlsAllInfo LsiCtlsInfo
	json.Unmarshal(CliCtlsAllInfoResult, &CtlsAllInfo)

	// Также сделаем выхлоп в стринге, далее пригодится
	CtlsInfoString := string(CliCtlsAllInfoResult)

	isNotOkUtil, _ := regexp.MatchString(`No Controller found`, CtlsInfoString)
	if isNotOkUtil {
		//если сразу не удалось, то заполняем поле Errors и возвращаем
		errText := fmt.Sprintf("Cannot get info with %s. No controllers found", CliName)
		lsiInfo.Errors = &errText
		return
	} else {
		for _, CtlInfo := range CtlsAllInfo.Controllers {
			StorageControllers = append(StorageControllers, StorageControllerInfo{
				Id:           CtlInfo.CtlResponseData.Basics.Num,
				Model:        CtlInfo.CtlResponseData.Basics.Model,
				SerialNumber: CtlInfo.CtlResponseData.Basics.SerialNumber,
				Firmware:     CtlInfo.CtlResponseData.Version.FirmwareVersion,
				State:        CtlInfo.CtlResponseData.Status.ControllerStatus})

			lsiInfo.StorageControllers = &StorageControllers

			if len(CtlInfo.CtlResponseData.VDList) > 0 {
				for _, VdInfo := range CtlInfo.CtlResponseData.VDList {
					LdId := LdDgReplaceRe.ReplaceAllString(VdInfo.DGVD, fmt.Sprintf("/c%d/v", CtlInfo.CtlResponseData.Basics.Num))
					LogicalDisks = append(LogicalDisks, LogicalDiskInfo{
						ControllerId: CtlInfo.CtlResponseData.Basics.Num,
						LdId:         LdId,
						DGVD:         VdInfo.DGVD,
						Consist:      VdInfo.Consist,
						Name:         VdInfo.Name,
						Size:         VdInfo.Size,
						State:        fixLsiLdstate(VdInfo.State),
						Type:         VdInfo.Type})
				}
			}

			lsiInfo.LogicalDisks = &LogicalDisks

			if len(CtlInfo.CtlResponseData.CVInfo) > 0 {
				for _, CvInfo := range CtlInfo.CtlResponseData.CVInfo {
					Batteries = append(Batteries, LsiCtlBatteryInfo{
						ControllerId: CtlInfo.CtlResponseData.Basics.Num,
						Type:         "CV",
						Model:        CvInfo.Model,
						State:        CvInfo.State,
						Temp:         commonfunctions.StrToInt(commonfunctions.ReplaceUnits(CvInfo.Temp))})
				}
				lsiInfo.Flags.BatteryExists = true
			}

			if len(CtlInfo.CtlResponseData.BBUInfo) > 0 {
				for _, BBUInfo := range CtlInfo.CtlResponseData.BBUInfo {
					Batteries = append(Batteries, LsiCtlBatteryInfo{
						ControllerId: CtlInfo.CtlResponseData.Basics.Num,
						Type:         "BBU",
						Model:        BBUInfo.Model,
						State:        BBUInfo.State,
						Temp:         commonfunctions.StrToInt(commonfunctions.ReplaceUnits(BBUInfo.Temp))})
				}
				lsiInfo.Flags.BatteryExists = true
			}

			lsiInfo.Batteries = &Batteries
		}

		HasSeenDisks, _ := regexp.MatchString(`Physical Device Information`, CtlsInfoString)
		if HasSeenDisks {
			lsiInfo.Flags.PdInfoFromCtl = true
			// чтобы метрика с дисками возвращала заглушку. пока ее триггер не отключит
			lsiInfo.Flags.PdInfoArgs = "stub"

			PhysicalDisks = GetPhysicalDisks(CtlsInfoString, "Response Data.Physical Device Information")
			lsiInfo.PhysicalDisks = &PhysicalDisks

			SmartMap = GetSmartMap(PhysicalDisks)
			lsiInfo.SmartMap = &SmartMap

		} else {
			EnclosuresNotExistsRegexp := regexp.MustCompile(`EID:Slt(.*)?\s:\d{1,}`)
			EnclosuresExistsRegexp := regexp.MustCompile(`EID:Slt(.*)?\d{1,}:\d{1,}`)
			if len(EnclosuresExistsRegexp.FindAllString(CtlsInfoString, -1)) > 0 {
				lsiInfo.Flags.PdInfoArgs = "enclosure"
			}
			if len(EnclosuresNotExistsRegexp.FindAllString(CtlsInfoString, -1)) > 0 {
				if lsiInfo.Flags.PdInfoArgs == "enclosure" {
					lsiInfo.Flags.PdInfoArgs = "combined"
				} else {
					lsiInfo.Flags.PdInfoArgs = "noenclosure"
				}
			}
		}

		HasCacheVault, _ := regexp.MatchString(`Cachevault_Info`, CtlsInfoString)
		HasBbu, _ := regexp.MatchString(`BBU_Info`, CtlsInfoString)

		if HasCacheVault {
			lsiInfo.Flags.BatteryInfoArgs = "cv"
		}
		if HasBbu {
			lsiInfo.Flags.BatteryInfoArgs = "bbu"
		}
		if HasCacheVault && HasBbu {
			lsiInfo.Flags.BatteryInfoArgs = "combined"
		}
	}
	return
}

func PrintCliOutput(CliDirPath, CliFullPath, CliName string,
	timeout int, UseSudo bool,
	CliParam string) string {

	CliParamSplitted := strings.Split(CliParam, " ")

	CliResult, err := commonfunctions.ExecWithContextTimeout(CliDirPath, CliFullPath, timeout, CliParamSplitted, UseSudo)
	if err != nil {
		return fmt.Sprintf("Cannot get info with %s\nError: %s", err, CliName)
	} else {
		return string(CliResult)
	}
}
