package main

import (
	"log"
	"net/http"

	"project/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	// Раздача статических файлов (CSS) из папки web/static.
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	// Главная страница перенаправляет на список.
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/expenses", http.StatusSeeOther)
	})

	mux.HandleFunc("GET /expenses", handlers.ShowList)
	mux.HandleFunc("GET /expenses/add", handlers.ShowForm)
	mux.HandleFunc("POST /expenses/add", handlers.AddExpense)

	log.Println("Сервер запущен: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
