package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// Club — структура, описывающая один клуб
type Club struct {
	ID       int
	Name     string
	City     string
	Stadium  string
	Played   int // сыграно матчей
	Wins     int // победы
	Draws    int // ничьи
	Losses   int // поражения
	GoalsFor int // забито
	GoalsAg  int // пропущено
	Points   int // очки
}

// db — глобальная переменная для подключения к базе
var db *sql.DB

// initDB открывает базу и создаёт таблицу, если её нет
func initDB() {
	var err error
	db, err = sql.Open("sqlite", "kpl.db")
	if err != nil {
		fmt.Println("Ошибка открытия базы:", err)
		return
	}

	// SQL-команда: создать таблицу clubs, если она ещё не существует
	createTable := `
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

	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println("Ошибка создания таблицы:", err)
		return
	}

	fmt.Println("База данных готова")
}

// getClubs возвращает все клубы, отсортированные по очкам
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
