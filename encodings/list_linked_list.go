package encodings

import (
	"time"

	re "github.com/SwanSpouse/redis_go/error"
	"github.com/SwanSpouse/redis_go/loggers"
	"github.com/SwanSpouse/redis_go/raw_type"
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

func (ll *ListLinkedList) LPush(vals []string) int {
	list := ll.GetValue().(*raw_type.List)
	var newList *raw_type.List
	for _, val := range vals {
		newList = list.ListAddNodeHead(val)
	}
	ll.SetValue(newList)
	return list.ListLength()
}

func (ll *ListLinkedList) RPush(values []string) int {
	list := ll.GetValue().(*raw_type.List)
	var newList *raw_type.List
	for _, val := range values {
		newList = list.ListAddNodeTail(val)
	}
	ll.SetValue(newList)
	return list.ListLength()
}

func (ll *ListLinkedList) LPop() string {
	list := ll.GetValue().(*raw_type.List)
	value := list.ListFirst().NodeValue()
	ll.SetValue(list.ListRemoveNode(list.ListFirst()))
	return value
}

func (ll *ListLinkedList) RPop() string {
	list := ll.GetValue().(*raw_type.List)
	value := list.ListLast().NodeValue()
	ll.SetValue(list.ListRemoveNode(list.ListLast()))
	return value
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
	if curNode := list.ListSearchKey(val); curNode == nil {
		return -1, nil
	} else {
		var newList *raw_type.List
		for _, newValue := range values {
			newList = list.ListInsertNode(curNode, newValue, insertFlag == RedisTypeListInsertAfter)
		}
		ll.SetValue(newList)
		return newList.ListLength(), nil
	}
}

/**
count > 0 : 从表头开始向表尾搜索，移除与 value 相等的元素，数量为 count 。
count < 0 : 从表尾开始向表头搜索，移除与 value 相等的元素，数量为 count 的绝对值。
count = 0 : 移除表中所有与 value 相等的值。
*/
func (ll *ListLinkedList) LRem(count int, key string) int {
	list := ll.GetValue().(*raw_type.List)

	succCount := 0
	if count == 0 {
		count = list.ListLength() + 1
	}
	if count > 0 {
		for node := list.ListFirst(); node != nil && count > 0; node = node.NodeNext() {
			if node.NodeValue() == key {
				count -= 1
				succCount += 1
				list = list.ListRemoveNode(node)
			}
		}
	} else {
		count = -count
		for node := list.ListLast(); node != nil && count > 0; node = node.NodePrev() {
			if node.NodeValue() == key {
				count -= 1
				succCount += 1
				list = list.ListRemoveNode(node)
			}
		}
	}
	if succCount > 0 {
		ll.SetValue(list)
	}
	return succCount
}

//TODO lmj
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

func (ll *ListLinkedList) GetAllMembers() []string {
	ret := make([]string, 0)
	if linkedList, ok := ll.GetValue().(*raw_type.List); !ok {
		return ret
	} else {
		for node := linkedList.ListFirst(); node != nil; node = node.NodeNext() {
			ret = append(ret, node.NodeValue())
		}
		return ret
	}
}

func (ll *ListLinkedList) Debug() {
	loggers.Info(ll.String())
}
