package services

import (
	"ClubTennis/repositories"

	"gorm.io/gorm"
)

type ServiceContainer struct {
	TokenService        *TokenService
	UserService         *UserService
	MatchService        *MatchService
	AnnouncementService *AnnouncementService
	PublicService       *PublicService
}

func SetupServices(db *gorm.DB) *ServiceContainer {
	tokenService := DefaultTokenService(repositories.NewTokenRepository())
	userService := NewUserService(db)
	matchService := NewMatchService(db)
	announcementService := NewAnnouncementService(db)
	publicService := NewPublicService(db)
	return &ServiceContainer{
		TokenService:        tokenService,
		UserService:         userService,
		MatchService:        matchService,
		AnnouncementService: announcementService,
		PublicService:       publicService,
	}
}
