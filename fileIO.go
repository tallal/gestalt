package main

import (
	"encoding/gob"
	"encoding/json"
	"os"
)

func loadPlayerFile(filename string) map[int]fflPlayer {
	// Open a RO file
	decodeFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer decodeFile.Close()

	// Create a decoder
	decoder := gob.NewDecoder(decodeFile)

	// Place to decode into
	playersFromFile := make(map[int]fflPlayer)

	// Decode -- We need to pass a pointer
	decoder.Decode(&playersFromFile)

	return playersFromFile
}

func savePlayerFile(filename string, plist map[int]fflPlayer) {
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

func loadTeamFile(filename string) team {
	// Open a RO file
	decodeFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer decodeFile.Close()

	// Create a decoder
	decoder := json.NewDecoder(decodeFile)

	//// Place to decode into
	//playersFromFile := make(map[int]fflPlayer)

	//// Decode -- We need to pass a pointer

	teamFromFile := team{}
	decoder.Decode(&teamFromFile)
	return teamFromFile
}
