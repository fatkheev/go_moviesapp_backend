package handler

import (
	"encoding/json"
	"filmoteca/internal/model"
	"filmoteca/internal/storage"
	"net/http"
	"strconv"
)

// AddActor обрабатывает POST запросы для добавления информации об актёре
func AddActor(w http.ResponseWriter, r *http.Request) {
	var actor model.Actor
	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		http.Error(w, "Ошибка при декодировании JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Логика для проверки, что имя и дата рождения не пустые
	if actor.Name == "" || actor.Birthdate.IsZero() {
		http.Error(w, "Имя и дата рождения не могут быть пустыми", http.StatusBadRequest)
		return
	}

	if err := storage.AddActor(actor); err != nil {
		http.Error(w, "Ошибка при добавлении актера в базу данных: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(actor)
}

// UpdateActor обрабатывает PUT запросы для изменения информации об актёре
func UpdateActor(actorId string, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(actorId) // Конвертируем actorId из строки в int
	if err != nil {
		http.Error(w, "Некорректный ID актёра", http.StatusBadRequest)
		return
	}

	var updatedActor model.Actor
	if err := json.NewDecoder(r.Body).Decode(&updatedActor); err != nil {
		http.Error(w, "Ошибка при декодировании JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// В этот момент мы уже знаем, что ID актёра есть и он корректный
	currentActor, err := storage.GetActorByID(id) // Используем id, а не updatedActor.ID
	if err != nil {
		http.Error(w, "Актёр не найден: "+err.Error(), http.StatusNotFound)
		return
	}

	// Применяем обновления к currentActor на основе данных из updatedActor
	if updatedActor.Name != "" {
		currentActor.Name = updatedActor.Name
	}
	if !updatedActor.Birthdate.IsZero() {
		currentActor.Birthdate = updatedActor.Birthdate
	}
	if updatedActor.Gender != "" { // Добавьте эту проверку
		currentActor.Gender = updatedActor.Gender
	}

	if err := storage.UpdateActor(currentActor); err != nil {
		http.Error(w, "Ошибка при обновлении актёра в базе данных: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currentActor) // Отправляем обновлённые данные
}

// DeleteActor обрабатывает DELETE запросы для удаления актёра по ID.
func DeleteActor(actorId string, w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(actorId) // Конвертируем actorId из строки в int
	if err != nil {
		http.Error(w, "Некорректный ID актёра", http.StatusBadRequest)
		return
	}

	// Вызываем функцию удаления актёра из пакета storage
	err = storage.DeleteActorByID(id)
	if err != nil {
		// Если актёр не найден, возвращаем 404 Not Found
		// В зависимости от реализации storage.DeleteActorByID, вы можете проверять здесь конкретную ошибку
		http.Error(w, "Ошибка при удалении актёра: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Если удаление прошло успешно, отправляем JSON-ответ
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Actor successfully deleted"}
	json.NewEncoder(w).Encode(response)
}
