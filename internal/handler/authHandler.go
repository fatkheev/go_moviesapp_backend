package handler

import (
	"encoding/json"
	"filmoteca/internal/model"
	"filmoteca/internal/storage"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
	"strings"
)

var jwtKey = []byte("kjhkhfskhfskjdh345345kjhkh")

type Claims struct {
	Username string `json:"username"`
	RoleID   int    `json:"role_id"`
	jwt.StandardClaims
}

// HashPassword генерирует хеш пароля с использованием bcrypt
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		 return "", err
	}
	return string(hashedPassword), nil
}

// ComparePasswords сравнивает хешированный пароль с введенным паролем
func ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// RegisterUser обрабатывает запросы на регистрацию пользователя.
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	// Декодирование пользователя из запроса
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		 log.Printf("Ошибка при декодировании пользователя: %v", err)
		 http.Error(w, err.Error(), http.StatusBadRequest)
		 return
	}

	log.Printf("Регистрация: принят пользователь %s с паролем %s", user.Username, user.Password)

	// Хеширование пароля
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		 log.Printf("Ошибка при хешировании пароля: %v", err)
		 http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		 return
	}

	log.Printf("Регистрация: хеш пароля для пользователя %s: %s", user.Username, hashedPassword)

	// Сохранение пользователя в базу данных
	user.Password = hashedPassword
	if err := storage.AddUser(user); err != nil {
		 log.Printf("Ошибка при добавлении пользователя в базу данных: %v", err)
		 http.Error(w, "Failed to register user", http.StatusInternalServerError)
		 return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "user created"})
	log.Printf("Регистрация успешна: пользователь %s", user.Username)
}


// LoginUser обрабатывает запросы на аутентификацию пользователя.
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var creds model.User
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		 log.Printf("Ошибка при декодировании учетных данных: %v", err)
		 http.Error(w, "Bad Request", http.StatusBadRequest)
		 return
	}

	log.Printf("Аутентификация: попытка входа для пользователя %s", creds.Username)

	user, err := storage.GetUserByUsername(creds.Username)
	if err != nil {
		 log.Printf("Ошибка получения пользователя из базы данных: %v", err)
		 http.Error(w, "User not found", http.StatusUnauthorized)
		 return
	}

	if !ComparePasswords(user.Password, creds.Password) {
		 log.Printf("Неверные учетные данные для пользователя %s", creds.Username)
		 http.Error(w, "Bad credentials", http.StatusUnauthorized)
		 return
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		 Username: user.Username,
		 RoleID:   user.RoleID,
		 StandardClaims: jwt.StandardClaims{
			  ExpiresAt: expirationTime.Unix(),
		 },
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		 http.Error(w, "Error while signing the token", http.StatusInternalServerError)
		 return
	}

	// Устанавливаем токен в куки
	http.SetCookie(w, &http.Cookie{
		 Name:    "token",
		 Value:   tokenString,
		 Expires: expirationTime,
	})

	// Возвращаем информацию о пользователе, токене и дате его истечения в формате JSON
	response := struct {
		 User      string    `json:"user"`
		 Token     string    `json:"token"`
		 ExpiresAt time.Time `json:"expires_at"`
	}{
		 User:      user.Username,
		 Token:     tokenString,
		 ExpiresAt: expirationTime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	log.Printf("Аутентификация успешна для пользователя %s. Токен выдан.", user.Username)
}

// AdminOnly является middleware, который проверяет, является ли пользователь администратором.
func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		 // Извлечение токена из запроса
		 tokenString := r.Header.Get("Authorization")
		 log.Printf("Попытка доступа с токеном: %s", tokenString)
		 tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		 // Парсинг токена
		 claims := &Claims{}
		 token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			  return jwtKey, nil
		 })

		 if err != nil || !token.Valid {
			  http.Error(w, "Invalid token", http.StatusUnauthorized)
			  return
		 }

		 // Проверка роли пользователя
		 if claims.RoleID != 1 { // предполагается, что ID роли админа - "1"
			  http.Error(w, "Access denied", http.StatusForbidden)
			  return
		 }

		 next(w, r)
	}
}