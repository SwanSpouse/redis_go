
### 五种 unix io 模型

Unix系统将所有外部设备都看做是一个文件来操作，对一个文件的读写操作会调用内核提供的系统命令，返回一个file descriptor (fd 文件描述符)。

在Unix系统中，I/O操作可以分为以下两个不同阶段：

* 用户程序发起I/O请求。

* 内核执行I/O请求。

如果在上述操作的过程中:

* 用户进程一直在等待（没有处理其他事情）称之为同步过程；反之为异步。

* 用户程序调用I/O请求后，当资源不可用时，内核不能够立即返回，称之为阻塞。(由于内核的原因阻塞了用户程序)；反之为非阻塞

通过两阶段状态的两两组合，可以得到四种模型：


|               | Blocking      |         Non-Blocking  |
| ------------- |:-------------:| -------------:        |
| Synchronous   | read/write    | read/write(O_NONBLOK) |
| Asynchronous  | io multiplexing(select/poll)  |   AIO |


#### 同步阻塞

用户进程调用API(read, write)会转化成一个I/O请求，一直等到I/O请求完成API调用才会完成。这意味着：在API调用期间用户程序是同步的；这个API调用会导致系统以阻塞的模式执行I/O，如果此时没有数据
则一直等待（放弃CPU主动挂起--Sleep状态）

![BlockingIO](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/socket/BlockingIO.png?raw=true)

#### 同步非阻塞

这种模式通过调用read,write的时候指定O_NONBLOCK参数。和"同步阻塞"的区别在于系统调用的时候它是以非阻塞的方式执行的，无论是否有数据都会立即返回。

![Non-blockingIO](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/socket/Non-blockingIO.png?raw=true)


#### 异步阻塞

同步模型最主要的问题是占用CPU，阻塞I/O会主动让出CPU但是用户空间的系统调用还是不会返回依然耗费CPU；

如果仔细分析同步模型霸占CPU的原因不难得出结论——都是在等待数据的到来。异步模式正是意识到这一点所以把I/O读取细化为订阅I/O事件，
实际I/O读写，在"订阅I/O事件"事件部分会主动让出CPU直到事件发生。异步模式下的I/O函数和同步模式下的
I/O函数是一样的。（都是read、write）唯一的区别在是异步模式"读"必有数据而同步模式则未必。

常见的异步阻塞函数包括: select, poll, epoll。
* 其中以select为例： 异步模式下我们的API调用分为两步，第一步是通过select订阅读写事件，这个函数会主动让出CPU直到事件发生（设置为Sleep状态，等待事件发生）；select一旦返回就证明可以开始读了，
    所以第二步是通过read读取数据（"读"必有数据）

* select/poll是顺序扫描fd是否就绪，而且支持的fd数量有限，因此它的使用受到了一些制约。Linux还提供了一个epoll系统调用，epoll使用基于时间驱动方式代替顺序扫描，因此性能更高。
当有fd就绪时，立即回调函数callback。

#### 异步阻塞之信号驱动

完美主义者看了上面的select之后会有点不爽————我还要"等待"读写时间，能不能有读写事件的时候主动通知我呢？借助"信号"机制我们可以实现这个，但是这个并不完美而且有点儿弄巧成拙的意思。

具体用法：通过fcntl函数设置一个F_GETFL|O_ASYNC，当有I/O事件的时候，操作系统会触发SIGIO信号，在程序里只需要绑定SIGIO信号的处理函数就可以了。但是这里面有个问题——信号处理函数由哪个进程执行呢？
答案是"属主"进程。操作系统只负责参数信号而实际的信号处理函数必须由用户空间的进程实现。（这就是设置F_SETOWN为当前进程PID的原因）。

信号驱动性能要比select、poll高，但是缺点是致命的 —— Linux中信号队列是有限制的，如果超过这个限制就完全无法读取数据。

![IOMultiplexing](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/socket/IOMultiplexing.png?raw=true)


#### 异步非阻塞

![AsynchronousIO](https://github.com/SwanSpouse/redis_go/blob/master/z_docs/socket/AsynchronousIO.png?raw=true)


这种模型是最"省事"的模型，系统调用完成之后就只要坐等数据就可以了。其实不然，问题出在实际上，Linux上的AIO两个实现版本，POSIX的实现最烂，性能很差而且是基于"事件驱动"还会出现"信号队列不足"的问题（所以它就
偷偷的创建线程，导致线程也不可控了）；一个是Linux自己实现的（readhat贡献）Native I/O。 Native I/O主要涉及到的两个函数 io_submit 设置需要I/O动作（读、写、数据大小、应用缓冲区）； io_getevents 等待I/O
动作完成。即便你的整个I/O行为是非阻塞的还是需要有一个办法知道数据是否读取、写入成功。