package redis_util

import (
	"github.com/gone-io/gone"
	"github.com/gone-io/gone/goner/logrus"
	"github.com/gone-io/gone/goner/redis"
	"math"
)

//go:gone
func NewRedisService() gone.Goner {
	return &redisService{}
}

type redisService struct {
	gone.Flag
	logrus.Logger `gone:"gone-logger"`
	Redis         redis.Pool `gone:"gone-redis-pool"`
}

func (s *redisService) GetConn() redis.Conn {
	return s.Redis.Get()
}

func (s *redisService) CloseConn(conn redis.Conn) {
	err := conn.Close()
	if err != nil {
		s.Logger.Errorf("Redis conn.Close() err:%v", err)
	}
}

func (s *redisService) SetWithEx(key, value string, expire int64) error {
	conn := s.GetConn()
	defer s.CloseConn(conn)
	_, err := conn.Do("SET", key, value, "EX", expire)
	return err
}

func (s *redisService) Decr(key string) error {
	conn := s.GetConn()
	defer s.CloseConn(conn)
	_, err := conn.Do("DECR", key)
	return err
}

func (s *redisService) DecrBy(key string, value int) error {
	conn := s.GetConn()
	defer s.CloseConn(conn)
	_, err := conn.Do("DECRBY", key, value)
	return err
}

func (s *redisService) Incr(key string) error {
	conn := s.GetConn()
	defer s.CloseConn(conn)
	_, err := conn.Do("INCR", key)
	return err
}

func (s *redisService) IncrBy(key string, value int) error {
	conn := s.GetConn()
	defer s.CloseConn(conn)
	_, err := conn.Do("INCRBY", key, value)
	return err
}

func (s *redisService) Get(key string) ([]byte, error) {
	conn := s.GetConn()
	defer s.CloseConn(conn)
	return redis.Bytes(conn.Do("GET", key))
}

func (s *redisService) HSet(key, field, value string) {
	conn := s.GetConn()
	defer s.CloseConn(conn)
	_, _ = conn.Do("HSET", key, field, value)
}

func (s *redisService) HGet(key, field string) ([]byte, error) {
	conn := s.GetConn()
	defer s.CloseConn(conn)
	return redis.Bytes(conn.Do("HGET", key, field))
}

func (s *redisService) HDel(key, field string) ([]byte, error) {
	conn := s.GetConn()
	defer s.CloseConn(conn)
	return redis.Bytes(conn.Do("HDEL", key, field))
}

func (s *redisService) HIncrBy(key, field string, incr int64) error {
	conn := s.GetConn()
	defer s.CloseConn(conn)
	_, err := redis.Int64(conn.Do("HINCRBY", key, field, incr))
	return err
}

func (s *redisService) ZSetWithEx(key, mem string, score int64, expire int64) error {
	conn := s.GetConn()
	defer s.CloseConn(conn)
	_, err := conn.Do("ZADD", key, score, mem)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", key, expire, "NX")
	return err
}

func (s *redisService) ZIncrbyWitEx(key string, mem string, incr, expire int64) error {
	conn := s.GetConn()
	defer s.CloseConn(conn)
	_, err := conn.Do("ZINCRBY", key, incr, mem)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", key, expire)

	return err
}

func (s *redisService) ZScore(key string, mem string) (int64, error) {
	conn := s.GetConn()
	defer s.CloseConn(conn)

	score, err := redis.Int64(conn.Do("ZSCORE", key, mem))
	s.Logger.Infof("get keys [%s] from redis is: %v", key, score)

	return score, err
}

func (s *redisService) ZRange(key string, start, end int64, withScores bool) ([]string, error) {
	conn := s.GetConn()
	defer s.CloseConn(conn)

	if withScores {
		return redis.Strings(conn.Do("ZRANGE", key, start, end, "WITHSCORES"))
	}

	return redis.Strings(conn.Do("ZRANGE", key, start, end))
}

func (s *redisService) ZRangeByScore(key string, start, end int64, withScores bool) ([]string, error) {
	conn := s.GetConn()
	defer s.CloseConn(conn)

	if withScores {
		return redis.Strings(conn.Do("ZRANGEBYSCORE", key, start, end, "WITHSCORES"))
	}

	return redis.Strings(conn.Do("ZRANGEBYSCORE", key, start, end))
}

func (s *redisService) ZRevrangeByScore(key string, start, end int64, withScores bool) ([]string, error) {
	conn := s.GetConn()
	defer s.CloseConn(conn)

	if withScores {
		return redis.Strings(conn.Do("ZREVRANGEBYSCORE", key, math.MaxInt64, 0, "WITHSCORES", "LIMIT", start, end))
	}

	return redis.Strings(conn.Do("ZREVRANGEBYSCORE", key, math.MaxInt64, 0, "LIMIT", start, end))
}

func (s *redisService) PfAdd(key string, element ...string) (int, error) {
	conn := s.GetConn()
	defer s.CloseConn(conn)

	return redis.Int(conn.Do("PFADD", key, element))
}

func (s *redisService) PfCount(key string) (int, error) {
	conn := s.GetConn()
	defer s.CloseConn(conn)

	return redis.Int(conn.Do("PFCOUNT", key))
}

func (s *redisService) Exists(key string) (bool, error) {
	conn := s.GetConn()
	defer s.CloseConn(conn)

	res, err := redis.Int(conn.Do("EXISTS", key))
	return res > 0, err
}
