package handlers

import (
	"expense-tracker/internal/models"
	"expense-tracker/internal/repository"
	"html/template"
	"net/http"
	"strconv"
)

type ExpenseHandler struct {
	Repo *repository.ExpenseRepository
	Tpl  *template.Template
}

// Структура для передачи в HTML-шаблон сразу всех данных: трат, категорий и выбранных фильтров
type IndexPageData struct {
	Expenses         []models.Expense
	Categories       []models.Category
	SelectedCategory string
	DateFrom         string
	DateTo           string
}

func NewExpenseHandler(repo *repository.ExpenseRepository, tpl *template.Template) *ExpenseHandler {
	return &ExpenseHandler{Repo: repo, Tpl: tpl}
}

func (h *ExpenseHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// Достаём параметры фильтра из адресной строки (?category=1&from=...&to=...)
		categoryID := r.URL.Query().Get("category")
		dateFrom := r.URL.Query().Get("from")
		dateTo := r.URL.Query().Get("to")

		list, err := h.Repo.GetAll(categoryID, dateFrom, dateTo)
		if err != nil {
			http.Error(w, "Ошибка базы данных при получении трат", http.StatusInternalServerError)
			return
		}

		categories, err := h.Repo.GetAllCategories()
		if err != nil {
			http.Error(w, "Ошибка базы данных при получении категорий", http.StatusInternalServerError)
			return
		}

		data := IndexPageData{
			Expenses:         list,
			Categories:       categories,
			SelectedCategory: categoryID,
			DateFrom:         dateFrom,
			DateTo:           dateTo,
		}

		h.Tpl.ExecuteTemplate(w, "layout.html", data)
		return
	}

	if r.Method == "POST" {
		amountStr := r.FormValue("amount")
		desc := r.FormValue("description")
		date := r.FormValue("date")
		catIDStr := r.FormValue("category_id")

		amount, _ := strconv.ParseFloat(amountStr, 64)
		categoryID, _ := strconv.Atoi(catIDStr)

		_ = h.Repo.Create(amount, desc, date, categoryID)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func (h *ExpenseHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		idStr := r.FormValue("id")
		id, _ := strconv.Atoi(idStr)

		_ = h.Repo.Delete(id)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
