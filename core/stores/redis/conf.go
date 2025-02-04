package redis

import (
	"errors"
	"time"
)

var (
	// ErrEmptyHost is an error that indicates no redis host is set.
	ErrEmptyHost = errors.New("empty redis host")
	// ErrEmptyType is an error that indicates no redis type is set.
	ErrEmptyType = errors.New("empty redis type")
	// ErrEmptyKey is an error that indicates no redis key is set.
	ErrEmptyKey = errors.New("empty redis key")
)

type (
	// A RedisConf is a redis config.
	RedisConf struct {
		Host     string `json:",env=REDIS_HOST"`
		Type     string `json:",default=node,options=node|cluster,env=REDIS_TYPE"`
		User     string `json:",optional,env=REDIS_USER"`
		Pass     string `json:",optional,env=REDIS_PASSWORD"`
		Tls      bool   `json:",optional,env=REDIS_TLS"`
		NonBlock bool   `json:",default=true"`
		// PingTimeout is the timeout for ping redis.
		PingTimeout time.Duration `json:",default=1s"`
	}

	// A RedisKeyConf is a redis config with key.
	RedisKeyConf struct {
		RedisConf
		Key string
	}
)

// NewRedis returns a Redis.
// Deprecated: use MustNewRedis or NewRedis instead.
func (rc RedisConf) NewRedis() *Redis {
	var opts []Option
	if rc.Type == ClusterType {
		opts = append(opts, Cluster())
	}
	if len(rc.User) > 0 {
		opts = append(opts, WithUser(rc.User))
	}
	if len(rc.Pass) > 0 {
		opts = append(opts, WithPass(rc.Pass))
	}
	if rc.Tls {
		opts = append(opts, WithTLS())
	}

	return newRedis(rc.Host, opts...)
}

// Validate validates the RedisConf.
func (rc RedisConf) Validate() error {
	if len(rc.Host) == 0 {
		return ErrEmptyHost
	}

	if len(rc.Type) == 0 {
		return ErrEmptyType
	}

	return nil
}

// Validate validates the RedisKeyConf.
func (rkc RedisKeyConf) Validate() error {
	if err := rkc.RedisConf.Validate(); err != nil {
		return err
	}

	if len(rkc.Key) == 0 {
		return ErrEmptyKey
	}

	return nil
}
