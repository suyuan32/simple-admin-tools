package redistest

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// CreateRedis returns an in process redis.Redis.
func CreateRedis(t *testing.T) *redis.Redis {
	r, _ := CreateRedisWithClean(t)
	return r
}

// CreateRedisWithClean returns an in process redis.Redis and a clean function.
func CreateRedisWithClean(t *testing.T) (r *redis.Redis, clean func()) {
	t.Helper()
	mr := miniredis.RunT(t)
	r, err := redis.NewRedis(redis.RedisConf{Host: mr.Addr(), Type: redis.NodeType})
	if err != nil {
		t.Fatal(err)
	}
	return r, mr.Close
}
