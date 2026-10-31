package main

import (
	"log"
	"net/http"

	"project/internal/handlers"
	"project/internal/repository"
)

func main() {
	// Открываем базу данных (файл expenses.db) и выполняем новые миграции из папки migrations.
	db, err := repository.OpenDB("expenses.db", "migrations")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewExpenseRepository(db)
	users := repository.NewUserRepository(db)
	h := handlers.NewExpenseHandler(repo)
	auth := handlers.NewAuthHandler(users)

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/expenses", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /register", auth.ShowRegister)
	mux.HandleFunc("POST /register", auth.Register)
	mux.HandleFunc("GET /login", auth.ShowLogin)
	mux.HandleFunc("POST /login", auth.Login)
	mux.HandleFunc("POST /logout", auth.Logout)

	mux.HandleFunc("GET /expenses", auth.RequireAuth(h.ShowList))
	mux.HandleFunc("GET /expenses/add", auth.RequireAuth(h.ShowForm))
	mux.HandleFunc("POST /expenses/add", auth.RequireAuth(h.AddExpense))
	mux.HandleFunc("GET /expenses/{id}/edit", auth.RequireAuth(h.ShowEditForm))
	mux.HandleFunc("POST /expenses/{id}/edit", auth.RequireAuth(h.UpdateExpense))
	mux.HandleFunc("POST /expenses/{id}/delete", auth.RequireAuth(h.DeleteExpense))

	log.Println("Сервер запущен: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
