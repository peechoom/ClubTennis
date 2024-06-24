package daemons

import (
	"ClubTennis/services"
	"log"
	"time"
)

type UserService = services.UserService

const AutoBackupDefaultFrequency = 48 * time.Hour

func AutoBackupDaemon(frequency time.Duration, autostart bool, es *EmailService, us *UserService) {
	if autostart {
		autostartSleep(23)
	}
	ticker := time.NewTicker(frequency)
	for range ticker.C {
		func() {
			users, err := us.FindAll()
			if err != nil {
				log.Println("error finding all users: " + err.Error())
				return
			}
			filename, err := services.UsersToSheet(users)
			if err != nil {
				log.Println("error making spreadsheet: " + err.Error())
				return
			}
			email := es.MakeBackupEmail(filename)
			if email == nil {
				log.Println("error making email: file not found or empty at location " + filename)
				return
			}
			if err = es.Send(email); err != nil {
				log.Print("error sending email: " + err.Error())
				return
			}
		}()
	}
}
