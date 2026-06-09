package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Club struct {
	ID       int
	Name     string
	City     string
	Stadium  string
	Played   int
	Wins     int
	Draws    int
	Losses   int
	GoalsFor int
	GoalsAg  int
	Points   int
}

type News struct {
	ID      int
	Title   string
	Date    string
	Content string
}

type Match struct {
	ID        int
	HomeTeam  string
	AwayTeam  string
	HomeScore int
	AwayScore int
	Date      string
	Stadium   string
	Played    bool
	Price     int
}

// User — структура пользователя
type User struct {
	ID       int
	Username string
	Password string // здесь хранится хеш пароля, не сам пароль
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "kpl.db")
	if err != nil {
		fmt.Println("Ошибка открытия базы:", err)
		return
	}

	createClubs := `
	CREATE TABLE IF NOT EXISTS clubs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		city TEXT,
		stadium TEXT,
		played INTEGER DEFAULT 0,
		wins INTEGER DEFAULT 0,
		draws INTEGER DEFAULT 0,
		losses INTEGER DEFAULT 0,
		goals_for INTEGER DEFAULT 0,
		goals_against INTEGER DEFAULT 0,
		points INTEGER DEFAULT 0
	);`
	_, err = db.Exec(createClubs)
	if err != nil {
		fmt.Println("Ошибка создания таблицы clubs:", err)
		return
	}

	createNews := `
	CREATE TABLE IF NOT EXISTS news (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		date TEXT,
		content TEXT
	);`
	_, err = db.Exec(createNews)
	if err != nil {
		fmt.Println("Ошибка создания таблицы news:", err)
		return
	}

	createMatches := `
	CREATE TABLE IF NOT EXISTS matches (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		home_team TEXT NOT NULL,
		away_team TEXT NOT NULL,
		home_score INTEGER DEFAULT 0,
		away_score INTEGER DEFAULT 0,
		date TEXT,
		stadium TEXT,
		played INTEGER DEFAULT 0,
		price INTEGER DEFAULT 0
	);`
	_, err = db.Exec(createMatches)
	if err != nil {
		fmt.Println("Ошибка создания таблицы matches:", err)
		return
	}

	createUsers := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`
	_, err = db.Exec(createUsers)
	if err != nil {
		fmt.Println("Ошибка создания таблицы users:", err)
		return
	}

	fmt.Println("База данных готова")
}

func getClubs() ([]Club, error) {
	rows, err := db.Query("SELECT id, name, city, stadium, played, wins, draws, losses, goals_for, goals_against, points FROM clubs ORDER BY points DESC, (goals_for - goals_against) DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clubs []Club
	for rows.Next() {
		var c Club
		err := rows.Scan(&c.ID, &c.Name, &c.City, &c.Stadium, &c.Played, &c.Wins, &c.Draws, &c.Losses, &c.GoalsFor, &c.GoalsAg, &c.Points)
		if err != nil {
			return nil, err
		}
		clubs = append(clubs, c)
	}
	return clubs, nil
}

func getClubByID(id int) (Club, error) {
	var c Club
	err := db.QueryRow(
		"SELECT id, name, city, stadium, played, wins, draws, losses, goals_for, goals_against, points FROM clubs WHERE id = ?",
		id,
	).Scan(&c.ID, &c.Name, &c.City, &c.Stadium, &c.Played, &c.Wins, &c.Draws, &c.Losses, &c.GoalsFor, &c.GoalsAg, &c.Points)
	return c, err
}

func getNews() ([]News, error) {
	rows, err := db.Query("SELECT id, title, date, content FROM news ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var newsList []News
	for rows.Next() {
		var n News
		err := rows.Scan(&n.ID, &n.Title, &n.Date, &n.Content)
		if err != nil {
			return nil, err
		}
		newsList = append(newsList, n)
	}
	return newsList, nil
}

func getNewsByID(id int) (News, error) {
	var n News
	err := db.QueryRow(
		"SELECT id, title, date, content FROM news WHERE id = ?",
		id,
	).Scan(&n.ID, &n.Title, &n.Date, &n.Content)
	return n, err
}

func getMatches() ([]Match, error) {
	rows, err := db.Query("SELECT id, home_team, away_team, home_score, away_score, date, stadium, played, price FROM matches ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []Match
	for rows.Next() {
		var m Match
		var playedInt int
		err := rows.Scan(&m.ID, &m.HomeTeam, &m.AwayTeam, &m.HomeScore, &m.AwayScore, &m.Date, &m.Stadium, &playedInt, &m.Price)
		if err != nil {
			return nil, err
		}
		m.Played = playedInt == 1
		matches = append(matches, m)
	}
	return matches, nil
}

func getMatchByID(id int) (Match, error) {
	var m Match
	var playedInt int
	err := db.QueryRow(
		"SELECT id, home_team, away_team, home_score, away_score, date, stadium, played, price FROM matches WHERE id = ?",
		id,
	).Scan(&m.ID, &m.HomeTeam, &m.AwayTeam, &m.HomeScore, &m.AwayScore, &m.Date, &m.Stadium, &playedInt, &m.Price)
	m.Played = playedInt == 1
	return m, err
}

// createUser добавляет нового пользователя (пароль уже хеширован)
func createUser(username, passwordHash string) error {
	_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, passwordHash)
	return err
}

// getUserByUsername находит пользователя по логину
func getUserByUsername(username string) (User, error) {
	var u User
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(&u.ID, &u.Username, &u.Password)
	return u, err
}
