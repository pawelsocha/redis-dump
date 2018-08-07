package main

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis"
)

//Dump get key dump from redis
func Dump(client *redis.Client, key string) (*KeyData, error) {
	data, err := client.Dump(key).Result()
	if err != nil {
		return nil, err
	}

	ttl, err := client.PTTL(key).Result()
	if err != nil {
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString([]byte(data))
	if ttl.Seconds() <= 0 {
		return &KeyData{
			Data:       encoded,
			Expiration: 0,
		}, nil
	}

	return &KeyData{
		Data:       encoded,
		Expiration: time.Now().Unix() + int64(ttl.Seconds()),
	}, nil
}

//DumpJson store data as json file in path
func DumpJson(path string, data interface{}) error {
	fd, err := os.Create(path)

	if err != nil {
		return err
	}
	defer fd.Close()
	err = json.NewEncoder(fd).Encode(data)

	return err
}
