package model

type Tour struct {
	ID          int64      `json:"id"`
	TourName    string     `json:"tour_name"`
	ProjectName string     `json:"project_name"`
	Start       string     `json:"start"`
	End         string     `json:"end"`
	Stops       []TourStop `json:"stops"`
	Contacts    []Contact  `json:"contacts"`
}

type TourInput struct {
	TourName    string `json:"tour_name"`
	ProjectName string `json:"project_name"`
	Start       string `json:"start"`
	End         string `json:"end"`
}

type TourUpdate struct {
	TourName    *string `json:"tour_name"`
	ProjectName *string `json:"project_name"`
	Start       *string `json:"start"`
	End         *string `json:"end"`
}

type TourCreate struct {
	UserID      int64
	TourName    string
	ProjectName string
	Start       string
	End         string
}
