package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

var league = make([]team, 16)
var fileInfo string

func main() {

	inf, err := os.Stat("./data/league.gob")
	if os.IsNotExist(err) {
		fmt.Println("ERROR: League GOB file does not exist...")
		os.Exit(1)
	}
	fileInfo = inf.ModTime().Format(http.TimeFormat)

	fmt.Println("Loading league GOB file...")
	league = loadLeagueFile("./data/league.gob")

	r := mux.NewRouter().StrictSlash(false)
	//r.NotFoundHandler = http.HandlerFunc(notFound)
	r.HandleFunc("/Team/{id}", teamHandler)
	r.HandleFunc("/", homeHandler)

	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)

	log.Println("Starting server on :3000")
	http.Handle("/", r)
	http.ListenAndServe(":3000", r)
}

type pageData struct {
	LastUpdated string
	Items       []team
}

func processByTemplate(rw http.ResponseWriter, r *http.Request) {
	t, err := template.New("home.html").ParseFiles("./tmpl/home.html")
	if err != nil {
		log.Println(err)
	}
	var data pageData

	data.LastUpdated = "19-Oct-2015@23:11"
	data.Items = league[:16]

	err = t.ExecuteTemplate(rw, "home.html", data) // merge & serve
	if err != nil {
		log.Println(err)
	}
}

func homeHandler(rw http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(rw, "<!DOCTYPE html>\n<html>\n<head>")
	fmt.Fprintln(rw, "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">")
	fmt.Fprintln(rw, "<link href=\"/static/site.css\" rel=\"stylesheet\">")
	fmt.Fprintln(rw, "<link rel=\"stylesheet\" href=\"http://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.css\">")
	fmt.Fprintln(rw, "<script src=\"http://code.jquery.com/jquery-1.11.3.min.js\"></script>")
	fmt.Fprintln(rw, "<script src=\"http://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.js\"></script>")
	fmt.Fprintln(rw, "</head>\n<body>")
	fmt.Fprintln(rw, "<div data-role=\"page\" id=\"pageone\">\n<div data-role=\"header\"><h1>FFL 2015/2016</h1></div>")

	fmt.Fprintf(rw, "<div class='smallSubheader'>Updated on:%s</div>\n", fileInfo)
	fmt.Fprintln(rw, "<table id=\"league\">")
	fmt.Fprintln(rw, "<tr>")
	fmt.Fprintln(rw, "<th>TEAM</th><th>OWNER</th><th>TOTAL POINTS</th>")
	fmt.Fprintln(rw, "</tr>")
	counter := 1
	for _, v := range league {
		if v.Owner != "" {

			fmt.Fprintf(rw, "<tr><td")
			if counter%2 == 0 {
				fmt.Fprintf(rw, " class=\"alt\"")
			}
			fmt.Fprintf(rw, "><a href=\"/Team/%s\">%s</a></td><td>%s</td><td>%d</td></tr>\n", v.Owner, v.Name, v.Owner, v.TotalScore)
			counter++
		}
	}
	fmt.Fprintln(rw, "</table>")
	fmt.Fprintln(rw, "<div data-role=\"footer\">\n<h4>Please report all errors to Tal, Mandeep & Suleman.\n For now the spreadsheet is the golden source</h4>\n</div>\n</div>\n</body>\n</html>")

}

func teamHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintln(rw, "<!DOCTYPE html>\n<html>\n<head>")
	fmt.Fprintln(rw, "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">")
	fmt.Fprintln(rw, "<link href=\"/static/site.css\" rel=\"stylesheet\">")
	fmt.Fprintln(rw, "<link rel=\"stylesheet\" href=\"http://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.css\">")
	fmt.Fprintln(rw, "<script src=\"http://code.jquery.com/jquery-1.11.3.min.js\"></script>")
	fmt.Fprintln(rw, "<script src=\"http://code.jquery.com/mobile/1.4.5/jquery.mobile-1.4.5.min.js\"></script>")
	fmt.Fprintln(rw, "</head>\n<body>")
	team, err := getTeam(vars["id"])
	if err == nil {
		fmt.Fprintf(rw, "<div data-role=\"page\" id=\"pageone\">\n<div data-role=\"header\"><h1>FFL Team - %s</h1></div>\n", vars["id"])
		fmt.Fprintf(rw, "<div class='smallSubheader'>Updated on:%s</div>\n", fileInfo)
		fmt.Fprintln(rw, "<p><a href=\"/\">BACK</a></p>")
		fmt.Fprintln(rw, "<table id=\"league\">")
		fmt.Fprintln(rw, "<tr><th>Player</th><th>Team</th><th>Total</th><th>Week Stats</th></tr>")

		processTeamByPos("Goalkeeper", team, rw)
		fmt.Fprintln(rw, "<tr><td colspan=5>Defence</td></tr>")
		processTeamByPos("Defender", team, rw)
		fmt.Fprintln(rw, "<tr><td colspan=5>Midfield</td></tr>")
		processTeamByPos("Midfielder", team, rw)
		fmt.Fprintln(rw, "<tr><td colspan=5>Attack</td></tr>")
		processTeamByPos("Striker", team, rw)

		fmt.Fprintln(rw, "</table>")

	} else {
		fmt.Fprintln(rw, "<h1>FFL 2015/2016</h1>")
		fmt.Fprintf(rw, "<h1>TEAM %s NOT FOUND - PLEASE RETURN</h1>\n", vars["id"])
	}

	fmt.Fprintln(rw, "<p><a href=\"/\">BACK</a></p>")
	fmt.Fprintln(rw, "<div data-role=\"footer\">\n<h4>Please report all errors to Tal,Mandeep & Suleman.</h4>\n</div>\n</div>\n</body>\n</html>")
}

func processTeamByPos(position string, t team, rw http.ResponseWriter) {
	for _, v := range getPlayerByPosition(position, t) {
		fmt.Fprintf(rw, "<tr><td>%s</td><td>%s</td><td>%d</td>", strings.Title(v.Name), v.Team, v.TotalPoints)
		fmt.Fprintf(rw, "<td><div data-role=\"collapsible\"><h1>Expand</h1><p><table>")
		fmt.Fprintf(rw, "<tr><th>Wk</th><th>Against</th><th>Points</th></tr>")

		var keys []int
		for k := range v.WeekStats {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, k := range keys {
			//if v.ScoringWeekStart v.WeekStats[k].Week
			//if player i
			if v.ScoringWeekStart == 0 && v.ScoringWeekEnd == 99 {
				fmt.Fprintf(rw, "<tr><td>%d</td><td>%s</td><td>%d</td></tr>", v.WeekStats[k].Week, v.WeekStats[k].Vs, v.WeekStats[k].Points)
			} else if v.ScoringWeekStart <= v.WeekStats[k].Week && v.ScoringWeekEnd >= v.WeekStats[k].Week {
				fmt.Fprintf(rw, "<tr><td>%d</td><td>%s</td><td>%d</td></tr>", v.WeekStats[k].Week, v.WeekStats[k].Vs, v.WeekStats[k].Points)
			} else {
				fmt.Fprintf(rw, "<tr><td class='nonscoring'>%d</td><td class='nonscoring'>%s</td><td class='nonscoring'>%d</td></tr>", v.WeekStats[k].Week, v.WeekStats[k].Vs, v.WeekStats[k].Points)
			}
		}

		fmt.Fprintln(rw, "</table></p></div></td>\n</tr>")
	}
}

func getPlayerByPosition(position string, t team) (players []fflPlayer) {
	//players = make(map[string]fflPlayer)
	for _, v := range t.Players {
		if v.Position == position {
			players = append(players, v)
		}
	}
	return players
}

func getTeam(managerName string) (t team, err error) {
	for _, v := range league {

		if v.Owner == managerName {

			return v, nil
		}
	}
	return team{}, errors.New("can't find the team selected")
}
