package database

type TList interface {
	LPush()
	RPush()
	LPop()
	RPop()
	LIndex()
	LLen()
	LInsert()
	LRem()
	LTrim()
	LSet()
}
