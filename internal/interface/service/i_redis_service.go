package service

import "github.com/gone-io/gone/goner/redis"

//go:generate sh -c "mockgen -package=mock -source=$GOFILE|gone mock -o mock/$GOFILE"
type IRedisService interface {
	GetConn() redis.Conn

	CloseConn(conn redis.Conn)

	SetWithEx(key, value string, expire int64) error

	Get(key string) ([]byte, error)

	Exists(key string) (bool, error)

	Incr(key string) error

	IncrBy(key string, value int) error

	Decr(key string) error

	DecrBy(key string, value int) error

	HSet(key, field, value string)

	HGet(key, field string) ([]byte, error)

	HDel(key, field string) ([]byte, error)

	HIncrBy(key, field string, incr int64) error

	ZSetWithEx(key, mem string, score int64, expire int64) error

	ZIncrbyWitEx(key string, mem string, incr, expire int64) error

	ZScore(key string, mem string) (int64, error)

	ZRange(key string, start, end int64, withScores bool) ([]string, error)

	ZRangeByScore(key string, min, max int64, withScores bool) ([]string, error)

	ZRevrangeByScore(key string, start, end int64, withScores bool) ([]string, error)

	PfAdd(key string, element ...string) (int, error)

	PfCount(key string) (int, error)
}
