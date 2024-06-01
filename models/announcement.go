package models

import (
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

type Announcement struct {
	ID        uint      `gorm:"primarykey"`                  // had to yoink the gorm.Model so i could put my own struct tag in
	CreatedAt time.Time `gorm:"index:,sort:desc,type:btree"` //always sorted desc
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Data      string         //the HTML for this announcement. Can include encoded images. encoded in base64
}

// TODO migrations, controllers
func NewAnnouncement(data string) *Announcement {
	if !IsCleanHTML(data) {
		return nil
	}
	return &Announcement{Data: data}

}

func IsCleanHTML(data string) bool {
	resultChan := make(chan bool)
	var wg sync.WaitGroup

	var unsafes = []string{
		"<script>",
		"</script>",
		"onload",
		"onerror",
		"onclick",
		"onmouseover",
		"onfocus",
		"onblur",
		"onsubmit",
		"onchange",
		"onkeypress",
		"onkeydown",
		"onkeyup",
		"onmousedown",
		"onmouseup",
		"style",
		"<iframe>",
		"<object>",
		"<embed>",
		"<applet>",
		"<meta http-equiv=\"refresh\">",
		"<form>",
		"<input>",
		"<button>",
		"<select>",
		"<textarea>",
		"href",
		"<!--",
	}

	for _, s := range unsafes {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			resultChan <- strings.Contains(data, s)
		}(s)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var ret bool = false
	for b := range resultChan {
		if b {
			ret = b
		}
	}
	return !ret
}
