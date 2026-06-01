package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/YASSERRMD/sql-sage/backend/internal/auth"
	"github.com/YASSERRMD/sql-sage/backend/internal/models"
	"github.com/YASSERRMD/sql-sage/backend/internal/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user is inactive")
	ErrInvalidToken       = errors.New("invalid refresh token")
)

type AuthService struct {
	users   *repositories.UserRepository
	refresh *repositories.RefreshTokenRepository
	jwt     *auth.Service
}

func NewAuthService(u *repositories.UserRepository, r *repositories.RefreshTokenRepository, j *auth.Service) *AuthService {
	return &AuthService{users: u, refresh: r, jwt: j}
}

type LoginResult struct {
	User         *models.User
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	if !u.IsActive {
		return nil, ErrUserInactive
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	pair, err := s.jwt.IssueTokens(u.ID, u.Role, u.Email)
	if err != nil {
		return nil, fmt.Errorf("issue tokens: %w", err)
	}
	if err := s.persistRefresh(ctx, u.ID, pair.RefreshToken); err != nil {
		return nil, err
	}
	return &LoginResult{
		User:         u,
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
	}, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*LoginResult, error) {
	hash := s.jwt.HashToken(refreshToken)
	rt, err := s.refresh.FindByHash(ctx, hash)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if rt.Revoked {
		return nil, ErrInvalidToken
	}
	if time.Unix(rt.ExpiresAt, 0).Before(time.Now()) {
		return nil, ErrInvalidToken
	}
	u, err := s.users.FindByID(ctx, rt.UserID)
	if err != nil {
		return nil, ErrInvalidToken
	}
	if !u.IsActive {
		return nil, ErrUserInactive
	}
	if err := s.refresh.Revoke(ctx, rt.ID); err != nil {
		return nil, err
	}
	pair, err := s.jwt.IssueTokens(u.ID, u.Role, u.Email)
	if err != nil {
		return nil, err
	}
	if err := s.persistRefresh(ctx, u.ID, pair.RefreshToken); err != nil {
		return nil, err
	}
	return &LoginResult{
		User:         u,
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		ExpiresIn:    pair.ExpiresIn,
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	hash := s.jwt.HashToken(refreshToken)
	rt, err := s.refresh.FindByHash(ctx, hash)
	if err != nil {
		return nil
	}
	return s.refresh.Revoke(ctx, rt.ID)
}

func (s *AuthService) Me(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.users.FindByID(ctx, id)
}

func (s *AuthService) persistRefresh(ctx context.Context, userID uuid.UUID, token string) error {
	rt := &models.RefreshToken{
		UserID:    userID,
		TokenHash: s.jwt.HashToken(token),
		ExpiresAt: time.Now().Add(s.jwt.RefreshTTL()).Unix(),
	}
	return s.refresh.Create(ctx, rt)
}
