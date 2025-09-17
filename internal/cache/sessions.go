package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const SessionTtlHours = time.Hour * 24 * 7

type SessionsCache struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) SessionsCache {
	return SessionsCache{
		rdb: rdb,
	}
}

func (c *SessionsCache) StoreSession(ctx context.Context, userID int, sessionID string) error {
	cmd := c.rdb.Set(fmt.Sprintf("session_%s", sessionID), userID, SessionTtlHours)
	if cmd.Err() != nil {
		return fmt.Errorf("error rdb.Set: %w", cmd.Err())
	}

	return nil
}

func (c *SessionsCache) GetUserIDBySessionID(ctx context.Context, sessionID string) (bool, int, error) {
	cmd := c.rdb.Get(fmt.Sprintf("session_%s", sessionID))
	if cmd.Err() != nil {
		if errors.Is(cmd.Err(), redis.Nil) {
			return false, 0, nil
		}
		return false, 0, fmt.Errorf("error rdb.Get: %w", cmd.Err())
	}

	userIdRaw := cmd.Val()
	userID, err := strconv.Atoi(userIdRaw)
	if err != nil {
		return false, 0, fmt.Errorf("error strconv.Atoi: %w", err)
	}

	return true, userID, nil
}
