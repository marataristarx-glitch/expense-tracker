package handlers

import (
	"expense-tracker/internal/repository"
	"html/template"
	"net/http"
	"strconv"
)

type ExpenseHandler struct {
	Repo *repository.ExpenseRepository
	Tpl  *template.Template
}

func NewExpenseHandler(repo *repository.ExpenseRepository, tpl *template.Template) *ExpenseHandler {
	return &ExpenseHandler{Repo: repo, Tpl: tpl}
}

// Один главный хендлер для главной страницы
func (h *ExpenseHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	// Если юзер просто зашел на сайт (GET запрос)
	if r.Method == "GET" {
		list, err := h.Repo.GetAll()
		if err != nil {
			http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
			return
		}
		// Отдаем html страницу и список трат
		h.Tpl.ExecuteTemplate(w, "layout.html", list)
		return
	}

	// Если юзер отправил форму добавления (POST запрос)
	if r.Method == "POST" {
		amountStr := r.FormValue("amount")
		desc := r.FormValue("description")
		date := r.FormValue("date")

		amount, _ := strconv.ParseFloat(amountStr, 64)

		_ = h.Repo.Create(amount, desc, date)

		// Возвращаем обратно
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

// Отдельный хендлер для удаления, потому что путь другой (/delete)
func (h *ExpenseHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	// На всякий случай проверяем, что это точно POST запрос
	if r.Method == "POST" {
		idStr := r.FormValue("id")
		id, _ := strconv.Atoi(idStr)

		_ = h.Repo.Delete(id)
	}
	
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
