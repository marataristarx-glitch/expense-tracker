package models

import "time"

// Expense описывает структуру одной траты (сумма, описание, дата)
type Expense struct {
	Amount      float64   // Сумма траты
	Description string    // Описание на что потрачено
	Date        time.Time // Дата совершения траты
}
