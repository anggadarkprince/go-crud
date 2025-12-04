package services

import (
	"context"
	"database/sql"
	"errors"
	"net/mail"
	"os"
	"strconv"
	"time"

	"github.com/anggadarkprince/crud-employee-go/dto"
	"github.com/anggadarkprince/crud-employee-go/exceptions"
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
			return nil, exceptions.ErrUserNotFound
		}
		return nil, err
	}

	if user.Status != "ACTIVATED" {
		return nil, exceptions.ErrUserInactive
	} 

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
		return nil, exceptions.ErrWrongPassword
	}

	return user, nil	
}


func (service *AuthService) Register(ctx context.Context, data *dto.RegisterUserRequest) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

	userModel := &models.User{
        Name: data.Name,
        Email: data.Email,
        Username: data.Username,
        Password: string(hashedPassword),
        UserType: "EXTERNAL",
        Status: "ACTIVATED",
    }
	user, err := service.userRepository.Create(ctx, userModel)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *AuthService) GenerateAuthToken(userId int, exp int64) (string, error)  {
    claims := jwt.MapClaims{
        "sub": strconv.Itoa(userId), // subject: user id
        "exp": exp,
        "iat": time.Now().Unix(), // issued at
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}