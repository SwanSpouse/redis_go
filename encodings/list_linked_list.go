package encodings

import (
	re "redis_go/error"
	"redis_go/loggers"
	"redis_go/raw_type"
	"time"
)

const (
	RedisTypeListInsertBefore = 0
	RedisTypeListInsertAfter  = 1
)

type ListLinkedList struct {
	RedisObject
}

func NewListLinkedList(ttl int) *ListLinkedList {
	var expireTime time.Time
	if ttl > 0 {
		expireTime = time.Now().Add(time.Duration(ttl) * time.Second)
	}
	return &ListLinkedList{
		RedisObject: RedisObject{
			objectType: RedisTypeList,
			encoding:   RedisEncodingLinkedList,
			ttl:        ttl,
			value:      raw_type.ListCreate(),
			expireTime: expireTime,
		},
	}
}

func (ll *ListLinkedList) LPush(val string) int {
	list := ll.GetValue().(*raw_type.List)
	newList := list.ListAddNodeHead(val)
	ll.SetValue(newList)
	return list.ListLength()
}

func (ll *ListLinkedList) RPush(val string) int {
	list := ll.GetValue().(*raw_type.List)
	newList := list.ListAddNodeTail(val)
	ll.SetValue(newList)
	return list.ListLength()
}

func (ll *ListLinkedList) LPop() string {
	list := ll.GetValue().(*raw_type.List)
	return list.ListFirst().NodeValue()
}

func (ll *ListLinkedList) RPop() string {
	list := ll.GetValue().(*raw_type.List)
	return list.ListLast().NodeValue()
}

func (ll *ListLinkedList) LIndex(index int) (string, error) {
	list := ll.GetValue().(*raw_type.List)
	if node := list.ListIndex(index); node == nil {
		return "", re.ErrNilValue
	} else {
		return node.NodeValue(), nil
	}
}

func (ll *ListLinkedList) LLen() int {
	list := ll.GetValue().(*raw_type.List)
	return list.ListLength()
}

func (ll *ListLinkedList) LInsert(insertFlag int, val string, values ...string) (int, error) {
	list := ll.GetValue().(*raw_type.List)
	if insertFlag != RedisTypeListInsertBefore && insertFlag != RedisTypeListInsertAfter {
		return 0, re.ErrSyntaxError
	}
	if curNode := list.ListSearchKey(val); curNode == nil {
		return -1, nil
	} else {
		newList := list.ListInsertNode(curNode, val, insertFlag == RedisTypeListInsertBefore)
		ll.SetValue(newList)
		return newList.ListLength(), nil
	}
}

func (ll *ListLinkedList) LRem(count int, key string) int {
	list := ll.GetValue().(*raw_type.List)

	succCount := 0
	for i := 0; i < count; i++ {
		if node := list.ListSearchKey(key); node == nil {
			continue
		} else {
			list = list.ListRemoveNode(node)
			succCount += 1
		}
	}
	if succCount > 0 {
		ll.SetValue(list)
	}
	return succCount
}

func (ll *ListLinkedList) LTrim(start int, stop int) error {
	//list := ll.GetValue().(*raw_type.List)
	//
	//if start >= list.ListLength() || -start >= list.ListLength() {
	//	return re.ErrNotIntegerOrOutOfRange
	//}
	//if stop >= list.ListLength() {
	//	stop = list.ListLength() - 1
	//} else if -stop >= list.ListLength() {
	//	stop = -(list.ListLength() - 1)
	//}
	//startNode := list.ListIndex(start)
	//stopNode := list.ListIndex(stop)
	return nil
}

func (ll *ListLinkedList) LSet(index int, val string) error {
	list := ll.GetValue().(*raw_type.List)
	if index >= list.ListLength() || -index >= list.ListLength() {
		return re.ErrNotIntegerOrOutOfRange
	}
	node := list.ListIndex(index)
	node.SetNodeValue(val)
	return nil
}

func (ll *ListLinkedList) LRange(start int, stop int) []string {
	list := ll.GetValue().(*raw_type.List)
	ret := make([]string, 0)

	if start >= list.ListLength() || -start >= list.ListLength() {
		return ret
	}
	if stop >= list.ListLength() {
		stop = list.ListLength() - 1
	} else if -stop >= list.ListLength() {
		stop = -(list.ListLength() - 1)
	}
	startNode := list.ListIndex(start)
	stopNode := list.ListIndex(stop)

	for node := startNode; node != stopNode.NodeNext(); node = node.NodeNext() {
		ret = append(ret, node.NodeValue())
	}
	return ret
}

func (ll *ListLinkedList) String() string {
	ret := "CURRENT_LIST:"
	if linkedList, ok := ll.GetValue().(*raw_type.List); !ok {
		return re.ErrWrongType.Error()
	} else {
		for node := linkedList.ListFirst(); node != nil; node = node.NodeNext() {
			ret += node.String()
		}
		return ret
	}
}

func (ll *ListLinkedList) Debug() {
	loggers.Info(ll.String())
}
