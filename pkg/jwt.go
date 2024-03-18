package pkg

import (
	"errors"
	"filmLibraryVk/internal/model/entity"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var privateKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

func GenerateJWT(user entity.User) (string, error) {
	expiration, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.Id,
		"role": user.RoleId,
		"iat":  time.Now().Unix(),
		"eat":  time.Now().Add(time.Second * time.Duration(expiration)).Unix(),
	})
	return token.SignedString(privateKey)
}

func ValidateJWT(w http.ResponseWriter, r *http.Request) error {
	token, err := getToken(w, r)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

func ValidateAdminRoleJWT(w http.ResponseWriter, r *http.Request) error {
	token, err := getToken(w, r)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := uint(claims["role"].(float64))
	if ok && token.Valid && userRole == 1 {
		return nil
	}
	log.Printf("Forbidden")
	return errors.New("invalid admin token provided")
}

func ValidateUserRoleJWT(w http.ResponseWriter, r *http.Request) error {
	token, err := getToken(w, r)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := uint(claims["role"].(float64))
	if ok && token.Valid && userRole == 2 || userRole == 1 {
		return nil
	}
	log.Printf("Forbidden")
	return errors.New("invalid author token provided")
}

func MockValidateJWT(w http.ResponseWriter, r *http.Request) error {
	tokenString := getTokenFromRequest(w, r)
	if tokenString == "ADMIN" || tokenString == "USER" {
		return nil
	}
	return errors.New("invalid token provided")
}

func MockValidateAdminRoleJWT(w http.ResponseWriter, r *http.Request) error {
	tokenString := getTokenFromRequest(w, r)
	if tokenString == "ADMIN" {
		return nil
	}
	log.Printf("Forbidden")
	return errors.New("invalid admin token provided")
}

func MockValidateUserRoleJWT(w http.ResponseWriter, r *http.Request) error {
	tokenString := getTokenFromRequest(w, r)
	if tokenString == "USER" || tokenString == "ADMIN" {
		return nil
	}
	log.Printf("Forbidden")
	return errors.New("invalid author token provided")
}

func getToken(w http.ResponseWriter, r *http.Request) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(w, r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return privateKey, nil
	})
	return token, err
}

func getTokenFromRequest(w http.ResponseWriter, r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

func EncodePassword(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func ComparePasswords(hashPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
