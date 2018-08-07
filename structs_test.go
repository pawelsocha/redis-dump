package main

import (
	"testing"
	"time"
)

func TestKeyNoExpiration(t *testing.T) {
	data := KeyData{
		Data:       "testdata",
		Expiration: 0,
	}

	if data.Expired() != false {
		t.Errorf("Key expired.")
	}
}

func TestKeyExpired(t *testing.T) {
	data := KeyData{
		Data:       "testdata",
		Expiration: time.Now().Add(time.Second * 1).Unix(),
	}

	if data.Expired() != false {
		t.Errorf("Key expired.")
	}

	time.Sleep(time.Second * 2)

	if data.Expired() == false {
		t.Errorf("Key should be expired.")
	}

	if data.GetExpiration().Seconds() != -1 {
		t.Errorf("Key expiration should be -1. Got: %d", data.GetExpiration())
	}
}

func TestKeyNotExpired(t *testing.T) {
	data := KeyData{
		Data:       "testdata",
		Expiration: time.Now().Add(time.Second * 10).Unix(),
	}

	time.Sleep(time.Second * 2)
	if data.Expired() != false {
		t.Errorf("Key expired.")
	}

	if data.GetExpiration().Seconds() != 8 {
		t.Errorf("Key expiration should be 8. Got: %d", data.GetExpiration())
	}
}
