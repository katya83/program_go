package repository

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite" // драйвер SQLite (подключается ради побочного эффекта)

	"project/internal/models"
)

// ExpenseRepository - вся работа с тратами и категориями. Только здесь есть SQL для трат.
type ExpenseRepository struct {
	db *sql.DB
}

// OpenDB открывает файл базы данных и выполняет все новые миграции из папки.
func OpenDB(dsn, migrationsDir string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn) // dsn - путь к файлу базы, например "expenses.db"
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite пишет в файл по одному, так не будет ошибки "database is locked"

	// В SQLite проверка внешних ключей по умолчанию выключена - включаем.
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	if err := runMigrations(db, migrationsDir); err != nil {
		return nil, err
	}
	return db, nil
}

// runMigrations выполняет файлы *.sql из папки по порядку (001, 002, 003, ...).
// Выполненные миграции запоминаются в таблице schema_migrations и повторно не запускаются.
func runMigrations(db *sql.DB, dir string) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS schema_migrations (name TEXT PRIMARY KEY)")
	if err != nil {
		return err
	}

	files, err := filepath.Glob(filepath.Join(dir, "*.sql")) // все .sql файлы папки
	if err != nil {
		return err
	}
	sort.Strings(files) // сортируем по имени: 001, 002, 003, ...

	for _, file := range files {
		name := filepath.Base(file)

		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE name = ?", name).Scan(&count)
		if err != nil {
			return err
		}
		if count > 0 {
			continue // эта миграция уже выполнена
		}

		script, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		// Транзакция: либо выполнится всё, либо ничего.
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		if _, err := tx.Exec(string(script)); err != nil {
			tx.Rollback()
			return err
		}
		if _, err := tx.Exec("INSERT INTO schema_migrations (name) VALUES (?)", name); err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// GetCategories возвращает все категории (для выпадающих списков).
func (r *ExpenseRepository) GetCategories() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

// Create добавляет новую трату (владелец - e.UserID).
func (r *ExpenseRepository) Create(e models.Expense) error {
	_, err := r.db.Exec(
		"INSERT INTO expenses (amount, description, date, category_id, user_id) VALUES (?, ?, ?, ?, ?)",
		e.Amount, e.Description, e.Date, e.CategoryID, e.UserID,
	)
	return err
}

// GetAll возвращает траты пользователя вместе с названием категории (JOIN) с учётом фильтра.
func (r *ExpenseRepository) GetAll(userID int, f models.ExpenseFilter) ([]models.Expense, error) {
	query := `SELECT e.id, e.amount, e.description, e.date, e.category_id, c.name
              FROM expenses e
              JOIN categories c ON c.id = e.category_id`

	// Первое условие есть всегда: берём только траты этого пользователя.
	// Остальные условия добавляем только для тех фильтров, которые заданы.
	conditions := []string{"e.user_id = ?"}
	args := []any{userID}

	if f.CategoryID != 0 {
		conditions = append(conditions, "e.category_id = ?")
		args = append(args, f.CategoryID)
	}
	if f.DateFrom != "" {
		conditions = append(conditions, "e.date >= ?")
		args = append(args, f.DateFrom)
	}
	if f.DateTo != "" {
		conditions = append(conditions, "e.date <= ?")
		args = append(args, f.DateTo)
	}

	query += " WHERE " + strings.Join(conditions, " AND ")
	query += " ORDER BY e.date DESC, e.id DESC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Expense
	for rows.Next() {
		var e models.Expense
		err := rows.Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.CategoryID, &e.CategoryName)
		if err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

// GetByID возвращает одну трату по номеру, но только если она принадлежит пользователю.
func (r *ExpenseRepository) GetByID(id, userID int) (models.Expense, error) {
	var e models.Expense
	err := r.db.QueryRow(
		"SELECT id, amount, description, date, category_id FROM expenses WHERE id = ? AND user_id = ?",
		id, userID,
	).Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.CategoryID)
	return e, err
}

// Update изменяет трату, но только если она принадлежит пользователю.
func (r *ExpenseRepository) Update(e models.Expense) error {
	_, err := r.db.Exec(
		"UPDATE expenses SET amount = ?, description = ?, date = ?, category_id = ? WHERE id = ? AND user_id = ?",
		e.Amount, e.Description, e.Date, e.CategoryID, e.ID, e.UserID,
	)
	return err
}

// Delete удаляет трату, но только если она принадлежит пользователю.
func (r *ExpenseRepository) Delete(id, userID int) error {
	_, err := r.db.Exec("DELETE FROM expenses WHERE id = ? AND user_id = ?", id, userID)
	return err
}
