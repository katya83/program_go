package repository

import (
	"database/sql"
	"time"

	"project/internal/models"
)

// UserRepository - работа с пользователями и их сессиями.
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser сохраняет нового пользователя (email должен быть уникальным).
func (r *UserRepository) CreateUser(email, passwordHash string) error {
	_, err := r.db.Exec(
		"INSERT INTO users (email, password_hash) VALUES (?, ?)",
		email, passwordHash,
	)
	return err
}

// GetByEmail ищет пользователя по email.
func (r *UserRepository) GetByEmail(email string) (models.User, error) {
	var u models.User
	err := r.db.QueryRow(
		"SELECT id, email, password_hash FROM users WHERE email = ?", email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash)
	return u, err
}

// CreateSession сохраняет новую сессию. expiresAt - время окончания (секунды Unix-времени).
func (r *UserRepository) CreateSession(sessionID string, userID int, expiresAt int64) error {
	_, err := r.db.Exec(
		"INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID, userID, expiresAt,
	)
	return err
}

// GetUserIDBySession возвращает номер пользователя, если сессия существует
// и срок её действия ещё не истёк. Иначе вернёт ошибку.
func (r *UserRepository) GetUserIDBySession(sessionID string) (int, error) {
	var userID int
	err := r.db.QueryRow(
		"SELECT user_id FROM sessions WHERE session_id = ? AND expires_at > ?",
		sessionID, time.Now().Unix(),
	).Scan(&userID)
	return userID, err
}

// DeleteSession удаляет сессию (выход из системы).
func (r *UserRepository) DeleteSession(sessionID string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}
