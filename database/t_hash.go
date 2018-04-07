package database

type THash interface {
	HSet()
	HGet()
	HExists()
	HDel()
	HLen()
	HGetAll()
}