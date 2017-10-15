package zomato

import (
	"encoding/json"
)

//Location is rad as fuck
type Location struct {
	Address   string `json:"address"`
	Locality  string `json:"locality"`
	City      string `json:"city"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Zipcode   string `json:"zipcode"`
	CountryID int    `json:"country_id"`
}

//Restaurant is rad as fuck
type Restaurant struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Location Location `json:"location"`
	URL      string   `json:"url"`
	Cuisines string   `json:"cuisines"`
}

//City is rad as fuck
type City struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	CountryID   string `json:"country_id"`
	CountryName string `json:"country_name"`
}

//SearchCityIDResponse is rad as fuck
type SearchCityIDResponse struct {
	LocationSuggestions []*City `json:"location_suggestions"`
}

//SearchByLatLongResponse is rad as fuck
type SearchByLatLongResponse struct {
	NTotalResult int
	Offset       int
	NResult      int
	Restaurants  []*Restaurant
}

type searchByLatLongResponseFormat struct {
	NTotalResult int `json:"results_found"`
	Offset       int `json:"results_start"`
	NResult      int `json:"results_shown"`
	Restaurants  []struct {
		Restaurant *Restaurant `json:"restaurant"`
	} `json:"restaurants"`
}

func (r *SearchByLatLongResponse) UnmarshalJSON(data []byte) error {
	var formatted searchByLatLongResponseFormat
	err := json.Unmarshal(data, &formatted)
	if err != nil {
		return err
	}
	r.NTotalResult = formatted.NTotalResult
	r.Offset = formatted.Offset
	r.NResult = formatted.NResult
	r.Restaurants = make([]*Restaurant, 0)
	for _, restaurant := range formatted.Restaurants {
		r.Restaurants = append(r.Restaurants, restaurant.Restaurant)
	}
	return nil
}
