package main

import (
	"encoding/json"
	"fmt"
	"github.com/soszin/hltv"
)

func main() {
	hltvClient := hltv.New()

	ranking, _ := hltvClient.GetPlayer(11893)
	jsonData, _ := json.Marshal(ranking)
	fmt.Printf("%s\n", jsonData)
}
