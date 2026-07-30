package spreadsheet

type NominatimData struct {
	Name        string  `json:"name"`
	Address     string  `json:"display_name"`
	AddressType string  `json:"addresstype"`
	Type        string  `json:"type"`
	Importance  float64 `json:"importance"`
	Latitude    string  `json:"lat"`
	Longitude   string  `json:"lon"`
}
