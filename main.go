package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"github.com/olekukonko/tablewriter"
)

type SteamJson struct{
	Headers []string `json:"headers"`
	Data [][]float64 `json:"data"`
}


var steamJson SteamJson
func main() {
	// units must be in MPa
	newRankine(0.01, 2)
}

func newRankine(pressureCondenser float64, pressureBoiler float64)  {
	var wp float64 // Work done at the pump
	var qh float64 // Heat added at the boiler
	var wt float64 // Work done at the turbine
	var x4 float64 // dryness fraction of steam when it leaves the boiler

	// At the pump
	v := getSteamPropsByPressure_sat(pressureCondenser)[13]
	wp = (v * (pressureBoiler - pressureCondenser)) * 1000
	h1 := getSteamPropsByPressure_sat(pressureCondenser)[7]
	h2 := h1 + wp
	
	// Boiler
	h3 := getSteamPropsByPressure_sat(pressureBoiler)[8]
	qh = h3 - h2

	// Turbine
	x4 = (getSteamPropsByPressure_sat(pressureBoiler)[11] - getSteamPropsByPressure_sat(pressureCondenser)[10]) / getSteamPropsByPressure_sat(pressureCondenser)[12]
	h4 := h1  + x4 * getSteamPropsByPressure_sat(pressureCondenser)[9]

	wt = h3 - h4

	// Condenser
	ql := h4 - h1

	// Efficiency of cycle
	nth := ((wt - wp)/ qh) * 100

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Process", "Inlet State", "Exit State", "Work done/Heat transfer"})
	data := [][]string{
		{"Pump", "P1 = " + strconv.FormatFloat(pressureCondenser, 'g', 5, 64) + "MPa", "h2 = " + strconv.FormatFloat(h2, 'g', 5, 64) + "MPa", "Work Input: " + strconv.FormatFloat(wp, 'g', 5, 64) + "KJ/KG"},
		{"Boiler", "P2 = " + strconv.FormatFloat(pressureBoiler, 'g', 5, 64) + "MPa", "No output", "Heat Added : " + strconv.FormatFloat(qh, 'g', 5, 64) + "KJ/KG"},
		{"Turbine", "State 3 known", "x4 = " + strconv.FormatFloat(x4, 'g', 5, 64) + "MPa " + "h4 = " + strconv.FormatFloat(h4, 'g', 5, 64) , "Work output : " + strconv.FormatFloat(wt, 'g', 5, 64) + "KJ/KG"},
		{"Condenser", "State 4 known", "Output known", "Heat rejected : " + strconv.FormatFloat(ql, 'g', 5, 64) + "KJ/KG"},
	}

	for _, v := range data {
		table.Append(v)
	}
	table.SetFooter([]string{"", "", "Efficiency of Cycle", strconv.FormatFloat(nth, 'g', 5, 64) + "%"})

	table.Render()

}

func getSteamPropsByPressure_sat(pressureValue float64) []float64 {
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