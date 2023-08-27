package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lekht/notepad/internal/auth"
	"github.com/lekht/notepad/internal/models"
	"github.com/lekht/notepad/pkg/speller"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrInternal   = errors.New("internal server error")
)

type notepadRoutes struct {
	u Usecase
	l Logger
}

func newNotepadRoutes(r *mux.Router, u Usecase, l Logger) {
	nr := &notepadRoutes{u, l}

	r.HandleFunc("/notepad", nr.getNotes).Methods(http.MethodGet)
	r.HandleFunc("/notepad", nr.newNote).Methods(http.MethodPost)
}

func (nr *notepadRoutes) newNote(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Token")
	userId := r.Header.Get("UserID")

	var note models.Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		nr.l.Error(err, "failed to decode request body")
		http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		nr.l.Error(err, "failed to convert user id")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		nr.l.Error(err, "failed to read UserID header")
		http.Error(w, ErrBadRequest.Error(), http.StatusBadRequest)
		return
	}

	note.UserId = id

	err = nr.u.NewNote(note, token)
	if err != nil {
		if errors.Is(err, auth.ErrNoUser) || errors.Is(err, speller.ErrSpelling) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		nr.l.Error(err, "failed to handle new note")
		http.Error(w, ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (nr *notepadRoutes) getNotes(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Token")
	userId := r.Header.Get("UserID")
	id, err := strconv.Atoi(userId)
	if err != nil {
		nr.l.Error(err, "failed to convert user id")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	notes, err := nr.u.Notes(id, token)
	if err != nil {
		if errors.Is(err, auth.ErrNoUser) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		nr.l.Error(err, "failed to get notes")
		http.Error(w, ErrInternal.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		http.Error(w, ErrInternal.Error(), http.StatusInternalServerError)
		return
	}
}
