package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<h1>Сайт футбольной аналитики</h1>")
		fmt.Fprintln(w, "<p>Скоро здесь будет полноценный портал про КПЛ.</p>")
	})

	fmt.Println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
