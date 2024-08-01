package main

import (
	"encoding/json"
	"fmt"
	"github.com/soszin/hltv"
)

func main() {
	hltvClient := hltv.New()

	player, _ := hltvClient.GetPlayer(11893)
	jsonData, _ := json.MarshalIndent(player, "", "    ")
	fmt.Printf("%s\n", jsonData)
}
