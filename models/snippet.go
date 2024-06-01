package models

import "gorm.io/gorm"

type Snippet struct {
	gorm.Model
	Category string `gorm:"index:,sort:desc"`
	Data     string // raw html snippet
}

func NewSnippet(category, data string) *Snippet {
	if !IsCleanHTML(data) {
		return nil
	}
	return &Snippet{Category: category, Data: data}
}
