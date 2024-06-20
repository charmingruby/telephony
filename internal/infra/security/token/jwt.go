package token

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTService struct {
	issuer    string
	secretKey string
}

func NewJWTService(secretKey, issuer string) *JWTService {
	return &JWTService{
		issuer:    issuer,
		secretKey: secretKey,
	}
}

type JWTClaim struct {
	jwt.StandardClaims
	// put fields here to be stored on token, if needed
}

func (j *JWTService) GenerateToken(userID int) (string, error) {
	tokenDuration := time.Duration(time.Minute * 60 * 24 * 7) //7 days

	claims := &JWTClaim{
		jwt.StandardClaims{
			Subject:   strconv.Itoa(userID),
			Issuer:    j.issuer,
			ExpiresAt: time.Now().Local().Add(tokenDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", nil
	}

	return tokenStr, nil
}

type Payload struct {
	UserID int `json:"user_id"`
}

func (j *JWTService) ValidateToken(token string) (*Payload, error) {
	t, err := jwt.Parse(token, j.isTokenValid)
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("unable to parse jwt claims")
	}

	userIDStr := claims["sub"].(string)
	userIDParsed, err := strconv.Atoi(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("unable to parse user_id")
	}

	payload := &Payload{
		UserID: userIDParsed,
	}

	return payload, err
}

func (j *JWTService) isTokenValid(t *jwt.Token) (interface{}, error) {
	if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
		return nil, fmt.Errorf("invalid token %v", t)
	}

	return []byte(j.secretKey), nil
}
