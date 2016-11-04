package redissession

import (
	"log"
	"testing"
	"time"
)

func TestSession(t *testing.T) {

	sessionConfig := &SessionConfig{
		RedisDB:   11,
		RedisHost: "127.0.0.1:6379",
		Prefix:    "customprefix",
		LifeTime:  5 * time.Second,
	}

	session := NewSession("redis", sessionConfig)

	session.SetSessionID("abcdefgh")

	session.Start()

	log.Println(session.Get("key_int"), session.Get("key_float"))

	log.Println(session.GetSessionID())

	session.Set("key_int", 1)
	session.Set("key_float", 1.131425926)
	session.Set("key_string", "abcdefgh")
	session.Set("key_map", map[string]interface{}{
		"a": 1,
		"b": 1.131425926,
		"c": "1231231",
	})
	session.Store()
}
