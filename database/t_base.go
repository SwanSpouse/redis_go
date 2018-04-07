package database

const (
	/**
	RedisTypeString  ->  RedisEncodingInt	 	: 使用整数值实现的字符串对象
	RedisTypeString  ->  RedisEncodingEmbStr        : 使用embstr编码的简单动态字符串实现的字符串对象
	RedisTypeString  ->  RedisEncodingRaw		: 使用简单动态字符串实现的字符串对象
	RedisTypeList    ->  RedisEncodingZipList	: 使用压缩列表实现的列表对象
	RedisTypeList    ->  RedisEncodingLinkedList	: 使用双端链表实现的列表对象
	RedisTypeHash    ->  RedisEncodingZipList	: 使用压缩链表实现的列表对象
	RedisTypeHash    ->  RedisEncodingHT		: 使用字典实现的哈希对象
	RedisTypeSet     ->  RedisEncodingIntSet	: 使用整数集合实现的集合对象
	RedisTypeSet     ->  RedisEncodingHT		: 使用字典实现的集合对象
	RedisTypeZSet    ->  RedisEncodingZipList	: 使用压缩链表实现的有序集合对象
	RedisTypeZSet    ->  RedisEncodingSkipList	: 使用跳跃表和字典实现的有序集合对象
	*/

	/* object type */
	RedisTypeString = "string"
	RedisTypeList   = "list"
	RedisTypeHash   = "hash"
	RedisTypeSet    = "set"
	RedisTypeZSet   = "zset"

	/* redis encoding type */
	RedisEncodingInt        = "int"
	RedisEncodingEmbStr     = "embstr"
	RedisEncodingRaw        = "raw"
	RedisEncodingHT         = "hashtable"
	RedisEncodingLinkedList = "linkedlist"
	RedisEncodingZipList    = "ziplist"
	RedisEncodingIntSet     = "intset"
	RedisEncodingSkipList   = "skiplist"
)

type TBase interface {
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
	GetValue() interface{}
	SetValue(interface{})
	IsExpired() bool
}
