package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func TestDumpStringKey(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	client.Set("test_string_1", "abc", 0)

	dumpdata, err := Dump(client, "test_string_1")
	if err != nil {
		t.Errorf("Dump returns error. Error: %s", err)
	}

	if dumpdata.Expiration != 0 {
		t.Errorf("Expiration should be zero.")
	}

	if dumpdata.Data == "" {
		t.Error("Empty data.")
	}
}

func TestDumpStringKeyExpiration(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	client.Set("test_string_2", "abc", time.Second*10)

	dumpdata, err := Dump(client, "test_string_2")
	if err != nil {
		t.Errorf("Dump returns error. Error: %s", err)
	}

	if dumpdata.Expiration == 0 {
		t.Error("Expiration time should be set.")
	}
}

func TestDumpToFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "dumptofiletest1")
	if err != nil {
		t.Errorf("Unable to create temp file. Error: %s", err)
	}
	defer os.RemoveAll(dir)

	row := make(Record)
	row["test_string_1_file"] = &KeyData{
		Data:       "foo",
		Expiration: 0,
	}

	path := filepath.Join(dir, "data.json")
	err = DumpJson(path, row)

	if err != nil {
		t.Errorf("Unable to store data in file. Error: %s", err)
	}

	fd, err := os.Open(path)
	if err != nil {
		t.Errorf("Unable to open data file. Error: %s", err)
	}
	defer fd.Close()

	data := make(Record)
	err = json.NewDecoder(fd).Decode(&data)

	if err != nil {
		t.Errorf("Unable to decode json data. Error: %s", err)
	}

	keydata, ok := data["test_string_1_file"]
	if !ok {
		t.Error("Key test_string_1_file does not exists in json data")
	}

	if keydata.Data != "foo" {
		t.Errorf("Data has invalid value. Expected foo, got: %s", keydata.Data)
	}
}
