package controllers

import "github.com/lekht/notepad/internal/models"

type Usecase interface {
	NewNote(n models.Note, token string) error
	Notes(userId int, token string) ([]models.Note, error)
}

type Logger interface {
	Error(message interface{}, args ...interface{})
}
