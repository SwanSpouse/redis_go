
#### unix I/O模型

Linux的内核将所有外部设备都看做一个文件来操作，对一个文件的读写操作会调用内核提供的系统命令，返回一个file descriptor (fd 文件描述符)。
而对一个socket的读写也会有相应的描述符，称为socketfd(socket描述符)，描述符就是一个数字，它指向内核中的一个结构（文件路径，数据区等一些属性）。

根据Unix网络编程对I/O模型的分类，Unix提供了5中I/O模型，分别如下。

![BlockingIO](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/socket/BlockingIO.png?raw=true)
* 阻塞I/O模型：最常用的I/O模型就是阻塞I/O模型，缺省情况下，所有文件操作都是阻塞的。我们以套接字接口为例来讲解此模型：在进程空间中调用recvfrom，气系统调用
直接到数据包到达且被复制到应用进程的缓冲区或者发生错误时才返回，在此期间一直会等待。进程在从调用recvfrom开始到它返回的整段时间内都是被阻塞的，因此被称为阻塞I/O模型。


![Non-blockingIO](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/socket/Non-blockingIO.png?raw=true)

* 非阻塞I/O模型：recvfrom从应用层到内核的时候，如果该缓冲区没有数据的话，就直接返回一个EWOULDBLOCK错误，一般都对非阻塞I/O模型进行轮询检查这个状态，看内核是不是有数据到来。

![IOMultiplexing](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/socket/IOMultiplexing.png?raw=true)

* I/O复用模型：Linux提供select/poll，进程通过将一个或者多个fd传递给select或poll系统调用，阻塞在select操作上，这样select/poll可以帮我们侦测多个fd是否处于就绪状态。
select/poll是顺序扫描fd是否就绪，而且支持的fd数量有限，因此它的使用受到了一些制约。Linux还提供了一个epoll系统调用，epoll使用基于时间驱动方式代替顺序扫描，因此性能更高。
当有fd就绪时，立即回调函数callback。

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


#### I/O模型中的同步/异步，阻塞/非阻塞

|               | Blocking      |         Non-Blocking  |
| ------------- |:-------------:| -------------:        |
| Synchronous   | read/write    | read/write(O_NONBLOK) |
| Asynchronous  | io multiplexing(select/poll)  |   AIO |


从内核角度看I/O操作分为两步：用户层API调用；内核层完成系统调用（发起I/O请求）。所以"异步/同步"是指API调用，"阻塞/非阻塞"是指内核完成I/O调用的模式。

同步是指函数完成之前会一直等待；阻塞是指系统调用的时候进程会被设置为Sleep状态直到等待的时间发生（比如有新的数据）。

同步阻塞：

* 用户空间调用API(read, write)会转化成一个I/O请求，一直等到I/O请求完成API调用才会完成。这意味着：在API调用期间用户程序是同步的；这个API调用会导致系统以阻塞的模式执行I/O，如果此时没有数据
则一直等待（放弃CPU主动挂起--Sleep状态）

同步非阻塞：

* 这种模式通过调用read,write的时候指定O_NONBLOCK参数。和"同步阻塞"的区别在于系统调用的时候它是以非阻塞的方式执行的，无论是否有数据都会立即返回。

异步阻塞：

* 同步模型最主要的问题是占用CPU，阻塞I/O会主动让出CPU但是用户空间的系统调用还是不会返回依然耗费CPU。如果仔细分析同步模型霸占
CPU的原因不难得出结论——都是在等待数据的到来。异步模式正是意识到这一点所以把I/O读取细化为订阅I/O事件，实际I/O读写，在"订阅I/O事件"事件部分会主动让出CPU直到事件发生。异步模式下的I/O函数和同步模式下的
I/O函数是一样的。（都是read、write）唯一的区别在是异步模式"读"必有数据而同步模式则未必。

* 常见的异步阻塞函数包括: select, poll, epoll。
    * 其中以select为例： 异步模式下我们的API调用分为两步，第一步是通过select订阅读写事件，这个函数会主动让出CPU直到事件发生（设置为Sleep状态，等待事件发生）；select一旦返回就证明可以开始读了，
    所以第二步是通过read读取数据（"读"必有数据）

异步阻塞之信号驱动

* 完美主义者看了上面的select之后会有点不爽————我还要"等待"读写时间，能不能有读写事件的时候主动通知我呢？借助"信号"机制我们可以实现这个，但是这个并不完美而且有点儿弄巧成拙的意思。

* 具体用法：通过fcntl函数设置一个F_GETFL|O_ASYNC，当有I/O事件的时候，操作系统会触发SIGIO信号，在程序里只需要绑定SIGIO信号的处理函数就可以了。但是这里面有个问题——信号处理函数由哪个进程执行呢？
答案是"属主"进程。操作系统只负责参数信号而实际的信号处理函数必须由用户空间的进程实现。（这就是设置F_SETOWN为当前进程PID的原因）。

* 信号驱动性能要比select、poll高，但是缺点是致命的 —— Linux中信号队列是有限制的，如果超过这个限制就完全无法读取数据。


异步非阻塞

* 这种模型是最"省事"的模型，系统调用完成之后就只要坐等数据就可以了。其实不然，问题出在实际上，Linux上的AIO两个实现版本，POSIX的实现最烂，性能很差而且是基于"事件驱动"还会出现"信号队列不足"的问题（所以它就
偷偷的创建线程，导致线程也不可控了）；一个是Linux自己实现的（readhat贡献）Native I/O。 Native I/O主要涉及到的两个函数 io_submit 设置需要I/O动作（读、写、数据大小、应用缓冲区）； io_getevents 等待I/O
动作完成。即便你的整个I/O行为是非阻塞的还是需要有一个办法知道数据是否读取、写入成功。

#### reference

* 《Netty权威指南》
* [透彻Linux(Unix)五种IO模型](https://mp.weixin.qq.com/s/vLTySPmujbAnR6wVbY9YnA)




