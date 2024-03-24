package main

import (
	"fmt"

	"git.zabbix.com/ap/plugin-support/plugin/container"
	"github.com/mykolq/zabbix_agent2_plugins/zabbix_agent2_plugins/lsi"
)

func main() {
	h, err := container.NewHandler("lsi")

	if err != nil {
		panic(fmt.Sprintf("failed to create plugin handler %s", err.Error()))
	}

	lsi.SetLogger(&h)

	err = h.Execute()
	if err != nil {
		panic(fmt.Sprintf("failed to execute plugin handler %s", err.Error()))
	}
}
