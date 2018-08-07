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

	fmt.Printf("Restoring %s into db %d on %s\n", dumpfile, db, host)
	if err := RestoreJson(client, dumpfile); err != nil {
		fmt.Printf("Can't restore redis from %s. Error: %s", dumpfile, err)
		return
	}
}
