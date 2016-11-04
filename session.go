/**
 a simple & fast session libs for redis engine.
 @author: Misko_Lee
 @date: 2016-11-04
**/

package redissession

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"

	"github.com/satori/go.uuid"
	"gopkg.in/redis.v3"
)

const (
	DefaultRedisHost = "127.0.0.1:6379"
	DefaultRedisDB   = 13
	DefaultPrefix    = "redissession-"
	DefaultLifeTime  = 3600 * time.Second
)

type SessionConfig struct {
	RedisHost string        //redis host.default is (127.0.0.1:6379)
	RedisDB   int           //redis db.default is 13
	Prefix    string        //redis session key prefix. default is `redissession-`
	LifeTime  time.Duration //session lifetime.defualt is 1 hour.
}

type Session interface {
	//SetSessionID set custom session id.default is uuid-v4 string.
	SetSessionID(sid string) bool
	//get session id
	GetSessionID() string
	Start() bool
	Destory() bool
	Set(key string, val interface{}) bool
	Get(key string) interface{}
	Delete(key string) bool
	Store() bool
}

//global redis connection pool
var redisConnectionPool *redis.Client

type RedisSession struct {
	config  *SessionConfig
	db      *redis.Client
	values  map[string]interface{}
	sid     string
	isStart bool
}

//NewSession get a redis session instance.
func NewSession(engine string, config *SessionConfig) Session {
	if engine == "redis" {
		redisSession := new(RedisSession)
		redisSession.config = config

		if len(redisSession.config.RedisHost) < 1 {
			redisSession.config.RedisDB = DefaultRedisDB
		}

		if len(redisSession.config.Prefix) < 1 {
			redisSession.config.Prefix = DefaultPrefix
		}
		if redisSession.config.LifeTime < 1*time.Second {
			redisSession.config.LifeTime = DefaultLifeTime
		}
		if redisSession.config.RedisDB < 0 {
			redisSession.config.RedisDB = DefaultRedisDB
		}
		redisSession.sid = uuid.NewV4().String()
		redisSession.init()
		return redisSession
	}
	return nil
}

func (session *RedisSession) init() {
	session.values = make(map[string]interface{})
	if redisConnectionPool != nil {
		session.db = redisConnectionPool
		return
	}

	opt := &redis.Options{}
	opt.Addr = session.config.RedisHost
	opt.DB = int64(session.config.RedisDB)
	opt.ReadTimeout = 30 * time.Second
	opt.WriteTimeout = 30 * time.Second
	opt.DialTimeout = 60 * time.Second
	opt.MaxRetries = 5
	opt.PoolTimeout = 120 * time.Second
	redisConnectionPool = redis.NewClient(opt)
	session.db = redisConnectionPool
}

//SetSessionID set custom session id
func (session *RedisSession) SetSessionID(sid string) bool {
	if len(sid) < 1 {
		return false
	}
	session.sid = sid
	return true
}

//GetSessionID get session id
func (session *RedisSession) GetSessionID() string {
	return session.sid
}

//Start.Call Start before Any thing.
func (session *RedisSession) Start() bool {

	cmd := session.db.Get(session.config.Prefix + session.sid)

	if err := cmd.Err(); err != nil && err != redis.Nil {
		log.Println("[RedisSession] Session Start Error", err)
		return false
	}

	vals := cmd.Val()

	//new session
	if len(vals) < 1 {
		session.isStart = true
		return true
	}

	session.isStart = true

	buffer := bytes.NewBuffer([]byte(vals))
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&session.values); err != nil {
		log.Println("[RedisSession] Session Start Error", cmd.Err())
		return false
	}
	return true
}

//Destory Destory all session data
func (session *RedisSession) Destory() bool {

	if !session.isStart {
		log.Println("[RedisSession] Cannt Destory Session Before Start")
		return false
	}

	cmd := session.db.Del(session.config.Prefix + session.sid)

	if cmd.Err() != nil {
		log.Println("[RedisSession] Session Destory Error", cmd.Err())
		return false
	}
	return true
}

//Set set a session data.
func (session *RedisSession) Set(key string, val interface{}) bool {
	session.values[key] = val
	return true
}

//Get. returns nil if data not exists.
func (session *RedisSession) Get(key string) interface{} {
	return session.values[key]
}

//Delete. delete a session key.
func (session *RedisSession) Delete(key string) bool {
	delete(session.values, key)
	return true
}

//Store. store session changes to redis.
func (session *RedisSession) Store() bool {
	if !session.isStart {
		log.Println("[RedisSession] Cannt Store Session Before Start")
		return false
	}
	buffer := bytes.NewBuffer([]byte(""))

	encoder := gob.NewEncoder(buffer)

	if err := encoder.Encode(session.values); err != nil {
		log.Println("[RedisSession] Session Store Error", err)
		return false
	}

	cmd := session.db.Set(session.config.Prefix+session.sid, buffer.String(), session.config.LifeTime)

	if cmd.Err() != nil {
		log.Println("[RedisSession] Session Store Error", cmd.Err())
		return false
	}
	return true
}

func init() {
	gob.Register(map[string]interface{}{})
}
