package cache

import (
	"remix-api/configs"
	cache2 "github.com/trumanwong/go-internal/cache"
)

var Cache *cache2.Cache

func Setup() {
	Cache = cache2.NewCache(
		configs.Config.Redis.Addr,
		configs.Config.Redis.Password,
		configs.Config.Cache.Prefix,
		configs.Config.Cache.Database,
	)
}
