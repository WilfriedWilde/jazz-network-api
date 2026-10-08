package model

type TourStop struct {
	ID          int64
	TourID      int64
	OrganizerID int64
	Type        string
	StopOrder   int
	Notes       string
}
