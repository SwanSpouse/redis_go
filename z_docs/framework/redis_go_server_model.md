
### Redis Go 新老server-client模型对比

#### 原有redis go server 模型:

![OldRedisGoModel](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/framework/old_redis_go_model.png?raw=true)

老模型描述:

1. 创建redis server, 同时监听9736端口。
2. Client向Server发起请求，请求建立TCP连接。
3. Server接收到Client的请求，同时建立TCP连接。
4. Server创建一个go routine, 将这个Client交给go routine来进行处理。
5. go routine轮询Client, 查看Client是否有未读数据，如果有未读数据，那么就进行处理。
6. go routine将处理好的结果发送给Client。继续轮询等待Client发送过来的命令。


老模型所存在的问题:

1. 在没有心跳的情况下，服务器端无法感知Client是否已经断开连接。
2. 每一个go routine的大部分时间都是在轮询Client的缓冲区查看是否有未处理的数据。这导致go routine的效率非常低。
3. go routine缺少回收机制，随着Client端的个数不断增加，go routine会越来越多


![NewRedisGoModel](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/framework/new_redis_go_model.png?raw=true)

新模型描述:

1. 创建redis server, 同时监听9736端口。
2. Client向Server发起请求，请求建立TCP连接。
3. Server接收到Client的请求，同时建立TCP连接。
4. 建立连接后，Server将Client放入Client队列（缓冲池）之中。
5. 有一个Scanner线程轮询Client队列中的每一个Client，处理处理超时和空闲超时的Client；如果Client缓冲区中有处理的请求，
则从task线程池中获取一个worker(go routine)来处理Client的这一个请求。
6. go routine处理完Client的命令之后，将数据发送给Client，将Client放回Client队列。


新模型能够解决的问题:

1. 引入了Scanner之后，可以对处理时间超时和空闲时间超时的Client进行处理，释放资源。
2. 将原来的Client-GoRoutine一一对应变成了Request-GoRoutine一一对应，能够提高GoRoutine的工作效率。
3. 引入了Client、GoRoutine缓存池，能够提高系统资源的利用率。



#### ps

虽然golang创建go routine的代价非常小，但是还要从节约内存资源的角度来思考问题。

如果说参考思想，这个的思想来源应该是I/O多路复用模型吧。
