package commonfunctions

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// We have all utilies executing with one timeot before cancelling it
func ExecWithContextTimeout(CliDirPath, CliFullPath string, Timeout int, CliParams []string, UseSudo bool) (result []byte, err error) {

	// Create a new context and add a timeout to it
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Timeout)*time.Second)
	defer cancel()

	if UseSudo {
		CliParams = append([]string{CliFullPath}, CliParams...)
		CliFullPath = "/usr/bin/sudo"
		CliDirPath = "/usr/bin/"
	}

	cmd := exec.CommandContext(ctx, CliFullPath, CliParams...)
	cmd.Dir = CliDirPath
	result, err = cmd.CombinedOutput()

	if err != nil {
		err = fmt.Errorf("util: %s, params: %v\n%v", CliFullPath, CliParams, err)
	}

	return
}

func StrToInt(str string) (integer int) {
	integer, err := strconv.Atoi(str)
	if err != nil {
		return 0
	} else {
		return integer
	}
}

func ReplaceUnits(str string) (changedstr string) {
	replaceRe := regexp.MustCompile(`(.*?)\s*(C)`)

	changedstr = replaceRe.ReplaceAllString(str, "$1")
	return
}

func FixPhysicalDiskInfo(model, vendor, ataProductId string, modelfamily *string) (modelFixed, vendorFixed string) {

	WDModelRegex := regexp.MustCompile(`^WUS.*`)
	SamsungModelRegex := regexp.MustCompile(`^MZ7LH.*`)
	HgstModelRegex := regexp.MustCompile(`^HUS.*`)
	SeagateModelRegex := regexp.MustCompile(`^ST\d{2,}.*`)
	ToshibaModelRegex := regexp.MustCompile(`^MG0.*`)
	IntelSsdModelRegex := regexp.MustCompile(`^SSDS.*`)
	MicronSsdModelRegex := regexp.MustCompile(`MTF\D+\d+\D+|.*Micron.*`)
	DellToshibaModelRegex := regexp.MustCompile(`^THN.*`)

	if (vendor == "" || vendor == "Unknown") && ataProductId != "" {
		vendor = strings.TrimSuffix(ataProductId, "(tm)")
	}

	if strings.Contains(vendor, "DELL") {
		if DellToshibaModelRegex.MatchString(model) {
			vendor = "Toshiba"
		}
		if SamsungModelRegex.MatchString(model) {
			vendor = "Samsung"
		}
	}

	if vendor == "" || vendor == "ATA" {

		if WDModelRegex.MatchString(model) {
			vendor = "Western Digital"
		}
		if SamsungModelRegex.MatchString(model) {
			vendor = "Samsung"
		}
		if HgstModelRegex.MatchString(model) {
			vendor = "HGST"
		}
		if SeagateModelRegex.MatchString(model) {
			vendor = "Seagate"
		}
		if MicronSsdModelRegex.MatchString(model) {
			vendor = "Micron"
			if strings.Contains(model, "_") {
				model = strings.TrimSpace(strings.TrimPrefix(strings.ReplaceAll(model, "_", " "), "Micron"))
			}
		}
		if IntelSsdModelRegex.MatchString(model) {
			vendor = "Intel"
		}

		if strings.Contains(model, " ") {
			vendor = strings.SplitN(model, " ", 2)[0]
			model = strings.SplitN(model, " ", 2)[1]
		}
	}

	if modelfamily != nil {
		if (vendor == "" || vendor == "Unknown" || vendor == "ATA") && strings.Contains(*modelfamily, " ") {
			vendor = strings.SplitN(*modelfamily, " ", 2)[0]
		}
	}

	if ToshibaModelRegex.MatchString(model) {
		vendor = "Toshiba"
	}
	if IntelSsdModelRegex.MatchString(model) {
		vendor = "Intel"
	}
	if HgstModelRegex.MatchString(model) {
		vendor = "HGST"
	}
	if SeagateModelRegex.MatchString(model) {
		vendor = "Seagate"
	}
	if SamsungModelRegex.MatchString(model) {
		vendor = "Samsung"
	}
	if WDModelRegex.MatchString(model) {
		vendor = "Western Digital"
	}

	if ToshibaModelRegex.MatchString(vendor) {
		model = vendor
		vendor = "Toshiba"
	}
	if IntelSsdModelRegex.MatchString(vendor) {
		model = vendor
		vendor = "Intel"
	}
	if HgstModelRegex.MatchString(vendor) {
		model = vendor
		vendor = "HGST"
	}
	if SeagateModelRegex.MatchString(vendor) {
		model = vendor
		vendor = "Seagate"
	}
	if SamsungModelRegex.MatchString(vendor) {
		model = vendor
		vendor = "Samsung"
	}
	if WDModelRegex.MatchString(vendor) {
		model = vendor
		vendor = "Western Digital"
	}

	if strings.Contains(vendor, "Crucial") && strings.Contains(model, "Crucial_") {
		model = strings.ReplaceAll(model, "Crucial_", "")
	}

	// если после всех экзекуций у нас в имени модели есть вендор, дропаем его
	if vendor != "" && strings.Contains(model, " ") {
		model = strings.SplitN(model, " ", 2)[1]
	}

	return strings.TrimSpace(model), strings.TrimSpace(vendor)
}

func SliceToString(inputSlice []string) (outString string) {

	endOfSlice := len(inputSlice) - 1

	for i, str := range inputSlice {

		if i == 0 {
			outString += "\"" + str + "\""
		}

		if i > 0 && i < endOfSlice {
			outString += ", \"" + str + "\""
		}

		if i == endOfSlice {
			outString += " or \"" + str + "\""
		}
	}

	return

}
