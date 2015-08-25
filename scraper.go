package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func scrapeAllPlayers() map[string]fflPlayer {

	players := make(map[string]fflPlayer)

	doc, err := goquery.NewDocument(os.Getenv("FFL-URL1"))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	p := fflPlayer{}
	doc.Find("tr").Each(func(i int, row *goquery.Selection) {
		s := row.Find("td.first")
		p.Name = strings.ToLower(strings.Trim(s.Text(), "' "))
		p.URL, _ = s.Find("a").Attr("href")
		if p.URL == "" {
			log.Printf("empty URL detected for %v", p.Name)
		}

		p.ID, _ = strconv.Atoi(strings.TrimPrefix(p.URL, os.Getenv("FFL-URL2")))
		p.Team = s.Next().Text()
		p.Index = p.Name + "|" + p.Team

		//fmt.Printf("Player Processed:\n Name:%s, Team:%v , index:%s, url:%s\n", p.Name, p.Team, p.Index, p.URL)
		players[p.Index] = p
	})

	return players

}

//https://fantasyfootball.telegraph.co.uk/premierleague/statistics/points/

func scrapePlayer(p *fflPlayer) {
	if p.URL == "" {
		log.Printf("empty URL detected for %v", p.Name)
		return
	}

	//fmt.Printf("Processing:%s URL:%s\n", p.Name, p.URL)

	p.WeekStats = make(map[string]weekStats)
	doc, err := goquery.NewDocument(p.URL)
	if err != nil {
		log.Print(p.URL)
		log.Fatal(err)
		return
	}

	//p.Team = doc.Find("#stats-team").Text()
	p.Value = doc.Find("#stats-value").Text()
	p.Position = doc.Find("#stats-position").Text()

	// the data differs depending on position
	switch p.Position {
	case "Midfielder", "Striker":
		p.WeekStats = processMidAndFwd(doc)
	case "Defender":
		p.WeekStats = processDef(doc)
	case "Goalkeeper":
		p.WeekStats = processGK(doc)
	}
	//fmt.Printf("weekstats:%v\n", p.WeekStats)
}

func processMidAndFwd(doc *goquery.Document) map[string]weekStats {

	p := make(map[string]weekStats)
	// need to skip the header so use .Next()
	doc.Find("#individual-player").Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {

		//fmt.Printf("Text:%s \n", s.Text())
		h, _ := s.Html()
		arr := strings.Split(strings.Trim(strings.Trim(h, "\n \t<td>"), "</"), "</td><td>")
		//fmt.Printf("item:%q \n", arr)
		//<td>1</td><td>Swansea</td><td>1</td><td>0</td><td>1</td><td>0</td><td>0</td><td>0</td><td>0</td><td>0</td><td>7</td>

		var w weekStats
		w.Week, _ = strconv.Atoi(arr[0])        //int
		w.Vs = arr[1]                           //string
		w.Goals, _ = strconv.Atoi(arr[2])       //int
		w.KeyContrib, _ = strconv.Atoi(arr[3])  //int
		w.StartingXI, _ = strconv.Atoi(arr[4])  //int
		w.Sub, _ = strconv.Atoi(arr[5])         //int
		w.YellowCard, _ = strconv.Atoi(arr[6])  //int
		w.RedCard, _ = strconv.Atoi(arr[7])     //int
		w.PenaltyMiss, _ = strconv.Atoi(arr[8]) //int
		w.OwnGoal, _ = strconv.Atoi(arr[9])     //int

		w.Points, _ = strconv.Atoi(arr[10]) //int
		//fmt.Println("Adding to map")
		key := strconv.Itoa(w.Week) + "-" + w.Vs
		p[key] = w
		//fmt.Printf("%s \n", element)
		//if (i+1)%11 == 0 {
		//	fmt.Printf("\n")
		//}
	})
	return p
}

func processDef(doc *goquery.Document) map[string]weekStats {
	p := make(map[string]weekStats)
	// need to skip the header so use .Next()
	doc.Find("#individual-player").Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {

		//fmt.Printf("Text:%s \n", s.Text())
		h, _ := s.Html()
		arr := strings.Split(strings.Trim(strings.Trim(h, "\n \t<td>"), "</"), "</td><td>")

		var w weekStats
		w.Week, _ = strconv.Atoi(arr[0])        //int
		w.Vs = arr[1]                           //string
		w.Goals, _ = strconv.Atoi(arr[2])       //int
		w.KeyContrib, _ = strconv.Atoi(arr[3])  //int
		w.StartingXI, _ = strconv.Atoi(arr[4])  //int
		w.Sub, _ = strconv.Atoi(arr[5])         //int
		w.YellowCard, _ = strconv.Atoi(arr[6])  //int
		w.RedCard, _ = strconv.Atoi(arr[7])     //int
		w.PenaltyMiss, _ = strconv.Atoi(arr[8]) //int
		w.OwnGoal, _ = strconv.Atoi(arr[9])     //int

		w.Conceeded, _ = strconv.Atoi(arr[10])      //int
		w.FullCleanSheet, _ = strconv.Atoi(arr[11]) //int
		w.PartCleanSheet, _ = strconv.Atoi(arr[12]) //int
		w.Points, _ = strconv.Atoi(arr[13])         //int
		key := strconv.Itoa(w.Week) + "-" + w.Vs
		p[key] = w
	})
	return p
}

func processGK(doc *goquery.Document) map[string]weekStats {
	p := make(map[string]weekStats)
	// need to skip the header so use .Next()
	doc.Find("#individual-player").Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {

		h, _ := s.Html()
		arr := strings.Split(strings.Trim(strings.Trim(h, "\n \t<td>"), "</"), "</td><td>")

		var w weekStats
		w.Week, _ = strconv.Atoi(arr[0])        //int
		w.Vs = arr[1]                           //string
		w.Goals, _ = strconv.Atoi(arr[2])       //int
		w.KeyContrib, _ = strconv.Atoi(arr[3])  //int
		w.StartingXI, _ = strconv.Atoi(arr[4])  //int
		w.Sub, _ = strconv.Atoi(arr[5])         //int
		w.YellowCard, _ = strconv.Atoi(arr[6])  //int
		w.RedCard, _ = strconv.Atoi(arr[7])     //int
		w.PenaltyMiss, _ = strconv.Atoi(arr[8]) //int

		w.PenaltySaved, _ = strconv.Atoi(arr[9]) //int

		w.OwnGoal, _ = strconv.Atoi(arr[10]) //int

		w.Conceeded, _ = strconv.Atoi(arr[11])      //int
		w.FullCleanSheet, _ = strconv.Atoi(arr[12]) //int
		w.PartCleanSheet, _ = strconv.Atoi(arr[13]) //int

		w.Points, _ = strconv.Atoi(arr[14]) //int
		key := strconv.Itoa(w.Week) + "-" + w.Vs
		p[key] = w
	})
	return p
}
