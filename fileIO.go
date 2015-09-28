package main

import (
	"bufio"
	"encoding/csv"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func loadLeagueFile(filename string) []team {
	// Open a RO file
	decodeFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer decodeFile.Close()

	// Create a decoder
	decoder := gob.NewDecoder(decodeFile)

	// Place to decode into
	leagueFromFile := make([]team, 16)

	// Decode -- We need to pass a pointer
	decoder.Decode(&leagueFromFile)

	return leagueFromFile
}

func saveLeagueFile(filename string, leuguelist []team) {
	// Create a file for IO
	encodeFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	// Since this is a binary format large parts of it will be unreadable
	encoder := gob.NewEncoder(encodeFile)

	// Write to the file
	if err := encoder.Encode(leuguelist); err != nil {
		panic(err)
	}
	encodeFile.Close()
}

func loadPlayerFile(filename string) map[string]fflPlayer {
	// Open a RO file
	decodeFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer decodeFile.Close()

	// Create a decoder
	decoder := gob.NewDecoder(decodeFile)

	// Place to decode into
	playersFromFile := make(map[string]fflPlayer)

	// Decode -- We need to pass a pointer
	decoder.Decode(&playersFromFile)

	return playersFromFile
}

func savePlayerFile(filename string, plist map[string]fflPlayer) {
	// Create a file for IO
	encodeFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	// Since this is a binary format large parts of it will be unreadable
	encoder := gob.NewEncoder(encodeFile)

	// Write to the file
	if err := encoder.Encode(plist); err != nil {
		panic(err)
	}
	encodeFile.Close()
}

func loadTeamFile(name string, bootstrap string, owner string, allPlayers map[string]fflPlayer) team {

	var team2load team

	team2load.Name = name
	team2load.Filename = ""
	team2load.Bootstrap = bootstrap
	team2load.Owner = owner

	//if we have a bootstrap file, use that to build up the team.
	if _, err := os.Stat(team2load.Bootstrap); err == nil {
		log.Printf("%v Bootstrap file exists; processing...\n", team2load.Bootstrap)
		team2load.Players = make(map[string]fflPlayer)
		// Open a RO file
		decodeFile, err := os.Open(team2load.Bootstrap)
		if err != nil {
			panic(err)
		}
		defer decodeFile.Close()

		scanner := bufio.NewScanner(decodeFile)
		for scanner.Scan() {
			//log.Printf(">%v<\n", scanner.Text())
			arr := strings.Split(strings.Trim(scanner.Text(), " "), "|")
			searchName := strings.ToLower(strings.Trim(arr[0], "' ")) + "|" + arr[1]

			//log.Printf("Searching for%v\n", searchName)
			//fmt.Println("PLAYER NAME:", strings.ToLower(strings.Trim(arr[0], "' ")))
			p, ok := allPlayers[searchName]

			if ok {
				p.ScoringWeekStart, _ = strconv.Atoi(arr[2])
				p.ScoringWeekEnd, _ = strconv.Atoi(arr[3])

				team2load.Players[arr[0]] = p
				//log.Printf("Added %v\n", p.Name)
			} else {
				log.Printf(">>PLAYER NOT FOUND %v\n", searchName)
				p, err := getMissingPlayer(searchName)
				if err == nil {
					team2load.Players[arr[0]] = p
				}
			}

		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading bootstrap file:", err)
		}

	} else if _, err := os.Stat(team2load.Filename); err == nil {
		log.Printf("filename found .. processing")

		// Create a decoder
		//decoder := json.NewDecoder(decodeFile)

		//// Place to decode into
		//playersFromFile := make(map[int]fflPlayer)

		//// Decode -- We need to pass a pointer

		//teamFromFile := team{}
		//	decoder.Decode(&teamFromFile)

		//if _, err := os.Stat("myplayers.txt"); err == nil {
		//	log.Printf("TEAM file exists; processing...\n")

		//	} else {
		//		log.Printf("TEAM file DOES NOT exist\n")
		//	}

	}

	return team2load
}

func getMissingPlayer(name string) (player fflPlayer, err error) {
	log.Printf("Searching missing players list for %s...\n", name)

	// Open a RO file
	decodeFile, err := os.Open("data\\missingPlayers.txt")
	if err != nil {
		panic(err)
	}
	defer decodeFile.Close()

	scanner := bufio.NewScanner(decodeFile)
	for scanner.Scan() {
		//log.Printf(">%v<\n", scanner.Text())
		arr := strings.Split(strings.Trim(scanner.Text(), " "), "|")
		searchName := strings.ToLower(strings.Trim(arr[0], "' ")) + "|" + arr[1]
		if name == searchName {
			player := fflPlayer{}
			player.ID = -1
			player.Index = searchName
			player.Name = arr[0]
			player.Team = arr[1]
			player.URL = ""
			player.Position = arr[2]
			for i := 2; i < len(arr); i++ {

				k, _ := strconv.Atoi(arr[i])
				player.TotalPoints += k
			}
			log.Printf("Found & processed %s...\n", player.Name)
			return player, nil
		}
	}
	return fflPlayer{}, errors.New("can't find the player")
}

func getTeams(f string) (records [][]string) {
	csvfile, err := os.Open(f)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return rawCSVdata
}
