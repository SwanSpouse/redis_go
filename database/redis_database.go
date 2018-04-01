package database

import (
	"sync/atomic"
)

var (
	databaseInc = uint64(0)
)

type Database struct {
	id      uint64                  // 数据库编号
	dict    map[string]*RedisObject // 数据库
	expires map[string]int64        // Key过期时间
}

func NewDatabase() *Database {
	return &Database{
		id:      atomic.AddUint64(&databaseInc, 1),
		dict:    make(map[string]*RedisObject),
		expires: make(map[string]int64),
	}
}

func (db *Database) SearchKeyInDB(key string) *RedisObject {
	if obj, ok := db.dict[key]; !ok {
		return nil
	} else {
		// TODO 增加是否过期的判断
		return obj
	}
}

func (db *Database) SetKeyInDB(key string, obj *RedisObject) {
	db.dict[key] = obj
}
