package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type SteamJson struct{
	Headers []string `json:"headers"`
	Data [][]float64 `json:"data"`
}


var steamJson SteamJson
func main() {
	file, err := os.Open("./tables/saturated_by_pressure.json")

	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()
	byteValue, _ := io.ReadAll(file)

	json.Unmarshal(byteValue, &steamJson)

	fmt.Println(steamJson.Headers)
}