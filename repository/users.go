package repository

var usersTableStatement = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		full_name TEXT NOT NULL
);`
