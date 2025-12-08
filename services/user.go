package services

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/anggadarkprince/crud-employee-go/dto"
	"github.com/anggadarkprince/crud-employee-go/exceptions"
	"github.com/anggadarkprince/crud-employee-go/models"
	"github.com/anggadarkprince/crud-employee-go/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository *repositories.UserRepository
	db *sql.DB
}

func NewUserService(
	userRepository *repositories.UserRepository,
	db *sql.DB,
) *UserService {
	return &UserService{
		userRepository: userRepository,
		db: db,
	}
}

func (service *UserService) UpdateAccount(ctx context.Context, data *dto.UpdateAccountRequest) (*models.User, error) {
	user, err := service.userRepository.GetById(ctx, data.Id)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.CurrentPassword))
    if err != nil {
		return nil, &exceptions.ValidationError{
			Message: "Current password is wrong",
		}
	}

	data.Avatar = user.Avatar.String
	if data.AvatarFile != nil {
        src, _ := data.AvatarFile.Open()
        defer src.Close()

		year := time.Now().Format("2006")
		month := time.Now().Format("01")

		uploadPath := fmt.Sprintf("avatars/%s/%s/%s", year, month, data.AvatarFile.Filename)
		uploadDir := filepath.Join("uploads", "avatars", year, month)
		err := os.MkdirAll(uploadDir, os.ModePerm)
		if err != nil {
			return nil, err
		}

		finalPath := filepath.Join(uploadDir, data.AvatarFile.Filename)

        dst, _ := os.Create(finalPath)
        defer dst.Close()

        _, err = io.Copy(dst, src)
		if err != nil {
			return nil, err
		}

		data.Avatar = uploadPath
    }

	hashedPassword := user.Password
	if data.Password != "" {
		newPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		hashedPassword = string(newPassword)
	}

	userModel := &models.User{
		Id: data.Id,
        Name: data.Name,
        Email: data.Email,
        Username: data.Username,
		Password: hashedPassword,
		Avatar: sql.NullString{String: data.Avatar, Valid: data.Avatar != ""},
    }
	user, err = service.userRepository.UpdateAccount(ctx, userModel)
	if err != nil {
		return nil, err
	}
	return user, nil
}