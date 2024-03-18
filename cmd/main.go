package main

import (
    "filmoteca/internal/handler"
    "filmoteca/internal/storage"
	 "filmoteca/middleware"
    "log"
    "net/http"
    "strings"
)

func main() {
	storage.InitDB()

	// Обертываем каждый обработчик функцией enableCORS для добавления заголовков CORS
	http.HandleFunc("/actors", middleware.EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Запрос на /actors с методом %s", r.Method)
		switch r.Method {
		case "POST":
			handler.AdminOnly(handler.AddActor)(w, r)
		case "GET":
			handler.GetActors(w, r)
		default:
			log.Printf("Метод %s не поддерживается для /actors", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/actors/", middleware.EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) != 2 {
			log.Printf("Некорректный запрос на %s: ожидался ID актёра", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		actorId := parts[1]

		switch r.Method {
		case "PUT":
			handler.AdminOnly(func(w http.ResponseWriter, r *http.Request) {
				handler.UpdateActor(actorId, w, r)
			})(w, r)
		case "DELETE":
			handler.AdminOnly(func(w http.ResponseWriter, r *http.Request) {
				handler.DeleteActor(actorId, w, r)
			})(w, r)
		default:
			log.Printf("Метод %s не поддерживается для %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	}))

	http.HandleFunc("/movies", middleware.EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Запрос на /movies с методом %s", r.Method)
		switch r.Method {
		case "POST":
			handler.AdminOnly(handler.AddMovie)(w, r)
		case "GET":
			handler.GetMovies(w, r)
		default:
			log.Printf("Метод %s не поддерживается для /movies", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))

	http.HandleFunc("/movies/", middleware.EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) != 2 {
			log.Printf("Некорректный запрос на %s: ожидался ID фильма", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		movieId := parts[1]

		switch r.Method {
		case "PUT":
			handler.AdminOnly(func(w http.ResponseWriter, r *http.Request) {
				handler.UpdateMovie(movieId, w, r)
			})(w, r)
		case "DELETE":
			handler.AdminOnly(func(w http.ResponseWriter, r *http.Request) {
				handler.DeleteMovie(movieId, w, r)
			})(w, r)
		default:
			log.Printf("Метод %s не поддерживается для %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	}))

	http.HandleFunc("/register", middleware.EnableCORS(handler.RegisterUser))
	http.HandleFunc("/login", middleware.EnableCORS(handler.LoginUser))

	log.Println("Сервер запущен на http://localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}