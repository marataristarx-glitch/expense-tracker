package models

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type Category struct {
	ID   int
	Name string
}

type Expense struct {
	ID           int
	UserID       int
	Amount       float64
	Description  string
	Date         string
	CategoryID   int
	CategoryName string
}
