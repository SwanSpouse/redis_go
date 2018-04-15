
#### 一次完整的客户端与服务器连接事件示例

假设一个Redis服务器正在运作，那么这个服务器的监听套接字的AE_READABLE事件应该正处于监听状态之下，而该事件所对应的处理器为
连接应答处理器。

如果这时有一个Redis客户端向服务器发起连接，那么监听套接字将产生AE_READABLE事件，出发连接应答处理器执行。处理器会对客户端的
连接请求进行应答，然后创建客户端套接字，以及客户端状态，并将客户端套接字的AE_READABLE事件与命令请求处理器进行关联，使得客户端
可以向主服务器发送命令请求。

之后，假设客户端向主服务器发送了一个命令请求，那么客户端套接字将产生AE_READABLE事件，已发命令处理器执行，处理器读取客户端的命令
的内容，然后传给相关程序去执行。

执行命令将产生相应的命令回复，为了将这些命令回复传送回给客户端，服务器会将客户端套接字的AE_WRITABLE事件与命令回复处理器进行
关联。当客户端尝试读取命令回复的时候，客户端套接字将产生AE_WRITABLE事件，触发命令回复处理器执行，当命令回复处理器将命令回复
全部写入到套接字之后，服务器就会接触客户端套接字的AE_WRITABLE事件与命令回复处理器之间的关联。



#### reference

* 《Redis设计与实现》黄建宏