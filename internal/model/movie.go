package model

import (
    "filmoteca/services" // Импортирование нового пакета services
)

type Movie struct {
    ID          int               `json:"id"`
    Title       string            `json:"title"`
    Description string            `json:"description,omitempty"`
    ReleaseDate services.JSONDate `json:"release_date"`
    Rating      float64           `json:"rating"`
    ActorIDs    []int             `json:"actor_ids,omitempty"`
}