package services

import (
	"crypto/rand"
	"fmt"
	"github.com/ChristophBe/desks/desks-api/util/configuration"
	"github.com/golang-jwt/jwt"
	"time"
)

const (
	Issuer = "de.christophb.rooms"
)

var singingKey []byte

type JwtTokenService interface {
	VerifyToken(tokenString string) (jwt.MapClaims, error)
	GenerateToken(claims jwt.MapClaims, duration time.Duration) (string, error)
}

type jwtTokenServiceImpl struct {
}

func (a jwtTokenServiceImpl) GenerateToken(claims jwt.MapClaims, duration time.Duration) (tokenString string, err error) {

	expirationTime := time.Now().Add(duration)
	claims["exp"] = expirationTime.Unix()
	claims["iss"] = Issuer

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString(getSigningKey())
	return
}

func (a jwtTokenServiceImpl) VerifyToken(tokenString string) (claims jwt.MapClaims, err error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return getSigningKey(), nil
	})
	if err != nil {
		return
	}

	var ok bool
	if claims, ok = token.Claims.(jwt.MapClaims); ok && token.Valid {
		//Check if Token Type is Valid
		issuer, ok := claims["iss"].(string)
		if !ok || issuer != Issuer {
			err = fmt.Errorf("unexpected issuer")
			return
		}
		return
	}
	err = fmt.Errorf("invalid token")

	return
}

func NewJwtTokenService() JwtTokenService {
	return jwtTokenServiceImpl{}
}

func getSigningKey() []byte {
	if singingKey == nil {
		singingKeyString := configuration.JwtSigningKey.GetValue()
		if len(singingKeyString) > 0 {
			singingKey = []byte(singingKeyString)
		} else {
			singingKey, _ = generateRandomBytes(64)
		}
	}
	return singingKey

}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
