package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	pc "steam-price-checker/steam-price-checker/pricechecker"
)

func main() {
	SetEnvironmentVariables()
	checker := pc.PriceChecker{Test: "aha"}
	checker.SetItemsToCheck()
	checker.GetAutographsPrices()
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
		value := strings.TrimSpace(envValue[1])
		environmentValues[key] = value
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	for key, value := range environmentValues {
		os.Setenv(key, value)
	}
}
