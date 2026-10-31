package handlers

import (
	"context"
	"net/http"
)

// contextKey - свой тип ключа, чтобы он не пересёкся с ключами других пакетов.
type contextKey string

const userIDKey contextKey = "userID"

// RequireAuth - middleware: оборачивает обработчик и пропускает дальше
// только запросы с действующей сессией, остальных отправляет на страницу входа.
func (h *AuthHandler) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookie)
		if err != nil { // cookie нет - пользователь не входил
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Сессия должна существовать и не быть просроченной.
		userID, err := h.users.GetUserIDBySession(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Кладём номер пользователя в контекст запроса и вызываем основной обработчик.
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

// userIDFromContext достаёт номер пользователя, который положил RequireAuth.
func userIDFromContext(r *http.Request) int {
	id, _ := r.Context().Value(userIDKey).(int)
	return id
}
