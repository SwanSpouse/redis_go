package raw_type

import "github.com/mitchellh/hashstructure"

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

type dictEntry struct {
	key   interface{}
	value interface{}
	next  *dictEntry
}

func NewDictEntry(has int, key, value interface{}, next *dictEntry) *dictEntry {
	return &dictEntry{
		key: key, value: value, next: next,
	}
}

func (de *dictEntry) getKey() interface{} {
	return de.key
}

func (de *dictEntry) getValue() interface{} {
	return de.value
}

func (de *dictEntry) hasCode() uint64 {
	hash, _ := hashstructure.Hash(de, nil)
	return hash
}
