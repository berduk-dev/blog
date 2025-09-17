package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func (s *Service) CreateSession(ctx context.Context, userID int) (string, error) {
	b := make([]byte, 32)
	rand.Read(b)
	sessionID := base64.URLEncoding.EncodeToString(b)[:32]

	err := s.sessionsManger.StoreSession(ctx, userID, sessionID)
	if err != nil {
		return "", fmt.Errorf("error sessionsManger.StoreSession: %w", err)
	}

	return sessionID, nil
}

func (s *Service) GetUserIDBySessionID(ctx context.Context, sessionID string) (bool, int, error) {
	ok, userID, err := s.sessionsManger.GetUserIDBySessionID(ctx, sessionID)
	if err != nil {
		return false, 0, fmt.Errorf("error sessionsManger.GetUserIDBySessionID: %w", err)
	}

	return ok, userID, nil
}
