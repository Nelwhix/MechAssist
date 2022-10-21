package main

import (
	"os"
	"encoding/json"
	"io"
	"log"
)

func GetSteamPropsByPressure_sat(pressureValue float64) []float64 {
	file, err := os.Open("./tables/saturated_by_pressure.json")

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)

	json.Unmarshal(byteValue, &steamJson)

	for i := 0; i < len(steamJson.Data); i++ {
		if pressureValue == steamJson.Data[i][0] {
			return steamJson.Data[i]
		}
	}

	return nil
}

func GetSteamPropsByTemp(tempValue float64) []float64 {
	file, err := os.Open("./tables/saturated_by_temperature.json")

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)

	json.Unmarshal(byteValue, &steamJson)

	for i := 0; i < len(steamJson.Data); i++ {
		if tempValue == steamJson.Data[i][0] {
			return steamJson.Data[i]
		}
	}

	return nil
}

func GetSteamProps_sup(pressureValue float64, tempValue float64) []float64 {
	file, err := os.Open("./tables/compressed_liquid_and_superheated_steam.json")

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	json.Unmarshal(byteValue, &steamJson)

	for i := 0; i < len(steamJson.Data); i++ {
		if pressureValue == steamJson.Data[i][0] && tempValue == steamJson.Data[i][1] {
			return steamJson.Data[i]
		}
	}
	return nil
}