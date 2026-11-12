package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"expense-tracker/internal/handlers"
	"expense-tracker/internal/repository"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "expenses.db")
	if err != nil {
		fmt.Printf("Не удалось открыть базу данных: %v\n", err)
		return
	}
	defer db.Close()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS sessions (
		session_id TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		expires_at TEXT NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE
	);

	INSERT OR IGNORE INTO categories (name) VALUES 
		('Еда'),
		('Транспорт'),
		('Жильё'),
		('Развлечения'),
		('Прочее');

	CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		amount REAL,
		description TEXT,
		date TEXT,
		category_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id),
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);`)
	if err != nil {
		fmt.Printf("Ошибка при создании таблиц в базе: %v\n", err)
		return
	}

	tpl := template.Must(template.ParseGlob("web/templates/*.html"))

	repo := repository.NewExpenseRepository(db)
	h := handlers.NewExpenseHandler(repo, tpl)

	mux := http.NewServeMux()

	// Открытые маршруты (доступны без входа)
	mux.HandleFunc("/register", h.HandleRegister)
	mux.HandleFunc("/login", h.HandleLogin)
	mux.HandleFunc("/logout", h.HandleLogout)

	// Защищённые маршруты (требуют авторизации через middleware)
	mux.HandleFunc("/", h.AuthMiddleware(h.HandleIndex))
	mux.HandleFunc("/delete", h.AuthMiddleware(h.HandleDelete))

	fmt.Println("Сервер Этапа 5 (авторизация и сессии) запущен на http://localhost:8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %v\n", err)
	}
}
