package main

type team struct {
	Name       string
	Filename   string
	Bootstrap  string
	Owner      string
	TotalScore int
	Players    map[string]fflPlayer
}
