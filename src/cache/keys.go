package cache

import (
	"errors"
)

var (
	NoKeyFoundError = errors.New("no such key found")
)

func GetKeyByPattern(pattern string) ([]byte, error) {
	ctx, _ := GetContext()

	keys, err := conn.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	if len(keys) != 1 {
		return nil, NoKeyFoundError
	}

	return []byte(keys[0]), nil
}
