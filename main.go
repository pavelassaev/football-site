package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func main() {
	// Инициализация базы данных
	initDB()
	seedClubs()
	defer db.Close()

	// Главная страница
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			fmt.Println("Ошибка шаблона:", err)
			return
		}
		tmpl.Execute(w, nil)
	})

	// Страница турнирной таблицы
	http.HandleFunc("/table", func(w http.ResponseWriter, r *http.Request) {
		clubs, err := getClubs()
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			fmt.Println("Ошибка getClubs:", err)
			return
		}

		funcMap := template.FuncMap{
			"add": func(a, b int) int { return a + b },
			"sub": func(a, b int) int { return a - b },
		}

		tmpl, err := template.New("table.html").Funcs(funcMap).ParseFiles("templates/table.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			fmt.Println("Ошибка шаблона:", err)
			return
		}
		tmpl.Execute(w, clubs)
	})

	// Страница со списком клубов
	http.HandleFunc("/clubs", func(w http.ResponseWriter, r *http.Request) {
		clubs, err := getClubs()
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			fmt.Println("Ошибка getClubs:", err)
			return
		}
		tmpl, err := template.ParseFiles("templates/clubs.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			fmt.Println("Ошибка шаблона:", err)
			return
		}
		tmpl.Execute(w, clubs)
	})

	// Страница одного клуба
	http.HandleFunc("/club", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Неверный id клуба", http.StatusBadRequest)
			return
		}

		club, err := getClubByID(id)
		if err != nil {
			http.Error(w, "Клуб не найден", http.StatusNotFound)
			fmt.Println("Ошибка getClubByID:", err)
			return
		}

		tmpl, err := template.ParseFiles("templates/club.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			fmt.Println("Ошибка шаблона:", err)
			return
		}
		tmpl.Execute(w, club)
	})

	// Раздаём статические файлы (CSS, картинки) из папки static
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
