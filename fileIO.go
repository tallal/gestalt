package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

func loadTeamFile(name string, filename string, bootstrap string, owner string, allPlayers map[string]fflPlayer) team {

	var team2load team

	team2load.Name = name
	team2load.Filename = filename
	team2load.Bootstrap = bootstrap
	team2load.Owner = owner

	//if we have a bootstrap file, use that to build up the team.
	if _, err := os.Stat(team2load.Bootstrap); err == nil {
		fmt.Printf("%v Bootstrap file exists; processing...\n", team2load.Bootstrap)
		team2load.Players = make(map[string]fflPlayer)
		// Open a RO file
		decodeFile, err := os.Open(team2load.Bootstrap)
		if err != nil {
			panic(err)
		}
		defer decodeFile.Close()

		scanner := bufio.NewScanner(decodeFile)
		for scanner.Scan() {
			//fmt.Printf(">%v<\n", scanner.Text())
			arr := strings.Split(strings.Trim(scanner.Text(), " "), "|")
			searchName := strings.ToLower(strings.Trim(arr[0], "' "))
			//fmt.Printf("Searching for%v\n", searchName)
			//fmt.Println("PLAYER NAME:", strings.ToLower(strings.Trim(arr[0], "' ")))
			p, ok := allPlayers[searchName]

			if ok {
				p.ScoringWeekStart, _ = strconv.Atoi(arr[1])
				p.ScoringWeekEnd, _ = strconv.Atoi(arr[2])

				team2load.Players[arr[0]] = p
			} else {
				fmt.Printf(">>NOT FOUND %v\n", searchName)
			}

		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "reading bootstrap file:", err)
		}

	} else if _, err := os.Stat(team2load.Filename); err == nil {
		fmt.Printf("filename found .. processing")

		// Create a decoder
		//decoder := json.NewDecoder(decodeFile)

		//// Place to decode into
		//playersFromFile := make(map[int]fflPlayer)

		//// Decode -- We need to pass a pointer

		//teamFromFile := team{}
		//	decoder.Decode(&teamFromFile)

		//if _, err := os.Stat("myplayers.txt"); err == nil {
		//	fmt.Printf("TEAM file exists; processing...\n")

		//	} else {
		//		fmt.Printf("TEAM file DOES NOT exist\n")
		//	}

	}

	return team2load

}
