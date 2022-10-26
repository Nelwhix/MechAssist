package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	//"fmt"
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
	// Pressure units must be in MPa, Temp units must be in C
	fmt.Println(GetSteamPropsByPressure_sat(0.0125))
	fmt.Println("Kj/Kg")
}


// If saturated vapor entered   the turbine
func newRankine_sat(pressureCondenser float64, pressureBoiler float64)  {
	var qh float64 // Heat added at the boiler
	var wt float64 // Work done at the turbine
	var x4 float64 // dryness fraction of steam when it leaves the boiler

	// At the pump
	h1 := pumpProcess(pressureCondenser, pressureBoiler)[0]
	h2 := pumpProcess(pressureCondenser, pressureBoiler)[1]
	wp := pumpProcess(pressureCondenser, pressureBoiler)[2]
	
	// Boiler
	h3 := GetSteamPropsByPressure_sat(pressureBoiler)[8]
	qh = h3 - h2

	// Turbine
	x4 = (GetSteamPropsByPressure_sat(pressureBoiler)[11] - GetSteamPropsByPressure_sat(pressureCondenser)[10]) / GetSteamPropsByPressure_sat(pressureCondenser)[12]
	h4 := h1  + x4 * GetSteamPropsByPressure_sat(pressureCondenser)[9]

	wt = h3 - h4

	// Condenser
	ql := h4 - h1

	// Efficiency of cycle
	nth := ((wt - wp)/ qh) * 100

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Process", "Inlet State", "Exit State", "Work done/Heat transfer"})
	data := [][]string{
		{"Pump", "P1 = " + strconv.FormatFloat(pressureCondenser, 'g', 5, 64) + "MPa", "h2 = " + strconv.FormatFloat(h2, 'g', 5, 64) + "KJ/KG", "Work Input: " + strconv.FormatFloat(wp, 'g', 5, 64) + "KJ/KG"},
		{"Boiler", "P2 = " + strconv.FormatFloat(pressureBoiler, 'g', 5, 64) + "MPa", "h3 = " +strconv.FormatFloat(h3, 'g', 5, 64) + "KJ/KG", "Heat Added : " + strconv.FormatFloat(qh, 'g', 5, 64) + "KJ/KG"},
		{"Turbine", "State 3 known", "x4 = " + strconv.FormatFloat(x4, 'g', 5, 64) + " h4 = " + strconv.FormatFloat(h4, 'g', 5, 64) + "KJ/KG", "Work output : " + strconv.FormatFloat(wt, 'g', 5, 64) + "KJ/KG"},
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
	inletSteamProps := GetSteamProps_sup(pressureBoiler, tempAfterBoiler)
	exitSteamProps := GetSteamPropsByPressure_sat(pressureCondenser)
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
		{"Pump", "P1 = " + strconv.FormatFloat(pressureCondenser, 'g', 5, 64) + "MPa", "h2 = " + strconv.FormatFloat(h2, 'g', 5, 64) + "KJ/KG", "Work Input: " + strconv.FormatFloat(wp, 'g', 5, 64) + "KJ/KG"},
		{"Turbine", "P3 = " + strconv.FormatFloat(pressureBoiler, 'g', 5, 64) + "MPa " + "T3 = " + strconv.FormatFloat(tempAfterBoiler, 'g', 5, 64) + "C", "x4 = " + strconv.FormatFloat(x4, 'g', 5, 64) + " h4 = " + strconv.FormatFloat(h4, 'g', 5, 64) + "KJ/KG" , "Work output : " + strconv.FormatFloat(wt, 'g', 5, 64) + "KJ/KG"},
		{"Boiler", "P2 = " + strconv.FormatFloat(pressureBoiler, 'g', 5, 64) + "MPa", "h3 = " +strconv.FormatFloat(h3, 'g', 5, 64) + "KJ/KG", "Heat Added : " + strconv.FormatFloat(qh, 'g', 5, 64) + "KJ/KG"},
	}

	table.AppendBulk(data)
	table.SetFooter([]string{"", "", "Efficiency of Cycle", strconv.FormatFloat(nth, 'g', 5, 64) + "%"})

	table.Render()
}

func pumpProcess(inletPressure float64, exitPressure float64) []float64 {
	var wp float64 // Work done at the pump
	steamProps := GetSteamPropsByPressure_sat(inletPressure)
	v := steamProps[Specific_Volume_f]
	wp = (v * (exitPressure - inletPressure)) * 1000
	h1 := steamProps[Enthalpy_f]
	h2 := h1 + wp

	return []float64{h1, h2, wp}
}

func newReheat(pressureBoiler float64, tempAfterBoiler float64, pressureAtReheat float64, pressureAfterLowTurbine float64) {
	var inletSteamProps []float64
	var exitSteamProps []float64

	// high pressure turbine
	inletSteamProps = GetSteamProps_sup(pressureBoiler, tempAfterBoiler)
	exitSteamProps = GetSteamPropsByPressure_sat(pressureAtReheat)
	h3 := inletSteamProps[Enthalpy_sup]
	s3 := inletSteamProps[Entropy_sup]
	sf4 := exitSteamProps[Entropy_f]
	sfg4 := exitSteamProps[Entropy_fg]

	x4 := (s3 - sf4) / sfg4
	h4 := exitSteamProps[Enthalpy_f] + x4 * exitSteamProps[Enthalpy_fg]

	whp := h3 - h4 // work done at high pressure turbine

	// Low Pressure turbine
	inletSteamProps = GetSteamProps_sup(pressureAtReheat, tempAfterBoiler)
	exitSteamProps = GetSteamPropsByPressure_sat(0.01)
	h5 := inletSteamProps[Enthalpy_sup]
	s5 := inletSteamProps[Entropy_sup]
	sf6 := exitSteamProps[Entropy_f]
	sfg6 := exitSteamProps[Entropy_fg]

	x6 := (s5 - sf6) / sfg6

	h6 := exitSteamProps[Enthalpy_f] + x6 * exitSteamProps[Enthalpy_fg]

	wt := whp + (h5 - h6)

	// Pump
	result := pumpProcess(pressureAfterLowTurbine, pressureBoiler)	
	h2 := result[1]
	wp := result[2]

	// Boiler
	qh := (h3 - h2) + (h5 - h4)

	nth := ((wt - wp) / qh) * 100

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Process", "Inlet State", "Exit State", "Work done/Heat transfer"})
	data := [][]string{
		{"High Pressure Turbine", "P3 = " + strconv.FormatFloat(pressureBoiler, 'g', 5, 64) + "MPa" + " T3 = " + strconv.FormatFloat(tempAfterBoiler, 'g', 5, 64) + "C", "h4 = " + strconv.FormatFloat(h4, 'g', 5, 64) + "MPa" + " x4 = " + strconv.FormatFloat(x4, 'g', 5, 64), "Work Output: " + strconv.FormatFloat(whp, 'g', 5, 64) + "KJ/KG"},
		{"Low Pressure Turbine", "P5 = " + strconv.FormatFloat(pressureAtReheat, 'g', 5, 64) + "MPa" + " T5 = " + strconv.FormatFloat(tempAfterBoiler, 'g', 5, 64) + "C", "h6 = " + strconv.FormatFloat(h6, 'g', 5, 64) + "MPa" + " x4 = " + strconv.FormatFloat(x6, 'g', 5, 64), "Work Output: " + strconv.FormatFloat((h5 - h6), 'g', 5, 64) + "KJ/KG"},
		{"Pump", "P1 = " + strconv.FormatFloat(pressureAfterLowTurbine, 'g', 5, 64) + "MPa", "h2 = " + strconv.FormatFloat(h2, 'g', 5, 64) + "MPa", "Work Input: " + strconv.FormatFloat(wp, 'g', 5, 64) + "KJ/KG"},
		{"Boiler", "P2 = " + strconv.FormatFloat(pressureBoiler, 'g', 5, 64) + "MPa " + "P4 = " + strconv.FormatFloat(pressureAtReheat, 'g', 5, 64) + "MPa", "Output from earlier steps", "Heat Added : " + strconv.FormatFloat(qh, 'g', 5, 64) + "KJ/KG"},
	}

	table.AppendBulk(data)
	table.SetFooter([]string{"", "", "Efficiency of Cycle", strconv.FormatFloat(nth, 'g', 5, 64) + "%"})

	table.Render()
}

func newRegenerative(pressureBoiler float64, tempAfterBoiler float64, pressureAtFeedwater float64, pressureAtCondenser float64) {
	// Low Pressure Pump
	result := pumpProcess(pressureAtCondenser, pressureAtFeedwater)
	h2 := result[1]
	wp1 := result[2]

	// Turbine
	inletTurbProps := GetSteamProps_sup(pressureBoiler, tempAfterBoiler)
	exitToFeed := GetSteamPropsByPressure_sat(pressureAtFeedwater)
	exitToCondenser := GetSteamPropsByPressure_sat(pressureAtCondenser)
	h5 := inletTurbProps[Enthalpy_sup]
	s5 := inletTurbProps[Entropy_sup]
	sf6 := exitToFeed[Entropy_f]
	sfg6 := exitToFeed[Entropy_fg]
	x6 := (s5 - sf6) / sfg6

	h6 := exitToFeed[Enthalpy_f] + (x6 * exitToFeed[Enthalpy_fg])

	sf7 := exitToCondenser[Entropy_f]
	sfg7 := exitToCondenser[Entropy_fg]
	x7 := (s5 - sf7) / sfg7

	h7 := exitToCondenser[Enthalpy_f] + (x7 * exitToCondenser[Enthalpy_fg])

	// Feedwater heater: y is the extraction fraction
	h3 := exitToFeed[Enthalpy_f]
	y := (h3 - h2) / (h6 - h2)
	wt := (h5 - h6) + ((1 - y) * (h6 - h7))

	// High Pressure Pump
	result2 := pumpProcess(0.4, 4)
	h4 := result2[1]
	wp2 := result2[2]

	// Boiler 
	qh := h5 - h4
	nth := ((wt - wp1 - wp2) / qh) * 100

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Process", "Inlet State", "Exit State", "Work done/Heat transfer"})
	data := [][]string{
		{"Low Pressure Pump", "P1 = " + strconv.FormatFloat(pressureAtCondenser, 'g', 5, 64) + "MPa", "h2 = " + strconv.FormatFloat(h2, 'g', 5, 64) + "MPa", "Work Input: " + strconv.FormatFloat(wp1, 'g', 5, 64) + "KJ/KG"},
		{"Turbine", "P5 = " + strconv.FormatFloat(pressureBoiler, 'g', 5, 64) + "MPa" + " T5 = " + strconv.FormatFloat(tempAfterBoiler, 'g', 5, 64) + "C", "h6 = " + strconv.FormatFloat(h6, 'g', 5, 64) + "MPa " + "x6 = " + strconv.FormatFloat(x6, 'g', 4, 64) + " h7 = " + strconv.FormatFloat(h7, 'g', 5, 64) + "MPa" + " x7 = " + strconv.FormatFloat(x7, 'g', 5, 64), "Work Output: " + strconv.FormatFloat(wt, 'g', 5, 64) + "KJ/KG"},
		{"FeedWater Heater", "States at Turbine and Low Pressure Pump", "y = " + strconv.FormatFloat(y, 'g', 4, 64), "N/A"},
		{"High Pressure Pump", "P3 = " + strconv.FormatFloat(pressureAtFeedwater, 'g', 5, 64) + "MPa", "h4 = " + strconv.FormatFloat(h4, 'g', 5, 64) + "MPa", "Work Input: " + strconv.FormatFloat(wp2, 'g', 5, 64) + "KJ/KG"},
		{"Boiler", "P4 = " + strconv.FormatFloat(pressureAtFeedwater, 'g', 5, 64) + "MPa ", "Same as Input into the Turbine", "Heat Added : " + strconv.FormatFloat(qh, 'g', 5, 64) + "KJ/KG"},
	}

	table.AppendBulk(data)
	table.SetFooter([]string{"", "", "Efficiency of Cycle", strconv.FormatFloat(nth, 'g', 5, 64) + "%"})

	table.Render()
}