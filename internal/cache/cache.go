package cache

import (
	cache2 "github.com/trumanwong/go-internal/cache"
	"remix-api/configs"
)

var Cache *cache2.Cache

func init() {
	Cache = cache2.NewCache(
		configs.Config.Redis.Addr,
		configs.Config.Redis.Password,
		configs.Config.Cache.Prefix,
		configs.Config.Cache.Database,
	)
}
