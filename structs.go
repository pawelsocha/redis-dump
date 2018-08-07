package main

import (
	"time"
)

type Record map[string]*KeyData

type KeyData struct {
	Data       string
	Expiration int64
}

//Expired this check if key should be expired
func (k KeyData) Expired() bool {
	if k.Expiration == 0 {
		return false
	}

	delta := k.GetExpiration()

	if delta <= 0 {
		return true
	}

	return false
}

//GetExpiration calculate key expiration from dump
func (k KeyData) GetExpiration() time.Duration {
	if k.Expiration == 0 {
		return time.Duration(0)
	}
	return time.Duration(k.Expiration-time.Now().Unix()) * time.Second
}
