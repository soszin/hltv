package main

import (
	"encoding/json"
	"fmt"
	"github.com/soszin/hltv"
)

func main() {
	hltvClient := hltv.New()

	ranking, _ := hltvClient.GetRanking()
	jsonData, _ := json.MarshalIndent(ranking, "", "    ")
	fmt.Printf("%s\n", jsonData)
}
