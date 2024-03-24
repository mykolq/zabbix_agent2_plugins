# Плагин для шаблона дисковых контроллеров LSI стека

## Параметры:

Метрика называется lsi и принимает следующие аргументы: 

- тип утилиты: perccli64 или storcli64
- callshowall: отдает всю доступную инфу, выводимую по команде ```storcli64 /call show all``` или ```perccli64 /call show all```
- calldisks: нужна если в выхлопе предыдущей команде недостаточно инфы для обнаружения ии мониторинга дисков (половина всех случаев). В случае этого аргумекнта должен быть
еще один аргумент, который покажет, нужно ли использовать что-то одно (```/call/eall/sall``` или ```/call/sall```), или использовать оба варианта чтобы отдать инфу по дискам.
- callbattery: нужна для вывода дополнительной инфы о батарее. тут может быть либо ```/call/bbu``` либо ```/call/cv show all```. В теории могут быть оба варианта, поэтому сделаем как с дисками выше.

## Структура вывода:

- с аргументом callshowall: 
```
type LsiInfo struct {
	StorageControllers *[]StorageControllerInfo `json:"StorageControllers,omitempty"`
	PhysicalDisks      *[]PhysicalDiskInfo      `json:"PhysicalDisks,omitempty"`
	LogicalDisks       *[]LogicalDiskInfo       `json:"LogicalDisks,omitempty"`
	Batteries          *[]LsiCtlBatteryInfo     `json:"Batteries,omitempty"`
	Flags              LsiTemplateFlags
	SmartMap           *map[string]string `json:"SmartMap,omitempty"`
	Errors             *string            `json:"Errors,omitempty"`
}
```
структура флагов
```
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
```
Флаги прописываются в соотвествующие макросы в шаблоне.

- с аргументом calldisks (если недстаточно инфы о дисках в выхлопе выше), дополнительно передаем enclosure, noenclosure, combined. enclosure = ```/call/eall/sall show all```,
noenclosure = ```/call/eall/sall show all```, combined по очереди запускает оба и добавляет в нужную структуру после обработки

```
type lsiPDInfo struct {
    PhysicalDisks      []PhysicalDiskInfo
}
```
- с аргументом callbattery, дополнительно передаем cv, bbu, combined.

```
type lsiBatteryInfo struct {
	CacheVaults        []CVInfo
	BBUs               []BBUInfo
}
```