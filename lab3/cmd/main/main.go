package main

import (
	"log"
	"net/http"

	"project/internal/handlers"
	"project/internal/repository"
)

func main() {
	// Открываем базу данных (файл expenses.db) и создаём таблицу из миграции.
	db, err := repository.OpenDB("expenses.db", "migrations/001_create_expenses.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewExpenseRepository(db)
	h := handlers.NewExpenseHandler(repo)

	mux := http.NewServeMux()

	// Раздача статических файлов (CSS) из папки web/static.
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Главная страница перенаправляет на список.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/expenses", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /expenses", h.ShowList)
	mux.HandleFunc("GET /expenses/add", h.ShowForm)
	mux.HandleFunc("POST /expenses/add", h.AddExpense)
	mux.HandleFunc("GET /expenses/{id}/edit", h.ShowEditForm)
	mux.HandleFunc("POST /expenses/{id}/edit", h.UpdateExpense)
	mux.HandleFunc("POST /expenses/{id}/delete", h.DeleteExpense)

	log.Println("Сервер запущен: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
