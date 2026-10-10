package model

type Tour struct {
	ID    int64      `json:"id"`
	Name  string     `json:"name"`
	Start string     `json:"start"`
	End   string     `json:"end"`
	Stops []TourStop `json:"stops"`
}

type TourInput struct {
	Name  string `json:"name"`
	Start string `json:"start"`
	End   string `json:"end"`
}

type TourCreate struct {
	UserID int64
	Name   string
	Start  string
	End    string
}
