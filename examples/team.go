package main

import (
	"encoding/json"
	"fmt"
	"github.com/soszin/hltv"
)

func main() {
	hltvClient := hltv.New()

	team, _ := hltvClient.GetTeam(9565)
	jsonData, _ := json.MarshalIndent(team, "", "    ")
	fmt.Printf("%s\n", jsonData)
}
