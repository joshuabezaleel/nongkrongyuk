package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

//SearchCityIDResponse is rad as fuck
type SearchCityIDResponse struct {
	LocationSuggestions []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"location_suggestions"`
}

//SearchByLatLongResponse is rad as fuck
type SearchByLatLongResponse struct {
	Restaurants []struct {
		Restaurant struct {
			R struct {
				ResID int `json:"res_id"`
			} `json:"R"`
			ID       string `json:"id"`
			Name     string `json:"name"`
			Location struct {
				City string `json:"city"`
			} `json:"location"`
			UserRating struct {
			} `json:"user_rating"`
		} `json:"restaurant"`
	} `json:"restaurants"`
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
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					// cityQuery := message.Text
					// apiURL := "https://developers.zomato.com/api/v2.1/cities?q=" + cityQuery
					// req, err := http.NewRequest("GET", apiURL, nil)
					// if err != nil {
					// 	log.Fatal(err)
					// }
					// req.Header.Set("Accept", "application/json")
					// req.Header.Set("User-Key", "d5a2d5a5ba29db0566b65335e27b5801")
					//
					// resp, err := http.DefaultClient.Do(req)
					// if err != nil {
					// 	log.Fatal(err)
					// }
					// defer resp.Body.Close()
					//
					// var SearchCityIDs SearchCityIDResponse
					// if err = json.NewDecoder(resp.Body).Decode(&SearchCityIDs); err != nil {
					// 	log.Println(err)
					// }
					//
					// if len(SearchCityIDs.LocationSuggestions) == 0 {
					// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("not found")).Do(); err != nil {
					// 		log.Print(err)
					// 	}
					// } else {
					// 	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(strconv.Itoa(SearchCityIDs.LocationSuggestions[0].ID))).Do(); err != nil {
					// 		log.Print(err)
					// 	}
					// }
					imageURL := "https://user-images.githubusercontent.com/7043511/31583356-630ca11c-b1c4-11e7-8109-16228f8a5c0b.png"
					template := linebot.NewCarouselTemplate(
						linebot.NewCarouselColumn(imageURL, "title carousel", "text carousel",
							linebot.NewMessageTemplateAction("label message template action", "text message template action"))
					)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("carousel alt text", template).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.LocationMessage:
					lat := message.Latitude
					lon := message.Longitude
					apiURL := "https://developers.zomato.com/api/v2.1/search?lat=" + fmt.Sprint(lat) + "&lon=" + fmt.Sprint(lon)
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

					var SearchByLatLongs SearchByLatLongResponse
					if err = json.NewDecoder(resp.Body).Decode(&SearchByLatLongs); err != nil {
						log.Println(err)
					}

					var lineMessage string
					for _, SearchByLatLong := range SearchByLatLongs.Restaurants {
						lineMessage = lineMessage + SearchByLatLong.Restaurant.Name + "\n"
					}
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(lineMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
