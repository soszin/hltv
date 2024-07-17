package hltv

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

type Team struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TeamDetails struct {
	Team
	LogoURL          string  `json:"logo"`
	Country          string  `json:"country"`
	WorldRanking     uint8   `json:"world_ranking"`
	WeeksInTop30     int     `json:"weeks_in_top30"`
	CoachName        string  `json:"coach_name"`
	AveragePlayerAge float32 `json:"average_player_age"`
}

func (client *Client) GetTeam(teamID int) (*TeamDetails, error) {
	res, err := client.fetch(fmt.Sprintf("%v/team/%v/teamX", client.baseURL, teamID))
	if err != nil {
		println(err.Error())
		return nil, err
	}

	document, err := goquery.NewDocumentFromReader(res.Body)
	worldRanking, _ := strconv.Atoi(strings.Replace(document.Find(".profile-team-stat:nth-child(1) a").Text(), "#", "", -1))
	top30, _ := strconv.Atoi(document.Find(".profile-team-stat:nth-child(2) span.right").Text())
	avgAge, _ := strconv.ParseFloat(document.Find(".profile-team-stat:nth-child(3) span.right").Text(), 32)

	var teamDetails TeamDetails
	teamDetails.ID = teamID
	teamDetails.Name = document.Find(".profile-team-name").Text()
	teamDetails.LogoURL = document.Find(".profile-team-logo-container img").AttrOr("src", "")
	teamDetails.Country = document.Find(".team-country img").AttrOr("title", "")
	teamDetails.WorldRanking = uint8(worldRanking)
	teamDetails.WeeksInTop30 = top30
	teamDetails.AveragePlayerAge = float32(avgAge)
	teamDetails.CoachName = strings.TrimSpace(document.Find(".profile-team-stat:nth-child(4) a").Text())

	return &teamDetails, nil
}
