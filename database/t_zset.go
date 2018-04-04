package database

type TZSet interface {
	ZAdd()
	ZCard()
	ZCount()
	ZRange()
	ZRevRange()
	ZRank()
	ZRevRank()
	ZRem()
	ZScore()
}
