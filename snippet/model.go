package snippet

import "time"

type Model struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type NewModelParams struct {
	Title          string
	Content        string
	DaystToExpires int
}
