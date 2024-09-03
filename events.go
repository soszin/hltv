package hltv

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
	"time"
)

var (
	ErrEventNameNotExist     = errors.New("event name not found")
	ErrStartDateNotExist     = errors.New("start date not found")
	ErrEndDateNotExist       = errors.New("end date not found")
	ErrStatusNotExist        = errors.New("status not found")
	ErrLocationNotExist      = errors.New("location not found")
	ErrNumberOfTeamsNotExist = errors.New("number of teams not found")
)

type eventHTMLParser struct {
	document *goquery.Document
}

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
	parser := eventHTMLParser{document}

	var event EventDetail

	event.ID = eventID
	event.Name, err = parser.getName()
	if err != nil {
		return nil, err
	}

	startDate, err := parser.getStartDate()
	if err != nil {
		return nil, err
	}

	endDate, err := parser.getEndDate()
	if err != nil {
		return nil, err
	}

	status, err := parser.getStatus()
	if err != nil {
		return nil, err
	}

	prizePool, err := parser.getPrizePool()
	if err != nil {
		return nil, err
	}

	location, err := parser.getLocation()
	if err != nil {
		return nil, err
	}

	numberOfTeams, err := parser.getNumberOfTeams()
	if err != nil {
		return nil, err
	}

	event.DateStart = startDate
	event.DateEnd = endDate
	event.Status = status
	event.PrizePool = prizePool
	event.Location = location
	event.NumberOfTeams = numberOfTeams

	document.Find(".placements .placement").Each(func(i int, placement *goquery.Selection) {
		placementTeam, err := getPlacementTeam(placement)
		if err != nil {
			return
		}
		event.Teams = append(event.Teams, placementTeam)
	})

	return &event, nil
}

func (parser eventHTMLParser) getName() (name string, err error) {
	name = parser.document.Find(".event-hub-title").Text()
	if name == "" {
		return "", ErrEventNameNotExist
	}

	return name, nil
}

func (parser eventHTMLParser) getStartDate() (time.Time, error) {
	startDateElement := parser.document.Find(".eventdate [data-unix]").First()
	startDateString := startDateElement.AttrOr("data-unix", "")
	if startDateString == "" {
		return time.Time{}, ErrStartDateNotExist
	}

	startUnixMilli, err := strconv.ParseInt(startDateString, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.UnixMilli(startUnixMilli), nil

}

func (parser eventHTMLParser) getEndDate() (time.Time, error) {
	endDateElement := parser.document.Find(".eventdate [data-unix]").Last()
	endDateString := endDateElement.AttrOr("data-unix", "")
	if endDateString == "" {
		return time.Time{}, ErrEndDateNotExist
	}

	endUnixMilli, err := strconv.ParseInt(endDateString, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.UnixMilli(endUnixMilli), nil

}

func (parser eventHTMLParser) getStatus() (string, error) {
	status := parser.document.Find(".event-hub-indicator").Text()
	if status == "" {
		return "", ErrStatusNotExist
	}

	return status, nil
}

func (parser eventHTMLParser) getPrizePool() (float32, error) {
	prizePoolReplacer := strings.NewReplacer("$", "", ",", "")
	prizePoolString := prizePoolReplacer.Replace(strings.TrimSpace(parser.document.Find(".prizepool").Last().Text()))
	prizePool, err := strconv.ParseFloat(prizePoolString, 32)
	if err != nil {
		return float32(0), err
	}

	return float32(prizePool), nil
}

func (parser eventHTMLParser) getLocation() (string, error) {
	location := parser.document.Find("tbody .location span").Text()
	if location == "" {
		return "", ErrLocationNotExist
	}
	return location, nil
}

func (parser eventHTMLParser) getNumberOfTeams() (int, error) {
	numberOfTeamsString := parser.document.Find(".teamsNumber").Last().Text()
	if numberOfTeamsString == "" {
		return 0, ErrNumberOfTeamsNotExist
	}
	numberOfTeams, err := strconv.Atoi(numberOfTeamsString)
	if err != nil {
		return 0, err
	}
	return numberOfTeams, nil
}

func getPlacementTeam(placement *goquery.Selection) (Placement, error) {
	prizeReplacer := strings.NewReplacer("$", "", ",", "")
	prizeSting := prizeReplacer.Replace(placement.Find(".prize").First().Text())
	prize, err := strconv.ParseFloat(prizeSting, 32)
	if err != nil {
		return Placement{}, err
	}

	teamElement := placement.Find(".team a")

	var team Placement
	team.ID = idFromURL(teamElement.AttrOr("href", ""), 2)
	team.Name = teamElement.Text()
	team.Position = placement.Children().Eq(1).Text()
	team.Prize = float32(prize)
	return team, nil
}
