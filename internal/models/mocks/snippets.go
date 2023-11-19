package mocks

import (
	"time"

	"poc.eduardo-luz.eu/internal/models"
)

var mockedSnippet = &models.Snippet{
	ID:      1,
	Title:   "Title",
	Content: "Content",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockedSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockedSnippet}, nil
}
