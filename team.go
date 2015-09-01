package main

import (
	"math/rand"
	"time"
)

type team struct {
	Name       string
	Filename   string
	Bootstrap  string
	Owner      string
	TotalScore int
	Players    map[string]fflPlayer
}

type fflPlayer struct {
	Index            string
	ID               int
	Name             string
	Team             string
	Value            string
	Position         string
	TotalPoints      int
	URL              string
	ScoringWeekStart int
	ScoringWeekEnd   int
	WeekStats        map[int]weekStats
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

func processTeam(t *team) {
	for k, v := range t.Players {
		duration := time.Duration(random(0, 2)) * time.Second
		time.Sleep(duration)
		//log.Printf("Pausing for %v\n", duration)
		scrapePlayer(&v)
		calculatePlayerScore(&v)
		t.Players[k] = v
		t.TotalScore = t.TotalScore + v.TotalPoints
		//log.Printf("%v %v %v %v %v\n", t.Players[k].Name, t.Players[k].TotalPoints, t.Players[k].ScoringWeekStart, t.Players[k].ScoringWeekEnd, t.Players[k].WeekStats)
	}

}

func calculatePlayerScore(p *fflPlayer) {
	totalScore := 0
	if p.ID == -1 {
		// missing player
		totalScore = p.TotalPoints
	} else if p.ScoringWeekStart == 0 && p.ScoringWeekEnd == 99 {
		for _, v := range p.WeekStats {
			totalScore = totalScore + v.Points
			//log.Printf("Player:%v %v %v\n", p.Name, v.Vs, v.Points)
		}
		p.TotalPoints = totalScore
	} else {
		for _, v := range p.WeekStats {
			if v.Week >= p.ScoringWeekStart && v.Week <= p.ScoringWeekEnd {
				totalScore = totalScore + v.Points
			}
			//log.Printf("%v %v %v\n", p.Name, v.Vs, v.Points)
		}
		p.TotalPoints = totalScore
	}
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
