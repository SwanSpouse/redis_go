package handlers

import (
	"fmt"
	"redis_go/client"
	"redis_go/database"
	re "redis_go/error"
	"redis_go/protocol"
	"strings"
)

// StringHandler可以处理的三种rawType
var stringEncodingTypeDict = map[string]bool{
	database.RedisEncodingInt:    true,
	database.RedisEncodingRaw:    true,
	database.RedisEncodingEmbStr: true,
}

type StringHandler struct{}

func (handler *StringHandler) Process(client *client.Client, command *protocol.Command) {
	switch strings.ToUpper(command.GetName()) {
	case "APPEND":
	case "BITCOUNT", "BITOP", "GETBIT", "SETBIT":
		client.ResponseWriter.AppendErrorf(re.ErrFunctionNotImplement)
	case "DECR":
	case "DECRBY":
	case "GET":
		handler.Get(client, command)
	case "GETRANGE":
	case "GETSET":
	case "INCR":
	case "INCRBY":
	case "INCRBYFLOAT":
	case "MGET":
	case "MSET":
	case "MSETNX":
	case "PSETEX":
	case "SET":
		handler.Set(client, command)
	case "SETEX":
	case "SETNX":
	case "SETRANGE":
	case "STRLEN":

	default:
		client.ResponseWriter.AppendErrorf("ERR unknown command %s", command.GetOriginName())
	}
	// 最后统一发送数据
	client.ResponseWriter.Flush()
}

func (handler *StringHandler) Set(client *client.Client, command *protocol.Command) {
	args := command.GetArgs()
	if len(args) < 2 {
		client.ResponseWriter.AppendErrorf(re.ErrWrongNumberOfArgs, command.GetOriginName())
		return
	}
	key := args[0]
	// TODO lmj 根据变量的值来判断创建什么样Encoding的StringObject
	if value, err := database.NewRedisStringObject(args[1]); err != nil || value == nil {
		client.ResponseWriter.AppendError(re.ErrUnknown.Error())
	} else {
		client.SelectedDatabase().SetKeyInDB(key, value)
		client.ResponseWriter.AppendOK()
	}
}

func (handler *StringHandler) Get(client *client.Client, command *protocol.Command) {
	args := command.GetArgs()
	// 判断参数个数是否合理
	if len(args) != 1 {
		client.ResponseWriter.AppendErrorf(re.ErrWrongNumberOfArgs, command.GetOriginName())
		return
	}
	// 获取args中的Key
	key := args[0]

	/* 获取key在数据库中对应的value(TBase:BaseType)
	 *      1. 处理database error
	 *      2. 验证BaseType的类型和编码方式
	 */
	baseType, err := client.SelectedDatabase().SearchKeyInDB(key)
	if err != nil {
		client.ResponseWriter.AppendErrorf("error type of %s", key)
		return
	}
	// 数据库中没有这个Key
	if baseType == nil {
		client.ResponseWriter.AppendNil()
		return
	}
	if _, ok := stringEncodingTypeDict[baseType.GetEncoding()]; !ok || baseType.GetObjectType() != database.RedisTypeString {
		client.ResponseWriter.AppendErrorf("error object type or encoding. type:%s, encoding:%s", baseType.GetObjectType(), baseType.GetEncoding())
		return
	}
	// 根据不同的encoding类型对数据进行处理
	switch baseType.GetEncoding() {
	case database.RedisEncodingInt, database.RedisEncodingEmbStr, database.RedisEncodingRaw:
		client.ResponseWriter.AppendInlineString(fmt.Sprintf("%s", baseType.GetValue()))
	}
}

func (handler *StringHandler) Incr(client *client.Client, command *protocol.Command) {
	args := command.GetArgs()
	if len(args) != 1 {
		client.ResponseWriter.AppendErrorf(re.ErrWrongNumberOfArgs, command.GetOriginName())
		return
	}
}
