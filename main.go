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

	// Создаём таблицу категорий, добавляем базовые категории и создаём таблицу расходов с внешним ключом
	_, err = db.Exec(`
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
		amount REAL,
		description TEXT,
		date TEXT,
		category_id INTEGER,
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
	mux.HandleFunc("/", h.HandleIndex)
	mux.HandleFunc("/delete", h.HandleDelete)

	fmt.Println("Сервер Этапа 4 (категории и фильтрация) запущен на http://localhost:8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %v\n", err)
	}
}
