package main

import (
	"database/sql" // Новый пакет для работы с  SQLite
	"fmt"          // Нужен для вывода сообщений в черное окно (ТЕРМИНАЛ).
	"html/template" // Нужен для чтения html файлов
	"net/http"     // Главный пакет для создания серверов.

	// Подключаем наши папки с базой данных
	"expense-tracker/internal/handlers"
	"expense-tracker/internal/repository"

	_ "modernc.org/sqlite" // Подключаем драйвер для базы данных
)

func main() {
	
	db, err := sql.Open("sqlite", "expenses.db")
	if err != nil {
		fmt.Printf("Не удалось открыть базу данных: %v\n", err)
		return
	}
	defer db.Close()

	// Сюда запихал создание таблицы, чтобы она сама создавалась при запуске, если пропала.
	//  работает xtnrf
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		amount REAL,
		description TEXT,
		date TEXT
	);`)
	if err != nil {
		fmt.Printf("Ошибка при создании таблицы в базе: %v\n", err)
		return
	}

	// Читаем все html файлы из папки шаблонов web  
	tpl := template.Must(template.ParseGlob("web/templates/*.html"))

	
	repo := repository.NewExpenseRepository(db)
	
	
	h := handlers.NewExpenseHandler(repo, tpl)

	// Он принимает запросы от браузера и решает, на какую страницу их отправить.
	mux := http.NewServeMux()

	// Этап 3: Главная страница сайта (переходим на http://localhost:8080)
	
	mux.HandleFunc("/", h.HandleIndex)         

	// кнопка удаления
	mux.HandleFunc("/delete", h.HandleDelete)

	fmt.Println("Сервер Этапа 3 (с SQLite) запущен на http://localhost:8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %v\n", err)
	}
}
