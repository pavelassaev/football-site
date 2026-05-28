package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	// Обработчик главной страницы
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Загружаем шаблон главной страницы
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
			fmt.Println("Ошибка шаблона:", err)
			return
		}
		// Отдаём страницу пользователю
		tmpl.Execute(w, nil)
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
