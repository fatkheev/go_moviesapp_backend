package storage

import (
	"database/sql"
	"filmoteca/internal/model"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
    // "golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

// DB - глобальная переменная для хранения соединения с базой данных
var DB *sql.DB

// InitDB инициализирует соединение с базой данных
func InitDB() {
	// Загрузка переменных окружения из файла .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки файла .env")
	}

	connectionString := fmt.Sprintf("host=%s port=5432 user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		// os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v", err)
	}

	fmt.Println("Соединение с базой данных успешно установлено.")
}

// Функция для получения актера по ID
func GetActorByID(id int) (model.Actor, error) {
	var actor model.Actor
	query := `SELECT actor_id, name, gender, birthdate FROM actors WHERE actor_id = $1`
	row := DB.QueryRow(query, id)
	err := row.Scan(&actor.ID, &actor.Name, &actor.Gender, &actor.Birthdate)
	if err != nil {
		return model.Actor{}, err
	}
	return actor, nil
}

// Функция для добавления актера:
func AddActor(actor model.Actor) error {
	query := `INSERT INTO actors (name, gender, birthdate) VALUES ($1, $2, $3)`
	_, err := DB.Exec(query, actor.Name, actor.Gender, actor.Birthdate)
	return err
}

// Функция для изменения данных актера:
func UpdateActor(actor model.Actor) error {
	// Составляем SQL-запрос с возможностью обновления всех полей актёра
	query := `
    UPDATE actors
    SET name = $2, gender = $3, birthdate = $4
    WHERE actor_id = $1
    `
	_, err := DB.Exec(query, actor.ID, actor.Name, actor.Gender, time.Time(actor.Birthdate))
	return err
}

// Функция для удаления актера по ID
func DeleteActorByID(id int) error {
	query := `DELETE FROM actors WHERE actor_id = $1`
	_, err := DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// Функция для добавления фильма
func AddMovie(movie model.Movie) (int, error) {
    // Сначала проверяем наличие всех актеров в базе данных
    for _, actorID := range movie.ActorIDs {
        var exists bool
        err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM actors WHERE actor_id = $1)", actorID).Scan(&exists)
        if err != nil {
            // Логирование ошибки запроса
            log.Printf("Ошибка при проверке актера с ID %d: %v", actorID, err)
            return 0, err
        }
        if !exists {
            // Если актер не найден, возвращаем ошибку
            errorMsg := fmt.Sprintf("Актер с ID %d не найден в базе данных", actorID)
            log.Println(errorMsg)
            return 0, fmt.Errorf(errorMsg)
        }
    }
    
    // Теперь добавляем фильм и связи с актерами как раньше
    var movieID int
    query := `INSERT INTO movies (title, description, release_date, rating) VALUES ($1, $2, $3, $4) RETURNING movie_id`
    err := DB.QueryRow(query, movie.Title, movie.Description, movie.ReleaseDate, movie.Rating).Scan(&movieID)
    if err != nil {
        return 0, err
    }

    // Добавление связей с актерами
    for _, actorID := range movie.ActorIDs {
        query = `INSERT INTO movies_actors (movie_id, actor_id) VALUES ($1, $2)`
        _, err = DB.Exec(query, movieID, actorID)
        if err != nil {
            return 0, err
        }
    }

    return movieID, nil
}


// Функция для обновления информации о фильме
func UpdateMovie(movie model.Movie) error {
	setParts := []string{}
	args := []interface{}{}
	argId := 1

	if movie.Title != "" {
		setParts = append(setParts, "title = $"+strconv.Itoa(argId))
		args = append(args, movie.Title)
		argId++
	}
	if movie.Description != "" {
		setParts = append(setParts, "description = $"+strconv.Itoa(argId))
		args = append(args, movie.Description)
		argId++
	}
	if !movie.ReleaseDate.IsZero() {
		setParts = append(setParts, "release_date = $"+strconv.Itoa(argId))
		args = append(args, movie.ReleaseDate)
		argId++
	}
	if movie.Rating != 0 {
		setParts = append(setParts, "rating = $"+strconv.Itoa(argId))
		args = append(args, movie.Rating)
		argId++
	}

	if len(setParts) > 0 {
		queryString := "UPDATE movies SET " + strings.Join(setParts, ", ") + " WHERE movie_id = $" + strconv.Itoa(argId)
		args = append(args, movie.ID)
		_, err := DB.Exec(queryString, args...)
		if err != nil {
			return err
		}
	}

	// Обработка обновления связей актёров, если список ActorIDs предоставлен
	if len(movie.ActorIDs) > 0 {
		_, err := DB.Exec("DELETE FROM movies_actors WHERE movie_id = $1", movie.ID)
		if err != nil {
			return err
		}

		for _, actorID := range movie.ActorIDs {
			_, err = DB.Exec("INSERT INTO movies_actors (movie_id, actor_id) VALUES ($1, $2)", movie.ID, actorID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Функция для удаления фильма
func DeleteMovie(movieID int) error {
	// Удаление связей фильма с актёрами
	_, err := DB.Exec("DELETE FROM movies_actors WHERE movie_id = $1", movieID)
	if err != nil {
		return err
	}

	// Удаление самого фильма
	_, err = DB.Exec("DELETE FROM movies WHERE movie_id = $1", movieID)
	if err != nil {
		return err
	}

	return nil
}

// Пример функции для получения списка всех актеров:
func GetActors() ([]model.Actor, error) {
	var actors []model.Actor
	query := `SELECT id, name, gender, birthdate FROM actors`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a model.Actor
		if err := rows.Scan(&a.ID, &a.Name, &a.Gender, &a.Birthdate); err != nil {
			return nil, err
		}
		actors = append(actors, a)
	}
	return actors, nil
}

// Функция для получения списка фильмов по заданным параметрам
func GetMovies(sortBy, searchQuery, searchActor string) ([]model.Movie, error) {
	var movies []model.Movie

	// Определение базового запроса
	baseQuery := `
    SELECT DISTINCT m.movie_id, m.title, m.description, m.release_date, m.rating
    FROM movies m
    LEFT JOIN movies_actors ma ON m.movie_id = ma.movie_id
    LEFT JOIN actors a ON ma.actor_id = a.actor_id
    WHERE 1=1
    `
	queryParams := []interface{}{}

	// Добавление условий поиска, если они заданы
	if searchQuery != "" {
		baseQuery += " AND m.title LIKE $1"
		queryParams = append(queryParams, "%"+searchQuery+"%")
	}

	if searchActor != "" {
		baseQuery += " AND a.name LIKE $" + fmt.Sprint(len(queryParams)+1)
		queryParams = append(queryParams, "%"+searchActor+"%")
	}

	// Добавление сортировки
	orderBy := " ORDER BY m.rating DESC" // Сортировка по умолчанию
	switch sortBy {
	case "title":
		orderBy = " ORDER BY m.title ASC"
	case "release_date":
		orderBy = " ORDER BY m.release_date DESC"
	case "rating":
		orderBy = " ORDER BY m.rating DESC"
	}
	baseQuery += orderBy

	// Выполнение запроса
	rows, err := DB.Query(baseQuery, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m model.Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.ReleaseDate, &m.Rating); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}

// Функция для получения списка актеров с фильмами
func GetActorsWithMovies() ([]model.Actor, error) {
    // Инициализация переменных и подготовка запроса
    var actors []model.Actor

    query := `
    SELECT a.actor_id, a.name, m.movie_id, m.title
    FROM actors a
    LEFT JOIN movies_actors ma ON a.actor_id = ma.actor_id
    LEFT JOIN movies m ON ma.movie_id = m.movie_id
    ORDER BY a.actor_id
    `
    rows, err := DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    currentActorID := 0
    var currentActor *model.Actor
    for rows.Next() {
        var actorID int
        var name string
        var movieID sql.NullInt64
        var movieTitle sql.NullString
        err := rows.Scan(&actorID, &name, &movieID, &movieTitle)
        if err != nil {
            return nil, err
        }
        
        if actorID != currentActorID {
            if currentActor != nil {
                actors = append(actors, *currentActor)
            }
            currentActor = &model.Actor{ID: actorID, Name: name, Movies: []model.Movie{}}
            currentActorID = actorID
        }
        if movieID.Valid { // Проверка на валидность movieID
            currentActor.Movies = append(currentActor.Movies, model.Movie{ID: int(movieID.Int64), Title: movieTitle.String})
        }
    }
    if currentActor != nil { // Добавление последнего актёра
        actors = append(actors, *currentActor)
    }

    return actors, nil
}

// AddUser добавляет нового пользователя в базу данных с хешированным паролем и id роли.
func AddUser(user model.User) error {
    _, err := DB.Exec("INSERT INTO users (username, password, role_id) VALUES ($1, $2, $3)", user.Username, user.Password, user.RoleID)
    if err != nil {
        log.Printf("Ошибка при добавлении пользователя: %v", err)
        return err
    }
    return nil
}

// GetUserByUsername получает пользователя по его имени
func GetUserByUsername(username string) (model.User, error) {
    var user model.User
    err := DB.QueryRow("SELECT user_id, username, password, role_id FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password, &user.RoleID)
    if err != nil {
        return model.User{}, err
    }
    return user, nil
}