package services

import (
	"ClubTennis/models"
	"fmt"
	"log"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"time"

	"github.com/jordan-wright/email"
)

type EmailService struct {
	dirpath       string
	senderAddress string
	auth          smtp.Auth
}

const max_email_recipients int = 99

// dirpath is the path to the directory that contains the email templates. returns nil if any template file not found
func NewEmailService(dirpath string, senderAddress string, password string) *EmailService {
	info, err := os.Stat(dirpath)
	if os.IsNotExist(err) || !info.IsDir() {
		return nil
	}
	for _, f := range templates {
		path := filepath.Join(dirpath, f)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			log.Printf("file %s not found", path)
			return nil
		}
	}

	return &EmailService{
		dirpath:       dirpath,
		senderAddress: senderAddress,
		auth:          smtp.PlainAuth("", senderAddress, password, "smtp.gmail.com"),
	}
}

// sends an email
func (s *EmailService) Send(e *email.Email) error {
	return e.Send("smtp.gmail.com:587", s.auth)
}

// makes an indicator/reciept email for the challenger and a notification for the challenged
func (s *EmailService) MakeChallengeEmails(challenger, challenged *User, message string) (cr *email.Email, cd *email.Email) {
	v, err := challengeEmailMap(challenger, challenged, message)
	if err != nil {
		return nil, nil
	}
	return s.makeChallengerEmail(challenger, challenged, v), s.makeChallengedEmail(challenger, challenged, v)
}

// the email that goes to the challenger
func (s *EmailService) makeChallengerEmail(challenger, challenged *User, v map[string]string) *email.Email {
	htmlBody, err := s.populateTemplate("challenger", v)
	if err != nil {
		return nil
	}

	return &email.Email{
		To:      []string{challenger.ContactEmail},
		From:    fmt.Sprintf("NC State Club Tennis <%s>", s.senderAddress),
		Cc:      []string{s.senderAddress},
		Subject: fmt.Sprintf("You successfully challenged %s %s", challenged.FirstName, challenged.LastName),
		HTML:    []byte(htmlBody),
		Text:    []byte(fmt.Sprintf("You successfully challenged %s %s (%s). You should expect an email from them soon regarding scheduling.", challenged.FirstName, challenged.LastName, challenged.ContactEmail)),
	}

}

// the email that goes to the challenged
func (s *EmailService) makeChallengedEmail(challenger, challenged *User, v map[string]string) *email.Email {
	htmlBody, err := s.populateTemplate("challenged", v)
	if err != nil {
		return nil
	}

	headers := textproto.MIMEHeader{}
	headers.Add("Reply-To", fmt.Sprintf("%s %s <%s>", challenger.FirstName, challenger.LastName, challenger.ContactEmail))

	return &email.Email{
		To:      []string{challenged.ContactEmail},
		From:    fmt.Sprintf("NC State Club Tennis <%s>", s.senderAddress),
		Cc:      []string{s.senderAddress},
		Headers: headers,
		Subject: "Club Tennis Challenge Match",
		HTML:    []byte(htmlBody),
		Text:    []byte(fmt.Sprintf("You have been challenged by %s %s (%s). Reply to this email to contact them for scheduling.", challenger.FirstName, challenger.LastName, challenger.ContactEmail)),
	}
}

// because an announcement may have >100 users, we make copies to send to users [0-99], [100,199], etc
func (s *EmailService) MakeAnnouncementEmail(ann *models.Announcement, recipients []User) []*email.Email {
	var ret []*email.Email

	//recursive, call MakeAnnouncementEmail on head (arbitrary len) and then proceed to tail
	if len(recipients) > max_email_recipients {
		cut := len(recipients) - max_email_recipients
		head := recipients[:cut]
		ret = s.MakeAnnouncementEmail(ann, head)

		recipients = recipients[cut:]
	}

	v, err := announcementEmailMap(ann)
	if err != nil {
		print(err.Error())
		ret = append(ret, nil)
		return ret
	}

	htmlBody, err := s.populateTemplate("announcement", v)
	if err != nil {
		print(err.Error())
		ret = append(ret, nil)
		return ret
	}
	e := s.stdHeader(mapSlice(recipients, func(u User) *User { return &u })...)
	e.HTML = []byte(htmlBody)
	e.Text = []byte("A new announcement has been posted.") //TODO improve

	if ann.Subject == "" {
		e.Subject = "Announcement from NC State Club Tennis"
	} else {
		e.Subject = ann.Subject
	}

	ret = append(ret, e)
	return ret
}

// makes an email to the challenged reminding them of the challenge match
func (s *EmailService) MakeChallengeReminder(challenger, challenged *User) *email.Email {
	v, err := challengeEmailMap(challenger, challenged, "")
	if v == nil || err != nil {
		panic("v")
	}
	htmlBody, err := s.populateTemplate("reminder", v)
	if err != nil {
		return nil
	}
	headers := textproto.MIMEHeader{}
	headers.Add("Reply-To", fmt.Sprintf("%s %s <%s>", challenger.FirstName, challenger.LastName, challenger.ContactEmail))

	return &email.Email{
		To:      []string{challenged.ContactEmail},
		From:    fmt.Sprintf("NC State Club Tennis <%s>", s.senderAddress),
		Cc:      []string{s.senderAddress},
		Headers: headers,
		Subject: "Challenge Match Reminder",
		HTML:    []byte(htmlBody),
		Text:    []byte(fmt.Sprintf("This is an automated reminder that you have one day left to play %s %s (%s) to defend your position in the ladder, or else the match is automatically considered a forefit.", challenger.FirstName, challenger.LastName, challenger.ContactEmail)),
	}
}

// makes an email telling the user that the forfeit has heppened
func (s *EmailService) MakeForfeitEmail(challenger, challenged *User) *email.Email {
	v, err := challengeEmailMap(challenger, challenged, "")
	if v == nil || err != nil {
		return nil
	}
	htmlBody, err := s.populateTemplate("forfeit", v)
	if err != nil {
		return nil
	}
	email := s.stdHeader(challenger, challenged)
	email.HTML = []byte(htmlBody)
	email.Text = []byte("This email is sent as indication that the challenger has won due to an auto-forfeit.")

	return email
}

func (s *EmailService) stdHeader(recipients ...*User) *email.Email {
	return &email.Email{
		To:   mapSlice(recipients, func(u *User) string { return u.ContactEmail }),
		From: fmt.Sprintf("Club Tennis <%s>", s.senderAddress),
		Cc:   []string{s.senderAddress},
	}
}

func mapSlice[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

func (s *EmailService) MakeBackupEmail(filename string) *email.Email {
	timestring := time.Now().Format("Monday, Jan 2")
	contents, err := os.ReadFile(filename)
	if err != nil || len(contents) == 0 {
		return nil
	}

	a := &email.Attachment{
		Filename:    "ladder.xlsx",
		ContentType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		Content:     contents,
		Header:      textproto.MIMEHeader{},
	}

	return &email.Email{
		To:      []string{s.senderAddress},
		From:    fmt.Sprintf("Club Tennis app <%s>", s.senderAddress),
		Subject: timestring + " Users Backup",
		Text: []byte("Attached is a spreadsheet backup of the club ladder and members from " + timestring +
			". This file can be used to reload the ladder in case the database needs resetting"),
		Attachments: []*email.Attachment{a},
	}
}
