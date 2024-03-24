# Zabbix Agent2 Plugins

# How to compile

- Install go
- Copy repo
- Execute command below, changing pathes

``` 
cd yourpath\zabbix_agent2_plugins
go mod tidy
go build ./cmd/zabbix-agent2-plugin-lsi.go
```

# How to use

- Copy binary file (for windows use exe extension) to zabbix agent2 path, for example, C:\zabbix_agent\loadableplugins
- In module config (lsi_win.conf or other), if needed, changing path to module binary file
`Plugins.lsi.System.Path=C:/zabbix_agent/loadableplugins/zabbix-agent2-plugin-lsi`
- Put module conf file (lsi_win.conf or other) to agent2 folder, for example, C:\zabbix_agent\zabbix_agent2.d\plugins.d
- Restart agent2

Don't forget to change sudoers file for using with linux
You can use compiled binary from releases

