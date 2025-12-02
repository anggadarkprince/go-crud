package services

import (
	"context"
	"database/sql"
	"errors"
	"net/mail"
	"os"
	"strconv"
	"time"

	appErrors "github.com/anggadarkprince/crud-employee-go/errors"
	"github.com/anggadarkprince/crud-employee-go/models"
	"github.com/anggadarkprince/crud-employee-go/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository *repositories.UserRepository
}

func NewAuthService(userRepository *repositories.UserRepository) *AuthService {
	return &AuthService{userRepository: userRepository}
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func (service *AuthService) Authenticate(ctx context.Context, username string, password string) (*models.User, error) {
	isEmail := true
	if _, err := mail.ParseAddress(username); err != nil {
        isEmail = false
    }
	
	var user *models.User
	var err error
    if isEmail {
		user, err = service.userRepository.GetByEmail(ctx, username)
	} else {
		user, err = service.userRepository.GetByUsername(ctx, username)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, appErrors.ErrUserNotFound
		}
		return nil, err
	}

	if user.Status != "ACTIVATED" {
		return nil, appErrors.ErrUserInactive
	} 

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
		return nil, appErrors.ErrWrongPassword
	}

	return user, nil	
}

func (service *AuthService) GenerateAuthToken(userId int) (string, error)  {
    claims := jwt.MapClaims{
        "sub": strconv.Itoa(userId), // subject: user id
        "exp": time.Now().Add(24 * time.Hour).Unix(), // expires in 24h
        "iat": time.Now().Unix(), // issued at
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}