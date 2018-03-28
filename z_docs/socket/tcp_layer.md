
##### redis go tcp layer

![server-client](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/socket/redis_go%20server_client.png?raw=true)

* Mock Client 使用RequestWriter通过WriterCmd把PING命令转换成
```shell
*1\r\n$4\r\nPING\r\n
```

* 把上述字符串通过TCP连接发送到Server端

* Server端通过RequestReader的ReadCmd读取上述数据，并封装成Command

* Redis Server 对Command进行解析，形成Response "PONG"

* 用ResponseWriter把PONG换成
```shell
*1\r\n$4\r\nPONG\r\n
```

* 把上述命令发送到Client端

* Client 解析Response