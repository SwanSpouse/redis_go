package raw_type

import (
	"github.com/mitchellh/hashstructure"
	"sync"
)

const (
	/**
	The largest possible table capacity.
	*/
	MaximumCapacity = 1 << 30

	/**
	  The default initial table capacity.  Must be a power of 2
	  (i.e., at least 1) and at most MaximumCapacity.
	*/
	DefaultCapacity = 1 << 4

	LoadFactory = 0.75

	// 不使用锁情况下最多尝试次数
	MaxScanRetries = 2
)

/************************************  dictEntry  ***************************************/
type dictEntry struct {
	Key   interface{}
	Value interface{} `hash:"ignore"` // 这个域不参与哈希值计算
	next  *dictEntry  `hash:"ignore"` // 这个域不参与哈希值计算
	hash  int         `hash:"ignore"` // 这个域不参与哈希值计算
}

func NewDictEntry(hash int, key, value interface{}, next *dictEntry) *dictEntry {
	if key == nil || value == nil {
		return nil
	}
	return &dictEntry{
		hash: hash, Key: key, Value: value, next: next,
	}
}

func (de *dictEntry) getKey() interface{} {
	return de.Key
}

func (de *dictEntry) getValue() interface{} {
	return de.Value
}

func (de *dictEntry) hasCode() uint64 {
	hash, _ := hashstructure.Hash(de, nil)
	return hash
}

func (de *dictEntry) equals(obj interface{}) bool {
	switch obj.(type) {
	case *dictEntry, dictEntry:
		if instance, ok := obj.(*dictEntry); ok && instance != nil {
			return de.Key == instance.Key && de.Value == instance.Value
		}
		if instance, ok := obj.(dictEntry); ok {
			return de.Key == instance.Key && de.Value == instance.Value
		}
	default:
		return false
	}
	return false
}

/************************************   segment   ***************************************/

type segment struct {
	table      []*dictEntry
	count      int64
	modCount   int64
	threshold  int64
	loadFactor float64
	locker     *sync.RWMutex // locker for segment
}

func newSegement(lf float64, threshold int64, table []*dictEntry) *segment {
	return &segment{
		loadFactor: lf,
		threshold:  threshold,
		table:      table,
		locker:     new(sync.RWMutex),
	}
}

/*
	put操作:(TODO)
		1. 首先在不加锁的情况下进行尝试，如果在put的过程中没有遇到并发修改，则顺利插入。
		2. 多线程情况下：如果在不加锁的情况下遇到了冲突导致插入失败，则进行尝试。如果一直都失败，并且达到MaxScanRetires的上限。则先lock再put

	现在的做法是直接上锁
*/
func (seg *segment) put(key, value interface{}) interface{} {
	seg.locker.Lock()
	defer seg.locker.Unlock()

	var oldValue interface{}
	index := hash(key) % len(seg.table)
	e := seg.table[index]
	for true {
		if e != nil {
			// 如果找到相同元素，先记录oldValue，再覆盖
			if hash(key) == e.hash && key == e.Key {
				oldValue = e.Value
				e.Value = value
				seg.modCount += 1
				break
			}
			e = e.next
		} else {
			// 如果遍历之后没有找到key相同的元素，则利用头插法插入新node
			node := NewDictEntry(hash(key), key, value, seg.table[index])
			// 如果节点容量达到了rehash的条件，那么就进行rehash
			if seg.count+1 > seg.threshold && len(seg.table) < MaximumCapacity {
				seg.rehash(node)
			} else {
				seg.table[index] = node
			}
			seg.modCount += 1
			seg.count += 1
			oldValue = nil
			break
		}
	}
	return oldValue
}

func (seg *segment) rehash(node *dictEntry) {

}

func (seg *segment) remove(key, value interface{}) interface{} {
	return nil
}

func (seg *segment) replace(key, oldValue, newValue interface{}) bool {
	seg.locker.Lock()
	defer seg.locker.Unlock()

	index := hash(key) % len(seg.table)
	for e := seg.table[index]; e != nil; e = e.next {
		// 查找目标元素
		if hash(key) == e.hash && key == e.Key && oldValue == e.Value {
			e.Value = newValue
			seg.modCount += 1
			return true
		}
	}
	return false
}

func (seg *segment) clear() {
	seg.locker.Lock()
	defer seg.locker.Lock()

	// 将所有entry置为nil
	for i := 0; i < len(seg.table); i++ {
		seg.table[i] = nil
	}
	seg.modCount += 1
	seg.count = 0
}

/************************************     dict    ***************************************/

type Dict struct {
	segments []*segment
}

/************************************   common   ***************************************/
func hash(value interface{}) int {
	hash, _ := hashstructure.Hash(value, nil)
	return int(hash)
}
