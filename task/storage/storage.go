package storage

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"task/base62"

	"github.com/gomodule/redigo/redis"
)

type Service interface {
	Save(string, time.Time) (string, error)
	Load(string) (string, error)
	LoadInfo(string) (*Item, error)
	Close() error
}

type Item struct {
	Id      uint64 `json:"id" redis:"id"`
	URL     string `json:"url" redis:"url"`
	Expires string `json:"expires" redis:"expires"`
	Visits  int    `json:"visits" redis:"visits"`
}

type db struct{ poll *redis.Pool }

func New(host, port, passwrod string) (Service, error) {
	poll := &redis.Pool {
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		},
	}

	return &db { poll }, nil
}

func (r *db) isUsed(id uint64) bool {
	conn := r.poll.Get()
	defer conn.Close()	

	exists, err := redis.Bool(conn.Do("EXISTS", "Shortener:"+strconv.FormatUint(id, 10)))
	if err != nil {
		return false
	}
	return exists
}

func (r *db) Save(url string, expires time.Time) (string, error) {
	conn := r.poll.Get()
	defer conn.Close()

	var id uint64

	for used := true; used; used = r.isUsed(id) {
		id = rand.Uint64()
	}

	shortLink := Item{id, url, expires.Format("2006-01-02 15:04:05.728046 +0300 EEST"), 0}

	_, err := conn.Do("HMSET", redis.Args{"Shortener:" + strconv.FormatUint(id, 10)}.AddFlat(shortLink)...)
	if err != nil {
		return "", err
	}

	_, err = conn.Do("EXPIREAT", "Shortener:"+strconv.FormatUint(id, 10), expires.Unix())
	if err != nil {
		return "", err
	}

	return base62.Encode(id), nil
}

func (r *db) Load(code string) (string, error) {
	conn := r.poll.Get()
	defer conn.Close()

	decodedId, err := base62.Decode(code)
	if err != nil {
		return "", err
	}

	urlString, err := redis.String(conn.Do("HGET", "Shortener:"+strconv.FormatUint(decodedId, 10), "url"))
	if err != nil {

	} else if len(urlString) == 0 {
		return "", nil //ErrNoLink
	}

	_, err = conn.Do("HINCRBY", "Shortener:"+strconv.FormatUint(decodedId, 10), "visits", 1)

	return urlString, nil
}

func (r *db) isAvailable(id uint64) bool {
	conn := r.poll.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", "Shortener:"+strconv.FormatUint(id, 10)))
	if err != nil {
		return false
	}

	return !exists
}

func (r *db) LoadInfo(code string) (*Item, error) {
	conn := r.poll.Get()
	defer conn.Close()

	decodedId, err := base62.Decode(code)
	if err != nil {
		return nil, err
	}

	values, err := redis.Values(conn.Do("HGETALL", "Shortener:"+strconv.FormatUint(decodedId, 10)))
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, nil  //ErrNoLink
	}

	var shortLink Item
	err = redis.ScanStruct(values, &shortLink)
	if err != nil {
		return nil, err
	}

	return &shortLink, nil
}

func (r *db) Close() error {
	return r.poll.Close()
}
