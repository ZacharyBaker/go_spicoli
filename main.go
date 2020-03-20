package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type challenge struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

type challengeResponse struct {
	Challenge string `json:"challenge"`
}

type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
	{
		ID:          "2",
		Title:       "Spicoli",
		Description: "tasty waves and a cool buzz",
	},
	{
		ID:          "3",
		Title:       "Three",
		Description: "three three threee",
	},
	{
		ID:          "4",
		Title:       "Four",
		Description: "four four four",
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Welcome home!")
	w.Write([]byte("Chicken boy"))
}

func spicoliHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hey bud, whats your problem")
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

func respondToChallenge(w http.ResponseWriter, r *http.Request) {
	var ch challenge
	err := json.NewDecoder(r.Body).Decode(&ch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := challengeResponse{
		ch.Challenge,
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(resp)
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink).Methods("GET")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/", respondToChallenge).Methods("POST")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
