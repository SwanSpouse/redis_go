package database

import (
	"redis_go/loggers"
	"redis_go/raw_type"
	"sync/atomic"
)

var (
	databaseInc = uint64(0)
)

type Database struct {
	id   uint64         // 数据库编号
	dict *raw_type.Dict // 数据库
}

func NewDatabase() *Database {
	return &Database{
		id:   atomic.AddUint64(&databaseInc, 1),
		dict: raw_type.NewDict(),
	}
}

// 获取Key在数据库中对应的Value
func (db *Database) SearchKeyInDB(key string) TBase {
	if obj := db.dict.Get(key); obj == nil {
		return nil
	} else {
		if tBase, ok := obj.(TBase); !ok || tBase.IsExpired() {
			loggers.Errorf("illegal value in database.dict or tBase is expired. key %s", key)
			return nil
		} else {
			return tBase
		}
	}
}

func (db *Database) SearchKeysInDB(keys []string) ([]TBase, error) {
	ret := make([]TBase, 0)
	for _, key := range keys {
		val := db.SearchKeyInDB(key)
		if val != nil {
			ret = append(ret, val)
		}
	}
	return ret, nil
}

// 将TBase写入到redis database
func (db *Database) SetKeyInDB(key string, obj TBase) {
	db.dict.Put(key, obj)
}

// 删除redis db 中的key
func (db *Database) RemoveKeyInDB(keys []string) int64 {
	var successCount int64
	for _, key := range keys {
		if oldValue := db.dict.RemoveKey(key); oldValue != nil {
			successCount += 1
		}
	}
	return successCount
}

// 获取数据库中所有的key
func (db *Database) GetAllKeys() []string {
	ret := make([]string, 0)
	for key := range db.dict.KeySet() {
		ret = append(ret, key.(string))
	}
	return ret
}

func (db *Database) DBSize() int {
	return db.dict.Size()
}
