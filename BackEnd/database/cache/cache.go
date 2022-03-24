package cache

import (
	"FINRepository/database"
	utilReflection "FINRepository/util/reflection"
	"context"
	"reflect"
	"sync"
	"time"
)

type cacheEntryKey struct {
	Type reflect.Type
	ID   database.ID
}

type cacheEntryValue struct {
	Obj    interface{}
	Expire time.Time
}

type DatabaseCache struct {
	parent    *DatabaseCache
	cache     map[cacheEntryKey]cacheEntryValue
	cacheLock sync.RWMutex
}

func NewDBCache(parent *DatabaseCache) *DatabaseCache {
	dbCache := &DatabaseCache{
		parent:    parent,
		cacheLock: sync.RWMutex{},
		cache:     make(map[cacheEntryKey]cacheEntryValue),
	}
	return dbCache
}

func (dbCache *DatabaseCache) GetDefaultExpire() time.Time {
	return time.Now().Add(time.Minute * 1)
}

func (dbCache *DatabaseCache) Add(obj interface{}) {
	pk := utilReflection.FindPrimaryKey(obj)
	if pk == 0 {
		return
	}
	if dbCache.parent != nil {
		go dbCache.parent.Add(obj)
	}
	dbCache.cacheLock.Lock()
	defer dbCache.cacheLock.Unlock()
	dbCache.cache[cacheEntryKey{reflect.TypeOf(obj), pk}] = cacheEntryValue{obj, dbCache.GetDefaultExpire()}
}

func (dbCache *DatabaseCache) Get(obj interface{}) *interface{} {
	pk := utilReflection.FindPrimaryKey(obj)
	if pk == 0 {
		return nil
	}
	return dbCache.GetByPK(obj, pk)
}

func (dbCache *DatabaseCache) GetByPK(obj interface{}, pk database.ID) *interface{} {
	key := cacheEntryKey{reflect.TypeOf(obj), pk}
	cached, ok := dbCache.cache[key]
	if !ok {
		return nil
	}
	if time.Now().After(cached.Expire) {
		delete(dbCache.cache, key)
		return nil
	}
	reflect.ValueOf(obj).Elem().Set(reflect.ValueOf(cached.Obj).Elem())
	return &obj
}

var globalDatabaseCache = NewDBCache(nil)

func CtxWithDBCache(ctx context.Context) context.Context {
	return context.WithValue(ctx, "dbCache", globalDatabaseCache)
}

func DBCacheFromCtx(ctx context.Context) *DatabaseCache {
	return ctx.Value("dbCache").(*DatabaseCache)
}

func CtxWithDBAuthCache(ctx context.Context) context.Context {
	return context.WithValue(ctx, "dbAuthCache", NewDBCache(DBCacheFromCtx(ctx)))
}

func DBAuthCacheFromCtx(ctx context.Context) *DatabaseCache {
	return ctx.Value("dbAuthCache").(*DatabaseCache)
}
