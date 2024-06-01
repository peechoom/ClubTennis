package models

import "gorm.io/gorm"

type Slide struct {
	gorm.Model
	SlideNum int    // up to 5
	Data     string // image data base64 w html src prefix
}
