package repository

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthRepo struct {
	client *redis.Client
}

func NewAuthRepo(client *redis.Client) *AuthRepo {
	return &AuthRepo{
		client: client,
	}
}

type RedisData struct {
	UserId primitive.ObjectID
	Email  string
	Role   string
	Ua     string
	Ip     string
	Exp    time.Duration
}

func (i RedisData) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

func UnMarshalBinary(str string) *RedisData {
	var data *RedisData
	json.Unmarshal([]byte(str), &data)
	return data
}

func (r *AuthRepo) CreateSession(token string, data RedisData) error {
	if err := r.client.Set(r.client.Context(), token, data, data.Exp).Err(); err != nil {
		return err
	}
	return nil
}

func (r *AuthRepo) GetDelSession(token string) (*RedisData, error) {
	cmd := r.client.GetDel(r.client.Context(), token)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	str, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	return UnMarshalBinary(str), nil
}

func (r *AuthRepo) RemoveSession(token string) error {
	if err := r.client.Del(r.client.Context(), token).Err(); err != nil {
		return err
	}
	return nil
}
