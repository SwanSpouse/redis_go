
#### unix I/O模型

Linux的内核将所有外部设备都看做一个文件来操作，对一个文件的读写操作会调用内核提供的系统命令，返回一个file descriptor (fd 文件描述符)。
而对一个socket的读写也会有相应的描述符，称为socketfd(socket描述符)，描述符就是一个数字，它指向内核中的一个结构（文件路径，数据区等一些属性）。

根据Unix网络编程对I/O模型的分类，Unix提供了5中I/O模型，分别如下。

![BlockingIO](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/socket/BlockingIO.png?raw=true)
* 阻塞I/O模型：最常用的I/O模型就是阻塞I/O模型，缺省情况下，所有文件操作都是阻塞的。我们以套接字接口为例来讲解此模型：在进程空间中调用recvfrom，气系统调用
直接到数据包到达且被复制到应用进程的缓冲区或者发生错误时才返回，在此期间一直会等待。进程在从调用recvfrom开始到它返回的整段时间内都是被阻塞的，因此被称为阻塞I/O模型。


![Non-blockingIO](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/socket/Non-blockingIO.png?raw=true)

* 非阻塞I/O模型：recvfrom从应用层到内核的时候，如果该缓冲区没有数据的话，就直接返回一个EWOULDBLOCK错误，一般都对非阻塞I/O模型进行轮询检查这个状态，看内核是不是有数据到来。

![IOMultiplexingIO](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/socket/IOMultiplexingIO.png?raw=true)

* I/O复用模型：Linux提供select/poll，进程通过将一个或者多个fd传递给select或poll系统调用，阻塞在select操作上，这样select/poll可以帮我们侦测多个fd是否处于就绪状态。
select/poll是顺序扫描fd是否就绪，而且支持的fd数量有限，因此它的使用受到了一些制约。Linux还提供了一个epoll系统调用，epoll使用基于时间驱动方式代替顺序扫描，因此性能更高。
当有fd就绪时，立即回调函数rollback。

![Signal-drivenIO](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/socket/Signal-drivenIO.png?raw=true)

* 信号驱动I/O模型：首先开启套接口信号驱动I/O功能，并通过系统调用sigaction执行一个信号处理函数（此系统调用立即返回，进程继续工作，它是非阻塞的）。当数据准备就绪时，就为该进程
生成一个SIGIO信号，通过信号回调通知应用程序调用recvfrom来读取数据，并通知主循环函数处理数据。

![AsynchronousIO](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/socket/AsynchronousIO.png?raw=true)

* 异步I/O：告知内核启动某个操作，并让内核在整个操作完成后（包括将数据从内核复制到用户自己的数据缓冲区）通知我们。这种模型与信号驱动模型的主要区别是：信号驱动I/O由内核通知我们何时可以
开始一个I/O操作；异步I/O模型由内核通知我们I/O操作何时完成。


#### I/O多路复用技术

目前支持I/O多路复用的系统调用有: select、pselect、poll、epoll，在Linux网络编程过程中，很长一段时间都是用select做轮询和网络时间通知。然而select的一些固有缺陷导致了它的应用受到
了很大的限制，最终Linux不得不在新的内核版本中寻找select的替代方案，最终选择了epoll。epoll和select的原理比较类似，为了克服select的缺点，epoll做了很多重大改进，总结如下：

* 支持一个进程打开的socket描述符(fd)不受限制(受限于操作系统的最大文件句柄数)
* I/O 效率不会随着FD数目的增加而线性下降。
* 使用mmap加速内核与用户空间的消息传递。
* epoll的API更加简单。



#### reference

* 《Netty权威指南》




