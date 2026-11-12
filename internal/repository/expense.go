package repository

import (
	"database/sql"
	"expense-tracker/internal/models"
	"time"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

// --- РАБОТА С ПОЛЬЗОВАТЕЛЯМИ И СЕССИЯМИ ---

func (r *ExpenseRepository) CreateUser(email string, passwordHash string) error {
	_, err := r.db.Exec("INSERT INTO users (email, password_hash) VALUES (?, ?)", email, passwordHash)
	return err
}

func (r *ExpenseRepository) GetUserByEmail(email string) (models.User, error) {
	var u models.User
	err := r.db.QueryRow("SELECT id, email, password_hash FROM users WHERE email = ?", email).
		Scan(&u.ID, &u.Email, &u.PasswordHash)
	return u, err
}

func (r *ExpenseRepository) CreateSession(sessionID string, userID int, expiresAt time.Time) error {
	_, err := r.db.Exec(
		"INSERT INTO sessions (session_id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID, userID, expiresAt.Format(time.RFC3339),
	)
	return err
}

func (r *ExpenseRepository) GetUserIDBySession(sessionID string) (int, error) {
	var userID int
	var expiresStr string

	err := r.db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_id = ?", sessionID).
		Scan(&userID, &expiresStr)
	if err != nil {
		return 0, err
	}

	// Проверяем, не истёк ли срок действия сессии
	expiresAt, err := time.Parse(time.RFC3339, expiresStr)
	if err != nil || time.Now().After(expiresAt) {
		_ = r.DeleteSession(sessionID)
		return 0, sql.ErrNoRows
	}

	return userID, nil
}

func (r *ExpenseRepository) DeleteSession(sessionID string) error {
	_, err := r.db.Exec("DELETE FROM sessions WHERE session_id = ?", sessionID)
	return err
}

// --- РАБОТА С КАТЕГОРИЯМИ И ТРАТАМИ ---

func (r *ExpenseRepository) GetAllCategories() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// Получаем траты ТОЛЬКО текущего пользователя (WHERE e.user_id = ?)
func (r *ExpenseRepository) GetAll(userID int, categoryID string, dateFrom string, dateTo string) ([]models.Expense, error) {
	query := `
		SELECT e.id, e.user_id, e.amount, e.description, e.date, e.category_id, c.name
		FROM expenses e
		JOIN categories c ON e.category_id = c.id
		WHERE e.user_id = ?
	`
	args := []any{userID}

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
		err := rows.Scan(&e.ID, &e.UserID, &e.Amount, &e.Description, &e.Date, &e.CategoryID, &e.CategoryName)
		if err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, nil
}

func (r *ExpenseRepository) Create(userID int, amount float64, desc string, date string, categoryID int) error {
	_, err := r.db.Exec(
		"INSERT INTO expenses (user_id, amount, description, date, category_id) VALUES (?, ?, ?, ?, ?)",
		userID, amount, desc, date, categoryID,
	)
	return err
}

// Удаляем трату только если она принадлежит этому пользователю
func (r *ExpenseRepository) Delete(id int, userID int) error {
	_, err := r.db.Exec("DELETE FROM expenses WHERE id = ? AND user_id = ?", id, userID)
	return err
}
