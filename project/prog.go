package main

import (
	"fmt"
	"net/http"
)

// Обработчик главной страницы
func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Добро пожаловать на главную страницу!")
}

// Обработчик страницы /about
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Это демонстрационный проект веб‑сервера на Go. \n"+
		"Сервер содержит три маршрута: главную страницу, /about и /ping.")
}

// Обработчик страницы /ping
func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "pong")
}
func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/ping", pingHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %v\n", err)
	}
}
