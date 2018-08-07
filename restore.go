package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
)

func Restore(client *redis.Client, key, data string, expire time.Duration) error {
	_, err := client.Restore(key, expire, data).Result()
	if err != nil {
		return err
	}
	return nil
}

func RestoreJson(client *redis.Client, path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	rows := make(Record)

	if err = json.NewDecoder(fd).Decode(&rows); err != nil {
		return err
	}

	for k, v := range rows {

		if v.Expired() {
			continue
		}

		//TODO: what now?
		data, err := base64.StdEncoding.DecodeString(v.Data)
		if err != nil {
			return err
		}

		if err = Restore(client, k, string(data), v.GetExpiration()); err != nil {
			return fmt.Errorf("Can't restore key %s. Error: %s", k, err)
		}
	}

	return nil
}
