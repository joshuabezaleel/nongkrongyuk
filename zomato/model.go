package zomato

//Location is rad as fuck
type Location struct {
	Address   string `json:"address"`
	Locality  string `json:"locality"`
	City      string `json:"city"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Zipcode   string `json:"zipcode"`
	CountryID string `json:"country_id"`
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
	ID          string `json:"id"`
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
	NTotalResult int           `json:"results_found"`
	Offset       int           `json:"results_start"`
	NResult      int           `json:"results_shown"`
	Restaurants  []*Restaurant `json:"restaurants"`
}
