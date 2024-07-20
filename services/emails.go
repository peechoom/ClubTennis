package services

import (
	"ClubTennis/models"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var templates = map[string]string{
	"challenger":   "challenger_email.html",          //email to the challenger
	"challenged":   "challenged_email.html",          //email to the challenged
	"announcement": "announcement_email.html",        //email containing an announcement
	"reminder":     "challenged_reminder_email.html", //email reminder to the challenged
	"forfeit":      "forfeit_email.html",             //email indicating forfeit
}

const sentinalOpen string = "&%"
const sentinalClose string = "%&"

// populates the template ("challenge", "announcement") using the values v that maps $template$ -> value
func (s *EmailService) populateTemplate(template string, v map[string]string) (string, error) {
	fpath := filepath.Join(s.dirpath, templates[template])
	f, err := os.ReadFile(fpath)
	if err != nil {
		return "", err
	}

	// template assumed to be complete
	if v == nil {
		return string(f), nil
	}

	// regex matches all single symbols encased with &%...%&
	expression := regexp.MustCompile(`&%([^%]*?)%&`)

	symbols := expression.FindAllIndex(f, -1)
	// no symbols found
	if symbols == nil {
		return string(f), nil
	}

	var b strings.Builder

	l := len(symbols)
	var fileIndex int = 0
	//i is index into symbols
	for i := 0; i < l; i++ {
		b.WriteString(string(f[fileIndex:symbols[i][0]]))
		b.WriteString(v[string(f[symbols[i][0]:symbols[i][1]])])
		fileIndex = symbols[i][1]
	}

	u := f[fileIndex:]
	b.WriteString(string(u))

	return b.String(), nil
}

// makes the values map for populating a challenge email template
func challengeEmailMap(challenger, challenged *User, message string) (map[string]string, error) {
	v := commonEmailMap()
	if v == nil {
		return nil, errors.New("env variables not defined")
	}

	v[sentinalOpen+"challenger_firstname"+sentinalClose] = challenger.FirstName
	v[sentinalOpen+"challenger_lastname"+sentinalClose] = challenger.LastName

	v[sentinalOpen+"challenged_firstname"+sentinalClose] = challenged.FirstName
	v[sentinalOpen+"challenged_lastname"+sentinalClose] = challenged.LastName

	v[sentinalOpen+"challenger_email"+sentinalClose] = challenger.ContactEmail
	v[sentinalOpen+"challenged_email"+sentinalClose] = challenged.ContactEmail

	v[sentinalOpen+"challenger_rank"+sentinalClose] = strconv.FormatUint(uint64(challenger.Rank), 10)
	v[sentinalOpen+"challenged_rank"+sentinalClose] = strconv.FormatUint(uint64(challenged.Rank), 10)

	v[sentinalOpen+"challenger_message"+sentinalClose] = message

	for k, s := range v {
		if s == "" {
			return nil, fmt.Errorf("%s not populated", k)
		}
	}
	return v, nil
}

// lowkey does nothing. announcements dont have any metadata atm
func announcementEmailMap(ann *models.Announcement) (map[string]string, error) {
	isToday := func(date time.Time) bool {
		now := time.Now()
		return date.Year() == now.Year() && date.YearDay() == now.YearDay()
	}

	v := commonEmailMap()
	if v == nil {
		return nil, errors.New("env variables not defined")
	}
	v[sentinalOpen+"announcement_body"+sentinalClose] = ann.Data

	var prefix string = ""
	if isToday(ann.CreatedAt) {
		prefix = "Today, "
	}
	v[sentinalOpen+"upload_date"+sentinalClose] = prefix + ann.CreatedAt.Format("Monday, January 2, 2006 at 3:04 PM")

	for k, s := range v {
		if s == "" {
			return nil, fmt.Errorf("%s not populated", k)
		}
	}

	return v, nil
}

func commonEmailMap() map[string]string {
	host := os.Getenv("SERVER_HOST")

	if host == "" {
		return nil
	}
	v := make(map[string]string)

	link := "https://" + host + "/static/"
	v[sentinalOpen+"email_header"+sentinalClose] = link + "email_header.webp"
	v[sentinalOpen+"challenger_title_link"+sentinalClose] = link + "challenger.webp"
	v[sentinalOpen+"challenged_title_link"+sentinalClose] = link + "challenged.webp"
	v[sentinalOpen+"tuffy_head_small_link"+sentinalClose] = link + "tuffy-small.webp"
	v[sentinalOpen+"wolfpack_regular_link"+sentinalClose] = link + "Wolfpack-Regular.woff2"
	v[sentinalOpen+"united_sans_regbold_link"+sentinalClose] = link + "UnitedSansReg-Bold.woff2"
	v[sentinalOpen+"gotham_narrow_bold_link"+sentinalClose] = link + "GothamNarrow-Bold.woff2"
	v[sentinalOpen+"gotham_narrow_bold_italic_link"+sentinalClose] = link + "GothamNarrow-BoldItalic.woff2"
	v[sentinalOpen+"gotham_narrow_book_link"+sentinalClose] = link + "GothamNarrow-Book.woff2"
	v[sentinalOpen+"gotham_narrow_book_italic_link"+sentinalClose] = link + "GothamNarrow-BookItalic.woff2"
	v[sentinalOpen+"gotham_narrow_medium_link"+sentinalClose] = link + "GothamNarrow-Medium.woff2"
	v[sentinalOpen+"gotham_narrow_medium_italic_link"+sentinalClose] = link + "GothamNarrow-MediumItalic.woff2"

	return v
}
