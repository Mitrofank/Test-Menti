package errorsx

import "errors"

var (
	ErrUserExists          = errors.New("user with this email already exists")
	ErrPasBeEmpty          = errors.New("the password field cannot be empt")
	ErrPasLength           = errors.New("password must be between 8 and 64 characters long")
	ErrPasAndLoginSame     = errors.New("login and password must not match")
	ErrIncorLoginOrPas     = errors.New("incorrect login or password")
	ErrCurrencyNotFound    = errors.New("currency not found error")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrRefreshTokenExpired = errors.New("token refresh time expired")
)
