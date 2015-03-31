package main

import (
	"fmt"
	"os"
)

type fflPlayer struct {
	Id          int
	Name        string
	Team        string
	Value       string
	Position    string
	TotalPoints int
	Url         string
	WeekStats   []weekStats
}
type weekStats struct {
	Week        int
	Vs          string
	Goals       int
	KeyContrib  int
	StartingXI  int
	Sub         int
	YellowCard  int
	RedCard     int
	PenaltyMiss int
	OwnGoal     int
	Points      int
}

func main() {

	players := make(map[int]fflPlayer)
	// first check if we have a local list of all the players
	if _, err := os.Stat("players.gob"); err == nil {
		fmt.Printf("Cache file exists; processing...\n")
		players = loadPlayerFile("players.gob")
	} else {
		fmt.Printf("Cache file does not exist; scraping...\n")
		players = scrapeAllPlayers()
		fmt.Printf("Scraping complete...Saving to cache file\n")
		savePlayerFile("players.gob", players)
	}

	fmt.Printf("Total players : %d\n", len(players))
	fmt.Printf("-------------------------\n")
	fmt.Printf("%v\n", players[3063])
	fmt.Printf("-------------------------\n")
	fmt.Println("Loading Team file...")

	var myteam team
	if _, err := os.Stat("myplayers.txt"); err == nil {
		fmt.Printf("TEAM file exists; processing...\n")
		myteam = loadTeamFile("myplayers.txt")
	} else {
		fmt.Printf("TEAM file DOES NOT exist\n")
	}
	fmt.Printf("", myteam)

}
