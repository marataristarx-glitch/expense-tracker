package main

import (
	"expense-tracker/internal/handlers"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Привязываем новые обработчики интерфейса
	mux.HandleFunc("/", handlers.ExpensesHandler)               // Главная страница со списком
	mux.HandleFunc("/expenses/add", handlers.AddExpenseHandler) // Обработка формы добавления

	fmt.Println("Сервер Этапа 2 запущен на http://localhost:8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %v\n", err)
	}
}
