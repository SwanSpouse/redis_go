package raw_type

import (
	"fmt"
	"github.com/mitchellh/hashstructure"
	"math"
	"redis_go/log"
	"reflect"
	"sync"
)

const (
	MaximumCapacity         = 1 << 30
	DefaultCapacity         = 1 << 4
	DefaultConcurrencyLevel = 16
	MaxSegments             = 1 << 16
	MinSegmentTableCapacity = 2
	LoadFactory             = 0.75
	MaxScanRetries          = 2 // 不使用锁情况下最多尝试次数
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
			return reflect.DeepEqual(de.Key, instance.Key) && reflect.DeepEqual(de.Value, instance.Value)
		}
		if instance, ok := obj.(dictEntry); ok {
			return reflect.DeepEqual(de.Key, instance.Key) && reflect.DeepEqual(de.Value, instance.Value)
		}
	default:
		return false
	}
	return false
}

func (de *dictEntry) String() string {
	return fmt.Sprintf("{%v=%v}", de.Key, de.Value)
}

/************************************   segment   ***************************************/

type segment struct {
	table      []*dictEntry
	count      int
	sizeMask   int
	modCount   int
	threshold  int
	loadFactor float64
	locker     *sync.RWMutex // locker for segment
}

func newSegment(capacity int, lf float64, threshold int) *segment {
	return &segment{
		table:      make([]*dictEntry, capacity),
		sizeMask:   capacity - 1,
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
func (seg *segment) put(hashCode int, key, value interface{}) interface{} {
	seg.locker.Lock()
	defer seg.locker.Unlock()

	var oldValue interface{}
	index := hashCode & seg.sizeMask
	e := seg.table[index]
	for true {
		if e != nil {
			// 如果找到相同元素，先记录oldValue，再覆盖
			if hashCode == e.hash && reflect.DeepEqual(key, e.Key) {
				oldValue = e.Value
				e.Value = value
				seg.modCount += 1
				break
			}
			e = e.next
		} else {
			// 如果遍历之后没有找到key相同的元素，则利用头插法插入新node
			node := NewDictEntry(hashCode, key, value, seg.table[index])
			// 如果节点容量达到了rehash的条件，那么就进行rehash
			if seg.count+1 > seg.threshold && len(seg.table) < MaximumCapacity {
				seg.rehash(node)
			} else {
				seg.table[index] = node
			}
			seg.modCount += 1
			seg.count += 1
			return nil
		}
	}
	return oldValue
}

/*
	先实现一个concurrentHashMap的版本

	TODO 然后再参考redis dict的rehash实现一个版本
*/
func (seg *segment) rehash(node *dictEntry) {
	oldTable := seg.table
	oldCapacity := len(oldTable)
	newCapacity := oldCapacity << 1
	threshold := int(float32(newCapacity) * LoadFactory)
	log.Debug("segment start rehash enlarge size from %d to %d", oldCapacity, newCapacity)
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
			newIdx := cur.hash & newTable.sizeMask
			entry := newTable.table[newIdx]
			newTable.table[newIdx] = NewDictEntry(cur.hash, cur.Key, cur.Value, entry)
		}
	}
	/**
	    oldTable中的节点迁移完成后，添加新节点(这里可以直接添加是因为进到这里的节点，肯定不会和table中的某个节点key相同)
	同时在添加结束后没有把table的count++, modCount++ 是因为这些操作都在put里面已经进行过了。
	*/
	newIdx := node.hash & newTable.sizeMask
	node.next = newTable.table[newIdx]
	newTable.table[newIdx] = node

	// 在rehash完成的时候切换成新的table
	seg.table = newTable.table
	seg.sizeMask = newCapacity - 1
	seg.threshold = threshold
	seg.loadFactor = LoadFactory
}

/*
删除segment中的特定元素,
	如果value == nil 则删除key相等的。
	如果value != nil 则删除key相等且value相等的。
*/
func (seg *segment) remove(hashCode int, key, value interface{}) interface{} {
	seg.locker.Lock()
	defer seg.locker.Unlock()
	var oldValue interface{}

	idx := hashCode & seg.sizeMask
	if seg.table[idx] == nil {
		return nil
	}
	var pre *dictEntry
	for cur := seg.table[idx]; cur != nil; cur = cur.next {
		if cur.hash == hashCode && reflect.DeepEqual(key, cur.Key) &&
			(value == nil || reflect.DeepEqual(value, cur.Value)) {
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
	@param hashCode: key对应的哈希值，节省计算
	@param key:
	@param oldValue: key对应的原有的值
	@param newValue: 需要替换成为的值
*/
func (seg *segment) replace(hashCode int, key, oldValue, newValue interface{}) bool {
	seg.locker.Lock()
	defer seg.locker.Unlock()

	index := hashCode & seg.sizeMask
	for e := seg.table[index]; e != nil; e = e.next {
		// 查找目标元素
		if hashCode == e.hash && reflect.DeepEqual(key, e.Key) && reflect.DeepEqual(oldValue, e.Value) {
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
	defer seg.locker.Unlock()

	// 将所有entry置为nil
	for i := 0; i < len(seg.table); i++ {
		seg.table[i] = nil
	}
	seg.modCount += 1
	seg.count = 0
}

func (seg *segment) printSegForDebug() {
	fmt.Printf("segement has %d slots and %d entries, threshold:%d\n", len(seg.table), seg.count, seg.threshold)
	for i := 0; i < len(seg.table); i++ {
		if seg.table[i] != nil {
			fmt.Printf("SEG[%d]\t->", i)
			for e := seg.table[i]; e != nil; e = e.next {
				fmt.Printf("%s->", e)
			}
			fmt.Printf("\n")
		}
	}
}

/************************************     dict    ***************************************/

type Dict struct {
	segments []*segment
}

func NewDict() *Dict {
	return NewDictWithCapacity(DefaultCapacity)
}

func NewDictWithCapacity(capacity int) *Dict {
	return NewDictWithCapacityAndConcurrencyLevel(capacity, DefaultConcurrencyLevel)
}

func NewDictWithCapacityAndConcurrencyLevel(capacity, concurrencyLevel int) *Dict {
	if capacity < 0 || concurrencyLevel <= 0 {
		panic("Illegal argument exception")
	}
	if concurrencyLevel > MaxSegments {
		concurrencyLevel = MaxSegments
	}
	if capacity > MaximumCapacity {
		capacity = MaximumCapacity
	}
	actualConcurrentLevel := 1
	for actualConcurrentLevel < concurrencyLevel {
		actualConcurrentLevel = actualConcurrentLevel << 1
	}
	expectedSegmentCap := capacity / actualConcurrentLevel
	if expectedSegmentCap*actualConcurrentLevel < capacity {
		expectedSegmentCap += 1
	}
	actualSegmentCap := MinSegmentTableCapacity
	for actualSegmentCap < expectedSegmentCap {
		actualSegmentCap <<= 1
	}
	log.Info("[SEGMENT] totalCapacity %+v, segment cap %+v, actualSegCap %+v", capacity, expectedSegmentCap, actualSegmentCap)
	segments := make([]*segment, concurrencyLevel)
	for i := 0; i < concurrencyLevel; i++ {
		segments[i] = newSegment(actualSegmentCap, LoadFactory, int(float32(expectedSegmentCap)*LoadFactory))
	}
	return &Dict{
		segments: segments,
	}
}

func (dict *Dict) IsEmpty() bool {
	sum := 0
	for i := 0; i < len(dict.segments); i++ {
		if dict.segments[i] != nil {
			if dict.segments[i].count != 0 {
				return false
			}
			sum += dict.segments[i].modCount
		}
	}
	if sum != 0 {
		for i := 0; i < len(dict.segments); i++ {
			if dict.segments[i] != nil {
				if dict.segments[i].count != 0 {
					return false
				}
				sum -= dict.segments[i].modCount
			}
		}
		if sum != 0 {
			return false
		}
	}
	return true
}

/*
	TODO lmj 改造成先尝试几次然后再加锁的逻辑
*/
func (dict *Dict) Size() int {
	for i := 0; i < len(dict.segments); i++ {
		dict.segments[i].locker.Lock()
	}
	defer func() {
		for i := 0; i < len(dict.segments); i++ {
			dict.segments[i].locker.Unlock()
		}
	}()

	var isOverflow bool
	var size int
	for i := 0; i < len(dict.segments); i++ {
		if dict.segments[i] == nil {
			continue
		}
		curSize := dict.segments[i].count
		if curSize < 0 || size+curSize < 0 {
			isOverflow = true
			break
		}
		size += curSize
	}
	if isOverflow {
		return MaximumCapacity
	}
	return size
}

/*
	根据hashCode先找到segment，再找dictEntry
*/
func (dict *Dict) Get(key interface{}) interface{} {
	hashCode := hash(key)
	segmentIdx := hashCode & (len(dict.segments) - 1)
	if dict.segments[segmentIdx] != nil {
		idx := hashCode & dict.segments[segmentIdx].sizeMask
		if dict.segments[segmentIdx].table[idx] != nil {
			for e := dict.segments[segmentIdx].table[idx]; e != nil; e = e.next {
				if e.hash == hashCode && reflect.DeepEqual(e.Key, key) {
					return e.Value
				}
			}
		}
	}
	return nil
}

/*
	和Get思路相同
*/
func (dict *Dict) ContainsKey(key interface{}) bool {
	return dict.Get(key) != nil
}

/*
	TODO 修改先尝试，再加锁
*/
func (dict *Dict) ContainsValue(value interface{}) bool {
	for i := 0; i < len(dict.segments); i++ {
		dict.segments[i].locker.RLock()
	}
	defer func() {
		for i := 0; i < len(dict.segments); i++ {
			dict.segments[i].locker.RUnlock()
		}
	}()

	for i := 0; i < len(dict.segments); i++ {
		for j := 0; dict.segments[i] != nil && j < len(dict.segments[i].table); j++ {
			for e := dict.segments[i].table[j]; e != nil; e = e.next {
				if reflect.DeepEqual(value, e.Value) {
					return true
				}
			}
		}
	}
	return false
}

func (dict *Dict) Contains(value interface{}) bool {
	return dict.ContainsValue(value)
}

func (dict *Dict) Put(key, value interface{}) {
	if key == nil || value == nil {
		panic("PUT key or value null pointer exception")
	}
	hashCode := hash(key)
	segmentIdx := hashCode & (len(dict.segments) - 1)
	dict.segments[segmentIdx].put(hashCode, key, value)
}

/*
	TODO lmj
*/
func (dict *Dict) PutAll(otherDict *Dict) {

}

func (dict *Dict) Remove(key, value interface{}) interface{} {
	if key == nil {
		panic("REMOVE key null pointer exception")
	}
	hashCode := hash(key)
	segmentIdx := hashCode & (len(dict.segments) - 1)
	return dict.segments[segmentIdx].remove(hashCode, key, value)
}

func (dict *Dict) RemoveKey(key interface{}) interface{} {
	return dict.Remove(key, nil)
}

func (dict *Dict) Replace(key, oldValue, newValue interface{}) interface{} {
	if key == nil {
		panic("REPLACE key null pointer exception")
	}
	hashCode := hash(key)
	segmentIdx := hashCode & (len(dict.segments) - 1)
	return dict.segments[segmentIdx].replace(hashCode, key, oldValue, newValue)
}

/*
清空Dict中所有的k-v
*/
func (dict *Dict) Clear() {
	for i := 0; i < len(dict.segments); i++ {
		if dict.segments[i] != nil {
			dict.segments[i].clear()
		}
	}
}

/*
	TODO lmj
*/
func (dict *Dict) KeySet() map[interface{}]bool {
	return nil
}

/*
	TODO lmj
*/
func (dict *Dict) Values() map[interface{}]bool {
	return nil
}

func (dict *Dict) printDictForDebug() {
	fmt.Printf("dict has %d segment and %d entries\n", len(dict.segments), dict.Size())
	for i := 0; i < len(dict.segments); i++ {
		if dict.segments[i] == nil {
			continue
		}
		fmt.Printf("==========   segement[%d]   ==========\n", i)
		dict.segments[i].printSegForDebug()
	}
}

/************************************     common   **************************************/
func hash(value interface{}) int {
	hashCode, _ := hashstructure.Hash(value, nil)
	return int(hashCode & uint64(math.MaxInt32))
}
