package main

import (
	"fmt"
	"os"
	"sort"
)

type fflPlayer struct {
	ID               int
	Name             string
	Team             string
	Value            string
	Position         string
	TotalPoints      int
	URL              string
	ScoringWeekStart int
	ScoringWeekEnd   int
	WeekStats        map[string]weekStats
}
type weekStats struct {
	Week           int
	Vs             string
	Goals          int
	KeyContrib     int
	StartingXI     int
	Sub            int
	YellowCard     int
	RedCard        int
	PenaltyMiss    int
	PenaltySaved   int
	OwnGoal        int
	Conceeded      int
	FullCleanSheet int
	PartCleanSheet int
	Points         int
}

type leagueType []team

func (a leagueType) Len() int           { return len(a) }
func (a leagueType) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a leagueType) Less(i, j int) bool { return a[i].TotalScore > a[j].TotalScore }

func main() {

	players := make(map[string]fflPlayer)
	// first check if we have a local list of all the players
	if _, err := os.Stat("data\\players.gob"); err == nil {
		fmt.Printf("Cache file exists; processing...\n")
		players = loadPlayerFile("data\\players.gob")
	} else {
		fmt.Printf("Cache file does not exist; scraping...\n")
		players = scrapeAllPlayers()
		fmt.Printf("Scraping complete...Saving to cache file\n")
		savePlayerFile("data\\players.gob", players)
	}

	fmt.Printf("Total ffl players available : %d\n", len(players))

	fmt.Println("Loading Team files into league...")

	league := make([]team, 20)

	adamTeam := loadTeamFile("ADAM", "data\\ADAM.json", "data\\ADAM.txt", "ADAM", players)
	processTeam(&adamTeam)
	league = append(league, adamTeam)

	charlieTeam := loadTeamFile("BLITZKRIEG TOTAL FOOTBALL ©", "data\\charlie.json", "data\\charlie.txt", "Charlie", players)
	processTeam(&charlieTeam)
	league = append(league, charlieTeam)

	ChrisGoodTeam := loadTeamFile("SLIGHTLY ATHLETICO", "data\\ChrisGood.json", "data\\ChrisGood.txt", "ChrisGood", players)
	processTeam(&ChrisGoodTeam)
	league = append(league, ChrisGoodTeam)

	darren := loadTeamFile("R.I.P ANDY COLE", "data\\darren.json", "data\\darren.txt", "darren", players)
	processTeam(&darren)
	league = append(league, darren)

	eTeam := loadTeamFile("BEGINNER'S LUCK", "data\\eddy.json", "data\\eddy.txt", "Eddy", players)
	processTeam(&eTeam)
	league = append(league, eTeam)

	Godber := loadTeamFile("YEAR OF THE RAT", "data\\Godber.json", "data\\Godber.txt", "Godber", players)
	processTeam(&Godber)
	league = append(league, Godber)

	howardTeam := loadTeamFile("VAN-HOOIJDONK.COM", "data\\howard.json", "data\\howard.txt", "Howard", players)
	processTeam(&howardTeam)
	league = append(league, howardTeam)

	mattTeam := loadTeamFile("THE BASH STREET KIDS", "data\\matt.json", "data\\matt.txt", "matt", players)
	processTeam(&mattTeam)
	league = append(league, mattTeam)

	mandeepTeam := loadTeamFile("WEDDED BLISS", "data\\mandeep.json", "data\\mandeep.txt", "mandeep", players)
	processTeam(&mandeepTeam)
	league = append(league, mandeepTeam)

	richardTeam := loadTeamFile("SNAKE IN THE GRASS", "data\\richard.json", "data\\richard.txt", "Richard", players)
	processTeam(&richardTeam)
	league = append(league, richardTeam)

	ryanTeam := loadTeamFile("TEN AND A HALF MEN IN FLIGHT", "data\\ryan.json", "data\\ryan.txt", "Ryan", players)
	processTeam(&ryanTeam)
	league = append(league, ryanTeam)

	steveTeam := loadTeamFile("Steve's Putney Pillagers", "data\\steve.json", "data\\steve.txt", "Steve", players)
	processTeam(&steveTeam)
	league = append(league, steveTeam)

	sulemanTeam := loadTeamFile("REBEL WITHOUT A CLAUSE", "data\\suleman.json", "data\\suleman.txt", "Suleman", players)
	processTeam(&sulemanTeam)
	league = append(league, sulemanTeam)

	talTeam := loadTeamFile("Tal's Terrible Thames Dittioners", "data\\tal.json", "data\\tal.txt", "Tal", players)
	processTeam(&talTeam)
	league = append(league, talTeam)

	tonyTeam := loadTeamFile("LOS TESTICULOS DE PERRO", "data\\tony.json", "data\\tony.txt", "Tony", players)
	processTeam(&tonyTeam)
	league = append(league, tonyTeam)

	yusufTeam := loadTeamFile("W.C. MILAN", "data\\yusuf.json", "data\\yusuf.txt", "Yusuf", players)
	processTeam(&yusufTeam)
	league = append(league, yusufTeam)

	sort.Sort(leagueType(league))

	for _, v := range league {

		if v.TotalScore > 0 {
			fmt.Printf("Team: %v, Owner: %v, TotalPoints: %v\n", v.Name, v.Owner, v.TotalScore)
		}
	}

}
func processTeam(t *team) {
	for _, p := range t.Players {
		scrapePlayer(&p)
		calculatePlayerScore(&p)
		t.TotalScore = t.TotalScore + p.TotalPoints
		//fmt.Printf("%v %v %v %v\n", p.Name, p.TotalPoints, p.ScoringWeekStart, p.ScoringWeekEnd)
	}

}

func calculatePlayerScore(p *fflPlayer) {
	totalScore := 0
	if p.ScoringWeekStart == 0 && p.ScoringWeekEnd == 999 {
		for _, v := range p.WeekStats {
			totalScore = totalScore + v.Points
			//fmt.Printf("%v %v %v\n", p.Name, v.Vs, v.Points)
		}
		p.TotalPoints = totalScore
	} else {
		for _, v := range p.WeekStats {
			if v.Week >= p.ScoringWeekStart && v.Week <= p.ScoringWeekEnd {
				totalScore = totalScore + v.Points
			}
			//fmt.Printf("%v %v %v\n", p.Name, v.Vs, v.Points)
		}
		p.TotalPoints = totalScore
	}
}
