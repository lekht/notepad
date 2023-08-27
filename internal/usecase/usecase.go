package usecase

import (
	"context"

	"github.com/lekht/notepad/internal/auth"
	"github.com/lekht/notepad/internal/models"
	"github.com/lekht/notepad/pkg/speller"
	"github.com/pkg/errors"
)

type NotepadUsecase struct {
	repo  NotepadRepository
	auth  Authenticator
	spell Speller
}

func New(nr NotepadRepository, a Authenticator, s Speller) *NotepadUsecase {
	return &NotepadUsecase{
		repo:  nr,
		auth:  a,
		spell: s,
	}
}

func (nu *NotepadUsecase) NewNote(n models.Note, token string) error {
	storedToken, err := nu.auth.Authorization(n.UserId)
	if err != nil {
		return errors.Wrap(err, "failed to authorize user")
	}
	if token != storedToken {
		return auth.ErrNoUser
	}

	ok, err := nu.spell.Check(n.Title, n.Body)
	if err != nil {
		return errors.Wrap(err, "failed to check text for spellig")
	}
	if !ok {
		return speller.ErrSpelling
	}

	err = nu.repo.CreateNote(context.Background(), n)
	if err != nil {
		return errors.Wrap(err, "failed to create note")
	}

	return nil
}

func (nu *NotepadUsecase) Notes(userId int, token string) ([]models.Note, error) {
	storedToken, err := nu.auth.Authorization(userId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to authorize user")
	}

	if token != storedToken {
		return nil, auth.ErrNoUser
	}

	notes, err := nu.repo.GetNotes(context.Background(), userId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get notes from db")
	}

	return notes, nil
}
