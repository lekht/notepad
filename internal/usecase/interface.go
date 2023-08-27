package usecase

import (
	"context"

	"github.com/lekht/notepad/internal/models"
)

type NotepadRepository interface {
	CreateNote(ctx context.Context, note models.Note) error
	GetNotes(ctx context.Context, userId int) ([]models.Note, error)
}

type Authenticator interface {
	Authorization(id int) (string, error)
}

type Speller interface {
	Check(title, body string) (bool, error)
}
