package database

import (
	"redis_go/client"
	"redis_go/protocol"
)

type TString interface {
	GetObjectType() string
	SetObjectType(string)
	GetEncoding() int
	SetEncoding(int)
	GetLRU() int
	SetLRU(int)
	GetRefCount() int
	IncrRefCount() int
	DecrRefCount() int
	GetTTL() int
	SetTTL(int)
	GetValue() string
	SetValue(interface{})
	Set(client *client.Client, command *protocol.Command) (int, error)
	Get(client *client.Client, command *protocol.Command) (string, error)
	Append()
	IncrByFloat()
	IncrBy()
	DecrBy()
	StrLen()
	SetRange()
	GetRange()
}
