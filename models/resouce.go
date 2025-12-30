package models

type Resource struct {
	UserID     int64
	ResourceID string
	Type       string
	Title      string
	URL        string
	Source     string
	Status     string
	Tags       []string
	Notes      string
	CreatedAt  string
}
