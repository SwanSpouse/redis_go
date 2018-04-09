package raw_type

import (
	"github.com/mitchellh/hashstructure"
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
)

/************************************  dictEntry  ***************************************/
type dictEntry struct {
	Key   interface{}
	Value interface{} `hash:"ignore"` // 表明在计算哈希值的时候，这个域不参与计算
	next  *dictEntry  `hash:"ignore"` // 表明在计算哈希值的时候，这个域不参与计算
}

func NewDictEntry(key, value interface{}, next *dictEntry) *dictEntry {
	return &dictEntry{
		Key: key, Value: value, next: next,
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
	
}

/************************************     dict    ***************************************/

type Dict struct {
}

