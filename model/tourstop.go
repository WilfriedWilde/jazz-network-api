package model

type TourStop struct {
	ID          int64   `json:"id"`
	TourID      int64   `json:"tour_id"`
	OrganizerID *int64  `json:"organizer_id"`
	Type        string  `json:"type"`
	Date        string  `json:"date"`
	Country     *string `json:"country"`
	City        *string `json:"city"`
	Notes       *string `json:"notes"`
}

type TourStopCreate struct {
	TourID      int64
	OrganizerID int64
	Type        string
	Date        string
	Country     string
	City        string
	Notes       string
}

type TourStopUpdate struct {
	OrganizerID *int64  `json:"organizer_id"`
	Type        *string `json:"type"`
	Country     *string `json:"country"`
	City        *string `json:"city"`
	Notes       *string `json:"notes"`
}
