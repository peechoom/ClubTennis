package services

import (
	"ClubTennis/repositories"
	"log"
	"os"

	"gorm.io/gorm"
)

type ServiceContainer struct {
	TokenService        *TokenService
	UserService         *UserService
	MatchService        *MatchService
	AnnouncementService *AnnouncementService
	PublicService       *PublicService
	EmailService        *EmailService
	ImageService        *ImageService
}

func SetupServices(db *gorm.DB, emailTemplatesDir string) *ServiceContainer {
	tokenService := DefaultTokenService(repositories.NewTokenRepository())
	userService := NewUserService(db)
	matchService := NewMatchService(db)
	announcementService := NewAnnouncementService(db)
	publicService := NewPublicService(db)
	ImageService := NewImageService(db)

	uname := os.Getenv("EMAIL_USERNAME")
	pw := os.Getenv("EMAIL_PASSWORD")

	emailService := NewEmailService(emailTemplatesDir, uname, pw)
	if emailService == nil {
		log.Fatalf("Could not load email service from directory %s, using email: %s, pw: %s", emailTemplatesDir, uname, pw)
	}
	return &ServiceContainer{
		TokenService:        tokenService,
		UserService:         userService,
		MatchService:        matchService,
		AnnouncementService: announcementService,
		PublicService:       publicService,
		EmailService:        emailService,
		ImageService:        ImageService,
	}
}
