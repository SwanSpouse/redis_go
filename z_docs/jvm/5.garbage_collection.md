## JVM 垃圾收集器

#### 概述

上图展示了7种作用不同分代的垃圾收集器。如果两个收集器之间存在连线，就说明他们可以搭配使用。虚拟机所处的区域，则表示它是属于新生代收集器还是老年代收集器。

没有最好的垃圾收集器，只有对具体应用最适合的收集器。

#### HotSpot 虚拟机架构图

Eden : S0 : S1  = 8 : 1 : 1

垃圾分代回收：
* 新生代 ： 复制
* 老年代 ： 标记-清除/ 整理

#### Serial收集器

* 单线程收集器，在收集器必须暂停其他所有的工作线程，直到它收集结束。-- "Stop the World"
* Serial收集器由于没有线程交互的开销，专心做垃圾收集自然可以获得最高的单线程收集效率。

#### Serial Old收集器

Serial Old收集器是Serial收集器的老年代版本，是**一个单线程收集器，使用“标记--整理”算法。**

#### ParNew 收集器

ParNew收集器其实就是Serial收集器的**多线程**版本。

#### Parallel Scavenge 收集器 ( 吞吐量有限收集器）

Parallel Scavenge 收集器是一个新生代收集器，使用复制算法。

Parallel Scavenge收集器的目标是达到一个可控制的吞吐量（Throughput）

吞吐量 = 运行用户代码的时间 / ( 运行用户代码的时间 + 垃圾收集的时间）

* -XX:MaxGCPauseMillis 最大垃圾收集停顿的时间
* -XX:MaxGCTimeRatio 吞吐量大小

#### Parallel Old 收集器

Parallel Old 收集器是 Parallel Scavenge收集器的老年代版本，使用**多线程和“标记--整理”算法。**

#### CMS收集器(Concurrent Mark Sweep)

CMS收集器是一种以获取最短回收停顿时间为目标的收集器。

基于“标记--清除“ 算法实现的。整个过程分为4个步骤：
* 初始标记 (CMS inital mark)  
* 并发标记 (CMS concurrent mark)
* 重新标记 (CMS remark)
* 并发清除 (CMS concurrent sweep)

初始标记、重新标记需要Stop The World。初始标记仅仅只是标记一下GC Roots能直接关联到的对象。速度很快。并发标记阶段就是进行GC Roots Tracing的过程，而重新标记阶段则是为了修正并发标记期间因用户程序继续运作儿导致标记产生的变动的哪一部分对象的标记记录，这个阶段的停顿时间一般会比初试标记阶段稍长一些，但是远比并发标记的时间短。

由于整个过程中耗时最长的并发标记和并发清除过程收集器线程都是可以与用户线程一起工作。所以，从总体上来说，CMS收集器的内存回收过程是与用户线程一起并发执行的。

![](http://images.cnblogs.com/cnblogs_com/swanspouse/795043/o_CMS.jpg)

CMS收集器有以下缺点：
* CMS收集器对CPU资源非常敏感。
* CMS收集器无法处理浮动垃圾(Floating Garbage)
* CMS收集器基于“标记 -- 清除”算法实现的，容易产生碎片。

#### G1 收集器(Garbage First)

G1 收集器与其他GC收集器相比，具备如下特点：
* 并发与并行：G1 能充分利用多CPU、多核环境下的硬件优势，使用多个CPU来缩短Stop-The-World停顿的时间。
* 分代收集：G1可以不需要其他收集器配合就能独立管理整个GC堆。
* 空间整合：基于“标记--整理”算法的收集器。不会产生内存空间碎片。
* 可预测的停顿：能够建立可预测的停顿时间模型，让使用者明确指定一个长度为N毫秒的时间片段内，消耗在垃圾收集上的时间不得超过N毫秒。

在G1之前的其他收集器进行收集的范围都是整个新生代或者老年代，而G1不再是这样。使用G1收集器时，Java堆的内存布局就与其他收集器有很大差别，它将整个Java堆划分为多个大小相等的独立区域(Region)，虽然还保留有新生代和老年代的概念，但新生代和老年代已经不再是物理隔离的了，它们都是Region的集合。

G1收集器之所以能建立可预测的停顿时间模型，是因为它可以计划地避免在整个Java堆中进行全区域的垃圾收集。G1跟踪各个Region里面的垃圾堆积的价值大小，在后台维护一个优先列表。每次根据允许的时间，优先回收价值最大的Region（Garbage-First名字的由来）。这种使用Region划分内存空间以及有优先级的区域回收方式，保证了G1收集器在有限的时间内可以获取尽可能高的收集效率。

G1收集器的运作大致可以划分为以下几个步骤：
* 初始标记 (Inital Marking)  
* 并发标记 (Concurrent Marking)
* 最终标记 (Final Marking)
* 筛选回收 (Live Data Counting and Evacuation)

初始标记阶段仅仅只是标记以下GC Roots能直接关联到的对象。并且修改TAMS(Next Top at Mark Start)的值，让下一阶段用户程序并发运行时，能在正确可用的Region中创建新对象，这阶段需要停顿线程，但耗时很短。并发标记阶段是从GC Root开始对对中对象进行可达性分析，找出存活的对象，这阶段耗时较长，但可以与用户程序并发执行。而最终标记阶段则是为了修正并发标记期间因用户程序继续运作而导致标记产生变动的哪一部分标记记录，虚拟机将这段时间对象变化记录在线程Remembered Set Logs里面，最终标记需要把Remembered Set Logs的数据合并到Remembered Set中，这阶段需要停顿线程，但是可以并行执行。最后在筛选回首阶段首先对各个Region的回收价值和成本进行排序，根据用户所期望的GC停顿时间来制定回收计划。


## 垃圾收集器参数总结

参数 | 参数描述
--------------------|---------
-XX:+UseSerialGC	| Jvm运行在Client模式下的默认值，打开此开关后，使用Serial + Serial Old的收集器组合进行内存回收
-XX:+UseParNewGC    | 打开此开关后，使用ParNew + Serial Old的收集器进行垃圾回收
-XX:+UseConcMarkSweepGC | 使用ParNew + CMS +  Serial Old的收集器组合进行内存回收，Serial Old作为CMS出现“Concurrent Mode Failure”失败后的后备收集器使用。
-XX:+UseParallelGC | Jvm运行在Server模式下的默认值，打开此开关后，使用Parallel Scavenge +  Serial Old的收集器组合进行回收
-XX:+UseParallelOldGC	| 使用Parallel Scavenge +  Parallel Old的收集器组合进行回收
-XX:SurvivorRatio| 	新生代中Eden区域与Survivor区域的容量比值，默认为8，代表Eden:Subrvivor = 8:1
-XX:PretenureSizeThreshold	| 直接晋升到老年代对象的大小，设置这个参数后，大于这个参数的对象将直接在老年代分配
-XX:MaxTenuringThreshold	| 晋升到老年代的对象年龄，每次Minor GC之后，年龄就加1，当超过这个参数的值时进入老年代
-XX:UseAdaptiveSizePolicy| 	动态调整java堆中各个区域的大小以及进入老年代的年龄
-XX:+HandlePromotionFailure	| 是否允许新生代收集担保，进行一次minor gc后, 另一块Survivor空间不足时，将直接会在老年代中保留
-XX:ParallelGCThreads	| 设置并行GC进行内存回收的线程数
-XX:GCTimeRatio | 	GC时间占总时间的比列，默认值为99，即允许1%的GC时间，仅在使用Parallel Scavenge 收集器时有效
-XX:MaxGCPauseMillis	| 设置GC的最大停顿时间，在Parallel Scavenge 收集器下有效
-XX:CMSInitiatingOccupancyFraction	| 设置CMS收集器在老年代空间被使用多少后出发垃圾收集，默认值为68%，仅在CMS收集器时有效，-XX:CMSInitiatingOccupancyFraction=70
-XX:+UseCMSCompactAtFullCollection | 由于CMS收集器会产生碎片，此参数设置在垃圾收集器后是否需要一次内存碎片整理过程，仅在CMS收集器时有效
-XX:+CMSFullGCBeforeCompaction | 设置CMS收集器在进行若干次垃圾收集后再进行一次内存碎片整理过程，通常与UseCMSCompactAtFullCollection参数一起使用
-XX:+UseFastAccessorMethods| 原始类型优化
-XX:+DisableExplicitGC| 是否关闭手动System.gc
-XX:+CMSParallelRemarkEnabled| 降低标记停顿
-XX:LargePageSizeInBytes| 内存页的大小不可设置过大，会影响Perm的大小，-XX:LargePageSizeInBytes=128m

java -Xmx3550m -Xms3550m -Xmn2g -Xss128k
* -Xmx3550m：设置JVM最大可用内存为3550M。
* -Xms3550m：设置JVM初始内存为3550m。此值可以设置与-Xmx相同，以避免每次垃圾回收完成后JVM重新分配内存。
* -Xmn2g：设置年轻代大小为2G。整个JVM内存大小=年轻代大小 + 年老代大小 + 永久代大小。永久代一般固定大小为64m，所以增大年轻代后，将会减小年老代大小。此值对系统性能影响较大，Sun官方推荐配置为整个堆的3/8。
* -Xss128k：设置每个线程的堆栈大小。JDK5.0以后每个线程堆栈大小为1M，以前每个线程堆栈大小为256K。更具应用的线程所需内存大小进行调整。在相同物理内存下，减小这个值能生成更多的线程。但是操作系统对一个进程内的线程数还是有限制的，不能无限生成，经验值在3000~5000左右。

## 内存分配与回收策略

#### 对象优先在Eden上分配

大多数情况下，新生对象在Eden上分配。当Eden区没有足够的空间进行分配时，虚拟机发生一次Minor GC。

* 新生代GC ( Minor GC )：指发生在新生代的垃圾收集动作，因为Java对象大多数具备朝生夕灭的特性，所以Minor GC 非常频繁，回收速度也比较快。

* 老年代GC ( Major GC / Full GC)：指发生在老年代的GC，出现了Major GC，经常会伴随至少一次的Minor GC。 Major GC的速度一般会比 Minor GC慢10倍以上。

#### 大对象直接进入老年代

所谓大对象是指，需要大量连续内存空间的Java对象，最典型的大对象就是那种很长的字符串以及数组。虚拟机提供一个 -XX:PretenureSizeThreshold参数，令大于这个设置值的对象直接在老年代分配。这样做的目的是避免在Eden区以及两个Survivor区之间发生大量的内存复制。

#### 长期存活的对象将进入老年代

每个对象有一个对象年龄计数器。如果对象在Eden出生并经过第一次Minor GC后仍然存活，并且能被Survivor容纳的话，将被移动到Survivor空间中，并且对象年龄设为1，对象在Survivor区每熬过一次Minor GC 年龄就增加1. 当年龄达到一定程度（默认15，可设置）就将会被晋升到老年代中。

#### 动态年龄判定
为了更好的适应不同内存的情况。如果在Survivor空间中相同年龄的所有对象大小的总和大于Survivor空间的一半，年龄大于或者等于该年龄对象就可以直接进入老年代。

#### 空间分配担保

在发生Minor GC之前，虚拟机会先检查老年代最大可用的连续空间是否大于新生代所有对象的总控件。如果这个条件成立，那么Minor GC可以确保是安全的。如果不成立，则虚拟机会查看HandlePromotionFailure设置的值是否允许担保失败。如果允许，那么会继续检查老年代最大可用的连续空间是否大于历次晋升到老年代对象的平均大小，如果大于，将尝试一次Minor GC，尽管这次Minor GC是有风险的；如果小于，或者HandlePromotionFailure设置不允许冒险，那么这时也要进行一次Full GC。
