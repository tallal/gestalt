package main

type team struct {
	name       string
	filename   string
	owner      string
	totalScore int
	players    map[int]fflPlayer
}
