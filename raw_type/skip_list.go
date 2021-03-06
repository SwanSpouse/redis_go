package raw_type

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	SkipListMaxLevel = 32
)

/*
	跳跃表层级定义
*/
type SkipLevel struct {
	forward *SkipNode // 前进指针
	span    int       // 这个层跨越节点的数量
}

/*
	跳跃表节点定义
*/
type SkipNode struct {
	obj      string      // 包含的对象
	score    float64     // 对象对应的分值
	backward *SkipNode   // 后退指针
	level    []SkipLevel // 跳跃表节点包含的层级
}

/*
	确定层数的方法：抛硬币，只要是正面就累加，直到遇见反面才停止，最后记录正面的次数并将其作为要添加新元素的层；
*/
func getNewSkipNodeLevel() int {
	levelCount := 1
	rand.Seed(time.Now().UnixNano())

	cur := rand.Intn(10)
	for ; cur < 5; cur = rand.Intn(10) {
		levelCount += 1
	}
	return levelCount
}

func (sn *SkipNode) GetValue() string {
	return sn.obj
}

func (sn *SkipNode) GetScore() float64 {
	return sn.score
}

func (sn *SkipNode) SetScore(score float64) {
	sn.score = score
}

func (sn *SkipNode) GetNextNode() *SkipNode {
	return sn.level[0].forward
}

func createSkipNode(obj string, score float64) *SkipNode {
	return &SkipNode{
		obj:   obj,
		score: score,
		level: make([]SkipLevel, getNewSkipNodeLevel()),
	}
}

func createHeaderSkipNode() *SkipNode {
	return &SkipNode{
		obj:   "H",
		level: make([]SkipLevel, SkipListMaxLevel),
	}
}

/*
	跳跃表定义
*/
type SkipList struct {
	header *SkipNode // 跳跃表头结点
	tail   *SkipNode // 跳跃表尾结点
	length int       // 跳跃表节点数量
	level  int       // 跳跃表内节点的最大层数
}

func NewSkipList() *SkipList {
	return &SkipList{
		header: createHeaderSkipNode(),
		tail:   nil,
		length: 0,
		level:  1,
	}
}

func (sl *SkipList) Insert(obj string, score float64) *SkipNode {
	updates := make([]*SkipNode, SkipListMaxLevel)
	rank := make([]int, SkipListMaxLevel)

	cur := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		if i != sl.level-1 {
			rank[i] = rank[i+1]
		}
		for cur.level[i].forward != nil &&
			(cur.level[i].forward.score < score || (cur.level[i].forward.score == score && strings.Compare(cur.level[i].forward.obj, obj) < 0)) {
			// 记录向前跨越了多少个元素
			rank[i] += cur.level[i].span
			cur = cur.level[i].forward
		}
		updates[i] = cur
	}
	// 创建新节点
	newNode := createSkipNode(obj, score)
	if len(newNode.level) > sl.level {
		for i := sl.level; i < len(newNode.level); i++ {
			rank[i] = 0
			updates[i] = sl.header
			updates[i].level[i].span = sl.length
		}
		sl.level = len(newNode.level)
	}
	// 更新level
	for i := 0; i < len(newNode.level); i++ {
		newNode.level[i].forward = updates[i].level[i].forward
		updates[i].level[i].forward = newNode

		// 前驱节点和新节点之间的距离是 rank[0]- rank[i]
		newNode.level[i].span = updates[i].level[i].span - (rank[0] - rank[i])
		updates[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	// 更新新节点上层的所有span
	for i := len(newNode.level); i < sl.level; i++ {
		updates[i].level[i].span += 1
	}

	// 设置后退指针
	if updates[0].backward != sl.header {
		newNode.backward = updates[0]
	}
	// 如果newNode不是尾指针，更新newNode下一个节点的backward指针；否则将newNode设置为尾节点
	if newNode.level[0].forward != nil {
		newNode.level[0].forward.backward = newNode
	} else {
		sl.tail = newNode
	}
	sl.length += 1
	return newNode
}

func (sl *SkipList) Length() int {
	return sl.length
}

// 这个函数要保证传进来的肯定是需要删除的
func (sl *SkipList) deleteNode(node *SkipNode, updates []*SkipNode) {
	// 更新forward指针
	for i := 0; i < sl.level; i++ {
		if updates[i].level[i].forward == node {
			// 有目标节点的层span =目标节点的span - 1（目标节点）
			updates[i].level[i].span += node.level[i].span - 1
			updates[i].level[i].forward = node.level[i].forward
		} else {
			// 没有目标节点的层直接span-1就好了
			updates[i].level[i].span -= 1
		}
	}
	if node.level[0].forward != nil {
		node.level[0].forward.backward = node.backward
	} else {
		sl.tail = node
	}
	// 收缩level
	for sl.level > 1 && sl.header.level[sl.level-1].forward == nil {
		sl.level -= 1
	}
	sl.length -= 1
}

func (sl *SkipList) Delete(obj string, score float64) int {
	updates := make([]*SkipNode, SkipListMaxLevel)

	cur := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for cur.level[i].forward != nil &&
			(cur.level[i].forward.score < score || (cur.level[i].forward.score == score && strings.Compare(cur.level[i].forward.obj, obj) < 0)) {
			cur = cur.level[i].forward
		}
		updates[i] = cur
	}

	cur = cur.level[0].forward
	// 只有当找到的时候才会进行删除
	if cur != nil && cur.score == score && strings.Compare(cur.obj, obj) == 0 {
		sl.deleteNode(cur, updates)
		return 1
	}
	return 0
}

type RangeSpec struct {
	Min, Max     float64
	MinEx, MaxEx bool
}

// 检查score是否比rangeSpec.min大
func valueGteMin(score float64, rng RangeSpec) bool {
	if rng.MinEx {
		return score > rng.Min
	}
	return score >= rng.Min
}

// 检查score是否比rangeSpec.max小
func valueLteMax(score float64, rng RangeSpec) bool {
	if rng.MaxEx {
		return score < rng.Max
	}
	return score <= rng.Max
}

// 检查SkipList中的元素是否在rangeSpec之中
func (sl *SkipList) IsInRange(rng RangeSpec) bool {
	if rng.Min > rng.Max || (rng.Min == rng.Max && (rng.MinEx || rng.MaxEx)) {
		return false
	}
	cur := sl.tail
	if cur == nil || !valueGteMin(cur.score, rng) {
		return false
	}

	cur = sl.header.level[0].forward
	if cur == nil || !valueLteMax(cur.score, rng) {
		return false
	}
	return true
}

// 返回在SkipList中第一个在给定区间中的元素
func (sl *SkipList) FirstInRange(rng RangeSpec) *SkipNode {
	if !sl.IsInRange(rng) {
		return nil
	}

	cur := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for cur.level[i].forward != nil && !valueGteMin(cur.level[i].forward.score, rng) {
			cur = cur.level[i].forward
		}
	}

	cur = cur.level[0].forward
	if !valueLteMax(cur.score, rng) {
		return nil
	}
	return cur
}

// 返回在SkipList中最后一个在给定区间中的元素
func (sl *SkipList) LastInRange(rng RangeSpec) *SkipNode {
	if !sl.IsInRange(rng) {
		return nil
	}
	// 自顶向下，查找符合条件的最后一个元素
	cur := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for cur.level[i].forward != nil && valueLteMax(cur.level[i].forward.score, rng) {
			cur = cur.level[i].forward
		}
	}

	if !valueGteMin(cur.score, rng) {
		return nil
	}
	return cur
}

func (sl *SkipList) DeleteRangeByScore(rng RangeSpec) int {
	updates := make([]*SkipNode, SkipListMaxLevel)

	cur := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for cur.level[i].forward != nil && ((rng.MinEx && cur.level[i].forward.score <= rng.Min) || (!rng.MinEx && cur.level[i].forward.score < rng.Min)) {
			cur = cur.level[i].forward
		}
		updates[i] = cur
	}
	/* Current node is the last with score < or <= Min. */
	cur = cur.level[0].forward

	removed := 0
	for cur != nil && ((rng.MaxEx && cur.score < rng.Max) || !rng.MaxEx && cur.score <= rng.Max) {
		// 保存后继节点，然后删除当前节点
		next := cur.level[0].forward
		sl.deleteNode(cur, updates)
		removed += 1
		cur = next
	}
	return removed
}

func (sl *SkipList) DeleteRangeByRank(start, end int) int {
	updates := make([]*SkipNode, SkipListMaxLevel)
	traversed := 0
	removed := 0

	cur := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for cur.level[i].forward != nil && (traversed+cur.level[i].span) < start {
			traversed += cur.level[i].span
			cur = cur.level[i].forward
		}
		updates[i] = cur
	}
	// 算上start节点
	traversed += 1

	cur = cur.level[0].forward
	for cur != nil && traversed <= end {
		nextNode := cur.level[0].forward
		sl.deleteNode(cur, updates)
		removed += 1
		traversed += 1
		cur = nextNode
	}
	return removed
}

// 获取obj的score在有序集中的排名。不存在返回0
func (sl *SkipList) GetRank(obj string, score float64) int {
	rank := 0
	cur := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for cur.level[i].forward != nil &&
			(cur.level[i].forward.score < score || (cur.level[i].forward.score == score && strings.Compare(cur.level[i].forward.obj, obj) <= 0)) {
			rank += cur.level[i].span
			cur = cur.level[i].forward
		}
	}

	if cur.score == score && strings.Compare(cur.obj, obj) == 0 {
		return rank
	}
	return 0
}

// 根据元素的排名来查找元素 第一个元素的排名是1
func (sl *SkipList) GetElementByRank(rank int) *SkipNode {
	traversed := 0
	cur := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for cur.level[i].forward != nil && cur.level[i].span+traversed <= rank {
			traversed += cur.level[i].span
			cur = cur.level[i].forward
		}
		if traversed == rank {
			return cur
		}
	}
	return nil
}

// debug skip list
func (sl *SkipList) DebugSkipListInfo() {
	fmt.Printf("current skip list info length:%d curLevel:%d\n", sl.length, sl.level)

	for curLevel := sl.level - 1; curLevel >= 0; curLevel -= 1 {
		for node := sl.header; node != nil; node = node.level[0].forward {
			if curLevel < len(node.level) {
				fmt.Printf("%d\t\t", node.level[curLevel].span)
			} else {
				fmt.Printf("E\t\t")
			}
		}
		fmt.Printf("\n")
	}

	for node := sl.header; node != nil; node = node.level[0].forward {
		fmt.Printf("%s\t\t", node.obj)
	}
	fmt.Println()

	for node := sl.header; node != nil; node = node.level[0].forward {
		fmt.Printf("%.1f\t\t", node.score)
	}
	fmt.Println()

	for node := sl.header; node != nil; node = node.level[0].forward {
		fmt.Printf("%d\t\t", len(node.level))
	}
	fmt.Println()

}

func printSkipNodeLevel(node *SkipNode, level int) string {
	ret := fmt.Sprintf("*->\t")
	for i := 0; i < node.level[level].span; i++ {
		ret += fmt.Sprintf("?->\t")
	}
	return ret
}
