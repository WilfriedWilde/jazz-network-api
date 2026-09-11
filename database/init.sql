CREATE TABLE IF NOT EXISTS organizers (
    id INTEGER PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
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
    latitude REAL NOT NULL CHECK (latitude >= -90 AND latitude <= 90),
    longitude REAL NOT NULL CHECK (latitude >= -180 AND latitude <= 180),
    public_visibility BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY,
    username TEXT NOT NULL CHECK (length(trim(username)) > 0),
    email TEXT NOT NULL CHECK (length(trim(email)) > 0),
    password_hash TEXT NOT NULL CHECK (length(trim(password_hash)) > 0),
    access_level INTEGER NOT NULL CHECK (access_level IN (0, 1))
);

CREATE TABLE IF NOT EXISTS sessions (
    uuid TEXT NOT NULL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS organizer_contacts (
    id INTEGER PRIMARY KEY,
    user_organizer_id INTEGER NOT NULL REFERENCES user_organizers(id) ON DELETE CASCADE,
    method TEXT NOT NULL CHECK (method IN ('email', 'phone', 'in_person')),
    status TEXT CHECK (status IN ('not_tried', 'contacted', 'confirmed', 'declined')),
    last_attempted_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,
    UNIQUE (user_organizer_id, method)
);

CREATE TABLE IF NOT EXISTS user_notes (
    id INTEGER PRIMARY KEY,
    user_organizer_id INTEGER NOT NULL REFERENCES user_organizers(id) ON DELETE CASCADE,
    body TEXT NOT NULL,
    created_at TEXT
);

CREATE TABLE IF NOT EXISTS tours (
    id INTEGER PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL CHECK (length(trim(name)) > 0),
    start TEXT NOT NULL CHECK (length(trim(start)) > 0),
    end TEXT NOT NULL CHECK (length(trim(end)) > 0)
);

CREATE TABLE IF NOT EXISTS tour_stops (
    id INTEGER PRIMARY KEY,
    tour_id INTEGER NOT NULL REFERENCES tours(id),
    organizer_id INTEGER REFERENCES organizers(id),
    type TEXT NOT NULL CHECK (type IN ('gig', 'free_day', 'travel_day')),
    stop_order INTEGER NOT NULL,
    notes TEXT,
    CHECK (
      (type = 'gig' AND organizer_id IS NOT NULL)
      OR (type != 'gig' AND organizer_id IS NULL)
    ),
    UNIQUE (tour_id, stop_order)
);

CREATE TABLE IF NOT EXISTS tour_stop_finances (
  id INTEGER PRIMARY KEY,
  tour_stop_id INTEGER NOT NULL REFERENCES tour_stops(id),
  category TEXT NOT NULL CHECK (category IN ('gig_fee', 'merch', 'food', 'transport', 'other')),
  amount REAL NOT NULL DEFAULT 0,
  notes TEXT
);

CREATE TABLE IF NOT EXISTS travels (
    id INTEGER PRIMARY KEY,
    from_stop_id INTEGER NOT NULL REFERENCES tour_stops(id),
    to_stop_id INTEGER NOT NULL REFERENCES tour_stops(id),
    mean TEXT NOT NULL CHECK (mean IN ('car', 'taxi', 'bus', 'train', 'plane')),
    cost REAL NOT NULL CHECK (cost >= 0),
    booking_status TEXT CHECK (booking_status IN ('pending', 'booked')),
    notes TEXT
);

CREATE TABLE IF NOT EXISTS accommodations (
    id INTEGER PRIMARY KEY,
    tour_stop_id INTEGER NOT NULL REFERENCES tour_stops(id),
    type TEXT NOT NULL CHECK (type IN ('hotel', 'hostel', 'airbnb', 'guesthouse', 'apartment', 'host_family', 'other')),
    cost REAL NOT NULL CHECK (cost >= 0),
    booking_status TEXT CHECK (booking_status IN ('pending', 'booked')),
    name TEXT,
    address TEXT,
    email TEXT,
    phone_number TEXT,
    notes TEXT
);