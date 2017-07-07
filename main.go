package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var timeToClap map[string]time.Time

func PostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside post handler")
	params := mux.Vars(r)
	teamName := params["team_name"]

	timeToClap[teamName] = time.Now().Add(time.Duration(5) * time.Second)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside GET handler")
	params := mux.Vars(r)
	teamName := params["team_name"]

	timeToClapForKubo, timerStarted := timeToClap[teamName]
	timeRemaining := time.Duration(0)

	if timerStarted {
		timeRemaining = time.Until(timeToClapForKubo)
	}

	if timeRemaining < 0 {
		delete(timeToClap, teamName)
		timerStarted = false
		timeRemaining = time.Duration(0)
	}

	fmt.Fprintf(w, "{ timer_started: %t, time_remaining: %f }", timerStarted, timeRemaining.Seconds())
}

func main() {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/{team_name}", PostHandler).Methods("POST")
	rtr.HandleFunc("/{team_name}", GetHandler).Methods("GET")

	http.Handle("/", rtr)

	timeToClap = map[string]time.Time{}

	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
