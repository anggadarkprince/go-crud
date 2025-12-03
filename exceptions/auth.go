package exceptions

import "errors"

var (
    ErrUserNotFound = errors.New("user not found")
    ErrUserInactive = errors.New("user not activated")
    ErrWrongPassword = errors.New("wrong password")
)