package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strconv"
	"strings"
)

func scrapeAllPlayers() map[int]fflPlayer {

	players := make(map[int]fflPlayer)

	doc, err := goquery.NewDocument("https://fantasyfootball.telegraph.co.uk/premierleague/PLAYERS/all")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	doc.Find("td.first").Each(func(i int, s *goquery.Selection) {

		p := fflPlayer{}
		p.Name = s.Text()
		p.Url, _ = s.Find("a").Attr("href")
		p.Id, _ = strconv.Atoi(strings.TrimPrefix(p.Url, "https://fantasyfootball.telegraph.co.uk/premierleague/statistics/points/"))
		players[p.Id] = p
	})

	return players

}

//https://fantasyfootball.telegraph.co.uk/premierleague/statistics/points/

func scrapePlayer(p fflPlayer) {

	doc, err := goquery.NewDocument(p.Url)
	if err != nil {
		log.Fatal(err)
		return
	}
	p.Team = doc.Find("#stats-team").Text()
	p.Value = doc.Find("#stats-value").Text()
	p.Position = doc.Find("#stats-position").Text()

	doc.Find("#individual-player").Find("td").Each(func(i int, s *goquery.Selection) {
		element := s.Text()
		fmt.Printf("%s ", element)
		if (i+1)%11 == 0 {
			fmt.Printf("\n")
		}
	})
}
