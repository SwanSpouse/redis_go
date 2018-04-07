package database

import (
	"sync/atomic"
)

var (
	databaseInc = uint64(0)
)

type Database struct {
	id      uint64           // 数据库编号
	dict    map[string]TBase // 数据库
	expires map[string]int64 // Key过期时间
}

func NewDatabase() *Database {
	return &Database{
		id:      atomic.AddUint64(&databaseInc, 1),
		dict:    make(map[string]TBase),
		expires: make(map[string]int64),
	}
}

// 获取Key在数据库中对应的Value
func (db *Database) SearchKeyInDB(key string) (TBase, error) {
	if obj, ok := db.dict[key]; !ok {
		return nil, nil
	} else {
		// TODO 增加是否过期的判断
		return obj, nil
	}
}

func (db *Database) SetKeyInDB(key string, obj TBase) {
	db.dict[key] = obj
}
