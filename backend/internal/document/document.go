package document

import "time"

type Document struct {
	ID         string
	CreateDate time.Time
	Title      string
	Text       string
}

func NewDocument(title string) *Document {
	return &Document{
		ID:         title,
		CreateDate: time.Now(),
		Title:      title,
		Text:       "",
	}
}
