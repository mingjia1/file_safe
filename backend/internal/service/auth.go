package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"password-timer-manager/internal/config"
	"password-timer-manager/internal/model"
	ptmerrors "password-timer-manager/internal/pkg/errors"
	"password-timer-manager/internal/repository"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	apiKeyRepo *repository.ApiKeyRepository
	cfg        *config.JWTConfig
}

func NewAuthService(userRepo *repository.UserRepository, apiKeyRepo *repository.ApiKeyRepository, cfg *config.JWTConfig) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		apiKeyRepo: apiKeyRepo,
		cfg:        cfg,
	}
}

type JWTClaims struct {
	UserID   uuid.UUID      `json:"user_id"`
	Username string         `json:"username"`
	Role     model.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*model.LoginResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, ptmerrors.ErrAuthFailed
	}

	if user.Status != model.UserStatusActive {
		return nil, ptmerrors.ErrAuthFailed
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ptmerrors.ErrPasswordWrong
	}

	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	s.userRepo.UpdateLastLogin(ctx, user.ID)

	return &model.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *AuthService) generateToken(user *model.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(time.Duration(s.cfg.ExpireHour) * time.Hour)

	claims := JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.Secret), nil
	})

	if err != nil {
		return nil, ptmerrors.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ptmerrors.ErrInvalidToken
}

func (s *AuthService) ValidateAPIKey(ctx context.Context, apiKey string) (*model.User, error) {
	keyHash := s.hashAPIKey(apiKey)

	user, err := s.apiKeyRepo.GetByHash(ctx, keyHash)
	if err != nil {
		return nil, ptmerrors.ErrInvalidAPIKey
	}

	if user.ApiKey.ExpiresAt != nil && user.ApiKey.ExpiresAt.Before(time.Now()) {
		return nil, ptmerrors.ErrInvalidAPIKey
	}

	return user.User, nil
}

func (s *AuthService) CreateAPIKey(ctx context.Context, userID uuid.UUID, name string, expiresIn int) (*model.ApiKeyResponse, error) {
	rawKey := s.generateAPIKey()
	keyHash := s.hashAPIKey(rawKey)

	apiKey := &model.ApiKey{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		KeyHash:   keyHash,
		CreatedAt: time.Now(),
	}

	if expiresIn > 0 {
		expiresAt := time.Now().Add(time.Duration(expiresIn) * 24 * time.Hour)
		apiKey.ExpiresAt = &expiresAt
	}

	if err := s.apiKeyRepo.Create(ctx, apiKey); err != nil {
		return nil, err
	}

	return &model.ApiKeyResponse{
		ID:        apiKey.ID,
		Name:      apiKey.Name,
		ApiKey:    rawKey,
		ExpiresAt: apiKey.ExpiresAt,
		CreatedAt: apiKey.CreatedAt,
	}, nil
}

func (s *AuthService) ListAPIKeys(ctx context.Context, userID uuid.UUID) ([]model.ApiKeyResponse, error) {
	apiKeys, err := s.apiKeyRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]model.ApiKeyResponse, len(apiKeys))
	for i, key := range apiKeys {
		responses[i] = model.ApiKeyResponse{
			ID:        key.ID,
			Name:      key.Name,
			ExpiresAt: key.ExpiresAt,
			CreatedAt: key.CreatedAt,
		}
	}

	return responses, nil
}

func (s *AuthService) DeleteAPIKey(ctx context.Context, userID, keyID uuid.UUID) error {
	return s.apiKeyRepo.Delete(ctx, userID, keyID)
}

func (s *AuthService) generateAPIKey() string {
	b := make([]byte, 32)
	for i := range b {
		b[i] = byte(time.Now().UnixNano() % 256)
	}
	hash := sha256.Sum256(b)
	return "ptm_sk_" + hex.EncodeToString(hash[:])
}

func (s *AuthService) hashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return hex.EncodeToString(hash[:])
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ptmerrors.ErrPasswordWrong
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	return s.userRepo.Update(ctx, user)
}

func (s *AuthService) Register(ctx context.Context, req *model.CreateUserRequest) (*model.UserResponse, error) {
	exists, err := s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ptmerrors.ErrUserExists
	}

	exists, err = s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ptmerrors.ErrUserExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		Status:       model.UserStatusActive,
		CreatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	response := user.ToResponse()
	return &response, nil
}
