package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type SteamJson struct{
	headers []string
	data [][]float64
}

var steamJson *SteamJson
func main() {
	file, err := os.Open("./tables/saturated_by_pressure_V1.3.json")

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	json.Unmarshal(byteValue, &steamJson)
	
}

// func getPropsByPressure_sat(pressureValue float64) int {
// 	file, err := os.Open("./tables/saturated_by_pressure_V1.3.json")

// 	if err != nil {
// 		log.Fatalf("Error opening file: %v", err)
// 	}
// 	defer file.Close()
// 	byteValue, _ := io.ReadAll(file)
// 	json.Unmarshal(byteValue, &steamJson)

// 	for i := 0; i < len(steamJson.data); i++ {
// 		for j := 0; j < len(steamJson.headers); j++ {
// 			if (steamJson.data[i])[j] == float32(pressureValue) {
// 				return len(steamJson.data)
// 			}
// 		}
// 	}
// 	return 1
// }