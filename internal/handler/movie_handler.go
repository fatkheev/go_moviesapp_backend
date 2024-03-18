package handler

import (
    "encoding/json"
    "net/http"
    "strconv"
    "filmoteca/internal/model"
    "filmoteca/internal/storage" // Импортируем пакет storage
)

// AddMovie обрабатывает POST запросы для добавления нового фильма
func AddMovie(w http.ResponseWriter, r *http.Request) {
    var movie model.Movie
    if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
        http.Error(w, "Ошибка при декодировании JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    movieID, err := storage.AddMovie(movie)
    if err != nil {
        http.Error(w, "Ошибка при добавлении фильма: "+err.Error(), http.StatusInternalServerError)
        return
    }
    movie.ID = movieID

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(movie)
}

// UpdateMovie обрабатывает PUT запросы для обновления информации о фильме
func UpdateMovie(movieId string, w http.ResponseWriter, r *http.Request) {
    // Конвертация movieId из строки в int
    id, err := strconv.Atoi(movieId)
    if err != nil {
        http.Error(w, "Invalid movie ID", http.StatusBadRequest)
        return
    }

    // Декодирование тела запроса в структуру Movie
    var movie model.Movie
    if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    movie.ID = id // Установка ID фильма

    // Вызов функции обновления фильма
    if err := storage.UpdateMovie(movie); err != nil {
        http.Error(w, "Failed to update movie: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(movie)
}

// DeleteMovie обрабатывает DELETE запросы для удаления фильма
func DeleteMovie(movieId string, w http.ResponseWriter, r *http.Request) {
    // Конвертация movieId из строки в int
    id, err := strconv.Atoi(movieId)
    if err != nil {
        http.Error(w, "Invalid movie ID", http.StatusBadRequest)
        return
    }

    // Вызов функции удаления фильма из пакета storage
    if err := storage.DeleteMovie(id); err != nil {
        http.Error(w, "Failed to delete movie: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    // Можно также отправить сообщение об успешном удалении, если это необходимо
    w.Write([]byte("Movie deleted successfully"))
}

// GetActors обрабатывает GET запросы для получения списка актеров
func GetActors(w http.ResponseWriter, r *http.Request) {
    actors, err := storage.GetActorsWithMovies()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(actors)
}

// GetMovies обрабатывает GET запросы для получения списка фильмов
func GetMovies(w http.ResponseWriter, r *http.Request) {
    sortBy := r.URL.Query().Get("sort")
    searchQuery := r.URL.Query().Get("title")
    searchActor := r.URL.Query().Get("actor")

    movies, err := storage.GetMovies(sortBy, searchQuery, searchActor)
    if err != nil {
        http.Error(w, "Failed to get movies: "+err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(movies)
}

