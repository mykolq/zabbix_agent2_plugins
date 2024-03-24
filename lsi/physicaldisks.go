package lsi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	comf "github.com/mykolq/zabbix_agent2_plugins/commonfunctions"
	"github.com/tidwall/gjson"
)

func fixLsiPdstate(inpstate string) (state string) {
	switch inpstate {
	case "Onln":
		return "Online"
	case "Ugood":
		return "Unconfigured-Good"
	case "Ubad":
		return "Unconfigured-Bad"
	case "Msng":
		return "Missing"
	case "Offln":
		return "Offline"
	case "F":
		return "Foreign"
	case "GHS":
		return "Global Hot Spare"
	case "DHS":
		return "Dedicated Hot Spare"
	case "Rbld":
		return "Rebuild"
	case "Cpybck":
		return "Copyback"

	default:
		return inpstate
	}
}

func GetLsiPdInfo(DataString, JsonPathConst, CliDirPath, CliFullPath, CliName, method string,
	timeout int,
	UseSudo bool) (lsiPdInfo LsiPdInfo, errStr error) {

	var PhysicalDisks = make([]PhysicalDiskInfo, 0)
	var SmartMap = make(map[string]string)

	DisksInfoCmdParamMap := map[string][]string{
		"enclosure":   {"/call/eall/sall", "show", "all", "j"},
		"noenclosure": {"/call/sall", "show", "all", "j"}}

	if method == "enclosure" {
		DisksInfoCmdResult, err := comf.ExecWithContextTimeout(CliDirPath, CliFullPath, timeout, DisksInfoCmdParamMap[method], UseSudo)
		errStr = err
		PhysicalDisks = GetPhysicalDisks(string(DisksInfoCmdResult), "Response Data")
	}

	if method == "noenclosure" {
		DisksInfoCmdResult, err := comf.ExecWithContextTimeout(CliDirPath, CliFullPath, timeout, DisksInfoCmdParamMap[method], UseSudo)

		errStr = err

		PhysicalDisks = GetPhysicalDisks(string(DisksInfoCmdResult), "Response Data")
	}

	if method == "combined" {
		for k, _ := range DisksInfoCmdParamMap {
			timeout = timeout / 2
			DisksInfoCmdResult, err := comf.ExecWithContextTimeout(CliDirPath, CliFullPath, timeout, DisksInfoCmdParamMap[k], UseSudo)
			if err != nil {
				errStr = fmt.Errorf("\n%v", err)
			}

			PhysicalDisks = append(PhysicalDisks, GetPhysicalDisks(string(DisksInfoCmdResult), "Response Data")...)
		}
	}

	if method == "stub" {
		return
	}

	if len(PhysicalDisks) == 0 {
		return
	}

	lsiPdInfo.PhysicalDisks = &PhysicalDisks

	SmartMap = GetSmartMap(PhysicalDisks)
	lsiPdInfo.SmartMap = &SmartMap

	return lsiPdInfo, errStr
}

func GetPhysicalDisks(DataString string, JsonPathConst string) (PhysicalDisks []PhysicalDiskInfo) {

	var DiskInfo PhysicalDiskInfo
	PhysicalDisks = make([]PhysicalDiskInfo, 0)

	SlotNameRegex := regexp.MustCompile(`Drive (\/c(\d{1,}).*\/s\d{1,})"`)
	DiskSizeRegex := regexp.MustCompile(`(\d{1,}\.*\d*)(\s*)(\w*).*`)
	SlotNamesArr := SlotNameRegex.FindAllString(DataString, -1)
	for _, SlotName := range SlotNamesArr {
		SlotNameNormalized := SlotNameRegex.ReplaceAllString(SlotName, "${1}")
		CtlNum, _ := strconv.ParseInt(SlotNameRegex.ReplaceAllString(SlotName, "${2}"), 10, 8)
		DiscAttrJsonPath := fmt.Sprintf(
			"Controllers.%d.%s.Drive %s - Detailed Information.Drive %s Device attributes",
			CtlNum, JsonPathConst, SlotNameNormalized, SlotNameNormalized)

		DiskAttrsMap := (gjson.Get(DataString, DiscAttrJsonPath)).Map()

		WWN := strings.TrimSpace(DiskAttrsMap["WWN"].String())

		if WWN != "NA" {
			DiskInfo.Slot, DiskInfo.SlotnameSource = SlotNameNormalized, "LSI"

			DiscStatesJsonPath := fmt.Sprintf(
				"Controllers.%d.%s.Drive %s - Detailed Information.Drive %s State",
				CtlNum, JsonPathConst, SlotNameNormalized, SlotNameNormalized)

			DiskInfo.Interface = fmt.Sprintf("%s", gjson.Get(DataString, fmt.Sprintf(
				"Controllers.%d.%s.Drive %s.0.Intf",
				CtlNum, JsonPathConst, SlotNameNormalized)))
			DiskInfo.Type = fmt.Sprintf("%s", gjson.Get(DataString, fmt.Sprintf(
				"Controllers.%d.%s.Drive %s.0.Med",
				CtlNum, JsonPathConst, SlotNameNormalized)))

			DiskInfo.State = fixLsiPdstate(fmt.Sprintf("%s", gjson.Get(DataString, fmt.Sprintf(
				"Controllers.%d.%s.Drive %s.0.State",
				CtlNum, JsonPathConst, SlotNameNormalized))))

			DiskStatesMap := (gjson.Get(DataString, DiscStatesJsonPath)).Map()

			DiskInfo.Size = DiskSizeRegex.ReplaceAllString(DiskAttrsMap["Raw size"].String(), "$1 $3")

			DiskInfo.Model, DiskInfo.Vendor = comf.FixPhysicalDiskInfo(DiskAttrsMap["Model Number"].String(),
				strings.TrimSpace(DiskAttrsMap["Manufacturer Id"].String()), "", nil)

			DiskInfo.SerialNumber = strings.TrimSpace(DiskAttrsMap["SN"].String())

			DiskInfo.Firmware = DiskAttrsMap["Firmware Revision"].String()

			DiskInfo.MediaErrCount, DiskInfo.OtherErrCount = DiskStatesMap["Media Error Count"].Uint(), DiskStatesMap["Other Error Count"].Uint()
			DiskInfo.PredictErrCount = DiskStatesMap["Predictive Failure Count"].Uint()

			DiskInfo.SmartFlag = strings.TrimSpace(DiskStatesMap["S.M.A.R.T alert flagged by drive"].String())
			DiskInfo.Firmware = DiskAttrsMap["Firmware Revision"].String()

			PhysicalDisks = append(PhysicalDisks, DiskInfo)
		}
	}
	return
}

func GetSmartMap(PhysicalDisks []PhysicalDiskInfo) map[string]string {

	SmartMap := make(map[string]string)

	for _, disk := range PhysicalDisks {
		SmartMap[disk.SerialNumber] = disk.Slot
	}

	return SmartMap
}
