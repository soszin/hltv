package hltv

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"strconv"
	"strings"
)

type Ranking struct {
	Position int   `json:"position"`
	Points   int   `json:"points"`
	Team     *Team `json:"team"`
}

func (client *Client) GetRanking() ([]Ranking, error) {
	res, err := client.fetch(fmt.Sprintf("%v/ranking/teams", client.baseURL))
	if err != nil {
		println(err.Error())
		return nil, err
	}
	defer res.Body.Close()

	var rankingList []Ranking

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	document.Find(".ranked-team").Each(func(i int, s *goquery.Selection) {

		name := s.Find(".name").Text()

		positionString := strings.Replace(s.Find(".position").Text(), "#", "", -1)
		position, _ := strconv.Atoi(positionString)

		var points int
		pointsRegexp := regexp.MustCompile(`\((\d+)\s.*\)$`)
		pointStrings := pointsRegexp.FindStringSubmatch(s.Find(".points").Text())
		if len(pointStrings) > 0 {
			points, _ = strconv.Atoi(pointStrings[1])
		} else {
			points = 0
		}

		var teamID int
		teamUrl, _ := s.Find(".more .moreLink").First().Attr("href")
		r := regexp.MustCompile(`/team/(\d+)/.*`)
		teamIdString := r.FindStringSubmatch(teamUrl)
		if len(teamIdString) > 0 {
			teamID, _ = strconv.Atoi(teamIdString[1])
		} else {
			teamID = 0
		}

		ranking := Ranking{
			Position: position,
			Points:   points,
		}

		ranking.Team = &Team{
			ID:   teamID,
			Name: name,
		}

		rankingList = append(rankingList, ranking)
	})

	return rankingList, nil
}
