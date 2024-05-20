package services

// uses a lot of code from https://github.com/JacobSNGoodwin/memrizr/blob/master/account/service/tokens.go#L22
import (
	"crypto/rsa"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// id token claims, is the user and registered claims
type idTokenClaims struct {
	UserID uint `json:"uID"`
	jwt.RegisteredClaims
}

// refresh token claims only needs the users id in the sql db
type refreshTokenClaims struct {
	UserID uint `json:"uID"`
	jwt.RegisteredClaims
}

// contains the signed string of the refresh token JWT as well as the ID and expiry for your convienience
type refreshTokenData struct {
	ID        uuid.UUID     // the uuid of the refresh token
	SS        string        // refresh token signed string
	ExpiresIn time.Duration // amount of time until this token expires
}

// returns the signed string representing the rs256-signed jwt token given the server's rsa secret and nanoseconds until expiry
func generateIDToken(u uint, key *rsa.PrivateKey, expires int64) (string, error) {
	currentTime := time.Now()
	claims := idTokenClaims{
		u,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Duration(expires) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			NotBefore: jwt.NewNumericDate(currentTime),
			//add more as needed
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return ss, nil
}

// generates a new refresh token using the secret symetric key and nanoseconds until expiry
func generateRefreshToken(u uint, key []byte, expires int64) (*refreshTokenData, error) {
	currentTime := time.Now()
	exp := currentTime.Add(time.Duration(expires) * time.Second)

	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	claims := refreshTokenClaims{
		u,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			NotBefore: jwt.NewNumericDate(currentTime),
			ID:        tokenID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(key))
	if err != nil {
		return nil, err
	}

	return &refreshTokenData{
		SS:        ss,
		ID:        tokenID,
		ExpiresIn: exp.Sub(currentTime),
	}, nil
}

// checks that an ID token is valid and returns the claims if it is
func validateIDToken(jwtString string, key *rsa.PublicKey) (*idTokenClaims, error) {
	claims := &idTokenClaims{}

	token, err := jwt.ParseWithClaims(jwtString, claims, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token not valid")
	}

	claims, ok := token.Claims.(*idTokenClaims)
	if !ok {
		return nil, errors.New("unknown claims type, cannot proceed")
	}

	return claims, nil
}

// checks that a refresh token is valid and returns the claims if it is
func validateRefreshToken(jwtString string, key []byte) (*refreshTokenClaims, error) {
	claims := &refreshTokenClaims{}

	token, err := jwt.ParseWithClaims(jwtString, claims, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token not valid")
	}

	claims, ok := token.Claims.(*refreshTokenClaims)
	if !ok {
		return nil, errors.New("unknown claims type, cannot proceed")
	}

	return claims, nil
}

// returns false if the claim times are not valid
func validateClaimTimes(claims jwt.Claims) (valid bool) {
	currentTime := time.Now()

	t, err := claims.GetExpirationTime()
	if err != nil || t.Time.Before(currentTime) {
		return false
	}

	t, err = claims.GetIssuedAt()
	if err != nil || currentTime.Before(t.Time) {
		return false
	}

	t, err = claims.GetNotBefore()
	if err != nil || currentTime.Before(t.Time) {
		return false
	}
	return true
}
