package repository

var notesTableStatement = `
	CREATE TABLE IF NOT EXISTS organizer_notes (
		id INTEGER PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		organizer_id INTEGER REFERENCES organizers(id),
		body TEXT NOT NULL,
		created_at TEXT
);`
