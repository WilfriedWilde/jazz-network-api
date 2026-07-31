package repository

type Organizer struct {
	ID             int      `json:"id"`
	Type           string   `json:"type"`
	Name           string   `json:"name"`
	Country        string   `json:"country"`
	City           string   `json:"ciy"`
	Address        string   `json:"address"`
	Websites       []string `json:"websites"`
	Emails         []string `json:"emails"`
	PhoneNumbers   []string `json:"phone_numbers"`
	ContactPersons []string `json:"contact_persons"`
	Schedule       string   `json:"schedule"`
	Notes          string   `json:"notes"`
	Latitude       float64  `json:"latitude"`
	Longitude      float64  `json:"longitude"`
}
