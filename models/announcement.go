package models

import (
	"time"

	"gorm.io/gorm"
)

type Announcement struct {
	ID        uint      `gorm:"primarykey"`                  // had to yoink the gorm.Model so i could put my own struct tag in
	CreatedAt time.Time `gorm:"index:,sort:desc,type:btree"` //always sorted desc
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Subject   string         `gorm:"-:all"` // the subject line for when this announcement is sent over email. not saved
	Data      string         //the HTML for this announcement. Can include encoded images. encoded in base64
}

func NewAnnouncement(data, subject string) *Announcement {
	if !IsCleanHTML(data) || len(subject) > 900 {
		return nil
	}
	return &Announcement{Data: data, Subject: subject}

}

const bytes_in_megabyte int = 1000000

// return the size of the buffer in megabytes
func (ann *Announcement) Size() int {
	return len(ann.Data) / bytes_in_megabyte
}
