package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// HomeData — данные для главной страницы
type HomeData struct {
	TopClubs   []Club
	LatestNews []News
}

func main() {
	initDB()
	seedClubs()
	seedNews()
	seedMatches()
	defer db.Close()

	// Главная страница
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		clubs, err := getClubs()
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			fmt.Println("Ошибка getClubs:", err)
			return
		}
		topClubs := clubs
		if len(topClubs) > 3 {
			topClubs = topClubs[:3]
		}

		newsList, err := getNews()
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			fmt.Println("Ошибка getNews:", err)
			return
		}
		latestNews := newsList
		if len(latestNews) > 3 {
			latestNews = latestNews[:3]
		}

		data := HomeData{TopClubs: topClubs, LatestNews: latestNews}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			fmt.Println("Ошибка шаблона:", err)
			return
		}
		tmpl.Execute(w, data)
	})

	// Турнирная таблица
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

	// Список клубов
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

	// Один клуб
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

	// Список новостей
	http.HandleFunc("/news", func(w http.ResponseWriter, r *http.Request) {
		newsList, err := getNews()
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			fmt.Println("Ошибка getNews:", err)
			return
		}
		tmpl, err := template.ParseFiles("templates/news.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			fmt.Println("Ошибка шаблона:", err)
			return
		}
		tmpl.Execute(w, newsList)
	})

	// Одна новость
	http.HandleFunc("/article", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Неверный id новости", http.StatusBadRequest)
			return
		}
		article, err := getNewsByID(id)
		if err != nil {
			http.Error(w, "Новость не найдена", http.StatusNotFound)
			fmt.Println("Ошибка getNewsByID:", err)
			return
		}
		tmpl, err := template.ParseFiles("templates/article.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			fmt.Println("Ошибка шаблона:", err)
			return
		}
		tmpl.Execute(w, article)
	})

	// Список матчей
	http.HandleFunc("/matches", func(w http.ResponseWriter, r *http.Request) {
		matches, err := getMatches()
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			fmt.Println("Ошибка getMatches:", err)
			return
		}
		tmpl, err := template.ParseFiles("templates/matches.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			fmt.Println("Ошибка шаблона:", err)
			return
		}
		tmpl.Execute(w, matches)
	})

	// Один матч
	http.HandleFunc("/match", func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Неверный id матча", http.StatusBadRequest)
			return
		}
		match, err := getMatchByID(id)
		if err != nil {
			http.Error(w, "Матч не найден", http.StatusNotFound)
			fmt.Println("Ошибка getMatchByID:", err)
			return
		}
		tmpl, err := template.ParseFiles("templates/match.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			fmt.Println("Ошибка шаблона:", err)
			return
		}
		tmpl.Execute(w, match)
	})

	// Статические файлы
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
