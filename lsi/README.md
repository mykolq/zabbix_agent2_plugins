# Plugin for zabbix agent 2 for lsi stack

## Install

Example of installing here (hust chaange you pathes, binary file for lsi has suffix lsi)

https://github.com/mykolq/zabbix_agent2_plugins 

Example of sudoers file for linux for this plugin:
```zabbix ALL= (ALL) NOPASSWD: /etc/zabbix/diskutils/avago/storcli64 /call show all j,/etc/zabbix/diskutils/avago/storcli64 /call/vall show all j,/etc/zabbix/diskutils/avago/storcli64 /call/eall/sall show all j,/etc/zabbix/diskutils/avago/storcli64 show all j,/etc/zabbix/diskutils/avago/storcli64 /call/cv show all j,/etc/zabbix/diskutils/avago/storcli64 /call/bbu show all j```

You have to create file like /etc/sudoers.d/sudoers_zabbix_lsiplugin and put it there.

After inslall pluging you can test it from zabbix server like this:
```zabbix_get -s yourservername -k 'lsi.allinfo[storcli64]'```
If all is ok you will see result like that:
```{"StorageControllers":[{"Id":0,"Model":"PERC H330 Adapter","SerialNumber":"9AP0551","Firmware":"4.300.01-8353","State":"Optimal"}],"LogicalDisks":[],"Batteries":[],"Flags":{"PdInfoFromCtl":false,"PdInfoArgs":"enclosure","BatteryExists":false,"BatteryInfoArgs":"stub"}}```
In this result you cannot see physical disks info cause storcli64 doesnt get enough info in this case. You can get physical disks info like that:
```zabbix_psk_get -s yourservername -k 'lsi.pdsinfo[storcli64, enclosure]'``` second parameter was has got from flags key above.
Result will be like that:
```{"PhysicalDisks":[{"Slot":"/c0/e32/s0","Model":"MTFDDAK240TCB","Vendor":"Micron","Type":"SSD","Interface":"SATA","Size":"223.570 GB","SerialNumber":"194324BFA1CE","State":"JBOD","MediaErrCount":0,"OtherErrCount":142043,"PredictErrCount":0,"SmartFlag":"No","SlotnameSource":"LSI","Firmware":" D0DE012"},{"Slot":"/c0/e32/s1","Model":"MG06ACA800EY","Vendor":"Toshiba","Type":"HDD","Interface":"SATA","Size":"7.277 TB","SerialNumber":"X9T0A0UMF1QF","State":"JBOD","MediaErrCount":0,"OtherErrCount":142043,"PredictErrCount":0,"SmartFlag":"No","SlotnameSource":"LSI","Firmware":"    GA06"},{"Slot":"/c0/e32/s2","Model":"MG06ACA800EY","Vendor":"Toshiba","Type":"HDD","Interface":"SATA","Size":"7.277 TB","SerialNumber":"X9T0A0TEF1QF","State":"JBOD","MediaErrCount":0,"OtherErrCount":142043,"PredictErrCount":0,"SmartFlag":"No","SlotnameSource":"LSI","Firmware":"    GA06"},{"Slot":"/c0/e32/s3","Model":"MG06ACA800EY","Vendor":"Toshiba","Type":"HDD","Interface":"SATA","Size":"7.277 TB","SerialNumber":"X9T0A0TZF1QF","State":"JBOD","MediaErrCount":0,"OtherErrCount":142043,"PredictErrCount":0,"SmartFlag":"No","SlotnameSource":"LSI","Firmware":"    GA06"},{"Slot":"/c0/e32/s4","Model":"MG06ACA800EY","Vendor":"Toshiba","Type":"HDD","Interface":"SATA","Size":"7.277 TB","SerialNumber":"X9U0A0N5F1QF","State":"JBOD","MediaErrCount":0,"OtherErrCount":142043,"PredictErrCount":0,"SmartFlag":"No","SlotnameSource":"LSI","Firmware":"    GA06"},{"Slot":"/c0/e32/s5","Model":"MG06ACA800EY","Vendor":"Toshiba","Type":"HDD","Interface":"SATA","Size":"7.277 TB","SerialNumber":"X9U0A0FDF1QF","State":"JBOD","MediaErrCount":0,"OtherErrCount":142043,"PredictErrCount":0,"SmartFlag":"No","SlotnameSource":"LSI","Firmware":"    GA06"},{"Slot":"/c0/e32/s6","Model":"MG06ACA800EY","Vendor":"Toshiba","Type":"HDD","Interface":"SATA","Size":"7.277 TB","SerialNumber":"X9U0A0JNF1QF","State":"JBOD","MediaErrCount":0,"OtherErrCount":142043,"PredictErrCount":0,"SmartFlag":"No","SlotnameSource":"LSI","Firmware":"    GA06"},{"Slot":"/c0/e32/s7","Model":"MG06ACA800EY","Vendor":"Toshiba","Type":"HDD","Interface":"SATA","Size":"7.277 TB","SerialNumber":"X9U0A0EMF1QF","State":"JBOD","MediaErrCount":0,"OtherErrCount":142043,"PredictErrCount":0,"SmartFlag":"No","SlotnameSource":"LSI","Firmware":"    GA06"}],"SmartMap":{"194324BFA1CE":"/c0/e32/s0","X9T0A0TEF1QF":"/c0/e32/s2","X9T0A0TZF1QF":"/c0/e32/s3","X9T0A0UMF1QF":"/c0/e32/s1","X9U0A0EMF1QF":"/c0/e32/s7","X9U0A0FDF1QF":"/c0/e32/s5","X9U0A0JNF1QF":"/c0/e32/s6","X9U0A0N5F1QF":"/c0/e32/s4"}}```
