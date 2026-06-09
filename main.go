package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type HomeData struct {
	TopClubs   []Club
	LatestNews []News
	Username   string
}

// MatchPageData — данные для страницы матча (матч + вошедший пользователь)
type MatchPageData struct {
	Match    Match
	Username string
}

func main() {
	initDB()
	seedClubs()
	seedNews()
	seedMatches()
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		clubs, err := getClubs()
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			return
		}
		topClubs := clubs
		if len(topClubs) > 3 {
			topClubs = topClubs[:3]
		}

		newsList, err := getNews()
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			return
		}
		latestNews := newsList
		if len(latestNews) > 3 {
			latestNews = latestNews[:3]
		}

		data := HomeData{TopClubs: topClubs, LatestNews: latestNews, Username: getCurrentUser(r)}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
	})

	http.HandleFunc("/table", func(w http.ResponseWriter, r *http.Request) {
		clubs, err := getClubs()
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			return
		}
		funcMap := template.FuncMap{
			"add": func(a, b int) int { return a + b },
			"sub": func(a, b int) int { return a - b },
		}
		tmpl, err := template.New("table.html").Funcs(funcMap).ParseFiles("templates/table.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, clubs)
	})

	http.HandleFunc("/clubs", func(w http.ResponseWriter, r *http.Request) {
		clubs, err := getClubs()
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			return
		}
		tmpl, err := template.ParseFiles("templates/clubs.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, clubs)
	})

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
			return
		}
		tmpl, err := template.ParseFiles("templates/club.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, club)
	})

	http.HandleFunc("/news", func(w http.ResponseWriter, r *http.Request) {
		newsList, err := getNews()
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			return
		}
		tmpl, err := template.ParseFiles("templates/news.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, newsList)
	})

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
			return
		}
		tmpl, err := template.ParseFiles("templates/article.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, article)
	})

	http.HandleFunc("/matches", func(w http.ResponseWriter, r *http.Request) {
		matches, err := getMatches()
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			return
		}
		tmpl, err := template.ParseFiles("templates/matches.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, matches)
	})

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
			return
		}
		data := MatchPageData{Match: match, Username: getCurrentUser(r)}
		tmpl, err := template.ParseFiles("templates/match.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
	})

	// Покупка билета
	http.HandleFunc("/buy", func(w http.ResponseWriter, r *http.Request) {
		username := getCurrentUser(r)
		// Если не вошёл — отправляем на страницу входа
		if username == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		idStr := r.URL.Query().Get("id")
		matchID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Неверный id матча", http.StatusBadRequest)
			return
		}

		user, err := getUserByUsername(username)
		if err != nil {
			http.Error(w, "Ошибка пользователя", http.StatusInternalServerError)
			return
		}

		err = buyTicket(user.ID, matchID)
		if err != nil {
			http.Error(w, "Ошибка покупки билета", http.StatusInternalServerError)
			return
		}

		// После покупки — на страницу "Мои билеты"
		http.Redirect(w, r, "/tickets", http.StatusSeeOther)
	})

	// Мои билеты
	http.HandleFunc("/tickets", func(w http.ResponseWriter, r *http.Request) {
		username := getCurrentUser(r)
		if username == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := getUserByUsername(username)
		if err != nil {
			http.Error(w, "Ошибка пользователя", http.StatusInternalServerError)
			return
		}

		tickets, err := getUserTickets(user.ID)
		if err != nil {
			http.Error(w, "Ошибка получения билетов", http.StatusInternalServerError)
			return
		}

		data := struct {
			Tickets  []Ticket
			Username string
		}{Tickets: tickets, Username: username}

		tmpl, err := template.ParseFiles("templates/tickets.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			username := r.FormValue("username")
			password := r.FormValue("password")

			if username == "" || password == "" {
				renderAuth(w, "register.html", "Заполните все поля")
				return
			}

			_, err := getUserByUsername(username)
			if err == nil {
				renderAuth(w, "register.html", "Этот логин уже занят")
				return
			}

			hash, err := hashPassword(password)
			if err != nil {
				renderAuth(w, "register.html", "Ошибка при регистрации")
				return
			}

			err = createUser(username, hash)
			if err != nil {
				renderAuth(w, "register.html", "Не удалось создать пользователя")
				return
			}

			loginUser(w, r, username)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		renderAuth(w, "register.html", "")
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			username := r.FormValue("username")
			password := r.FormValue("password")

			user, err := getUserByUsername(username)
			if err != nil {
				renderAuth(w, "login.html", "Неверный логин или пароль")
				return
			}

			if !checkPassword(password, user.Password) {
				renderAuth(w, "login.html", "Неверный логин или пароль")
				return
			}

			loginUser(w, r, username)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		renderAuth(w, "login.html", "")
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		logoutUser(w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

type AuthData struct {
	Error string
}

func renderAuth(w http.ResponseWriter, page, errMsg string) {
	tmpl, err := template.ParseFiles("templates/" + page)
	if err != nil {
		http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, AuthData{Error: errMsg})
}
