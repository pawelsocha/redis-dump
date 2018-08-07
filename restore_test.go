package main

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func TestRestore(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	client.Set("test_string_1", "abc", 0)
	dumpdata, err := Dump(client, "test_string_1")
	if err != nil {
		t.Errorf("Can't dump key test_string_1. Error: %s", err)
	}

	defer client.Del("test_string_1")

	decoded, _ := base64.StdEncoding.DecodeString(dumpdata.Data)
	err = Restore(client, "test_string_1r", string(decoded), 0)
	if err != nil {
		t.Errorf("Restore returns an error. Error: %s", err)
	}

	defer client.Del("test_string_1r")

	data, err := client.Get("test_string_1r").Result()
	if err != nil {
		t.Errorf("Can't get key test_string_1r. Error: %s", err)
	}

	if data != "abc" {
		t.Errorf("Key test_string_1r should have abc. Got: %s", data)
	}
}

func TestRestoreWithExpire(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	client.Set("test_string_2", "abc", time.Second*10)
	dumpdata, err := Dump(client, "test_string_2")

	if err != nil {
		t.Errorf("Can't dump key test_string_2. Error: %s", err)
	}

	defer client.Del("test_string_2")

	if dumpdata.GetExpiration() <= 9 && dumpdata.GetExpiration() >= 10 {
		t.Errorf("Key test_string_2 should be set to 9 or 10. Got: %d", dumpdata.GetExpiration())
	}
	decoded, _ := base64.StdEncoding.DecodeString(dumpdata.Data)
	err = Restore(client, "test_string_2r", string(decoded), dumpdata.GetExpiration())

	if err != nil {
		t.Errorf("Can't restore test_string_2. Error: %s", err)
	}

	defer client.Del("test_string_2r")
}

func TestRestoreFromFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "restorefromfiletest1")
	if err != nil {
		t.Errorf("Unable to create temp file. Error: %s", err)
	}
	defer os.RemoveAll(dir)

	row := make(Record)

	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	client.Set("test_restore_string_1_file", "foobar", time.Second*10)
	dumpdata, err := Dump(client, "test_restore_string_1_file")

	if err != nil {
		t.Errorf("Can't dump key test_string_2. Error: %s", err)
	}

	row["test_restore_string_1_file"] = dumpdata

	path := filepath.Join(dir, "data.json")
	err = DumpJson(path, row)
	if err != nil {
		t.Errorf("Unable to store data in file. Error: %s", err)
	}

	client.Del("test_restore_string_1_file")

	err = RestoreJson(client, path)
	if err != nil {
		t.Errorf("Unable to restore data in redis. Error: %s", err)
	}

	data, err := client.Get("test_restore_string_1_file").Result()
	if err != nil {
		t.Errorf("Unable to get test_restore_string_1_file. Error: %s", err)
	}

	if data != "foobar" {
		t.Errorf("Data from key test_restore_string_1_file should be foobar. Got: %s", data)
	}

	ttl, _ := client.TTL("test_restore_string_1_file").Result()

	if ttl.Seconds() < 9 {
		t.Errorf("TTL for key test_restore_string_1_file should be 9. Got: %0.4f", ttl.Seconds())
	}
}

//TODO: add BUSYKEY test
