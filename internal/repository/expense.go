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

// Получаем все категории из базы для выпадающих списков
func (r *ExpenseRepository) GetAllCategories() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// Получаем траты с учётом фильтров 
func (r *ExpenseRepository) GetAll(categoryID string, dateFrom string, dateTo string) ([]models.Expense, error) {
	
	query := `
		SELECT e.id, e.amount, e.description, e.date, e.category_id, c.name
		FROM expenses e
		JOIN categories c ON e.category_id = c.id
		WHERE 1=1
	`
	var args []any

	if categoryID != "" {
		query += " AND e.category_id = ?"
		args = append(args, categoryID)
	}
	if dateFrom != "" {
		query += " AND e.date >= ?"
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		query += " AND e.date <= ?"
		args = append(args, dateTo)
	}

	query += " ORDER BY e.id DESC"

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
	return list, nil
}

// Добавляем новую трату вместе с ID категории
func (r *ExpenseRepository) Create(amount float64, desc string, date string, categoryID int) error {
	_, err := r.db.Exec(
		"INSERT INTO expenses (amount, description, date, category_id) VALUES (?, ?, ?, ?)",
		amount, desc, date, categoryID,
	)
	return err
}

// Удаление по айдишнику
func (r *ExpenseRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM expenses WHERE id = ?", id)
	return err
}
