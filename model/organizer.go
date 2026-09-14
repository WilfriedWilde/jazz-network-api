package model

import "strings"

type Organizer struct {
	ID               int64
	Type             string
	Name             string
	Country          string
	City             string
	Address          string
	Websites         string
	Emails           string
	PhoneNumbers     string
	ContactPersons   string
	Schedule         string
	Notes            string
	Latitude         float64
	Longitude        float64
	PublicVisibility bool
}

type OrganizerInput struct {
	Type             string   `json:"type"`
	Name             string   `json:"name"`
	Country          string   `json:"country"`
	City             string   `json:"city"`
	Address          string   `json:"address"`
	Websites         []string `json:"websites"`
	Emails           []string `json:"emails"`
	PhoneNumbers     []string `json:"phone_numbers"`
	ContactPersons   []string `json:"contact_persons"`
	Schedule         string   `json:"schedule"`
	Notes            string   `json:"notes"`
	Latitude         *float64 `json:"latitude"`
	Longitude        *float64 `json:"longitude"`
	PublicVisibility *bool    `json:"public_visibility"`
}

type OrganizerCreate struct {
	UserID           *int64
	Type             string
	Name             string
	Country          string
	City             string
	Address          string
	Websites         string
	Emails           string
	PhoneNumbers     string
	ContactPersons   string
	Schedule         string
	Notes            string
	Latitude         float64
	Longitude        float64
	PublicVisibility bool
}

type OrganizerJSON struct {
	ID             int64    `json:"id"`
	Type           string   `json:"type"`
	Name           string   `json:"name"`
	Country        string   `json:"country"`
	City           string   `json:"city"`
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

func (o *OrganizerInput) ToOrganizerCreate() OrganizerCreate {
	return OrganizerCreate{
		Type:             o.Type,
		Name:             o.Name,
		Country:          o.Country,
		City:             o.City,
		Address:          o.Address,
		Websites:         strings.Join(o.Websites, ";"),
		Emails:           strings.Join(o.Emails, ";"),
		PhoneNumbers:     strings.Join(o.PhoneNumbers, ";"),
		ContactPersons:   strings.Join(o.ContactPersons, ";"),
		Schedule:         o.Schedule,
		Notes:            o.Notes,
		Latitude:         *o.Latitude,
		Longitude:        *o.Longitude,
		PublicVisibility: *o.PublicVisibility,
	}
}

func (o *Organizer) ToOrganizerJSON() OrganizerJSON {
	return OrganizerJSON{
		ID:             o.ID,
		Type:           o.Type,
		Name:           o.Name,
		Country:        o.Country,
		City:           o.City,
		Address:        o.Address,
		Websites:       strings.Split(o.Websites, ";"),
		Emails:         strings.Split(o.Emails, ";"),
		PhoneNumbers:   strings.Split(o.PhoneNumbers, ";"),
		ContactPersons: strings.Split(o.ContactPersons, ";"),
		Schedule:       o.Schedule,
		Notes:          o.Notes,
		Latitude:       o.Latitude,
		Longitude:      o.Longitude,
	}
}
