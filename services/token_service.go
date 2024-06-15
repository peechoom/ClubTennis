package services

import (
	"ClubTennis/models"
	"ClubTennis/repositories"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mazen160/go-random"
)

// ID tokens are valid for 6 hours
const DefaultIDTokenLifetime int64 = 10 //6 * int64(time.Hour/time.Second)

// refresh tokens are valid for 3 days
const DefaultRefreshTokenLifetime int64 = 3 * 24 * int64(time.Hour/time.Second)

type TokenService struct {
	repo *repositories.TokenRepository

	PrivateKey           *rsa.PrivateKey // the private key for id tokens
	PublicKey            *rsa.PublicKey  //the public key for id tokens
	RefreshKey           []byte          // the symetric key for refresh keys
	IDTokenLifetime      int64           //seconds that an id token can live
	RefreshTokenLifetime int64           //seconds that a refresh token can live
}

func NewTokenService(repo *repositories.TokenRepository,
	privateKey *rsa.PrivateKey,
	publicKey *rsa.PublicKey,
	refreshKey []byte,
	iDLifetime int64,
	refreshLifetime int64) *TokenService {
	return &TokenService{
		repo:                 repo,
		PrivateKey:           privateKey,
		PublicKey:            publicKey,
		RefreshKey:           refreshKey,
		IDTokenLifetime:      iDLifetime,
		RefreshTokenLifetime: refreshLifetime,
	}
}

// just needs a repo, generates random rsa and symmetric keys by itself
func DefaultTokenService(repo *repositories.TokenRepository) *TokenService {
	const rsaBitSize int = 2048
	priv, err := rsa.GenerateKey(rand.Reader, rsaBitSize)
	if err != nil {
		return nil
	}

	sym, _ := random.Bytes(64)
	if len(sym) == 0 {
		return nil
	}

	return &TokenService{
		repo:                 repo,
		PrivateKey:           priv,
		PublicKey:            priv.Public().(*rsa.PublicKey),
		RefreshKey:           sym,
		IDTokenLifetime:      DefaultIDTokenLifetime,
		RefreshTokenLifetime: DefaultRefreshTokenLifetime,
	}
}

// generates new refresh and id tokens for the user once they log in, used for singing in
func (ts *TokenService) GetNewTokenPair(userID uint, prevTokenID string) (*models.TokenPair, error) {
	var err error
	var idString string
	var refresh *refreshTokenData
	userIDString := strconv.FormatUint(uint64(userID), 10)

	if prevTokenID != "" {
		err := ts.repo.DeleteRefreshToken(userIDString, prevTokenID)
		if err != nil {
			return nil, err
		}
	}

	idString, err = generateIDToken(userID, ts.PrivateKey, ts.IDTokenLifetime)
	if err != nil {
		return nil, err
	}

	refresh, err = generateRefreshToken(userID, ts.RefreshKey, ts.RefreshTokenLifetime)
	if err != nil {
		return nil, err
	}

	if userIDString == "" {
		return nil, errors.New("bruh, the uhhhh string conversion failed. how did you manage that???")
	}

	ts.repo.SetRefreshToken(userIDString, refresh.ID.String(), refresh.ExpiresIn)

	return &models.TokenPair{
		IDToken:      models.IDToken{SS: idString},
		RefreshToken: models.RefreshToken{ID: refresh.ID, UserID: userID, SS: refresh.SS},
	}, nil
}

// validates the raw id token string and returns the signed in user id if valid
func (ts *TokenService) ValidateIDToken(tokenString string) (uint, error) {
	var claims *idTokenClaims
	var err error
	claims, err = validateIDToken(tokenString, ts.PublicKey)
	if err != nil {
		return 0, err
	}
	if !validateClaimTimes(claims) {
		return 0, jwt.ErrTokenExpired
	}
	return claims.UserID, nil
}

// validates the raw refresh token string and returns it decrypted if valid
func (ts *TokenService) ValidateRefreshToken(tokenString string) (*models.RefreshToken, error) {
	var claims *refreshTokenClaims
	var err error
	claims, err = validateRefreshToken(tokenString, ts.RefreshKey)
	if err != nil {
		return nil, err
	}

	if !validateClaimTimes(claims) {
		return nil, jwt.ErrTokenExpired
	}

	tokenUUID, err := uuid.Parse(claims.ID)
	if err != nil {
		return nil, err
	}

	return &models.RefreshToken{
		ID:     tokenUUID,
		UserID: claims.UserID,
		SS:     tokenString,
	}, nil
}

func (ts *TokenService) DeleteRefreshToken(userID uint, tokenID string) {
	str := strconv.FormatUint(uint64(userID), 10)
	ts.repo.DeleteRefreshToken(str, tokenID)
}
