// auth/auth.go

package auth

import (
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// JwtWrapper wraps the signing key and the issuer
type JwtWrapper struct {
	SecretKey       string
	Issuer          string
	ExpirationHours int64
}

// JwtClaim adds email as a claim to the token
type JwtClaim struct {
	Roll_no string
	jwt.StandardClaims
}

// GenerateToken generates a jwt token
func (JWT_Wrap *JwtWrapper) GenerateToken(roll_no string) (signedToken string, err error) {
	JWT_Claim := &JwtClaim{
		Roll_no: roll_no,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(JWT_Wrap.ExpirationHours)).Unix(),
			Issuer:    JWT_Wrap.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWT_Claim)

	signedToken, err = token.SignedString([]byte(JWT_Wrap.SecretKey))
	if err != nil {
		return
	}

	return
}

//ValidateToken validates the jwt token
func (JWT_Wrap *JwtWrapper) ValidateToken(signedToken string) (JWT_Claim *JwtClaim, err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(JWT_Wrap.SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	JWT_Claim, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("Couldn't parse JWT_Claim")
		return
	}

	if JWT_Claim.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return
	}

	return

}
