package storage

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"task/base62" // base62.go | at task/base62  

	"github.com/gomodule/redigo/redis" // library 
)

// Create the Service interface 
type Service interface {
	// All Methods
	Save(string, time.Time) (string, error) // For save
	Load(string) (string, error) // For load
	LoadInfo(string) (*Item, error) // For load the information
	Close() error // For close
}

// Create the Item struct for get the id and url and ...
type Item struct {
	Id      uint64 `json:"id" redis:"id"`
	URL     string `json:"url" redis:"url"`
	Expires string `json:"expires" redis:"expires"`
	Visits  int    `json:"visits" redis:"visits"`
}

// DB or Redis connection
type db struct{ poll *redis.Pool }

// Implementing New
func New(host, port, passwrod string) (Service, error) { // Return Service with errors
	poll := &redis.Pool {
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		},
	}

	return &db { poll }, nil
}

// Implementing isUsed
func (r *db) isUsed(id uint64) bool { // Return a boolean
	conn := r.poll.Get() // New connection 
	defer conn.Close()	

	exists, err := redis.Bool(conn.Do("EXISTS", "Shortener:"+strconv.FormatUint(id, 10)))
	if err != nil { // We check that if *err* has an error, we can return that error 
		return false
	}
	return exists
}

// Implementing Save 
func (r *db) Save(url string, expires time.Time) (string, error) { // Return a string and error
	conn := r.poll.Get() // New connection
	defer conn.Close()

	var id uint64 // The id

	for used := true; used; used = r.isUsed(id) {
		id = rand.Uint64()
	}
	
	// The shortLink variable
	shortLink := Item{id, url, expires.Format("2006-01-02 15:04:05.728046 +0300 EEST"), 0}

	_, err := conn.Do("HMSET", redis.Args{"Shortener:" + strconv.FormatUint(id, 10)}.AddFlat(shortLink)...)
	if err != nil { // We check that if *err* has an error, we can return that error
		return "", err
	}

	_, err = conn.Do("EXPIREAT", "Shortener:"+strconv.FormatUint(id, 10), expires.Unix())
	if err != nil { // We check that if *err* has an error, we can return that error
		return "", err
	}

	return base62.Encode(id), nil // Return the encoded id and error
}

// Implementing Load
func (r *db) Load(code string) (string, error) { // Return a string and error 
	conn := r.poll.Get() // New connection
	defer conn.Close()

	decodedId, err := base62.Decode(code)
	if err != nil { // We check that if *err* has an error, we can return that error
		return "", err
	}

	urlString, err := redis.String(conn.Do("HGET", "Shortener:"+strconv.FormatUint(decodedId, 10), "url"))
	if err != nil { // We check that if *err* has an error, we can return that error
		return "", err
	} else if len(urlString) == 0 {
		return "", nil // ErrNoLink
	}

	_, err = conn.Do("HINCRBY", "Shortener: " + strconv.FormatUint(decodedId, 10), "visits", 1)

	return urlString, nil // Return the urlString 
}

// Implementing isAvailable
func (r *db) isAvailable(id uint64) bool { // Return a boolean
	conn := r.poll.Get() // New connection
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", "Shortener: " + strconv.FormatUint(id, 10)))
	if err != nil { // We check that if *err* has an error, we can return that error as false
		return false
	}

	return !exists
}

// Implementing LoadInfo
func (r *db) LoadInfo(code string) (*Item, error) { // Return *Item and error
	conn := r.poll.Get() // New connection
	defer conn.Close()

	decodedId, err := base62.Decode(code)
	if err != nil { // We check that if *err* has an error, we can return that error
		return nil, err
	}

	values, err := redis.Values(conn.Do("HGETALL", "Shortener:"+strconv.FormatUint(decodedId, 10)))
	if err != nil { // We check that if *err* has an error, we can return that error
		return nil, err
	} else if len(values) == 0 {
		return nil, nil  // ErrNoLink
	}

	var shortLink Item // Create the shortLink variable 
	err = redis.ScanStruct(values, &shortLink)
	if err != nil { // We check that if *err* has an error, we can return that error
		return nil, err
	}

	return &shortLink, nil // Return the reference of shortLink and error
}

// Implementing Close
func (r *db) Close() error { 
	return r.poll.Close() // Return the poll.Close() from DB
}

/* ('Sami Ghasemi) */
