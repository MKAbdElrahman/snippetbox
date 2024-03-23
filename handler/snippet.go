package handler

import "net/http"

type SnippetHandler struct{}

func NewSnippetHandler() *SnippetHandler {
	return &SnippetHandler{}
}

func (h *SnippetHandler) View(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display Snippet with Id = " + r.PathValue("id")))
}

func (h *SnippetHandler) ViewCreateForm(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a Form for Creating a New Snippet"))
}
