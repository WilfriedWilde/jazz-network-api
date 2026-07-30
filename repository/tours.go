package repository

var toursTableStatement = `
	CREATE TABLE IF NOT EXISTS tours (
		id INTEGER PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		name TEXT NOT NULL
);`
