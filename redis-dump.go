package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       db,
	})

	keys, err := client.Keys("*").Result()

	if err != nil {
		fmt.Errorf("Can't get keys. Error: %s", err)
		return
	}

	mapdata := make(Record)

	for _, key := range keys {
		data, err := Dump(client, key)
		if err != nil {
			fmt.Errorf("Can't dump key %s. Error: %s", key, err)
			return
		}
		mapdata[key] = data
	}

	err = DumpJson(dumpfile, mapdata)

	if err != nil {
		fmt.Errorf("Can't store data in dumpfile %s. Error: %s", dumpfile, err)
		return
	}
}
