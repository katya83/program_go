package models

// Expense - одна трата.
type Expense struct {
	ID           int     // номер траты (его выдаёт база данных)
	Amount       float64 // сумма
	Description  string  // описание
	Date         string  // дата (строкой, например "2026-10-10")
	CategoryID   int     // номер категории
	CategoryName string  // название категории (приходит из базы через JOIN)
	UserID       int     // номер владельца траты
}

type ExpenseFilter struct {
	CategoryID int    // 0 - все категории
	DateFrom   string // "" - без нижней границы даты
	DateTo     string // "" - без верхней границы даты
}
