package repository

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite" // драйвер SQLite

	"project/internal/models"
)

// ExpenseRepository - вся работа с таблицей expenses. Только здесь есть SQL.
type ExpenseRepository struct {
	db *sql.DB
}

// OpenDB открывает файл базы данных и выполняет файл миграции.
func OpenDB(dsn, migrationFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn) // dsn - путь к файлу базы, например "expenses.db"
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite пишет в файл по одному, так не будет ошибки "database is locked"

	script, err := os.ReadFile(migrationFile) // читаем SQL-скрипт из файла
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(string(script)); err != nil { // создаём таблицу
		return nil, err
	}
	return db, nil
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// Create добавляет новую трату.
func (r *ExpenseRepository) Create(e models.Expense) error {
	_, err := r.db.Exec(
		"INSERT INTO expenses (amount, description, date) VALUES (?, ?, ?)",
		e.Amount, e.Description, e.Date,
	)
	return err
}

// GetAll возвращает все траты.
func (r *ExpenseRepository) GetAll() ([]models.Expense, error) {
	rows, err := r.db.Query("SELECT id, amount, description, date FROM expenses ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Expense
	for rows.Next() { // переходим к следующей строке результата
		var e models.Expense
		if err := rows.Scan(&e.ID, &e.Amount, &e.Description, &e.Date); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

// GetByID возвращает одну трату по номеру.
func (r *ExpenseRepository) GetByID(id int) (models.Expense, error) {
	var e models.Expense
	err := r.db.QueryRow(
		"SELECT id, amount, description, date FROM expenses WHERE id = ?", id,
	).Scan(&e.ID, &e.Amount, &e.Description, &e.Date)
	return e, err
}

// Update изменяет существующую трату.
func (r *ExpenseRepository) Update(e models.Expense) error {
	_, err := r.db.Exec(
		"UPDATE expenses SET amount = ?, description = ?, date = ? WHERE id = ?",
		e.Amount, e.Description, e.Date, e.ID,
	)
	return err
}

// Delete удаляет трату по номеру.
func (r *ExpenseRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM expenses WHERE id = ?", id)
	return err
}
