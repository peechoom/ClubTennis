package services

import (
	"ClubTennis/repositories"

	"gorm.io/gorm"
)

type ServiceContainer struct {
	TokenService *TokenService
	UserService  *UserService
	MatchService *MatchService
}

func SetupServices(db *gorm.DB) *ServiceContainer {
	tokenService := DefaultTokenService(repositories.NewTokenRepository())
	userService := NewUserService(db)
	matchService := NewMatchService(db)
	return &ServiceContainer{
		TokenService: tokenService,
		UserService:  userService,
		MatchService: matchService}
}
