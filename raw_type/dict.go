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
	count      int
	modCount   int
	threshold  int
	loadFactor float64
	locker     *sync.RWMutex // locker for segment
}

func newSegment(capacity int, lf float64, threshold int) *segment {
	return &segment{
		table:      make([]*dictEntry, capacity),
		locker:     new(sync.RWMutex),
		loadFactor: lf,
		threshold:  threshold,
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

/*
	先实现一个concurrentHashMap的版本

	TODO 然后再参考redis dict的rehash实现一个版本
*/
func (seg *segment) rehash(node *dictEntry) {
	seg.locker.Lock()
	defer seg.locker.Lock()

	oldTable := seg.table
	oldCapacity := len(seg.table)
	newCapacity := oldCapacity << 1
	threshold := newCapacity * LoadFactory

	newTable := newSegment(newCapacity, LoadFactory, threshold)
	// 将老table中的数据迁移到新table中去
	for i := 0; i < oldCapacity; i++ {
		/**
		rehash:
			1. 如果oldTable对应序号上无元素，则不需要进行操作
			2. 如果oldTable对应序号上有元素，利用头插法依次挪动oldTable某个序号上的所有node到新table
		*/
		if oldTable[i] == nil {
			continue
		}
		for cur := oldTable[i]; cur != nil; cur = cur.next {
			newIdx := hash(cur.Key) % len(newTable.table)
			entry := newTable.table[newIdx]
			newTable.table[newIdx] = NewDictEntry(hash(cur.Key), cur.Key, cur.Value, entry)
		}
	}
	/**
		oldTable中的节点迁移完成后，添加新节点(这里可以直接添加是因为进到这里的节点，肯定不会和table中的某个节点key相同)
	同时在添加结束后没有把table的count++, modCount++ 是因为这些操作都在put里面已经进行过了。
	*/
	idx := hash(node.Value) % len(newTable.table)
	node.next = newTable.table[idx]
	newTable.table[idx] = node
	// 在rehash完成的时候切换成新的table
	seg.table = newTable.table
}

// 删除segment中的特定元素
func (seg *segment) remove(key, value interface{}) interface{} {
	seg.locker.Lock()
	defer seg.locker.Unlock()
	var oldValue interface{}

	idx := hash(key) % len(seg.table)
	if seg.table[idx] == nil {
		return nil
	}
	var pre *dictEntry
	for cur := seg.table[idx]; cur != nil; cur = cur.next {
		if cur.hash == hash(key) && cur.Key == key && cur.Value == value {
			if pre == nil {
				seg.table[idx] = cur.next
			} else {
				pre.next = cur.next
			}
			oldValue = cur.Value
			seg.count -= 1
			seg.modCount += 1
			break
		}
		pre = cur
	}
	return oldValue
}

/**
用newValue替换segment中key对应的oldValue
	@param key:
	@param oldValue: key对应的原有的值
	@param newValue: 需要替换成为的值
*/
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

// 清空整个segment
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
