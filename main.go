package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joshuabezaleel/nongkrongyuk/zomato"
	"github.com/line/line-bot-sdk-go/linebot"
)

//CarouselColumn is rad as fuck
type CarouselColumn struct {
	imageURL string
	title    string
	text     string
	labelURI string
	linkURI  string
}

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
							linebot.NewMessageTemplateAction("label MTA", "text MTA")),
					)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("carousel alt text", template)).Do(); err != nil {
						log.Print(err)
					}
				case *linebot.LocationMessage:
					lat := message.Latitude
					lon := message.Longitude
					restaurants, err := zomatoService.SearchRestaurantsByLatLong(lat, lon, 1, 5)

					// var lineMessage string
					var CarouselColumnFiller []*linebot.CarouselColumn
					// var URITemplate linebot.URITemplateAction
					// linebot.NewURITemplateAction(label, uri)
					for i, restaurant := range restaurants {
						if i > 5 {
							break
						}
						OneCarouselColumn := new(linebot.CarouselColumn)
						OneCarouselColumn.ThumbnailImageURL = "https://user-images.githubusercontent.com/7043511/31583356-630ca11c-b1c4-11e7-8109-16228f8a5c0b.png"
						OneCarouselColumn.Title = restaurant.Name
						OneCarouselColumn.Text = restaurant.Cuisines
						// URITemplate.Label = "Go to Zomato!"
						// URITemplate.URI = "" + restaurant.URL
						// OneCarouselColumn.Actions = URITemplate
						CarouselColumnFiller = append(CarouselColumnFiller, OneCarouselColumn)
					}
					// var zaky [5]linebot.CarouselTemplate

					template := linebot.NewCarouselTemplate(CarouselColumnFiller...)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("carousel alt text", template)).Do(); err != nil {
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
