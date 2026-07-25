package spreadsheet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func updateOrganizerData(organizer *Organizer) error {
	data, err := getNominatimData(organizer)
	if err != nil {
		return err
	}

	lat, err := strconv.ParseFloat(data.Latitude, 64)
	if err != nil {
		return err
	}

	lon, err := strconv.ParseFloat(data.Longitude, 64)
	if err != nil {
		return err
	}

	organizer.Name = data.Name
	organizer.Address = data.Address
	organizer.Latitude = lat
	organizer.Longitude = lon

	return nil
}

func getNominatimData(organizer *Organizer) (NominatimData, error) {
	// build url query
	var queryParts []string

	switch organizer.Type {
	case "venue":
		queryParts = []string{organizer.Name, organizer.City, organizer.Country}
	case "festival":
		queryParts = []string{organizer.City, organizer.Country}
	default:
		queryParts = []string{organizer.City, organizer.Country}
	}

	searchQuery := strings.Join(queryParts, " ")

	params := url.Values{}
	params.Set("q", searchQuery)
	params.Set("format", "json")
	params.Set("accept-language", "en")

	urlQuery := "https://nominatim.openstreetmap.org/search?" + params.Encode()

	// create request and set user agent (nominatim usage policy requires it)
	request, err := http.NewRequest("GET", urlQuery, nil)
	if err != nil {
		return NominatimData{}, fmt.Errorf("error creating request: %s", err.Error())
	}

	request.Header.Set("User-Agent", "jazz-network-api/1.0 (wildewilfried@gmail.com)")

	// get data
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return NominatimData{}, fmt.Errorf("error fetching data from nominatim: %s", err.Error())
	}

	var geodata []NominatimData

	err = json.NewDecoder(response.Body).Decode(&geodata)
	response.Body.Close()
	if err != nil {
		return NominatimData{}, fmt.Errorf("error decoding nominatim response body: %s", err.Error())
	}

	var nominatimData NominatimData
	var unknowType string

	if len(geodata) == 0 {
		return NominatimData{}, fmt.Errorf("no match for %s, %s, %s", organizer.Name, organizer.City, organizer.Country)
	} else if len(geodata) > 1 { // if multiple results get city or town, else just the single result
		for _, d := range geodata {
			if d.Type == "city" || d.Type == "town" {
				nominatimData = d
				break
			} else {
				unknowType = d.Type
			}
		}
	} else {
		nominatimData = geodata[0]
	}

	if nominatimData == (NominatimData{}) {
		return nominatimData, fmt.Errorf("failed to find the correct address type %s for %s, %s, %s", unknowType, organizer.Name, organizer.City, organizer.Country)
	}

	return nominatimData, nil
}
