package handlers

import (
	"fmt"
	"redis_go/client"
	"redis_go/database"
	re "redis_go/error"
	"strconv"
)

// StringHandler可以处理的三种rawType
var stringEncodingTypeDict = map[string]bool{
	database.RedisEncodingInt:    true,
	database.RedisEncodingRaw:    true,
	database.RedisEncodingEmbStr: true,
}

type StringHandler struct{}

func (handler *StringHandler) Process(client *client.Client) {
	if client.Cmd == nil {
		client.ResponseError("ERR nil command")
		return
	}
	switch client.Cmd.GetName() {
	case "APPEND":
	case "BITCOUNT", "BITOP", "GETBIT", "SETBIT":
		client.ResponseError(re.ErrFunctionNotImplement)
	case "DECR":
		handler.Decr(client)
	case "DECRBY":
	case "GET":
		handler.Get(client)
	case "GETRANGE":
	case "GETSET":
	case "INCR":
		handler.Incr(client)
	case "INCRBY":
	case "INCRBYFLOAT":
	case "MGET":
	case "MSET":
	case "MSETNX":
	case "PSETEX":
	case "SET":
		handler.Set(client)
	case "SETEX":
	case "SETNX":
	case "SETRANGE":
	case "STRLEN":

	default:
		client.ResponseError("ERR unknown command %s", client.Cmd.GetOriginName())
	}
	// 最后统一发送数据
	client.Flush()
}

func (handler *StringHandler) Set(client *client.Client) {
	args := client.Cmd.GetArgs()
	if len(args) < 2 {
		client.ResponseError(re.ErrWrongNumberOfArgs, client.Cmd.GetOriginName())
		return
	}
	key := args[0]
	// TODO lmj 根据变量的值来判断创建什么样Encoding的StringObject
	if value, err := database.NewRedisStringObject(args[1]); err != nil || value == nil {
		client.ResponseError(re.ErrUnknown.Error())
	} else {
		client.SelectedDatabase().SetKeyInDB(key, value)
		client.ResponseOK()
	}
}

func (handler *StringHandler) Get(client *client.Client) {
	args := client.Cmd.GetArgs()
	// 判断参数个数是否合理
	if len(args) != 1 {
		client.ResponseError(re.ErrWrongNumberOfArgs, client.Cmd.GetOriginName())
		return
	}
	// 获取args中的Key
	key := args[0]
	// 获取key在数据库中对应的value(TBase:BaseType)
	baseType := client.SelectedDatabase().SearchKeyInDB(key)
	if !handler.isStringObjectAndEncodingValid(client, baseType) {
		return
	}
	// 根据不同的encoding类型对数据进行处理
	switch baseType.GetEncoding() {
	case database.RedisEncodingInt, database.RedisEncodingEmbStr, database.RedisEncodingRaw:
		client.Response(fmt.Sprintf("%s", baseType.GetValue()))
	}
}

func (handler *StringHandler) Incr(client *client.Client) {
	args := client.Cmd.GetArgs()
	if len(args) != 1 {
		client.ResponseError(re.ErrWrongNumberOfArgs, client.Cmd.GetOriginName())
		return
	}
	key := args[0]
	baseType := client.SelectedDatabase().SearchKeyInDB(key)
	if !handler.isStringObjectAndEncodingValid(client, baseType) {
		return
	}
	switch baseType.GetEncoding() {
	case database.RedisEncodingEmbStr, database.RedisEncodingRaw:
		value := baseType.GetValue().(string)
		if valueInt, err := strconv.ParseInt(value, 10, 64); err != nil {
			client.ResponseError(re.ErrNotIntegerOrOutOfRange)
			return
		} else {
			baseType.SetValue(strconv.FormatInt(valueInt+1, 10))
		}
	case database.RedisEncodingInt:
		value := baseType.GetValue().(int64)
		baseType.SetValue(value + 1)
	}
	client.ResponseOK()
}

func (handler *StringHandler) Decr(client *client.Client) {
	args := client.Cmd.GetArgs()
	if len(args) != 1 {
		client.ResponseError(re.ErrWrongNumberOfArgs, client.Cmd.GetOriginName())
		return
	}
	key := args[0]
	baseType := client.SelectedDatabase().SearchKeyInDB(key)
	if !handler.isStringObjectAndEncodingValid(client, baseType) {
		return
	}
	switch baseType.GetEncoding() {
	case database.RedisEncodingEmbStr, database.RedisEncodingRaw:
		value := baseType.GetValue().(string)
		if valueInt, err := strconv.ParseInt(value, 10, 64); err != nil {
			client.ResponseError(re.ErrNotIntegerOrOutOfRange)
			return
		} else {
			baseType.SetValue(strconv.FormatInt(valueInt-1, 10))
		}
	case database.RedisEncodingInt:
		value := baseType.GetValue().(int64)
		baseType.SetValue(value - 1)
	}
	client.ResponseOK()
}

/**
	对baseType的类型是否为string进行校验。
首先判断baseType是否为空，再判断baseType的Encoding是否为string的Encoding, baseType的Type是否为RedisTypeString
*/
func (handler *StringHandler) isStringObjectAndEncodingValid(client *client.Client, baseType database.TBase) bool {
	if baseType == nil {
		client.Response(nil)
		return false
	}
	if _, ok := stringEncodingTypeDict[baseType.GetEncoding()]; !ok || baseType.GetObjectType() != database.RedisTypeString {
		client.ResponseError("error object type or encoding. type:%s, encoding:%s", baseType.GetObjectType(), baseType.GetEncoding())
		return false
	}
	return true
}
