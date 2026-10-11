package models

// Expense - одна трата.
type Expense struct {
	ID          int     // номер траты (его выдаёт база данных)
	Amount      float64 // сумма
	Description string  // описание
	Date        string  // дата (строкой, например "2026-10-10")
}
