package lsi

import (
	"encoding/json"
	"fmt"
	"runtime"

	"golang.org/x/exp/slices"

	comf "github.com/mykolq/zabbix_agent2_plugins/commonfunctions"

	"git.zabbix.com/ap/plugin-support/conf"
	"git.zabbix.com/ap/plugin-support/log"
	"git.zabbix.com/ap/plugin-support/plugin"
)

type PluginOptions struct {
	plugin.SystemOptions `conf:"optional,name=System"`
	Timeout              int    `conf:"optional,range=1:30"`
	UtilsPath            string `conf:"optional"`
	UtilDefault          string `conf:"optional"`
}

// стандартный функционал плагинов. либы подключены выше
type Plugin struct {
	plugin.Base
	options PluginOptions
}

var impl Plugin

func SetLogger(logger log.Logger) {
	impl.Logger = logger
}

func (p *Plugin) Configure(global *plugin.GlobalOptions, options interface{}) {
	if err := conf.Unmarshal(options, &p.options); err != nil {
		p.Errf("cannot unmarshal configuration options: %s", err)
	}
	if p.options.Timeout == 0 {
		p.options.Timeout = global.Timeout
	}
	if p.options.UtilsPath == "" {
		p.options.UtilsPath = "C:/zabbix_agent/diskutils/avago"
	}
	if p.options.UtilDefault == "" {
		p.options.UtilDefault = "storcli64"
	}
}

func (p *Plugin) Validate(options interface{}) error {
	var o PluginOptions
	return conf.Unmarshal(options, &o)
}

func (p *Plugin) Export(key string, params []string, ctx plugin.ContextProvider) (result interface{}, err error) {
	p.Infof("received request to handle %s key with %d parameters", key, len(params))

	if len(params) == 0 {
		return nil, fmt.Errorf("This plugin must have one parameter minimum (name of cli)")
	} else {
		timeout := p.options.Timeout

		UtilPath := p.options.UtilsPath
		UtilsList := []string{"storcli64", "perccli64", "storcli"}
		Util := params[0]

		if !(slices.Contains(UtilsList, Util)) {
			Util = p.options.UtilDefault
		}

		UseSudo := false
		OsType := runtime.GOOS

		if OsType != "windows" {
			UseSudo = true
		}

		switch key {

		case "lsi.allinfo":
			lsiInfo, err := GetLsiInfo(fmt.Sprintf("%s/", UtilPath), fmt.Sprintf("%s/%s", UtilPath, Util), Util, timeout, UseSudo)
			if err != nil {
				return nil, err
			}
			lsiinfoJson, err := json.Marshal(lsiInfo)

			if err != nil {
				return nil, err
			}
			return string(lsiinfoJson), nil

		case "lsi.pdsinfo":
			if len(params) < 2 {
				return nil, fmt.Errorf("Key lsipdinfo must use with two parameters: cliname + control flag (enclosure , noenclosure, combined)")
			} else {

				flagsList := []string{"enclosure", "noenclosure", "combined", "stub"}
				flag := params[1]

				if !(slices.Contains(flagsList, flag)) {
					p.Infof("received second parameter is not in known list. Will use combined flag")
					flag = "combined"
				}

				lsiPdInfo, err := GetLsiPdInfo("", "", fmt.Sprintf("%s/", UtilPath), fmt.Sprintf("%s/%s", UtilPath, Util), Util, flag, timeout, UseSudo)
				if err != nil {
					return nil, err
				}

				lsiPdInfoJson, err := json.Marshal(lsiPdInfo)

				if err != nil {
					return nil, err
				}
				return string(lsiPdInfoJson), nil
			}

		case "lsi.batteriesinfo":
			if len(params) < 2 {
				return nil, fmt.Errorf("Key lsipdinfo must use with two parameters: cliname + control flag (cv , bbu, combined)")
			} else {

				flagsList := []string{"cv", "bbu", "combined", "stub"}
				flag := params[1]

				if !(slices.Contains(flagsList, flag)) {
					p.Infof("received second parameter is not in known list. Will use combined flag")
					flag = "combined"
				}

				lsiBatteryInfo, _ := GetLsiBatteryInfo(fmt.Sprintf("%s/", UtilPath), fmt.Sprintf("%s/%s", UtilPath, Util), Util, flag, timeout, UseSudo)
				/*if err != nil {
					return nil, err
				}*/

				lsiBatteryInfoJson, _ := json.Marshal(lsiBatteryInfo)

				/*if err != nil {
					return nil, err[0]
				}*/
				return string(lsiBatteryInfoJson), nil
			}

		case "lsi.clioutput":

			FlagsList := []string{"/call/eall/sall show all j",
				"/call/sall show all j",
				"/call show all j",
				"/call/vall show all j",
				"/call/bbu show all j",
				"/call/cv show all j"}

			if len(params) < 2 {
				errText := fmt.Sprintf(`Key lsipdinfo must use with two parameters: cliname + one of list: %s`,
					comf.SliceToString(FlagsList))
				return nil, fmt.Errorf("%s", errText)
			} else {
				flag := params[1]
				if !(slices.Contains(FlagsList, flag)) {
					errText := fmt.Sprintf("You must use only allowed parameters from list: %s",
						comf.SliceToString(FlagsList))
					return nil, fmt.Errorf("%s", errText)
				}

				CliOutput := PrintCliOutput(fmt.Sprintf("%s/", UtilPath), fmt.Sprintf("%s/%s", UtilPath, Util), Util, timeout, UseSudo, flag)

				return CliOutput, nil
			}

		default:
			return nil, plugin.UnsupportedMetricError

		}
	}
}

func init() {
	plugin.RegisterMetrics(&impl, "lsi", "lsi.allinfo", "All lsi info.",
		"lsi.pdsinfo", "All lsi physical disks info.",
		"lsi.batteriesinfo", "All lsi batteries info.",
		"lsi.clioutput", "Just print cli result.")
}
