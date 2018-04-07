
![framework](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/framework/framework.png?raw=true)

### log
log处理

### tcp
tcp层，和tcp通信打交道，负责数据的接收和发送。

### database

数据库，底层的数据库实现

定义五种redis基本类型

### raw_type

redis 五种基本类型的具体某个encoding的实现

string | list | hash | set | zset

### conf
配置，关于配置的一些常量。

### protocol
redis协议层，把数据包装成redis协议格式；把redis协议格式的数据解析成command

### client
redis客户端，向服务器发送请求

### handler
处理redis命令。不同的handler处理不同的redis命令

### server
redis server负责和客户端打交道，执行客户端命令，返回客户端需要的数据。

### mock
测试