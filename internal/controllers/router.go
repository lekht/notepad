package controllers

import (
	"github.com/gorilla/mux"
)

func NewRouter(u Usecase, l Logger) *mux.Router {

	r := mux.NewRouter()

	newNotepadRoutes(r, u, l)

	return r
}
