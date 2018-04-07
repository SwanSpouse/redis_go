package database

import (
	"sync/atomic"
)

var (
	databaseInc = uint64(0)
)

type Database struct {
	id   uint64           // 数据库编号
	dict map[string]TBase // 数据库
}

func NewDatabase() *Database {
	return &Database{
		id:   atomic.AddUint64(&databaseInc, 1),
		dict: make(map[string]TBase),
	}
}

// 获取Key在数据库中对应的Value
func (db *Database) SearchKeyInDB(key string) TBase {
	if obj, ok := db.dict[key]; ok && !obj.IsExpired() {
		return obj
	}
	return nil
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
	db.dict[key] = obj
}

// 删除redis db 中的key
func (db *Database) RemoveKeyInDB(keys []string) int64 {
	var successCount int64
	for _, key := range keys {
		if _, ok := db.dict[key]; ok {
			delete(db.dict, key)
			successCount += 1
		}
	}
	return successCount
}
