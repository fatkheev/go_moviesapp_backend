package storage

import "filmoteca/internal/model"

type Storage interface {
    AddActor(actor model.Actor) error
}