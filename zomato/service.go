package zomato

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	apiURL = "https://developers.zomato.com/api/v2.1"
)

type Service struct {
	apiKey string
}

func NewService(apiKey string) *Service {
	return &Service{
		apiKey: apiKey,
	}
}

func (s *Service) get(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Key", s.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func (s *Service) SearchCityByName(city string) ([]*City, error) {
	url := fmt.Sprintf("%s/cities?q=%s", apiURL, url.QueryEscape(city))
	response, err := s.get(url)
	if err != nil {
		return nil, err
	}

	var SearchCityIDs SearchCityIDResponse
	if err = json.Unmarshal(response, &SearchCityIDs); err != nil {
		return nil, err
	}

	return SearchCityIDs.LocationSuggestions, nil
}

func (s *Service) SearchRestaurantsByLatLong(latitude, longitude float64) ([]*Restaurant, error) {
	url := fmt.Sprintf("%s/search?lat=%g&lon=%g", apiURL, latitude, longitude)
	response, err := s.get(url)
	if err != nil {
		return nil, err
	}

	var SearchByLatLongs SearchByLatLongResponse
	if err = json.Unmarshal(response, &SearchByLatLongs); err != nil {
		return nil, err
	}

	return SearchByLatLongs.Restaurants, nil
}
