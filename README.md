# Zabbix Agent2 Plugins

Плагины для zabbiz agent 2

# Как скомпилировать

- Установить go
- Скопировать репозиторий
- выполнить команды ниже, изменив пути на свои

``` 
cd yourpath\zabbix_agent2_plugins
go mod tidy
go build ./cmd/zabbix-agent2-plugin-lsi.go
```

# Как использовать 

- Копируем получившийся бинарь в папку с агентом, например, C:\zabbix_agent\loadableplugins
- В конфиге модуля правим параметр пути к бинарю, если надо, например в файле lsi.conf
`Plugins.Irstinfo.System.Path=C:/zabbix_agent/loadableplugins/zabbix-agent2-plugin-intelrst`
- Кладем конфиг в папку C:\zabbix_agent\zabbix_agent2.d\plugins.d
- Перезапускаем агента
