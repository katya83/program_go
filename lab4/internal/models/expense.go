package models

// Expense - одна трата.
type Expense struct {
	ID           int     // номер траты
	Amount       float64 // сумма
	Description  string  // описание
	Date         string  // дата
	CategoryID   int     // номер категории
	CategoryName string  // название категории
}

// ExpenseFilter - параметры фильтрации списка.
// Пустое значение (0 или "") означает "этот фильтр не задан".
type ExpenseFilter struct {
	CategoryID int    // 0 - все категории
	DateFrom   string // "" - без нижней границы даты
	DateTo     string // "" - без верхней границы даты
}
