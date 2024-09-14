package auth

import (
	"encoding/json"
	"main/cache"
	"main/models"
)

func getUserKeyById(id models.ModelID) ([]byte, error) {
	pattern := getUserKeyPatternById(id)
	return cache.GetKeyByPattern(pattern)
}

func getUserKeyByToken(token *models.AuthToken) ([]byte, error) {
	pattern := getUserKeyPatternByToken(token)
	return cache.GetKeyByPattern(pattern)
}

func getUserByKey(key string) (*models.User, error) {
	ctx, _ := cache.GetContext()
	conn := cache.GetConnection()

	userJSON, err := conn.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var user models.User
	err = json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func setUserByKey(user *models.User, key string) error {
	ctx, _ := cache.GetContext()
	conn := cache.GetConnection()

	userJSON, err := json.Marshal(&user)
	if err != nil {
		return err
	}
	return conn.Set(ctx, key, userJSON, 0).Err()
}

func delUserByKey(key string) error {
	ctx, _ := cache.GetContext()
	conn := cache.GetConnection()

	return conn.Del(ctx, key).Err()
}

func GetUserByID(id models.ModelID) (*models.User, error) {
	key, err := getUserKeyById(id)
	if err != nil {
		return nil, err
	}

	return getUserByKey(string(key))
}

func GetTokenByUserID(id models.ModelID) (*models.AuthToken, error) {
	key, err := getUserKeyById(id)
	if err != nil {
		return nil, err
	}
	return parseTokenByKey(string(key))
}

func GetUserByToken(token *models.AuthToken) (*models.User, error) {
	key, err := getUserKeyByToken(token)
	if err != nil {
		return nil, err
	}

	return getUserByKey(string(key))
}

func SetUserByToken(user *models.User, token *models.AuthToken) error {
	previousKey, err := getUserKeyById(user.ID)
	if previousKey != nil && err == nil {
		err = DelUser(user)
	}

	key := createUserKeyByToken(user, token)
	return setUserByKey(user, key)
}

func DelUser(user *models.User) error {
	key, err := getUserKeyById(user.ID)
	if err != nil {
		return err
	}
	return delUserByKey(string(key))
}
