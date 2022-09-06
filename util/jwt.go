package util

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWT struct {
	UserName   string
	UserEmail  string
	UserRole   string
	signature  string
	timeOffset int
}

type jwtClaims struct {
	Email string
	Name  string
	Role  string
	jwt.StandardClaims
}

func JWT_New(name string, email string, role string, timeOffset int) *JWT {
	fmt.Println(getSignature())
	return &JWT{
		signature:  getSignature(),
		timeOffset: timeOffset,
		UserName:   name,
		UserEmail:  email,
		UserRole:   role,
	}
}

func JWT_NewSignatureOnly() *JWT {
	return &JWT{
		signature: getSignature(),
	}
}

func getSignature() string {
	return os.Getenv("JWT_SIGNATURE")
}

func (j *JWT) GenerateToken() (string, error) {
	exp, err := strconv.Atoi(os.Getenv("JWT_EXPIRED")) // in minutes

	if err != nil {
		exp = 3 * 60 //default 3 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims{
		j.UserEmail,
		j.UserName,
		j.UserRole,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(exp) * time.Minute).Unix(),
		},
	})

	return token.SignedString([]byte(j.signature))
}

func (j *JWT) VerifyToken(token string) (*JWT, error) {
	verifyToken, err := jwt.ParseWithClaims(token, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.signature), nil
	})

	if err != nil {
		return nil, err
	}

	if claim, ok := verifyToken.Claims.(*jwtClaims); ok && verifyToken.Valid {
		j.UserEmail = claim.Email
		j.UserName = claim.Name
		j.UserRole = claim.Role
	} else {
		return nil, errors.New("Claim token error")
	}

	return j, nil
}

func (j *JWT) CheckingThisIsAdmin() error {
	fmt.Println(strings.ToLower(j.UserRole))
	if strings.ToLower(j.UserRole) == "admin" {
		return nil
	}
	return errors.New("Your not admin")
}
