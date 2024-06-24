package daemons

import (
	"ClubTennis/services"
	"log"
	"time"
)

const (
	// Check if the match is about to expire once every 6 hours
	WarningDefaultFrequency time.Duration = 6 * time.Hour
	// Auto-forefit the challenged player every day (I would prefer at a set time like 3am or smthn)
	ExpiredDefaultFrequency time.Duration = 24 * time.Hour
	// Auto-delete matches that are too old to matter every 2 days
	DeletionDefaultFrequency time.Duration = 48 * time.Hour
)

type MatchService = services.MatchService
type EmailService = services.EmailService

func MatchWarningDaemon(frequency time.Duration, autostart bool, ms *MatchService, es *EmailService) {
	if autostart {
		autostartSleep(4)
	}
	ticker := time.NewTicker(frequency)
	for range ticker.C {
		matches, err := ms.FindByNearlyExpired()
		if err == nil { //IF THERE IS NO ERROR
			for i := 0; i < len(matches); i++ {
				if matches[i].LateNotifSent {
					continue
				}
				email := es.MakeChallengeReminder(matches[i].Challenger(), matches[i].Challenged())

				err := es.Send(email)
				if err != nil {
					matches[i].LateNotifSent = true
				}
				ms.Save(&(matches[i]))
			}
		} else {
			log.Print("Match Expiration Daemon error: " + err.Error())
		}
	}
}

func MatchExpiredDaemon(frequency time.Duration, autostart bool, ms *MatchService, es *EmailService) {
	if autostart {
		autostartSleep(2)
	}
	ticker := time.NewTicker(frequency)
	for range ticker.C {
		matches, err := ms.FindByExpired()
		if err == nil { //IF THERE IS NO ERROR
			for i := 0; i < len(matches); i++ {
				matches[i].SubmitScore(6, 0)
				es.Send(es.MakeForfeitEmail(matches[i].Challenger(), matches[i].Challenged()))
				ms.Save(&(matches[i]))
			}
		} else {
			log.Print("Match Expiration Daemon error: " + err.Error())
		}
	}
}

func MatchDeletionDaemon(frequency time.Duration, autostart bool, ms *MatchService) {
	if autostart {
		autostartSleep(0)
	}
	ticker := time.NewTicker(frequency)
	for range ticker.C {
		err := ms.DeleteOldMatches()
		if err != nil {
			log.Println("Match Deletion Daemon error: " + err.Error())
		}
	}
}

// will return the next time it is {hour} o'clock
func autostartSleep(hour int) {
	location := time.FixedZone("UTC-8", -4*60*60)

	now := time.Now()
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, location)
	if targetTime.Before(now) {
		targetTime = targetTime.Add(24 * time.Hour)
	}
	d := targetTime.Sub(now)
	log.Println("daemon thread sleeping for " + d.String())
	time.Sleep(d)
}
