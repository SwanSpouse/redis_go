package handlers

import (
	"fmt"
	"redis_go/client"
	"redis_go/database"
	"redis_go/encodings"
	re "redis_go/error"
)

// StringHandler可以处理的三种rawType
var stringEncodingTypeDict = map[string]bool{
	encodings.RedisEncodingInt:    true,
	encodings.RedisEncodingRaw:    true,
	encodings.RedisEncodingEmbStr: true,
}

type StringHandler struct{}

func (handler *StringHandler) Process(client *client.Client) {
	if client.Cmd == nil {
		client.ResponseError("ERR nil command")
		return
	}
	switch client.Cmd.GetName() {
	case "APPEND":
		handler.Append(client)
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

func (handler *StringHandler) Append(client *client.Client) {
	args := client.Cmd.GetArgs()
	if len(args) < 2 {
		client.ResponseError(re.ErrWrongNumberOfArgs, client.Cmd.GetOriginName())
		return
	}
	// 获取args中的Key
	key := args[0]
	// 获取key在数据库中对应的value(TBase:BaseType)
	baseType := client.SelectedDatabase().SearchKeyInDB(key)
	var sb database.TString
	var ok bool
	if sb, ok = handler.convertTBastToStringObject(client, baseType); !ok || sb == nil {
		client.ResponseError(re.ErrWrongType)
		return
	}
	client.Response(sb.Append(args[1]))
}

func (handler *StringHandler) Set(client *client.Client) {
	args := client.Cmd.GetArgs()
	if len(args) < 2 {
		client.ResponseError(re.ErrWrongNumberOfArgs, client.Cmd.GetOriginName())
		return
	}
	client.SelectedDatabase().SetKeyInDB(args[0], database.NewRedisStringObject(args[1]))
	client.ResponseOK()
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
	var sb database.TString
	var ok bool
	if sb, ok = handler.convertTBastToStringObject(client, baseType); !ok || sb == nil {
		client.ResponseError(re.ErrWrongType)
		return
	}
	client.Response(fmt.Sprintf("%s", sb.GetValue()))
}

func (handler *StringHandler) Incr(client *client.Client) {
	args := client.Cmd.GetArgs()
	if len(args) != 1 {
		client.ResponseError(re.ErrWrongNumberOfArgs, client.Cmd.GetOriginName())
		return
	}
	key := args[0]
	baseType := client.SelectedDatabase().SearchKeyInDB(key)
	var sb database.TString
	var ok bool
	if sb, ok = handler.convertTBastToStringObject(client, baseType); !ok || sb == nil {
		client.ResponseError(re.ErrWrongType)
		return
	}
	if ret, err := sb.Incr(); err != nil {
		client.ResponseError(err.Error())
	} else {
		client.Response(ret)
	}
}

func (handler *StringHandler) Decr(client *client.Client) {
	args := client.Cmd.GetArgs()
	if len(args) != 1 {
		client.ResponseError(re.ErrWrongNumberOfArgs, client.Cmd.GetOriginName())
		return
	}
	key := args[0]
	baseType := client.SelectedDatabase().SearchKeyInDB(key)
	var sb database.TString
	var ok bool
	if sb, ok = handler.convertTBastToStringObject(client, baseType); !ok || sb == nil {
		client.ResponseError(re.ErrWrongType)
		return
	}
	if ret, err := sb.Decr(); err != nil {
		client.ResponseError(err.Error())
	} else {
		client.Response(ret)
	}
}

/**
	对baseType的类型是否为string进行校验。
首先判断baseType是否为空，再判断baseType的Encoding是否为string的Encoding, baseType的Type是否为RedisTypeString
*/
func (handler *StringHandler) convertTBastToStringObject(client *client.Client, baseType database.TBase) (database.TString, bool) {
	if baseType == nil {
		client.Response(nil)
		return nil, false
	}
	if _, ok := stringEncodingTypeDict[baseType.GetEncoding()]; !ok || baseType.GetObjectType() != encodings.RedisTypeString {
		client.ResponseError("error object type or encoding. type:%s, encoding:%s", baseType.GetObjectType(), baseType.GetEncoding())
		return nil, false
	}
	if stringObject, ok := baseType.(database.TString); !ok {
		client.ResponseError("error object type or encoding. type:%s, encoding:%s", baseType.GetObjectType(), baseType.GetEncoding())
		return stringObject, true
	}
	return nil, false
}
