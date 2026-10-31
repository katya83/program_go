-- Миграция 003: пользователи, сессии и привязка трат к пользователю.

CREATE TABLE IF NOT EXISTS users (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);
 
CREATE TABLE IF NOT EXISTS sessions (
    session_id TEXT PRIMARY KEY,
    user_id    INTEGER NOT NULL REFERENCES users(id),
    expires_at INTEGER NOT NULL
);
 
ALTER TABLE expenses ADD COLUMN user_id INTEGER REFERENCES users(id);
