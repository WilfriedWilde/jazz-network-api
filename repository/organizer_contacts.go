package repository

var contactsTableStatement = `
	CREATE TABLE IF NOT EXISTS organizer_contacts (
		id INTEGER PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		organizer_id INTEGER REFERENCES organizers(id),
		method TEXT NOT NULL CHECK (method IN ('email', 'phone', 'in_person')),
		status TEXT CHECK (status IN ('not_tried', 'contacted', 'confirmed', 'declined')),
		last_attempted_at TEXT,
		notes TEXT
);`
