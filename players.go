package hltv

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Player struct {
	// w Go częściej się używa ID
	// https://google.github.io/styleguide/go/decisions#initialisms
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

// playerId -> playerID
// https://google.github.io/styleguide/go/decisions#initialisms

func (client *Client) GetPlayer(playerId int) (*PlayerDetails, error) {
	res, err := Fetch(fmt.Sprintf("%v/player/%v/playerX", client.baseUrl, playerId))
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

	var playerDetails PlayerDetails

	age, _ := strconv.Atoi(strings.ReplaceAll(document.Find(".playerAge span[itemprop='text']").Text(), " years", ""))
	// brak error handlingu

	playerDetails.Id = playerId
	playerDetails.Name = document.Find(".playerRealname").Text()
	playerDetails.Nickname = document.Find(".playerNickname").Text()
	playerDetails.YearOfBirth = time.Now().Year() - age
	playerDetails.Image = document.Find(".bodyshot-img").AttrOr("src", "")

	teamUrl := document.Find(".playerTeam a").AttrOr("href", "")
	teamName := document.Find(".playerTeam a").Text()

	if teamUrl != "" {
		// ten regexp nie musi być kompilowany przy każdym wywołaniu tej metody
		// mógłbyś go wrzucić jako globalną zmienną
		// regexp.Compile jest kosztowny
		teamUrlRegexp := regexp.MustCompile(`/team/(\d+)/(.*)`)
		teamUrlMatches := teamUrlRegexp.FindStringSubmatch(teamUrl)
		// według mnie tutaj nawet nie potrzebujesz regexpa tak na prawdę
		// strings.Split by wystarczył
		teamId, _ := strconv.Atoi(teamUrlMatches[1])
		// brak error handlingu
		playerDetails.CurrentTeam = &Team{
			Id:   teamId,
			Name: teamName,
		}
	}

	return &playerDetails, nil
}
