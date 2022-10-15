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

const (
      Internal_Energy_f = 4
      Internal_Energy_fg = 6
      Enthalpy_f = 7
      Enthalpy_g = 8 
      Enthalpy_fg = 9
      Entropy_f = 10 
      Entropy_g = 11 
      Entropy_fg = 12 
      Specific_Volume_f = 13
      Specific_Volume_g = 14
	  Internal_energy_sup = 4
	  Enthalpy_sup = 5
	  Entropy_sup = 6
)


var steamJson SteamJson
func main() {
	// units must be in MPa
	newRankine_sup(0.01, 4, 400)
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

func getSteamProps_sup(pressureValue float64, tempValue float64) []float64 {
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

// If saturated vapor entered the turbine
func newRankine_sat(pressureCondenser float64, pressureBoiler float64)  {
	var qh float64 // Heat added at the boiler
	var wt float64 // Work done at the turbine
	var x4 float64 // dryness fraction of steam when it leaves the boiler

	// At the pump
	h1 := pumpProcess(pressureCondenser, pressureBoiler)[0]
	h2 := pumpProcess(pressureCondenser, pressureBoiler)[1]
	wp := pumpProcess(pressureCondenser, pressureBoiler)[2]
	
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
		{"Boiler", "P2 = " + strconv.FormatFloat(pressureBoiler, 'g', 5, 64) + "MPa", "h3 = " +strconv.FormatFloat(h3, 'g', 5, 64) + "MPa", "Heat Added : " + strconv.FormatFloat(qh, 'g', 5, 64) + "KJ/KG"},
		{"Turbine", "State 3 known", "x4 = " + strconv.FormatFloat(x4, 'g', 5, 64) + "MPa " + "h4 = " + strconv.FormatFloat(h4, 'g', 5, 64) , "Work output : " + strconv.FormatFloat(wt, 'g', 5, 64) + "KJ/KG"},
		{"Condenser", "State 4 known", "Output known", "Heat rejected : " + strconv.FormatFloat(ql, 'g', 5, 64) + "KJ/KG"},
	}

	for _, v := range data {
		table.Append(v)
	}
	table.SetFooter([]string{"", "", "Efficiency of Cycle", strconv.FormatFloat(nth, 'g', 5, 64) + "%"})

	table.Render()

}

//if superheated steam entered the turbine
func newRankine_sup(pressureCondenser float64, pressureBoiler float64, tempAfterBoiler float64) {	
	var qh float64 // Heat added at the boiler
	var wt float64 // Work done at the turbine
	var x4 float64 // dryness fraction of steam when it leaves the boiler

	// Pump
	h2 := pumpProcess(pressureCondenser, pressureBoiler)[1]
	wp := pumpProcess(pressureCondenser, pressureBoiler)[2]
	
	// Turbine
	inletSteamProps := getSteamProps_sup(4, 400)
	exitSteamProps := getSteamPropsByPressure_sat(0.01)
	h3 := inletSteamProps[Enthalpy_sup]
	s3 := inletSteamProps[Entropy_sup]
	sf4 := exitSteamProps[Entropy_f]
	sfg4 := exitSteamProps[Entropy_fg]
	x4 = (s3 - sf4) / sfg4

	h4 := exitSteamProps[Enthalpy_f] + (x4 * exitSteamProps[Enthalpy_fg])
	wt = h3 - h4

	// Boiler 
	qh = h3 - h2

	nth := ((wt - wp) / qh) * 100

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Process", "Inlet State", "Exit State", "Work done/Heat transfer"})
	data := [][]string{
		{"Pump", "P1 = " + strconv.FormatFloat(pressureCondenser, 'g', 5, 64) + "MPa", "h2 = " + strconv.FormatFloat(h2, 'g', 5, 64) + "MPa", "Work Input: " + strconv.FormatFloat(wp, 'g', 5, 64) + "KJ/KG"},
		{"Turbine", "P3 = " + strconv.FormatFloat(pressureBoiler, 'g', 5, 64) + "MPa " + "T3 = " + strconv.FormatFloat(tempAfterBoiler, 'g', 5, 64) + "C", "x4 = " + strconv.FormatFloat(x4, 'g', 5, 64) + " h4 = " + strconv.FormatFloat(h4, 'g', 5, 64) + "MPa" , "Work output : " + strconv.FormatFloat(wt, 'g', 5, 64) + "KJ/KG"},
		{"Boiler", "P2 = " + strconv.FormatFloat(pressureBoiler, 'g', 5, 64) + "MPa", "h3 = " +strconv.FormatFloat(h3, 'g', 5, 64) + "MPa", "Heat Added : " + strconv.FormatFloat(qh, 'g', 5, 64) + "KJ/KG"},
	}

	table.AppendBulk(data)
	table.SetFooter([]string{"", "", "Efficiency of Cycle", strconv.FormatFloat(nth, 'g', 5, 64) + "%"})

	table.Render()
}

func pumpProcess(pressureCondenser float64, pressureBoiler float64) []float64 {
	var wp float64 // Work done at the pump
	v := getSteamPropsByPressure_sat(pressureCondenser)[13]
	wp = (v * (pressureBoiler - pressureCondenser)) * 1000
	h1 := getSteamPropsByPressure_sat(pressureCondenser)[7]
	h2 := h1 + wp

	return []float64{h1, h2, wp}
}

func newReheat(pressureBoiler float64, tempAfterBoiler float64, pressureAtReheat float64, pressureAtLowTurbine float64) {
	var inletSteamProps []float64
	var exitSteamProps []float64

	// high pressure turbine
	inletSteamProps = getSteamProps_sup(pressureBoiler, tempAfterBoiler)
	exitSteamProps = getSteamPropsByPressure_sat(pressureAtReheat)
	h3 := inletSteamProps[Enthalpy_sup]
	s3 := inletSteamProps[Entropy_sup]
	sf4 := exitSteamProps[Entropy_f]
	sfg4 := exitSteamProps[Entropy_fg]

	x4 := (s3 - sf4) / sfg4
	h4 := exitSteamProps[Enthalpy_f] + x4 * exitSteamProps[Enthalpy_fg]

	whp := h3 - h4 // work done at high pressure turbine

	// Low Pressure turbine
	inletSteamProps = getSteamProps_sup(pressureAtReheat, tempAfterBoiler)
	h5 := inletSteamProps[Enthalpy_sup]
	s5 := inletSteamProps[Entropy_sup]




}