package main

import (
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// store — хранилище сессий. Строка-ключ используется для подписи cookie.
// В реальном проекте её прячут, для учебного — сгодится так.
var store = sessions.NewCookieStore([]byte("kpl-analytics-secret-key-2026"))

// hashPassword превращает пароль в безопасный хеш
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// checkPassword сверяет введённый пароль с хешем из базы
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// getCurrentUser возвращает логин вошедшего пользователя (или пустую строку)
func getCurrentUser(r *http.Request) string {
	session, _ := store.Get(r, "session")
	username, ok := session.Values["username"].(string)
	if !ok {
		return ""
	}
	return username
}

// loginUser сохраняет пользователя в сессию (выдаёт "пропуск")
func loginUser(w http.ResponseWriter, r *http.Request, username string) {
	session, _ := store.Get(r, "session")
	session.Values["username"] = username
	session.Save(r, w)
}

// logoutUser удаляет пользователя из сессии
func logoutUser(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	delete(session.Values, "username")
	session.Save(r, w)
}
