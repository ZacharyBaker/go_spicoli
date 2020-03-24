package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type item struct {
	Channel string `json:"channel"`
}
type event struct {
	Type    string `json:"type"`
	Text    string `json:"text"`
	EventTs string `json:"event_ts"`
	User    string `json:"user"`
	Channel string `json:"channel"`
}

// Event comment
type Event struct {
	Token     string `json:"token"`
	TeamID    string `json:"team_id"`
	APIAppID  string `json:"api_app_id"`
	Event     event  `json:"event"`
	Type      string `json:"type"`
	EventID   string `json:"event_id"`
	EventTime int    `json:"event_time"`
}

type response struct {
	text    string
	channel string
}

func handleEvent(w http.ResponseWriter, r *http.Request) {
	var e Event
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(e)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("☄ HTTP status code returned!"))

	if e.Event.Type == "app_mention" {
		fmt.Println("app_mention is what is happening", e.Event.Text)
		url := "https://slack.com/api/chat.postMessage"
		client := &http.Client{}
		t := os.Getenv("BOT_TOKEN")
		bt := "Bearer " + t

		var jsonStr = []byte(fmt.Sprintf(`
			{
				"text":"hey bud",
				"channel": "%s",
				"blocks": [
					{
						"type": "section",
						"text": {
							"type": "mrkdwn",
							"text": ":surfer: \n You know, I've been thinking..."
						}
					},
					{
						"type": "image",
						"image_url": "https://media.giphy.com/media/yQaYsWfVTPyZW/giphy.gif",
						"alt_text": "Spicoli Philosophizing"
					}
				]
			}`, e.Event.Channel))
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", bt)

		fmt.Println("REQ:::", req)

		resp, err := client.Do(req)
		fmt.Println("RESP", resp)
		if err != nil {
			fmt.Println("err:::", err)
			panic(err)
		}
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Healthy Boy")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("☄ Healthy and happy"))
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/slack/event", handleEvent).Methods("POST")
	router.HandleFunc("/health-check", healthCheck).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
