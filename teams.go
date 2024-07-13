package hltv

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Team struct {
	// w Go częściej się używa ID
	// https://google.github.io/styleguide/go/decisions#initialisms
	Id   int    `json:"id"`
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

// teamId -> teamID
// https://google.github.io/styleguide/go/decisions#initialisms

func (client *Client) GetTeam(teamId int) (*TeamDetails, error) {
	res, err := Fetch(fmt.Sprintf("%v/team/%v/teamX", client.baseUrl, teamId))
	if err != nil {
		// println nie powinien być używany, w docsach jest info
		// " it is not guaranteed to stay in the language."
		println(err.Error())
		return nil, err
	}
	// response.Body nie jest zamykany
	// defer res.Body.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)
	// brak error handlingu
	worldRanking, _ := strconv.Atoi(strings.Replace(document.Find(".profile-team-stat:nth-child(1) a").Text(), "#", "", -1))
	// brak error handlingu
	top30, _ := strconv.Atoi(document.Find(".profile-team-stat:nth-child(2) span.right").Text())
	// brak error handlingu
	avgAge, _ := strconv.ParseFloat(document.Find(".profile-team-stat:nth-child(3) span.right").Text(), 32)
	// brak error handlingu

	var teamDetails TeamDetails
	teamDetails.Id = teamId
	teamDetails.Name = document.Find(".profile-team-name").Text()
	teamDetails.LogoURL = document.Find(".profile-team-logo-container img").AttrOr("src", "")
	teamDetails.Country = document.Find(".team-country img").AttrOr("title", "")
	teamDetails.WorldRanking = uint8(worldRanking)
	teamDetails.WeeksInTop30 = top30
	teamDetails.AveragePlayerAge = float32(avgAge)
	teamDetails.CoachName = strings.TrimSpace(document.Find(".profile-team-stat:nth-child(4) a").Text())

	return &teamDetails, nil
}
