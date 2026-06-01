package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/YASSERRMD/sql-sage/backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	Type   string `json:"typ"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

type Service struct {
	cfg *config.Config
}

func NewService(cfg *config.Config) *Service {
	return &Service{cfg: cfg}
}

func (s *Service) IssueTokens(userID uuid.UUID, role, email string) (*TokenPair, error) {
	access, err := s.sign(userID, role, email, "access", s.cfg.JWTAccessTTL)
	if err != nil {
		return nil, err
	}
	refresh, err := s.randomToken()
	if err != nil {
		return nil, err
	}
	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    int(s.cfg.JWTAccessTTL.Seconds()),
	}, nil
}

func (s *Service) sign(uid uuid.UUID, role, email, typ string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: uid.String(),
		Role:   role,
		Email:  email,
		Type:   typ,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			Subject:   uid.String(),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *Service) ParseAccess(tokenStr string) (*Claims, error) {
	c := &Claims{}
	t, err := jwt.ParseWithClaims(tokenStr, c, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected alg: %s", t.Method.Alg())
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !t.Valid {
		return nil, errors.New("invalid token")
	}
	if c.Type != "access" {
		return nil, errors.New("not an access token")
	}
	return c, nil
}

func (s *Service) RandomToken() (string, error) { return s.randomToken() }

func (s *Service) HashToken(tok string) string {
	sum := sha256.Sum256([]byte(tok))
	return hex.EncodeToString(sum[:])
}

func (s *Service) RefreshTTL() time.Duration { return s.cfg.JWTRefreshTTL }

func (s *Service) randomToken() (string, error) {
	b := make([]byte, 48)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
