# redissession
--
    import "github.com/weixinhost/redissession"

## Usage

```go
const (
       	DefaultRedisHost = "127.0.0.1:6379"
       	DefaultRedisDB   = 13
       	DefaultPrefix    = "redissession-"
       	DefaultLifeTime  = 3600 * time.Second
)
```

####type RedisSession

```go
type RedisSession struct {
}
```


#### func (*RedisSession) Delete

```go
func (session *RedisSession) Delete(key string) bool
```

Delete. delete a session key.

#### func (*RedisSession) Destory

```go
func (session *RedisSession) Destory() bool
```

Destory Destory all session data

#### func (*RedisSession) Get

```go
func (session *RedisSession) Get(key string) interface{}
```

Get. returns nil if data not exists.

#### func (*RedisSession) GetSessionID

```go
func (session *RedisSession) GetSessionID() string
```

GetSessionID get session id

#### func (*RedisSession) Set

```go
func (session *RedisSession) Set(key string, val interface{}) bool
```

Set set a session data.

#### func (*RedisSession) SetSessionID

```go
func (session *RedisSession) SetSessionID(sid string) bool
```

SetSessionID set custom session id

#### func (*RedisSession) Start

```go
func (session *RedisSession) Start() bool
```

Start.Call Start before Any thing.

#### func (*RedisSession) Store

```go
func (session *RedisSession) Store() bool
```

Store. store session changes to redis.

####type Session

```go
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
```


#### func  NewSession

```go
func NewSession(engine string, config *SessionConfig) Session
```

NewSession get a redis session instance.

####type SessionConfig

```go
type SessionConfig struct {
       	RedisHost string        //redis host.default is (127.0.0.1:6379)
       	RedisDB   int           //redis db.default is 13
       	Prefix    string        //redis session key prefix. default is `redissession-`
       	LifeTime  time.Duration //session lifetime.defualt is 1 hour.
}
```