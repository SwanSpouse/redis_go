
#### 一次完整的客户端与服务器连接事件示例

假设一个Redis服务器正在运作，那么这个服务器的监听套接字的AE_READABLE事件应该正处于监听状态之下，而该事件所对应的处理器为
连接应答处理器。

```c

# redis.c 将acceptTcpHandler函数注册为AE_READABLE事件的应答器
if (server.ipfd > 0 && aeCreateFileEvent(server.el,server.ipfd,AE_READABLE, acceptTcpHandler, NULL) == AE_ERR)
    redisPanic("Unrecoverable error creating server.ipfd file event.");

```

如果这时有一个Redis客户端向服务器发起连接，那么监听套接字将产生AE_READABLE事件，出发连接应答处理器执行。处理器会对客户端的
连接请求进行应答，然后创建客户端套接字，以及客户端状态，并将客户端套接字的AE_READABLE事件与命令请求处理器进行关联，使得客户端
可以向主服务器发送命令请求。

```c

# networking.c client建立和服务端的联系之后，会把套接字的AE_READABLE事件与请求处理器进行关联，这样发送过来的命令就可以被readQueryFromClient获取到，
# 通过processInputBuffer(networking.c)来执行客户端发来的命令。

if (aeCreateFileEvent(server.el,fd,AE_READABLE, readQueryFromClient, c) == AE_ERR) {
    close(fd);
    zfree(c);
    return NULL;
}
```

之后，假设客户端向主服务器发送了一个命令请求，那么客户端套接字将产生AE_READABLE事件，已发命令处理器执行，处理器读取客户端的命令
的内容，然后传给相关程序去执行。

```c

aeCreateFileEvent(server.el, c->fd, AE_WRITABLE, sendReplyToClient, c) == AE_ERR)
```

执行命令将产生相应的命令回复，为了将这些命令回复传送回给客户端，服务器会将客户端套接字的AE_WRITABLE事件与命令回复处理器进行关联。
当客户端尝试读取命令回复的时候，客户端套接字将产生AE_WRITABLE事件，触发命令回复处理器执行，当命令回复处理器将命令回复
全部写入到套接字之后，服务器就会接解除户端套接字的AE_WRITABLE事件与命令回复处理器之间的关联。

```c
// networking.c 回复全部发送完毕，删除事件处理器
if (c->bufpos == 0 && listLength(c->reply) == 0) {
    c->sentlen = 0;
    aeDeleteFileEvent(server.el,c->fd,AE_WRITABLE);

    /* Close connection after entire reply has been sent. */
    // 如果状态为“回复完毕之后关闭”，那么关闭客户端
    if (c->flags & REDIS_CLOSE_AFTER_REPLY) freeClient(c);
}
```

#### 流程图

![server-client-communication](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/redis/server_client_communication.png?raw=true)


#### 思考

因为goroutine创建的成本比较低，是不是可以不用I/O多路复用的思想来解决这个问题。直接创建goroutine，同步阻塞的方式来处理就好了？

#### reference

* 《Redis设计与实现》黄建宏