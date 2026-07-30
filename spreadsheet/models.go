package spreadsheet

type Organizer struct {
	Type           string
	Name           string
	Country        string
	City           string
	Address        string
	Websites       []string
	Emails         []string
	PhoneNumbers   []string
	ContactPersons []string
	Schedule       string
	Notes          string
	Latitude       float64
	Longitude      float64
}

type NominatimData struct {
	Name        string  `json:"name"`
	Address     string  `json:"display_name"`
	AddressType string  `json:"addresstype"`
	Type        string  `json:"type"`
	Importance  float64 `json:"importance"`
	Latitude    string  `json:"lat"`
	Longitude   string  `json:"lon"`
}
