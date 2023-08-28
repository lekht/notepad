package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/lekht/notepad/internal/models"
	"github.com/lekht/notepad/pkg/postgres"
	"github.com/pkg/errors"
)

const (
	sqlAddNote  = `INSERT INTO storage.notes (title, body, user_id, created_at) VALUES ($1,$2,$3,$4);`
	sqlGetNotes = `SELECT id, title, body, user_id, created_at FROM storage.notes WHERE user_id = $1`
)

type NotesRepository struct {
	*postgres.PostgreDB
}

func New(pg *postgres.PostgreDB) *NotesRepository {
	return &NotesRepository{pg}
}

// Метод создания новой записи в базе данных
func (nr *NotesRepository) CreateNote(ctx context.Context, note models.Note) error {
	_, err := nr.Pool.Exec(ctx, sqlAddNote, note.Title, note.Body, note.UserId, note.CreatedAt)
	if err != nil {
		return errors.Wrap(err, "failed to insert note")
	}

	return nil
}

// Метод получения строк из базы данных
func (nr *NotesRepository) GetNotes(ctx context.Context, userId int) ([]models.Note, error) {
	rows, err := nr.Pool.Query(ctx, sqlGetNotes, userId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get notes' rows")
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		err = rows.Scan(&note.Id, &note.Title, &note.Body, &note.UserId, &note.CreatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan rows")
		}
		notes = append(notes, note)
	}

	if err = rows.Err(); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Wrap(err, "rows.Err()")
		}
	}

	return notes, nil
}
