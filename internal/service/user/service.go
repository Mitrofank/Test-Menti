package user

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/MitrofanK/Test-Menti/internal/errorsx"
	"github.com/MitrofanK/Test-Menti/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateSession(ctx context.Context, refresh models.RefreshSession) error
	GetSession(ctx context.Context, refreshTokenHash string) (models.RefreshSession, error)
	DeleteSession(ctx context.Context, sessionID int) error
	GetUserByID(ctx context.Context, id int) (models.User, error)
}

type Token interface {
	GenerateToken(userID int, roleID int) (string, error)
}

type UserService struct {
	repo  Repository
	token Token
}

func NewService(repo Repository, token Token) *UserService {
	return &UserService{
		repo:  repo,
		token: token,
	}
}

func (s *UserService) SignUp(ctx context.Context, email, password string) (int, error) {
	if _, err := s.repo.GetUserByEmail(ctx, email); err == nil {
		return 0, errorsx.ErrUserExists
	}

	if password == email {
		return 0, errorsx.ErrPasAndLoginSame
	}

	// проверка на состав пороля

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	var newUser = models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		RoleID:       2,
	}

	id, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (s *UserService) generateTokens(ctx context.Context, email string) (models.RefreshSession, string, string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return models.RefreshSession{}, "", "", errorsx.ErrIncorLoginOrPas
	}

	accessToken, err := s.token.GenerateToken(user.ID, user.RoleID)
	if err != nil {
		return models.RefreshSession{}, "", "", fmt.Errorf("token generation error: %w", err)
	}

	originalRefreshToken := uuid.NewString()
	refreshTokenHash := createHash(originalRefreshToken)
	refreshTokenTTL := 30 * 24 * time.Hour

	session := models.RefreshSession{
		UserID:       user.ID,
		RefreshToken: refreshTokenHash,
		UserAgent:    "", // TODO: Получить из контекста
		IPAddress:    "", // TODO: Получить из контекста
		ExpiresAt:    time.Now().Add(refreshTokenTTL),
	}

	return session, accessToken, originalRefreshToken, nil
}

func (s *UserService) SignIn(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", errorsx.ErrIncorLoginOrPas
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", errorsx.ErrIncorLoginOrPas
	}

	session, accessToken, refreshToken, err := s.generateTokens(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("token generation error: %w", err)
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return "", "", fmt.Errorf("create session error: %w", err)
	}

	return accessToken, refreshToken, nil
}

func createHash(token string) string {
	hasher := sha256.New()
	hasher.Write([]byte(token))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func (s *UserService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	refreshTokenHash := createHash(refreshToken)

	session, err := s.repo.GetSession(ctx, refreshTokenHash)
	if err != nil {
		return "", "", errorsx.ErrInvalidRefreshToken
	}

	if time.Now().After(session.ExpiresAt) {
		return "", "", errorsx.ErrRefreshTokenExpired
	}

	if err := s.repo.DeleteSession(ctx, session.ID); err != nil {
		return "", "", err
	}

	user, err := s.repo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return "", "", err
	}

	newSession, newAccessToken, newRefreshToken, err := s.generateTokens(ctx, user.Email)
	if err != nil {
		return "", "", fmt.Errorf("token generation error: %w", err)
	}

	if err := s.repo.CreateSession(ctx, newSession); err != nil {
		return "", "", fmt.Errorf("create session error: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}
