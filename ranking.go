package hltv

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Ranking struct {
	Position int   `json:"position"`
	Points   int   `json:"points"`
	Team     *Team `json:"team"`
}

func (client *Client) GetRanking() ([]Ranking, error) {
	res, err := Fetch(fmt.Sprintf("%v/ranking/teams", client.baseUrl))
	if err != nil {
		// println nie powinien być używany, w docsach jest info
		// " it is not guaranteed to stay in the language."
		println(err.Error())
		return nil, err
	}
	// response.Body nie jest zamykany
	// defer res.Body.Close()

	// tutaj można by było preallokować pamięć
	// przy użyciu document.Find(".ranked-team").Length() i make
	var rankingList []Ranking

	document, err := goquery.NewDocumentFromReader(res.Body)

	document.Find(".ranked-team").Each(func(i int, s *goquery.Selection) {

		name := s.Find(".name").Text()

		positionString := strings.Replace(s.Find(".position").Text(), "#", "", -1)
		position, _ := strconv.Atoi(positionString)
		// brak error handlingu

		var points int
		// ten regexp nie musi być kompilowany przy każdym wywołaniu tej metody
		// mógłbyś go wrzucić jako globalną zmienną
		// regexp.Compile jest kosztowny
		pointsRegexp := regexp.MustCompile(`\((\d+)\s.*\)$`)
		pointStrings := pointsRegexp.FindStringSubmatch(s.Find(".points").Text())
		if len(pointStrings) > 0 {
			points, _ = strconv.Atoi(pointStrings[1])
		} else {
			points = 0
		}

		var teamId int
		teamUrl, _ := s.Find(".more .moreLink").First().Attr("href")
		// ten regexp nie musi być kompilowany przy każdym wywołaniu tej metody
		// mógłbyś go wrzucić jako globalną zmienną
		// regexp.Compile jest kosztowny
		r := regexp.MustCompile(`/team/(\d+)/.*`)
		teamIdString := r.FindStringSubmatch(teamUrl)
		if len(teamIdString) > 0 {
			teamId, _ = strconv.Atoi(teamIdString[1])
		} else {
			teamId = 0
		}

		ranking := Ranking{
			Position: position,
			Points:   points,
		}

		ranking.Team = &Team{
			Id:   teamId,
			Name: name,
		}

		rankingList = append(rankingList, ranking)
	})

	return rankingList, nil
}
