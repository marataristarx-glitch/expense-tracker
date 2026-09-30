package handlers

import (
	"expense-tracker/internal/models"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

// Временное хранилище в оперативной памяти (срез структур)
var ExpensesList = []models.Expense{}

// ExpensesHandler отображает страницу со списком трат
func ExpensesHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Парсим файлы шаблонов (layout и контент страницы)
	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблонов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	// Передаем срез ExpensesList в шаблон для отрисовки таблицы
	tmpl.ExecuteTemplate(w, "layout", ExpensesList)
}

// AddExpenseHandler обрабатывает отправку формы (POST)
func AddExpenseHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Разбираем данные формы
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка разбора формы", http.StatusBadRequest)
		return
	}

	// Извлекаем значения по атрибуту name из HTML
	amountStr := r.FormValue("amount")
	description := r.FormValue("description")
	dateStr := r.FormValue("date")

	// Конвертируем строку суммы в число с плавающей точкой
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, "Некорректная сумма", http.StatusBadRequest)
		return
	}

	// Конвертируем строку даты в тип time.Time
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		date = time.Now() // если дата сбоит, ставим текущую
	}

	// Создаем объект траты
	newExpense := models.Expense{
		Amount:      amount,
		Description: description,
		Date:        date,
	}

	// Добавляем созданную трату в наш срез в памяти
	ExpensesList = append(ExpensesList, newExpense)

	// Паттерн Post/Redirect/Get: делаем редирект обратно на главную страницу,
	// чтобы при обновлении страницы форма не отправлялась повторно.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
