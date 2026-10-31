package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"project/internal/repository"
)

const sessionCookie = "session_id" // имя cookie с идентификатором сессии

// authPage - данные для страниц входа и регистрации.
type authPage struct {
	Error string
	Email string // чтобы не терять введённый email при ошибке
}

// AuthHandler - регистрация, вход и выход.
type AuthHandler struct {
	users *repository.UserRepository
}

func NewAuthHandler(users *repository.UserRepository) *AuthHandler {
	return &AuthHandler{users: users}
}

// newSessionID создаёт случайный идентификатор сессии (64 символа).
func newSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil { // crypto/rand - криптостойкий генератор
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// ShowRegister показывает форму регистрации.
func (h *AuthHandler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "register.html", authPage{})
}

// Register создаёт нового пользователя.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	password := r.FormValue("password")

	if email == "" || len(password) < 6 || len(password) > 72 {
		w.WriteHeader(http.StatusBadRequest)
		tmpl.ExecuteTemplate(w, "register.html", authPage{
			Error: "Введите email и пароль длиной от 6 до 72 символов",
			Email: email,
		})
		return
	}

	// Превращаем пароль в хэш. Сам пароль дальше нигде не сохраняется.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	if err := h.users.CreateUser(email, string(hash)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		tmpl.ExecuteTemplate(w, "register.html", authPage{
			Error: "Не удалось зарегистрироваться: возможно, такой email уже занят",
			Email: email,
		})
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ShowLogin показывает форму входа.
func (h *AuthHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login.html", authPage{})
}

// Login проверяет пароль, создаёт сессию и устанавливает cookie.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	password := r.FormValue("password")

	// Ищем пользователя и сравниваем пароль с хэшем.
	// Сообщение одинаковое и для неверного email, и для неверного пароля.
	user, err := h.users.GetByEmail(email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		tmpl.ExecuteTemplate(w, "login.html", authPage{Error: "Неверный email или пароль", Email: email})
		return
	}

	// Создаём сессию на 24 часа.
	sessionID, err := newSessionID()
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	expires := time.Now().Add(24 * time.Hour)
	if err := h.users.CreateSession(sessionID, user.ID, expires.Unix()); err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	// Отправляем браузеру cookie с идентификатором сессии.
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    sessionID,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,                 // JavaScript страницы не сможет прочитать cookie
		SameSite: http.SameSiteLaxMode, // cookie не отправляется с запросов с других сайтов
	})
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

// Logout удаляет сессию на сервере и cookie в браузере.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookie); err == nil {
		h.users.DeleteSession(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   sessionCookie,
		Value:  "",
		Path:   "/",
		MaxAge: -1, // просим браузер удалить cookie
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
