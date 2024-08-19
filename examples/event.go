package main

import (
	"encoding/json"
	"fmt"
	"github.com/soszin/hltv"
	"log"
)

func main() {
	hltvClient := hltv.New()

	event, err := hltvClient.GetEvent(7436)
	if err != nil {
		log.Fatalln(err)
	}
	jsonData, _ := json.MarshalIndent(event, "", "    ")
	fmt.Printf("%s\n", jsonData)
}
