package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"project/internal/models"
	"project/internal/repository"
)

// Загружаем все шаблоны один раз при запуске программы.
var tmpl = template.Must(template.ParseFiles(
	"web/templates/layout.html",
	"web/templates/list.html",
	"web/templates/add.html",
	"web/templates/edit.html",
))

// ExpenseHandler хранит репозиторий, через который обработчики работают с базой.
type ExpenseHandler struct {
	repo *repository.ExpenseRepository
}

func NewExpenseHandler(repo *repository.ExpenseRepository) *ExpenseHandler {
	return &ExpenseHandler{repo: repo}
}

// readForm читает поля формы и собирает из них трату.
func readForm(r *http.Request) (models.Expense, error) {
	r.ParseForm()
	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		return models.Expense{}, err
	}
	return models.Expense{
		Amount:      amount,
		Description: r.FormValue("description"),
		Date:        r.FormValue("date"),
	}, nil
}

// ShowList показывает список трат.
func (h *ExpenseHandler) ShowList(w http.ResponseWriter, r *http.Request) {
	expenses, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Ошибка базы данных"+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "list.html", expenses)
}

// ShowForm показывает пустую форму добавления.
func (h *ExpenseHandler) ShowForm(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "add.html", nil)
}

// AddExpense сохраняет новую трату в базу.
func (h *ExpenseHandler) AddExpense(w http.ResponseWriter, r *http.Request) {
	e, err := readForm(r)
	if err != nil {
		http.Error(w, "Сумма должна быть числом", http.StatusBadRequest)
		return
	}
	if err := h.repo.Create(e); err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

// ShowEditForm показывает форму редактирования с текущими значениями.
func (h *ExpenseHandler) ShowEditForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id")) // берём {id} из адреса
	if err != nil {
		http.NotFound(w, r)
		return
	}
	e, err := h.repo.GetByID(id)
	if err != nil {
		http.NotFound(w, r) // такой траты нет
		return
	}
	tmpl.ExecuteTemplate(w, "edit.html", e)
}

// UpdateExpense сохраняет изменения траты.
func (h *ExpenseHandler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	e, err := readForm(r)
	if err != nil {
		http.Error(w, "Сумма должна быть числом", http.StatusBadRequest)
		return
	}
	e.ID = id
	if err := h.repo.Update(e); err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

// DeleteExpense удаляет трату.
func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}
