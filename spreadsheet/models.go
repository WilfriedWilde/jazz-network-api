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
	Status         bool
}

type NominatimData struct {
	Name      string `json:"name"`
	Address   string `json:"display_name"`
	Type      string `json:"addresstype"`
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}
