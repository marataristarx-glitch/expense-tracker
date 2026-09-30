package main

import (
	"fmt"
	"net/http"
)

// Обработчик для главной страницы
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что пользователь зашел именно на "/", а не на случайный адрес
	if r.URL.Path != "/" {
		http.NotFound(w, r) // Возвращает статус 404 Not Found
		return
	}
	w.WriteHeader(http.StatusOK) // Отправляем статус 200 OK
	fmt.Fprint(w, "Добро пожаловать в приложение Учёт личных трат!")
}

// Обработчик для страницы /about
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK) // Отправляем статус 200 OK
	fmt.Fprint(w, "Это семестровое практическое задание по веб-разработке на языке Go.")
}

// Обработчик для страницы /ping
func pingHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод. Если это НЕ GET-запрос — возвращаем ошибку 405 Method Not Allowed
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается. Используйте только GET.", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK) // Отправляем статус 200 OK
	fmt.Fprint(w, "pong")
}

func main() {
	// Создаем встроенный маршрутизатор (ServeMux)
	mux := http.NewServeMux()

	// Регистрируем маршруты и привязываем к ним наши функции
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/about", aboutHandler)
	mux.HandleFunc("/ping", pingHandler)

	fmt.Println("Сервер запущен на http://localhost:8080")

	// Запускаем веб-сервер на порту 8080
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %v\n", err)
	}
}
