CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    'date' CHAR(8) NOT NULL,
    title TEXT NOT NULL,
    comment TEXT,
    'repeat' VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS scheduler_date_index ON scheduler ('date');