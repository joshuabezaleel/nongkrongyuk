package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

type SearchCityIDResponse struct {
	LocationSuggestions []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		// CountryID            int    `json:"country_id"`
		// CountryName          string `json:"country_name"`
		// ShouldExperimentWith int    `json:"should_experiment_with"`
		// DiscoveryEnabled     int    `json:"discovery_enabled"`
		// HasNewAdFormat       int    `json:"has_new_ad_format"`
		// IsState              int    `json:"is_state"`
		// StateID              int    `json:"state_id"`
		// StateName            string `json:"state_name"`
		// StateCode            string `json:"state_code"`
	} `json:"location_suggestions"`
	// Status   string `json:"status"`
	// HasMore  int    `json:"has_more"`
	// HasTotal int    `json:"has_total"`
}

func main() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		events, err := bot.ParseRequest(r)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("joshua gateng")).Do(); err != nil {
				// 	log.Print(err)
				// }
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					cityQuery := message.Text
					apiURL := "https://developers.zomato.com/api/v2.1/cities?q=" + cityQuery
					req, err := http.NewRequest("GET", apiURL, nil)
					if err != nil {
						log.Fatal(err)
					}
					req.Header.Set("Accept", "application/json")
					req.Header.Set("User-Key", "d5a2d5a5ba29db0566b65335e27b5801")

					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						log.Fatal(err)
					}
					defer resp.Body.Close()

					var SearchCityIDs SearchCityIDResponse
					if err := json.NewDecoder(resp.Body).Decode(&SearchCityIDs); err != nil {
						log.Println(err)
					}

					if len(SearchCityIDs.LocationSuggestions) == 0 {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("not found")).Do(); err != nil {
							log.Print(err)
						}
					} else {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(strconv.Itoa(SearchCityIDs.LocationSuggestions[0].ID))).Do(); err != nil {
							log.Print(err)
						}
					}
					// for _, SearchCityID := range SearchCityIDs.LocationSuggestions {
					// 	if SearchCityID.ID != nil {
					// 		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("jobel")).Do(); err != nil {
					// 			log.Print(err)
					// 		}
					// 	} else {
					// 		if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("not found")).Do(); err != nil {
					// 			log.Print(err)
					// 		}
					// 	}
					// }
					// if message.Text == "a" {
					// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("jobel")).Do(); err != nil {
					// 		log.Print(err)
					// 	}
					// } else {
					// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
					// 		log.Print(err)
					// 	}
					// }
				}
			}
		}
	})

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
