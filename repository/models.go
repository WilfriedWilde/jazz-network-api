package repository

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

type Select struct {
	ColumnList []string
	Distinct   bool
	JoinTables []string
	RowFilter  map[string]any
}
