package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type original struct {
	URL string `json:"url"`
}

type images struct {
	Original original `json:"original"`
}

type gif struct {
	URL    string `json:"url"`
	Slug   string `json:"slug"`
	Images images `json:"images"`
}

type gifResponse struct {
	Data []gif `json:"data"`
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

func getRandomGif(t string) string {
	// get token from heroku for the api
	to := os.Getenv("GIF_API_TOKEN")

	// change the event text to work in the api request
	// remove "<@UVBE8EDMZ> " from string as well (including first space)
	r := regexp.MustCompile(`<(.+)>\s+`)
	q := r.ReplaceAllString(t, "")

	// replace any spaces with "+" signs
	q = strings.ReplaceAll(q, " ", "+")

	// make the call to the giphy api
	url := fmt.Sprintf("http://api.giphy.com/v1/gifs/search?api_key=%s&limit=10&q=%s", to, q)

	fmt.Println("url:::", url)

	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("err::: gif::: GET request:::", err)
		panic(err)
	}

	// capture the response
	var gr gifResponse
	gerr := json.NewDecoder(resp.Body).Decode(&gr)
	if gerr != nil {
		fmt.Println("err::: gif::: response decoding:::", err)
		panic(err)
	}

	fmt.Println("response from gif api :::", gr)

	// check how many are in the array
	// choose one of those returned randomly
	fmt.Println(len(gr.Data))

	l := len(gr.Data)

	rand.Seed(time.Now().UnixNano())
	min := 0
	max := l - 1
	ra := rand.Intn(max-min+1) + min

	fmt.Println("random number::", ra)

	// return the url
	g := gr.Data[ra]

	return g.Images.Original.URL
}

func handleAppMention(e Event) {
	gif := getRandomGif(e.Event.Text)
	fmt.Println("gif", gif)

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
						"text": ":surfer: _Powered By GIPHY_ thanks bud"
					}
				},
				{
					"type": "image",
					"image_url": "%s",
					"alt_text": "Spicoli Philosophizing"
				}
			]
		}`, e.Event.Channel, gif))
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bt)

	_, err := client.Do(req)
	if err != nil {
		fmt.Println("err:::", err)
		panic(err)
	}
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
		handleAppMention(e)
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
