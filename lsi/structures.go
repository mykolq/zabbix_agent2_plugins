package lsi

import "encoding/json"

type LsiCommandStatus struct {
	Status      string `json:"Status"`
	Description string `json:"Description"`
	Controller  int    `json:"Controller"`
}

type LsiCtlBasics struct {
	Num          int    `json:"Controller"`
	Model        string `json:"Model"`
	SerialNumber string `json:"Serial Number"`
}

type LsiCtlVersions struct {
	FirmwareVersion string `json:"Firmware Version"`
	BiosVersion     string `json:"Bios Version"`
	DriverName      string `json:"Driver Name"`
	DriverVersion   string `json:"Driver Version"`
}

type LsiCtlstatus struct {
	ControllerStatus string `json:"Controller Status"`
}

type LsiCtlResponseData struct {
	Basics  LsiCtlBasics    `json:"Basics"`
	Version LsiCtlVersions  `json:"Version"`
	Status  LsiCtlstatus    `json:"Status"`
	VDList  []LsiLdInfo     `json:"VD LIST"`
	BBUInfo []LsiCtlBattery `json:"BBU_Info"`
	CVInfo  []LsiCtlBattery `json:"Cachevault_Info"`
}

type LsiCtlBattery struct {
	ControllerId int
	Type         string
	Model        string `json:"Model"`
	State        string `json:"State"`
	Temp         string `json:"Temp"`
}

type LsiCtlBatteryInfo struct {
	ControllerId int
	Type         string
	Model        string
	State        string
	Temp         int
}

type LsiResponse struct {
	CommandStatus   LsiCommandStatus   `json:"Command Status"`
	CtlResponseData LsiCtlResponseData `json:"Response Data"`
}

type LsiCtlsInfo struct {
	Controllers []LsiResponse `json:"Controllers"`
}

type LsiBatteryProperty struct {
	Property string          `json:"Property"`
	Value    json.RawMessage `json:"Value"`
}

type LsiCacheVaultResponseData struct {
	CacheVaultInfo []LsiBatteryProperty `json:"Cachevault_Info"`
	FirmwareStatus []LsiBatteryProperty `json:"Firmware_Status"`
	GasGaugeStatus []LsiBatteryProperty `json:"GasGaugeStatus"`
	DesignInfo     []LsiBatteryProperty `json:"Design_Info"`
}

type CVInfo struct {
	ControllerId int
	State        string
	Subtype      string
	Temperature  int
	SerialNumber string
	NeedReplace  string
}

type BBUInfo struct {
	ControllerId int
	State        string
	Subtype      string
	//CeLsius
	Temperature int
	// microvolts (mV)
	Voltage int
	// amperage (current) mA
	Amperage             int
	ChargingStatus       string
	NeedReplace          string
	RemainingCapacityLow string
	TemperatureState     string
	VoltageState         string
}

type LsiBBUResponseData struct {
	BBUInfo         []LsiBatteryProperty `json:"BBU_Info"`
	BBUStatus       []LsiBatteryProperty `json:"BBU_Firmware_Status"`
	BBUCapacityInfo []LsiBatteryProperty `json:"BBU_Capacity_Info"`
}

type LsiInfo struct {
	StorageControllers *[]StorageControllerInfo `json:"StorageControllers,omitempty"`
	PhysicalDisks      *[]PhysicalDiskInfo      `json:"PhysicalDisks,omitempty"`
	LogicalDisks       *[]LogicalDiskInfo       `json:"LogicalDisks,omitempty"`
	Batteries          *[]LsiCtlBatteryInfo     `json:"Batteries,omitempty"`
	Flags              LsiTemplateFlags
	SmartMap           *map[string]string `json:"SmartMap,omitempty"`
	Errors             *string            `json:"Errors,omitempty"`
}

type LsiPdInfo struct {
	PhysicalDisks *[]PhysicalDiskInfo `json:"PhysicalDisks,omitempty"`
	SmartMap      *map[string]string  `json:"SmartMap,omitempty"`
}

type LsiTemplateFlags struct {
	//получаем ли инфу от call show all по дискам в достаточном длля мониторинга объеме
	PdInfoFromCtl bool
	// с enclosure, без или комбинированный вариант
	PdInfoArgs string
	// а есть ли вообще батарейки, чтобы зазря не делать попыток получить дополнительную инфу о батарейках
	BatteryExists bool
	// аналогично по батарее
	BatteryInfoArgs string
}

type LsiLdInfo struct {
	Type    string `json:"TYPE"`
	State   string `json:"State"`
	Size    string `json:"Size"`
	Name    string `json:"Name"`
	Consist string `json:"Consist"`
	DGVD    string `json:"DG/VD"`
}

type LogicalDiskInfo struct {
	ControllerId int
	Consist      string
	Name         string
	Type         string
	Size         string
	State        string
	DGVD         string
	LdId         string
}

type PhysicalDiskInfo struct {
	Slot            string
	Model           string
	Vendor          string
	Type            string
	Interface       string
	Size            string
	SerialNumber    string
	State           string
	MediaErrCount   uint64
	OtherErrCount   uint64
	PredictErrCount uint64
	SmartFlag       string
	SlotnameSource  string
	Firmware        string
}

type StorageControllerInfo struct {
	Id           int
	Model        string
	SerialNumber string
	Firmware     string
	State        string
}

type BatteryInfo struct {
	BBU []BBUInfo
	CV  []CVInfo
}

type LsiCtlBbuResponseData struct {
	BBUInfo   LsiBBUResponseData
	CtlBasics CtlBasics
}

type LsiBbuResponse struct {
	Controllers []LsiCtlsBbuInfo `json:"Controllers"`
}

type LsiCtlsBbuInfo struct {
	CtlResponseData LsiCtlBbuResponseData `json:"Response Data"`
	CommandStatus   LsiCommandStatus      `json:"Command Status"`
}

type LsiCvResponse struct {
	Controllers []LsiCtlsCvInfo `json:"Controllers"`
}

type LsiCtlsCvInfo struct {
	CommandStatus   LsiCommandStatus          `json:"Command Status"`
	CtlResponseData LsiCacheVaultResponseData `json:"Response Data"`
}

type CtlBasics struct {
	Basics LsiCtlBasics `json:"Basics"`
}
