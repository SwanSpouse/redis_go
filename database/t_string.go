package database

type TString interface {
	GetObjectType() string
	SetObjectType(string)
	GetEncoding() string
	SetEncoding(string)
	GetLRU() int
	SetLRU(int)
	GetRefCount() int
	IncrRefCount() int
	DecrRefCount() int
	GetTTL() int
	SetTTL(int)
	GetValue() string
	SetValue(interface{})
	Set() (int, error)
	Get() (string, error)
	Append()
	IncrByFloat()
	IncrBy()
	DecrBy()
	StrLen()
	SetRange()
	GetRange()
}
