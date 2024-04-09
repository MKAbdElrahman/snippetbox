package service

import (
	"github.com/mkabdelrahman/snippetbox/db"
	"github.com/mkabdelrahman/snippetbox/model"
)

type SnippetService struct {
	store *db.SnippetStore
}

func NewSnippetService(store *db.SnippetStore) *SnippetService {
	return &SnippetService{
		store: store,
	}
}

func (s *SnippetService) GetById(id int) (*model.Snippet, error) {
	return s.store.GetById(id)
}

func (s *SnippetService) Insert(params model.NewSnippetParams) (*model.Snippet, error) {
	return s.store.Insert(params)
}

func (s *SnippetService) GetLatestSnippets(limit int) ([]model.Snippet, error) {
	return s.store.GetLatestSnippets(limit)
}
