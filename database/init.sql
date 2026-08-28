CREATE TABLE IF NOT EXISTS organizers (
    id INTEGER PRIMARY KEY,
    type TEXT NOT NULL CHECK (length(trim(type)) > 0),
    name TEXT NOT NULL CHECK (length(trim(name)) > 0),
    country TEXT NOT NULL CHECK (length(trim(country)) > 0),
    city TEXT NOT NULL CHECK (length(trim(city)) > 0),
    address TEXT NOT NULL CHECK (length(trim(address)) > 0),
    websites TEXT,
    emails TEXT,
    phone_numbers TEXT,
    contact_persons TEXT,
    schedule TEXT,
    notes TEXT,
    latitude REAL,
    longitude REAL
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL CHECK (length(trim(username)) > 0),
    email TEXT NOT NULL CHECK (length(trim(email)) > 0),
    password_hash TEXT NOT NULL (length(trim(password_hash)) > 0)
);

CREATE TABLE IF NOT EXISTS sessions (
    uuid TEXT NOT NULL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS organizer_contacts (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    organizer_id INTEGER REFERENCES organizers(id),
    method TEXT NOT NULL CHECK (method IN ('email', 'phone', 'in_person')),
    status TEXT CHECK (status IN ('not_tried', 'contacted', 'confirmed', 'declined')),
    last_attempted_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,
    UNIQUE (user_id, organizer_id, method)
);

CREATE TABLE IF NOT EXISTS user_notes (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    organizer_id INTEGER REFERENCES organizers(id),
    body TEXT NOT NULL,
    created_at TEXT
);

CREATE TABLE IF NOT EXISTS tours (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    name TEXT NOT NULL CHECK (length(trim(name)) > 0)
);

CREATE TABLE IF NOT EXISTS tour_stops (
    id INTEGER PRIMARY KEY,
    tour_id INTEGER REFERENCES tours(id),
    organizer_id INTEGER REFERENCES organizers(id),
    stop_order INTEGER NOT NULL, 
    notes TEXT,
    UNIQUE (tour_id, organizer_id)
);