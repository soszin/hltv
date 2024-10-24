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
	ErrEventNameNotExist         = errors.New("event name not found")
	ErrStartDateNotExist         = errors.New("start date not found")
	ErrEndDateNotExist           = errors.New("end date not found")
	ErrStatusNotExist            = errors.New("status not found")
	ErrLocationNotExist          = errors.New("location not found")
	ErrNumberOfTeamsNotExist     = errors.New("number of teams not found")
	ErrPlacementIdNotExist       = errors.New("placement id not found")
	ErrPlacementTeamNameNotExist = errors.New("placement team name not found")
	ErrPlacementPrizeNotExist    = errors.New("placement prize not found")
	ErrPlacementPositionNotExist = errors.New("placement position not found")
)

type eventHTMLParser struct {
	document *goquery.Document
}

type eventPlacementHTMLParser struct {
	document *goquery.Selection
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

	placements, err := parser.getPlacementTeam()
	if err != nil {
		return nil, err
	}

	event.DateStart = startDate
	event.DateEnd = endDate
	event.Status = status
	event.PrizePool = prizePool
	event.Location = location
	event.NumberOfTeams = numberOfTeams
	event.Teams = placements

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

func (parser eventHTMLParser) getPlacementTeam() ([]Placement, error) {
	var firstError error
	var placements []Placement
	parser.document.Find(".placements .placement").Each(func(i int, placement *goquery.Selection) {
		if firstError != nil {
			return
		}
		placementParser := eventPlacementHTMLParser{placement}
		var team Placement

		teamID, err := placementParser.getId()
		if err != nil {
			if firstError == nil {
				firstError = err
			}
			return
		}

		teamName, err := placementParser.getTeamName()
		if err != nil {
			if firstError == nil {
				firstError = err
			}
			return
		}

		prize, err := placementParser.getPrize()
		if err != nil {
			if firstError == nil {
				firstError = err
			}
			return
		}

		position, err := placementParser.getPosition()
		if err != nil {
			if firstError == nil {
				firstError = err
			}
			return
		}

		team.ID = teamID
		team.Name = teamName
		team.Prize = prize
		team.Position = position

		placements = append(placements, team)
	})

	return placements, firstError
}

func (parser eventPlacementHTMLParser) getId() (int, error) {
	teamURL := parser.document.Find(".team a").AttrOr("href", "")
	if teamURL == "" {
		return 0, ErrPlacementIdNotExist
	}
	teamID, err := idFromURL(teamURL, 2)

	if err != nil {
		return 0, err
	}

	return teamID, nil
}

func (parser eventPlacementHTMLParser) getTeamName() (string, error) {
	teamName := parser.document.Find(".team a").Text()
	if teamName == "" {
		return "", ErrPlacementTeamNameNotExist
	}

	return teamName, nil
}

func (parser eventPlacementHTMLParser) getPrize() (float32, error) {
	prizeReplacer := strings.NewReplacer("$", "", ",", "")
	prizeString := parser.document.Find(".prize").First().Text()

	if prizeString == "" {
		return float32(0), ErrPlacementPrizeNotExist
	}

	prizeString = prizeReplacer.Replace(prizeString)
	prize, err := strconv.ParseFloat(prizeString, 32)

	if err != nil {
		return float32(0), err
	}

	return float32(prize), nil
}

func (parser eventPlacementHTMLParser) getPosition() (string, error) {
	position := parser.document.Children().Eq(1).Text()
	if position == "" {
		return "", ErrPlacementPositionNotExist
	}

	return position, nil
}
