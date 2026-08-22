CREATE TABLE IF NOT EXISTS organizers (
    id INTEGER PRIMARY KEY,
    type TEXT NOT NULL,
    name TEXT NOT NULL,
    country TEXT NOT NULL,
    city TEXT NOT NULL,
    address TEXT NOT NULL,
    websites TEXT,
    emails TEXT,
    phone_numbers TEXT,
    contact_persons TEXT,
    schedule TEXT,
    notes TEXT,
    latitude REAL,
    longitude REAL
);

CREATE TABLE IF NOT EXISTS organizer_contacts (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    organizer_id INTEGER REFERENCES organizers(id),
    method TEXT NOT NULL CHECK (method IN ('email', 'phone', 'in_person')),
    status TEXT CHECK (status IN ('not_tried', 'contacted', 'confirmed', 'declined')),
    last_attempted_at TEXT,
    notes TEXT
)

CREATE TABLE IF NOT EXISTS organizer_notes (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    organizer_id INTEGER REFERENCES organizers(id),
    body TEXT NOT NULL,
    created_at TEXT
)

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    full_name TEXT NOT NULL
)

CREATE TABLE IF NOT EXISTS tours (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    name TEXT NOT NULL
)

CREATE TABLE IF NOT EXISTS tour_stops (
    id INTEGER PRIMARY KEY,
    tour_id INTEGER REFERENCES tours(id),
    organizer_id INTEGER REFERENCES organizers(id),
    notes TEXT
)