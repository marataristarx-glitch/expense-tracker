package main

import (
	// Подключаем нашу папку internal/handlers, где лежит логика сайта.
	
	"expense-tracker/internal/handlers"
	
	"fmt"      // Нужен для вывода сообщений в черное окно (консоль).
	"net/http" // Главный пакет для работы с сетью и создания серверов.
)

func main() {
	
	// Он принимает запросы от браузера и решает, на какую страницу их отправить.
	mux := http.NewServeMux()

	// этап 1: Главная страница сайта (переходим на http://localhost:8080)
	// Запускает функцию ExpensesHandler, которая показывает форму и таблицу трат.
	mux.HandleFunc("/", handlers.ExpensesHandler)         

	// этап 2: Служебный адрес для приема данных из формы (метод POST).
	// Запускает функцию AddExpenseHandler, которая считывает сумму, описание и сохраняет их.
	mux.HandleFunc("/expenses/add", handlers.AddExpenseHandler) 

	
	fmt.Println("Сервер Этапа 2 запущен на http://localhost:8080")

	
	
	err := http.ListenAndServe(":8080", mux)
	
	
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %v\n", err)
	}
}


