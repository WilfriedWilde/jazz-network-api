package repository

var stopsTableStatement = `
	CREATE TABLE IF NOT EXISTS tour_stops (
		id INTEGER PRIMARY KEY,
		tour_id INTEGER REFERENCES tours(id),
		organizer_id INTEGER REFERENCES organizers(id),
		notes TEXT
)`
