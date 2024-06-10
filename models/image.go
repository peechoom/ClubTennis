package models

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"log"

	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	Data      []byte //raw image binary
	FileName  string `gorm:"index:unique,length:40"` // {hash}.{png/jpeg/gif}
	Extension string // file extension
}

// creates new image with auto generated filename. fileExtension must be of the form "png"
func NewImage(base64Data string, fileExtension string) *Image {
	imageData := make([]byte, base64.StdEncoding.DecodedLen(len(base64Data)))
	n, err := base64.StdEncoding.Decode(imageData, []byte(base64Data))
	if err != nil {
		log.Print(err.Error())
		return nil
	}
	imageData = imageData[:n]
	return &Image{Data: imageData, FileName: GenerateFileName(imageData, fileExtension), Extension: fileExtension}
}

func GenerateFileName(data []byte, fileExtension string) string {
	hash := md5.Sum(data)
	hexName := hex.EncodeToString(hash[:])
	return (hexName + "." + fileExtension)
}
