package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type StateStore struct {
	client *redis.Client
}

func NewStateStore(client *redis.Client) *StateStore {
	return &StateStore{client: client}
}

func (s *StateStore) SetState(ctx context.Context, userID int64, state string, ttl time.Duration) error {
	key := fmt.Sprintf("fsm:state:%d", userID)
	return s.client.Set(ctx, key, state, ttl).Err()
}

func (s *StateStore) GetState(ctx context.Context, userID int64) (string, error) {
	key := fmt.Sprintf("fsm:state:%d", userID)
	state, err := s.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return state, nil
}

func (s *StateStore) ClearState(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("fsm:state:%d", userID)
	return s.client.Del(ctx, key).Err()
}

func (s *StateStore) SetData(ctx context.Context, userID int64, key string, data any, ttl time.Duration) error {
	fullKey := fmt.Sprintf("fsm:data:%d:%s", userID, key)

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal fsm data: %w", err)
	}

	return s.client.Set(ctx, fullKey, bytes, ttl).Err()
}

func (s *StateStore) GetData(ctx context.Context, userID int64, key string, dest any) error {
	fullKey := fmt.Sprintf("fsm:data:%d:%s", userID, key)

	bytes, err := s.client.Get(ctx, fullKey).Bytes()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, dest)
}


func (s *StateStore) ClearData(ctx context.Context, userID int64, key string) error {
	fullKey := fmt.Sprintf("fsm:data:%d:%s", userID, key)
	return s.client.Del(ctx, fullKey).Err()
}