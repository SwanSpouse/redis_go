
![framework](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/framework/framework.png?raw=true)

### client
redis客户端，向服务器发送请求

### conf
配置，关于配置的一些常量。

### database
数据库，底层的数据库实现

同时定义五种redis基本类型(string,list,hash,set,zset)

### encodings

redis 五种基本类型的具体某个encoding的实现

* string类型:
    * string_int.go 使用整数值实现的字符串对象
    * string_raw.go 使用golang string实现的字符串对象
    * string_sds.go 使用simple dynamic string实现的字符串对象

* list类型:

* hash类型:

* set类型:

* zset类型:

### err
常见err的定义


### handler

处理redis命令。不同的handler处理不同的redis命令

* connection_handler.go 处理Connection相关的命令
* key_handler.go 处理key相关的命令
* string_handler.go 处理字符串相关的命令
* list_handler.go 处理列表相关的命令
* hash_handler.go 处理哈希表相关的命令
* set_handler.go 处理集合相关的命令
* sorted_set_handler.go 处理有序集合相关的命令
* pub_sub_handler.go 处理发布、订阅相关的命令
* transaction_handler.go 处理事务相关的命令

### log
log处理

### mock
用于进行mock测试

### protocol
redis协议层，把数据包装成redis协议格式；把redis协议格式的数据解析成command

### raw_type
各种encodings的低层实现。

* dict.go hash表的低层实现
* simple_dynamic_string.go 字符串对象的低层实现

### server
redis server负责和客户端打交道，执行客户端命令，返回客户端需要的数据。

### tcp
tcp层，和tcp通信打交道，负责数据的接收和发送。

### z_docs
相关文档和说明
