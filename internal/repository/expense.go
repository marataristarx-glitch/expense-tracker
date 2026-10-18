package repository

import (
	"database/sql"
	"expense-tracker/internal/models"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// Получаем вообще все траты из базы чтобы вывести в таблицу
func (r *ExpenseRepository) GetAll() ([]models.Expense, error) {
	// Сортирую по id  чтобы новые траты были сверху, так удобнее
	rows, err := r.db.Query("SELECT id, amount, description, date FROM expenses ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Expense
	for rows.Next() {
		var e models.Expense
		// Сканируем все поля по очереди. Главное не перепутать порядок колонок
		err := rows.Scan(&e.ID, &e.Amount, &e.Description, &e.Date)
		if err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}


func (r *ExpenseRepository) Create(amount float64, desc string, date string) error {
	// 
	_, err := r.db.Exec("INSERT INTO expenses (amount, description, date) VALUES (?, ?, ?)", amount, desc, date)
	return err
}

// Удаление по айдишнику. Работает вроде четко
func (r *ExpenseRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM expenses WHERE id = ?", id)
	return err
}
