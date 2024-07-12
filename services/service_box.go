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
	SnippetService      *SnippetService
	ResetService        *ResetService
}

func SetupServices(db *gorm.DB, emailTemplatesDir string) *ServiceContainer {
	tokenService := DefaultTokenService(repositories.NewTokenRepository())
	if tokenService == nil {
		panic("could not create token service")
	}
	userService := NewUserService(db)
	if userService == nil {
		panic("could not create user service")
	}
	matchService := NewMatchService(db)
	if matchService == nil {
		panic("could not create match service")
	}
	announcementService := NewAnnouncementService(db)
	if announcementService == nil {
		panic("could not create announcement service")
	}
	publicService := NewPublicService(db)
	if publicService == nil {
		panic("could not create public service")
	}
	imageService := NewImageService(db)
	if imageService == nil {
		panic("could not create image service")
	}
	snippetService := NewSnippetService(db)
	if snippetService == nil {
		panic("could not create snippet service")
	}
	resetService := NewResetService(db)
	if resetService == nil {
		panic("could not create reset service")
	}

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
		ImageService:        imageService,
		SnippetService:      snippetService,
		ResetService:        resetService,
	}
}
