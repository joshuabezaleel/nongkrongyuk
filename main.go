package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joshuabezaleel/nongkrongyuk/zomato"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	zomatoService := zomato.NewService("ZOMATO_API_KEY")

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
					cities, err := zomatoService.SearchCityByName(cityQuery)
					if err != nil {
						log.Printf("Error while SearchCityByName: %s", err.Error())
						break
					}

					if len(cities) == 0 {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("not found")).Do(); err != nil {
							log.Print(err)
						}
					} else {
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(cities[0].ID)).Do(); err != nil {
							log.Print(err)
						}
					}
				case *linebot.LocationMessage:
					lat := message.Latitude
					lon := message.Longitude
					restaurants, err := zomatoService.SearchRestaurantsByLatLong(lat, lon)

					var lineMessage string
					for _, restaurant := range restaurants {
						lineMessage = lineMessage + restaurant.Name + "\n"
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
