package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"project/internal/models"
)

// Загружаем все три шаблона один раз при запуске программы.
// template.Must остановит программу с ошибкой, если файл не найден.
var tmpl = template.Must(template.ParseFiles(
	"web/templates/layout.html",
	"web/templates/list.html",
	"web/templates/add.html",
))

// ShowList показывает список трат.
func ShowList(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "list.html", models.GetAll())
}

// ShowForm показывает пустую форму добавления.
func ShowForm(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "add.html", nil)
}

// AddExpense принимает данные формы и сохраняет трату.
func AddExpense(w http.ResponseWriter, r *http.Request) {
	// 1. Разбираем данные, которые прислал браузер.
	r.ParseForm()

	// 2. Достаём значения полей по их name из HTML-формы.
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		http.Error(w, "Сумма должна быть числом", http.StatusBadRequest)
		return
	}
	description := r.FormValue("description")
	date := r.FormValue("date")

	// 3. Сохраняем трату в список.
	models.Add(models.Expense{
		Amount:      amount,
		Description: description,
		Date:        date,
	})

	// 4. Перенаправляем на список (а не рисуем страницу сами).
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}
