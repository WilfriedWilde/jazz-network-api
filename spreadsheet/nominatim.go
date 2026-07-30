package spreadsheet

import (
	"cmp"
	"encoding/json"
	"fmt"
	r "jazz-network-api/repository"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"
)

var venueType = []string{"music_venue", "restaurant", "hotel", "nightclub", "arts_centre", "theatre", "bar", "pub", "events_venue", "community_centre", "cooperative"}
var placeTypes = []string{"city", "town", "village", "hamlet", "municipality", "borough", "suburb", "park"}

func updateOrganizerData(organizer *r.Organizer) error {
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

	if organizer.Type == "venue" {
		organizer.Name = data.Name
	}

	organizer.Address = data.Address
	organizer.Latitude = lat
	organizer.Longitude = lon

	return nil
}

func getNominatimData(organizer *r.Organizer) (NominatimData, error) {
	// build url query
	var queryParts []string
	params := url.Values{}

	switch organizer.Type {
	case "venue":
		queryParts = []string{organizer.Name, organizer.City, organizer.Country}
	default:
		queryParts = []string{organizer.City, organizer.Country}
		params.Set("featureType", "settlement")
	}

	searchQuery := strings.Join(queryParts, " ")

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
	var typeMatches []NominatimData

	if len(geodata) == 0 {
		return NominatimData{}, fmt.Errorf("no match for %s, %s, %s", organizer.Name, organizer.City, organizer.Country)
	} else if len(geodata) > 1 { // if multiple results get city or town, else just the single result
		for _, d := range geodata {
			switch organizer.Type {
			case "venue":
				if slices.Contains(venueType, d.Type) {
					typeMatches = append(typeMatches, d)
				}
			default:
				if slices.Contains(placeTypes, d.AddressType) {
					typeMatches = append(typeMatches, d)
				}
			}

		}

		// pick type match with highest importance
		if len(typeMatches) > 0 {
			slices.SortFunc(typeMatches, func(a, b NominatimData) int {
				return cmp.Compare(b.Importance, a.Importance)
			})

			nominatimData = typeMatches[0]
		}
	} else {
		nominatimData = geodata[0]
	}

	if nominatimData == (NominatimData{}) {
		return nominatimData, fmt.Errorf("failed to find the correct type for %s, %s, %s", organizer.Name, organizer.City, organizer.Country)
	}

	return nominatimData, nil
}
