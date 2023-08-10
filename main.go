package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	pc "steam-price-checker/steam-price-checker/pricechecker"
	sw "steam-price-checker/steam-price-checker/sheetswriter"
)

func main() {
	SetEnvironmentVariables()
	checker := pc.PriceChecker{}
	checker.SetItemsToCheck()
	checker.SetPrices()
	writer := sw.SheetsWriter{}
	writer.Init(&checker)
	writer.InsertColumn(1)
	writer.WriteData()
}

func SetEnvironmentVariables() {
	file, err := os.Open(".env")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	environmentValues := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		envValue := strings.Split(line, "=")
		key := strings.TrimSpace(envValue[0])
		if key == "SHEETSKEY" {
			value := fmt.Sprintf("%v==", strings.TrimSpace(envValue[1]))
			environmentValues[key] = value
		} else {
			value := strings.TrimSpace(envValue[1])
			environmentValues[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	for key, value := range environmentValues {
		os.Setenv(key, value)
	}
}
