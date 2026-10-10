package models

// Expense - одна трата.
type Expense struct {
	Amount      float64 // сумма
	Description string  // описание
	Date        string  // дата (строкой, например "2026-10-10")
}

// expenses - список трат в памяти (обычный срез).
// Имя с маленькой буквы, значит, снаружи пакета он недоступен.
var expenses []Expense

// Add добавляет трату в конец списка.
func Add(e Expense) {
	expenses = append(expenses, e)
}

// GetAll возвращает все траты.
func GetAll() []Expense {
	return expenses
}
