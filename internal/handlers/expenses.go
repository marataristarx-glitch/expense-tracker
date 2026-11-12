package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"expense-tracker/internal/models"
	"expense-tracker/internal/repository"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const userIDKey contextKey = "userID"

type ExpenseHandler struct {
	Repo *repository.ExpenseRepository
	Tpl  *template.Template
}

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

// Генерация случайного токена сессии
func generateSessionID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Middleware авторизации: пускает дальше только залогиненных пользователей
func (h *ExpenseHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userID, err := h.Repo.GetUserIDBySession(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Кладём userID в контекст запроса, чтобы достать его в обработчике
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

// Регистрация
func (h *ExpenseHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.Tpl.ExecuteTemplate(w, "register.html", nil)
		return
	}

	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Хэшируем пароль перед сохранением в БД
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			h.Tpl.ExecuteTemplate(w, "register.html", "Ошибка обработки пароля")
			return
		}

		err = h.Repo.CreateUser(email, string(hash))
		if err != nil {
			h.Tpl.ExecuteTemplate(w, "register.html", "Пользователь с таким email уже существует")
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// Вход
func (h *ExpenseHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.Tpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		user, err := h.Repo.GetUserByEmail(email)
		if err != nil {
			h.Tpl.ExecuteTemplate(w, "login.html", "Неверный email или пароль")
			return
		}

		// Сравниваем введённый пароль с хэшем из базы
		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
		if err != nil {
			h.Tpl.ExecuteTemplate(w, "login.html", "Неверный email или пароль")
			return
		}

		// Создаём сессию на 24 часа и ставим cookie
		sessionID := generateSessionID()
		expiresAt := time.Now().Add(24 * time.Hour)
		_ = h.Repo.CreateSession(sessionID, user.ID, expiresAt)

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Expires:  expiresAt,
			Path:     "/",
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

// Выход
func (h *ExpenseHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		_ = h.Repo.DeleteSession(cookie.Value)
	}

	// Удаляем cookie в браузере
	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Главная страница (доступна только после AuthMiddleware)
func (h *ExpenseHandler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)

	if r.Method == "GET" {
		categoryID := r.URL.Query().Get("category")
		dateFrom := r.URL.Query().Get("from")
		dateTo := r.URL.Query().Get("to")

		list, err := h.Repo.GetAll(userID, categoryID, dateFrom, dateTo)
		if err != nil {
			http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
			return
		}

		categories, err := h.Repo.GetAllCategories()
		if err != nil {
			http.Error(w, "Ошибка загрузки категорий", http.StatusInternalServerError)
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

		_ = h.Repo.Create(userID, amount, desc, date, categoryID)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func (h *ExpenseHandler) HandleDelete(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)

	if r.Method == "POST" {
		idStr := r.FormValue("id")
		id, _ := strconv.Atoi(idStr)

		_ = h.Repo.Delete(id, userID)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
