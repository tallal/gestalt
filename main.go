package main

import (
	"log"
	"os"
	"sort"
)

type leagueType []team

func (a leagueType) Len() int           { return len(a) }
func (a leagueType) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a leagueType) Less(i, j int) bool { return a[i].TotalScore > a[j].TotalScore }

var league = make([]team, 16)

func main() {

	if os.Getenv("FFL-URL1") == "" || os.Getenv("FFL-URL2") == "" {
		log.Println("ERROR: REQUIRED URLS undefined")
		os.Exit(1)
	}

	players := make(map[string]fflPlayer)
	// first check if we have a local list of all the players
	if _, err := os.Stat("data\\players.gob"); err == nil {
		log.Printf("Player cache file exists; processing...\n")
		players = loadPlayerFile("data\\players.gob")
	} else {
		log.Printf("Cache file does not exist; scraping...\n")
		players = scrapeAllPlayers()
		log.Printf("Scraping complete...Saving to cache file\n")
		savePlayerFile("data\\players.gob", players)
	}

	log.Printf("Total ffl players available : %d\n", len(players))
	log.Println("Loading Team files into league...")

	teams := getTeams("data\\teams.csv")

	for _, each := range teams {
		t := loadTeamFile(each[0], each[1], each[2], players)
		processTeam(&t)
		league = append(league, t)
	}

	sort.Sort(leagueType(league))

	log.Println("Saving league GOB file...")
	saveLeagueFile("data\\league.gob", league)
}
