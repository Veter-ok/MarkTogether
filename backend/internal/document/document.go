package document

import "time"

type Document struct {
	ID         string
	createDate time.Time
	title      string
	text       string
}

func NewDocument(title string) *Document {
	return &Document{
		ID:         title,
		createDate: time.Now(),
		title:      title,
		text:       "",
	}
}
