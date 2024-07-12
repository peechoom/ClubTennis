package services

import (
	"ClubTennis/repositories"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/h2non/bimg"
	"gorm.io/gorm"
)

// service for all public-facing resources
type PublicService struct {
	snippetRepo *repositories.SnippetRepository
}

// max of 10 mb
const MAX_IMG_SIZE int = 10 * 1024 * 1024
const SLIDE_COUNT = 5

func NewPublicService(db *gorm.DB) *PublicService {
	s := PublicService{snippetRepo: repositories.NewSnippetRepository(db)}
	if (&s).ensureSlideCount() != nil {
		return nil
	}
	return &s
}

// ensures there are 5 slides. fills in missing slides with nc state wallpaper
func (s *PublicService) ensureSlideCount() error {
	for i := 1; i <= SLIDE_COUNT; i++ {
		staticFilename := "static/slide" + strconv.FormatInt(int64(i), 10) + ".webp"
		persistFilename := os.Getenv("SERVER_FILES_MOUNTPOINT") + "/slide" + strconv.FormatInt(int64(i), 10) + ".webp"

		_, err := os.Stat(persistFilename)

		if errors.Is(err, os.ErrNotExist) {
			if e := ConvertToWebp("static/wallpaper.jpg", staticFilename); e != nil {
				return e
			}
		} else {
			if err = func() (err error) {
				var in, out *os.File
				if in, err = os.Open(persistFilename); err != nil {
					return err
				}
				defer in.Close()
				if out, err = os.Create(staticFilename); err != nil {
					return err
				}
				defer out.Close()
				_, err = io.Copy(out, in)
				return
			}(); err != nil {
				return err
			}
		}
	}
	return nil
}

// expects filename like /tmp/file and output filename like static/image.webp. must be webp
func ConvertToWebp(filename, outputName string) error {
	dir := filepath.Dir(outputName)
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return err
	}

	data, err := bimg.Read(filename)
	if err != nil {
		return err
	}
	newData, err := bimg.NewImage(data).Convert(bimg.WEBP)
	if err != nil {
		return err
	}
	if bimg.NewImage(newData).Type() != "webp" {
		return err
	}

	bimg.Write(outputName, newData)
	return nil
}
