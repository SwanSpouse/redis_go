
![framework](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/framework/framework.png?raw=true)


### 当前设计存在的问题：

server和handler的这个设计不好，没有完全解耦。

server_handler里面要用到好多server的信息，所以现在已经把server_handler放到了server的pkg中。

handler command client 之间的矛盾：

    * 本想着能够通过设计把handler command client很好的拆分开。但是写着写着发现又合回去了。

    * command由自己对应的handler来进行处理。handler是用来处理client的，然后client又拥有command队列。这就循环引用了。

    * 所以抽象出来了一个base_handler来分离handler和client & command之间的耦合。但是client和command好像很难拆分开了。

server database client 之间的矛盾：

    * database是单独出来的。client中应该包含它选择的database，不能只包含index，不然每次获取数据都需要server的参与。

    * 现在只有在创建client的时候，server会把client选择的db告诉它。剩下的都交给handler来处理了。所以client就失去了换database的机会。这很尴尬。

server handler client 之间的矛盾：

    * 接受到客户端的命令之后。不好把命令直接交给cmmand对应的handler来进行处理。这样当client需要找server的时候，就联系不到了。

最终的解决办法可能就是把handler和server进行合并。但是这又是我不想看到的。看来设计模式不过关啊，设计设计着就耦合到一起了。

### client
client.go: redis客户端，向服务器发送请求
command: redis command，定义command结构，其中包括command名称、参数个数、处理当前command的方法等。

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
    * list_linked_list.go 使用双端链表实现的列表对象

* hash类型:
    * dict.go 使用字典实现的哈希对象

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

### loggers
log相关，这里做了几个针对性的处理

* 把所有的\r\n都打印出来，不显示成回车符。
* 自定义info, debug, error的不同级别log的颜色。
* 打印指定代码的go routineID，方便定位问题。

### mock
用于进行mock测试

* 模拟redis client向server发送命令，并验证得到的回复是否为期望的回复。

### protocol
redis协议层，把数据包装成redis协议格式；把redis协议格式的数据解析成command

### raw_type
各种encodings的低层实现。

* dict.go hash的低层实现
* simple_dynamic_string.go 字符串对象的低层实现
* list.go list的底层实现

### server
redis server负责和客户端打交道，执行客户端命令，返回客户端需要的数据。

### tcp
tcp层，和tcp通信打交道，负责数据的接收和发送。

### z_docs
相关文档和说明
