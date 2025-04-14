package userservice


import (
	"context"
	"errors"
	"time"
)
func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if u.Password != password {
		return "", errors.New("invalid credentials")
	}

	return s.jwtManager.GenerateToken(string(u.Role), time.Hour)
}
