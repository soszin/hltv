package hltv

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Placement struct {
	Team
	Position string  `json:"position"`
	Prize    float32 `json:"prize"`
}

type EventDetail struct {
	Event
	DateStart     time.Time   `json:"date_start"`
	DateEnd       time.Time   `json:"date_end"`
	Status        string      `json:"status"`
	PrizePool     float32     `json:"prize_pool"`
	Location      string      `json:"location"`
	NumberOfTeams int         `json:"number_of_teams"`
	Teams         []Placement `json:"teams"`
}

func (client *Client) GetEvent(eventID int) (*EventDetail, error) {
	res, err := client.fetch(fmt.Sprintf("%v/events/%v/matchX", client.baseURL, eventID))
	if err != nil {
		return nil, err
	}

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var event EventDetail

	eventDates := document.Find(".eventdate [data-unix]")

	start, _ := strconv.ParseInt(eventDates.First().AttrOr("data-unix", ""), 10, 64)
	end, _ := strconv.ParseInt(eventDates.Last().AttrOr("data-unix", ""), 10, 64)
	numberOfTeams, _ := strconv.Atoi(document.Find(".teamsNumber").Last().Text())

	prizePoolReplacer := strings.NewReplacer("$", "", ",", "")
	prizePoolString := prizePoolReplacer.Replace(strings.TrimSpace(document.Find(".prizepool").Last().Text()))
	prizePool, _ := strconv.ParseFloat(prizePoolString, 32)

	event.ID = eventID
	event.Name = document.Find(".event-hub-title").Text()
	event.DateStart = time.UnixMilli(start)
	event.DateEnd = time.UnixMilli(end)
	event.Status = document.Find(".event-hub-indicator").Text()
	event.PrizePool = float32(prizePool)
	event.Location = document.Find("tbody .location span").Text()
	event.NumberOfTeams = numberOfTeams

	document.Find(".placements .placement").Each(func(i int, placement *goquery.Selection) {
		event.Teams = append(event.Teams, getPlacementTeam(placement))
	})

	return &event, nil
}

func getPlacementTeam(placement *goquery.Selection) (team Placement) {

	priceReplacer := strings.NewReplacer("$", "", ",", "")
	prizeSting := priceReplacer.Replace(placement.Find(".prize").First().Text())
	prize, _ := strconv.ParseFloat(prizeSting, 32)

	teamElement := placement.Find(".team a")
	team.ID = idFromURL(teamElement.AttrOr("href", ""), 2)
	team.Name = teamElement.Text()
	team.Position = placement.Children().Eq(1).Text()
	team.Prize = float32(prize)
	return
}
