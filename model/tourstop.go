package model

type TourStop struct {
	ID          int64   `json:"id"`
	TourID      int64   `json:"tour_id"`
	OrganizerID *int64  `json:"organizer_id"`
	Type        string  `json:"type"`
	StopOrder   int     `json:"stop_order"`
	Notes       *string `json:"notes"`
}

type TourStopCreate struct {
	TourID      int64
	OrganizerID int64
	Type        string
	StopOrder   int
	Notes       string
}
