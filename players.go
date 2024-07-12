package hltv

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Player struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Nickname    string `json:"nickname"`
	YearOfBirth int    `json:"year_of_birth"`
}

type PlayerDetails struct {
	Player
	Image       string `json:"image"`
	CurrentTeam *Team  `json:"current_team"`
}

func (client *Client) GetPlayer(playerId int) (*PlayerDetails, error) {
	res, err := Fetch(fmt.Sprintf("%v/player/%v/playerX", client.baseUrl, playerId))
	if err != nil {
		println(err.Error())
		return nil, err
	}

	document, err := goquery.NewDocumentFromReader(res.Body)

	var playerDetails PlayerDetails

	age, _ := strconv.Atoi(strings.ReplaceAll(document.Find(".playerAge span[itemprop='text']").Text(), " years", ""))

	playerDetails.Id = playerId
	playerDetails.Name = document.Find(".playerRealname").Text()
	playerDetails.Nickname = document.Find(".playerNickname").Text()
	playerDetails.YearOfBirth = time.Now().Year() - age
	playerDetails.Image = document.Find(".bodyshot-img").AttrOr("src", "")

	teamUrl := document.Find(".playerTeam a").AttrOr("href", "")
	teamName := document.Find(".playerTeam a").Text()

	if teamUrl != "" {
		teamUrlRegexp := regexp.MustCompile(`/team/(\d+)/(.*)`)
		teamUrlMatches := teamUrlRegexp.FindStringSubmatch(teamUrl)
		teamId, _ := strconv.Atoi(teamUrlMatches[1])
		playerDetails.CurrentTeam = &Team{
			Id:   teamId,
			Name: teamName,
		}
	}

	return &playerDetails, nil
}
