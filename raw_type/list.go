package raw_type

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/SwanSpouse/redis_go/loggers"
)

/*==============================      ListNode       =================================*/

type ListNode struct {
	pre   *ListNode
	next  *ListNode
	value string
}

func newListNode(value string) *ListNode {
	return &ListNode{value: value}
}

func (node *ListNode) NodePrev() *ListNode {
	return node.pre
}

func (node *ListNode) NodeNext() *ListNode {
	return node.next
}

func (node *ListNode) NodeValue() string {
	return node.value
}

func (node *ListNode) SetNodeValue(val string) {
	node.value = val
}

func (node *ListNode) String() string {
	if node == nil {
		return fmt.Sprint("[LNODE:nil]")
	}
	return fmt.Sprintf("[LNODE:%+v]", node.value)
}

/*==============================      List       =================================*/

type List struct {
	head   *ListNode
	tail   *ListNode
	length int
	Locker *sync.Mutex
}

func ListCreate() *List {
	return &List{
		head: nil, tail: nil, length: 0, Locker: new(sync.Mutex),
	}
}

func (list *List) ListLength() int {
	return list.length
}

func (list *List) ListFirst() *ListNode {
	return list.head
}

func (list *List) ListLast() *ListNode {
	return list.tail
}

func (list *List) ListAddNodeHead(value string) *List {
	list.Locker.Lock()
	defer list.Locker.Unlock()

	node := newListNode(value)
	if list.length == 0 {
		list.head = node
		list.tail = node
		node.next = nil
		node.pre = nil
	} else {
		node.pre = nil
		node.next = list.head
		list.head.pre = node
		list.head = node
	}

	list.length += 1
	return list
}

func (list *List) ListAddNodeTail(value string) *List {
	list.Locker.Lock()
	defer list.Locker.Unlock()

	node := newListNode(value)
	if list.length == 0 {
		list.head = node
		list.tail = node
		node.next = nil
		node.pre = nil
	} else {
		node.pre = list.tail
		node.next = nil
		list.tail.next = node
		list.tail = node
	}

	list.length += 1
	return list
}

func (list *List) ListInsertNode(oldNode *ListNode, value string, after bool) *List {
	list.Locker.Lock()
	defer list.Locker.Unlock()

	if oldNode == nil {
		return nil
	}

	newNode := newListNode(value)
	if after {
		newNode.pre = oldNode
		newNode.next = oldNode.next
		if list.tail == oldNode {
			list.tail = newNode
		}
	} else {
		newNode.next = oldNode
		newNode.pre = oldNode.pre
		if list.head == oldNode {
			list.head = newNode
		}
	}

	if newNode.pre != nil {
		newNode.pre.next = newNode
	}
	if newNode.next != nil {
		newNode.next.pre = newNode
	}

	list.length += 1
	return list
}

func (list *List) ListSearchKey(key string) *ListNode {
	list.Locker.Lock()
	defer list.Locker.Unlock()

	for node := list.head; node != nil; node = node.next {
		if node.value == key && reflect.DeepEqual(node.value, key) {
			return node
		}
	}
	return nil
}

func (list *List) ListIndex(index int) *ListNode {
	list.Locker.Lock()
	defer list.Locker.Unlock()

	var node *ListNode
	if index < 0 {
		index = -index - 1
		node = list.tail
		for ; node != nil && index > 0; node = node.pre {
			index -= 1
		}
	} else {
		node = list.head
		for ; node != nil && index > 0; node = node.next {
			index -= 1
		}
	}
	return node
}

func (list *List) ListRemoveNode(oldNode *ListNode) *List {
	list.Locker.Lock()
	defer list.Locker.Unlock()

	if oldNode == nil {
		return list
	}

	// 处理前驱节点
	if oldNode.pre != nil {
		oldNode.pre.next = oldNode.next
	} else {
		list.head = oldNode.next
	}
	// 处理后续节点
	if oldNode.next != nil {
		oldNode.next.pre = oldNode.pre
	} else {
		list.tail = oldNode.pre
	}

	list.length -= 1
	return list
}

func (list *List) ListRotate() *List {
	list.Locker.Lock()
	defer list.Locker.Unlock()

	if list.length <= 1 {
		return list
	}
	newTail := list.tail

	list.tail = newTail.pre
	list.tail.next = nil

	list.head.pre = newTail
	newTail.pre = nil
	newTail.next = list.head
	list.head = newTail
	return list
}

func (list *List) PrintListForDebug() {
	if list == nil {
		loggers.Info("current list length is nil")
		return
	}
	loggers.Info("current list length is %d, headNode: %s, TailNode: %s", list.length, list.head, list.tail)
	var msg string
	for node := list.head; node != nil; node = node.next {
		msg += fmt.Sprintf("%+v==>", node.value)
	}
	loggers.Info("LIST:%s", msg)
}

///*==============================      ListIter       =================================*/
//
//const (
//	RedisListIteratorDirectionStartHead = 0
//	RedisListIteratorDirectionStartTail = 1
//)
//
//type ListIter struct {
//	next      *ListNode
//	direction int
//}
//
//func ListGetIterator(list *List, direction int) *ListIter {
//	if list == nil {
//		return nil
//	}
//	iter := &ListIter{direction: direction}
//	if direction == RedisListIteratorDirectionStartHead {
//		iter.next = list.head
//	} else if direction == RedisListIteratorDirectionStartTail {
//		iter.next = list.tail
//	}
//	return iter
//}
//
//func ListRewind(list *List, iter *ListIter) {
//	iter.direction = RedisListIteratorDirectionStartHead
//	iter.next = list.head
//}
//
//func ListRewindTail(list *List, iter *ListIter) {
//	iter.direction = RedisListIteratorDirectionStartTail
//	iter.next = list.tail
//}
//
//func (it *ListIter) ListNext() *ListNode {
//	cur := it.next
//	if cur != nil {
//		if it.direction == RedisListIteratorDirectionStartHead {
//			it.next = cur.next
//		} else {
//			it.next = cur.pre
//		}
//	}
//	return cur
//}
