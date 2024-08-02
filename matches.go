package hltv

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

type Match struct {
	ID    int         `json:"id"`
	Teams []MatchTeam `json:"teams"`
	Event event       `json:"event"`
	Maps  []matchMap  `json:"maps"`
}

type MatchTeam struct {
	ID      int           `json:"id"`
	Name    string        `json:"name"`
	Logo    string        `json:"logo"`
	Score   int           `json:"score"`
	Players []playerStats `json:"players"`
}

type playerStats struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Nickname string  `json:"nickname"`
	Kills    int     `json:"kills"`
	Deaths   int     `json:"deaths"`
	ADR      float64 `json:"adr"`
	KAST     float64 `json:"kast"`
	Rating   float64 `json:"rating"`
}

type matchMap struct {
	Name  string    `json:"name"`
	Teams []mapTeam `json:"teams"`
}

type mapTeam struct {
	Name    string    `json:"name"`
	Results mapResult `json:"results"`
}

type mapResult struct {
	Total int `json:"total"`
	Ct    int `json:"ct"`
	Tt    int `json:"tt"`
}

type event struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
	URL  string `json:"url"`
}

func (client *Client) GetMatch(matchID int) (*Match, error) {
	res, err := client.fetch(fmt.Sprintf("%v/matches/%v/matchX", client.baseURL, matchID))
	if err != nil {
		return nil, err
	}

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var match Match

	match.ID = matchID

	teamContainer := document.Find(".standard-box.teamsBox").Children()
	playerContainer := document.Find(".stats-content").First().Find(".totalstats")
	teamA := getTeam(teamContainer.First(), playerContainer.First())
	teamB := getTeam(teamContainer.Last(), playerContainer.Last())
	match.Teams = append(match.Teams, teamA, teamB)

	match.Event = getEvent(document.Find(".matchSidebarEvent"), teamContainer)

	document.Find(".mapholder").Each(func(i int, mapContainer *goquery.Selection) {
		match.Maps = append(match.Maps, getMatchMap(mapContainer))
	})

	return &match, nil
}

func getEvent(eventContainer *goquery.Selection, teamContainer *goquery.Selection) (event event) {

	eventLink := teamContainer.Find(".event a")
	eventURL := eventLink.AttrOr("href", "")
	event.ID = idFromURL(eventURL, 2)
	event.Name = eventLink.Text()
	event.Logo = eventContainer.Find(".matchSidebarEventLogo").AttrOr("src", "")
	event.URL = fmt.Sprintf("%v/"+strings.Trim(eventURL, "/"), BaseURL)

	return
}

func getTeam(teamContainer *goquery.Selection, playerContainer *goquery.Selection) (team MatchTeam) {

	teamURL := playerContainer.Find("a").AttrOr("href", "")
	team.ID = idFromURL(teamURL, 2)
	team.Name = teamContainer.Find(".teamName").Text()
	team.Logo = teamContainer.Find(".logo").AttrOr("src", "")
	team.Score, _ = strconv.Atoi(teamContainer.Find(".teamName").Parent().Next().Text())

	playerContainer.Find(".players").Parent().Each(func(i int, playerRow *goquery.Selection) {
		if i == 0 {
			return
		}
		team.Players = append(team.Players, getStats(playerRow))
	})

	return
}

func getStats(playerRow *goquery.Selection) (stats playerStats) {
	playerURL := playerRow.Find("a").AttrOr("href", "")
	stats.ID = idFromURL(playerURL, 2)

	nameParts := strings.Split(playerRow.Find(".statsPlayerName").First().Text(), "'")
	stats.Name = strings.TrimSpace(nameParts[0]) + nameParts[2]
	stats.Nickname = playerRow.Find(".player-nick").Text()

	kd := strings.Split(playerRow.Find(".kd").Text(), "-")
	kills, _ := strconv.Atoi(kd[0])
	stats.Kills = kills
	deaths, _ := strconv.Atoi(kd[1])
	stats.Deaths = deaths

	adr, _ := strconv.ParseFloat(playerRow.Find(".adr").Text(), 64)
	stats.ADR = adr

	kastString := playerRow.Find(".kast").Text()
	stats.KAST, _ = strconv.ParseFloat(kastString, 64)

	rating, _ := strconv.ParseFloat(playerRow.Find(".rating").Text(), 64)
	stats.Rating = rating
	return
}

func getMatchMap(mapContainer *goquery.Selection) (matchMap matchMap) {
	matchMap.Name = mapContainer.Find(".mapname").Text()

	names := mapContainer.Find(".results-teamname")
	scores := mapContainer.Find(".results-team-score")

	var teamA, teamB *mapTeam

	halfScores := mapContainer.Find(".results-center-half-score").Children()

	teamA = getMapTeam(names.First().Text(), scores.First().Text(), halfScores.Eq(1), halfScores.Eq(5))
	teamB = getMapTeam(names.Last().Text(), scores.Last().Text(), halfScores.Eq(3), halfScores.Eq(7))

	matchMap.Teams = append(matchMap.Teams, *teamA, *teamB)

	return
}

func getMapTeam(name string, score string, halfScores ...*goquery.Selection) *mapTeam {
	total, _ := strconv.Atoi(score)

	var ct, tt int

	for _, half := range halfScores {
		if half.HasClass("ct") {
			ct, _ = strconv.Atoi(half.Text())
		} else {
			tt, _ = strconv.Atoi(half.Text())
		}
	}

	return &mapTeam{
		Name: name,
		Results: mapResult{
			Total: total,
			Ct:    ct,
			Tt:    tt,
		},
	}
}
