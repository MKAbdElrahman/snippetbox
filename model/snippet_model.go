package model

import "time"

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type NewSnippetParams struct {
	Title   string
	Content string
	Expires int
}
