package main

import (
	"encoding/json"
	"fmt"
	"github.com/soszin/hltv"
)

func main() {
	hltvClient := hltv.New()

	ranking, _ := hltvClient.GetRanking()
	jsonData, _ := json.Marshal(ranking)
	fmt.Printf("%s\n", jsonData)
}
