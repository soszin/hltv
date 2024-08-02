package main

import (
	"encoding/json"
	"fmt"
	"github.com/soszin/hltv"
	"log"
)

func main() {
	hltvClient := hltv.New()

	player, err := hltvClient.GetMatch(2373786)
	if err != nil {
		log.Fatalln(err)
	}
	jsonData, _ := json.MarshalIndent(player, "", "    ")
	fmt.Printf("%s\n", jsonData)
}
