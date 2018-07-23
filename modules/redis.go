package modules

import (
	"github.com/go-redis/redis"
	"strings"
	"time"
	"strconv"
)

type Redis struct {
	Config map[string]string //配置文件
	Conn   redis.Client      //连接
	Info   map[string]string //redis info
}

func (R *Redis) Connect() {
	db, _ := strconv.Atoi(R.Config["redis_db"])
	client := redis.NewClient(&redis.Options{
		Addr:     R.Config["redis_host"] + ":" + R.Config["redis_port"],
		Password: R.Config["redis_password"], // no password set
		DB:       db,                         // use default DB
	})
	R.Conn = *client //redis 连接
}

func (R *Redis) Key(key string) (string) {
	newKey := R.Config["redis_pre"] + "_" + key
	return newKey
}

func (R *Redis) Set(key, value string, expiration int) error {
	key = R.Key(key)
	set_err := R.Conn.Set(key, value, time.Duration(1000000000*expiration)).Err()
	return set_err
}

func (R *Redis) Get(key string) (string, error) {
	key = R.Key(key)
	value, err := R.Conn.Get(key).Result()
	return value, err
}

func (R *Redis) Del(key string) error {
	key = R.Key(key)
	_, err := R.Conn.Del(key).Result()
	return err
}

//刷新key过期时间
func (R *Redis) ExPire(key string,expiration int) error {
	key = R.Key(key)
	err := R.Conn.Expire(key,time.Duration(1000000000*expiration)).Err()
	return err
}

func (R *Redis) ServerStatus() {
	R.Info = make(map[string]string)
	for _, v := range strings.Split(R.Conn.Info().String(), "\r\n") {
		if value := strings.Split(v, ":"); len(value) == 2 {
			R.Info[value[0]] = value[1]
		}
	}
}
