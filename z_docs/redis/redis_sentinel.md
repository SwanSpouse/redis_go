

### redis sentinel 提纲


Sentinel (哨兵) 是Redis 高可用性解决方案：有一个或者多个Sentinel实例组成的Sentinel系统可以件事任意多个主服务器，
以及这些主服务器属下的所有从服务器，并在被监视的主服务器进入下线状态时，自动将下线的主服务器属下的某个从服务器升级为新的主服务器，
然后由新的主服务器代替已下线的主服务器继续处理命令请求。


sentinel本质上只是一个运行在特殊模式下的Redis服务器

sentinel启动步骤：
* 初始化服务器
* 将普通Redis服务器使用的代码换成Sentinel专用的代码
* 初始化 Sentinel状态
* 根据给定的配置文件，初始化Sentinel的监视主服务器列表
* 创建连向主服务器的网络连接


#### 创建向主服务器的网络连接

Sentinel会创建两个连向主服务器的异步网络连接

* 一个是命令连接，这个连接专门用于向主服务器发送命令，并接收命令回复。
* 一个是订阅连接，这个连接专门用于订阅主服务器的__sentinel__:hello频道。

#### 获取主服务器的信息

Sentinel默认会以每10s一次的频率通过命令连接向被监视的主服务器发送INFO命令，并通过分析INFO命令的回复来获取主服务器的当前信息。

* 一方面是主服务器本身的信息，run_id域记录的服务器运行ID以及role域记录的服务器角色。
* 另一方面是关于主服务器属下的所有从服务器的信息。Sentinel无需用户提供从服务器的地址信息，就可以自动发现从服务器。

#### 获取从服务器的信息

Sentinel发现主服务器有新的从服务器出现时，Sentinel除了会为这个新的从服务器创建相应的实例结构之外，Sentinel还会创建连接到从服务器的命令连接和订阅连接。

* 创建命令连接之后，Sentinel在默认情况下会10s一次的频率通过命令连接发送INFO命令。同时更新从服务器的实例结构。

#### 向主服务器和从服务器发送信息

默认情况下，Sentinel会以2s一次的频率通过命令连接向所有被监视的主服务器和从服务器发送以下格式的命令

```shell
PUBLISH __sentinel__:hello "<s_ip>,<s_port>,<s_runid>,<s_epoch>,<m_name>,<m_ip>,<m_port><m_epoch>"
```

* 其中s_开头的参数记录的是Sentinel本身的信息。
* 其中m_开头的参数记录的是主服务器的信息。

#### 接受来自主服务器和从服务器的频道信息

当Sentinel与一个主服务器或者从服务器建立起订阅连接之后，Sentinel会通过订阅连接，向服务器发送一下命令：

```shell
SUBSCRIBE __sentinel__:hello 
```
Sentinel对__sentinel__:hello 频道的订阅会一直持续到Sentinel与服务器的连接断开为止。

对于监视同一个服务器的多个Sentinel来说，一个Sentinel发送的信息会被其他Sentinel接收到，这些信息会被用于更新其他Sentinel对发送信息Sentinel的认知。

#### 创建连向其他Sentinel的命令连接

当Sentinel通过频道信息发现一个新的Sentinel时，它不仅会为新的Sentinel在sentinels字典中创建相应的实例结构，还会创建一个连向新Sentinel的命令连接，而新的Sentinel同样会创建向这个Sentinel的命令连接。

无须为运行的每个 Sentinel 分别设置其他 Sentinel 的地址， 因为 Sentinel 可以通过发布与订阅功能来自动发现正在监视相同主服务器的其他 Sentinel ， 这一功能是通过向频道 __sentinel__:hello 发送信息来实现的。

#### 检测主观下线状态

默认情况下，Sentinel会以每秒一次的频率向所有与它创建了命令连接的实例（包括主服务器，从服务器，其他Sentinel在内）发送PING命令，并通过返回值来判断实例是否在线。

* 配置文件中 down-after-milliseconds选项指定了Sentinel判断实例进入主观下线所需要的时间长度。

#### 检查客观下线状态

当一个Sentinel将一个主服务器判断为主观下线之后，为了确认这个主服务器是否真的下线了，它会向同样监视这一主服务器的其他Sentinel进行询问。当Sentinel从其他Sentinel那里接收到足够数量的已下线判断之后，
Sentinel就会将主服务器判定为客观下线。

* 使用SENTINEL is-master-down-by-addr <ip> <port> <current_epoch> <run_id> 命令询问其他Sentinel是否同意服务器已下线

* 接收到SENTINEL is-master-down-by-addr命令时，会检查目标服务器是否已下线，然后对命令进行回复。

* 判断客观下线的条件：根据配置文件配置的数量。sentinel monitor master 127.0.0.1 6379 2，包括当前Sentinel在内，如果有两个Sentinel认为服务器已经下线，当前Sentinel就将服务器判断为客观下线。

#### 主观下线到客观下线

从主观下线状态切换到客观下线状态并没有使用严格的法定人数算法（strong quorum algorithm）， 而是使用了流言协议： 如果 Sentinel 在给定的时间范围内， 从其他 Sentinel 那里接收到了足够数量的主服务器下线报告， 那么 Sentinel 就会将主服务器的状态从主观下线改变为客观下线。 如果之后其他 Sentinel 不再报告主服务器已下线， 那么客观下线状态就会被移除。

客观下线条件只适用于主服务器： 对于任何其他类型的 Redis 实例， Sentinel 在将它们判断为下线前不需要进行协商， 所以从服务器或者其他 Sentinel 永远不会达到客观下线条件。

只要一个 Sentinel 发现某个主服务器进入了客观下线状态， 这个 Sentinel 就可能会被其他 Sentinel 推选出， 并对失效的主服务器执行自动故障迁移操作。

#### redis领头sentinel选举规则

所有在线的Sentinel都有被选为领头sentinel的资格，换句话说，监视同一个主服务器的多个在线Sentinel中的任意一个都有可能成为领头Sentinel。

每次进行领头Sentinel选举之后，不论选举是否成功，所有sentinel的配置纪元（configuration epoch）的值都会自增一次。配置纪元实际上就是一个计数器，并没有什么特别的。

在一个配置纪元里面，所有Sentinel都有一次将某个Sentinel设置为局部领头Sentinel的机会，并且局部领头Sentinel一旦设置，在这个纪元里面就不能再更改。

每个发现主服务器进入客观下线的Sentinel都会要求其他Sentinel将自己设置为局部领头Sentinel。

当一个Sentinel（源Sentinel）向另一个Sentinel（目标Sentinel）发送SENTINEL is-master-down-by-addr命令，并且命令中的runid参数不是*符号而是源Sentinel的运行ID时，表示源Sentinel要求目标Sentinel将前者设置为后者的局部领头Sentinel。


Sentinel设置局部领头Sentinel的规则是先到先得：最先向目标Sentinel发送设置要求的源Sentinel将成为目标Sentinel的局部领头Sentinel，而之后接受到所有设置要求
都会被目标Sentinel拒绝。

目标Sentinel在接受到SENTINEL is-master-down-by-addr命令之后，将向源Sentinel返回一条命令回复，回复中的leader_runid参数和leader_epoch参数分别记录了目标
Sentinel的局部领头Sentinel的运行ID和配置纪元。

源Sentinel在接收到目标Sentinel返回的命令之后，会检查回复中的leader_epoch参数的值和自己的配置纪元是否相同，如果相同的话，那么源Sentinel继续取出回复中的leader_runid参数，
如果leader_runid参数的值和源Sentinel的运行ID一致，那么表示目标Sentinel将源Sentinel设置成了局部领头Sentinel。

如果有某个Sentinel被半数以上的Sentinel设置成了局部领头Sentinel，那么这个Sentinel成为领头Sentinel。在一个由10个Sentinel组成的Sentinel系统里面，只要有大于等于10/2+1=6个
Sentinel将某个Sentinel设置为局部领头Sentinel，那么被设置的那个Sentinel就会成为领头Sentinel。

因为领头Sentinel的产生需要半数以上的Sentinel支持，并且每个Sentinel在每个配置纪元里面只能设置一次局部领头Sentinel，所以在一个配置纪元里面，只会出现一个领头Sentinel。

如果在给定时限内，没有一个Sentinel被选举为领头Sentinel，那么各个Sentinel将在一段时间之后再次进行选举，直到选举出领头Sentinel为止。

#### 故障转移

在选举产生领头Sentinel之后，领头Sentinel将对已下线的主服务器执行故障转移操作。

* 在已下线的主服务器属下的所有从服务器里面，挑选出一个从服务器，并将其转换为主服务器。

* 让已下线主服务器属下的所有从服务器改为复制新的主服务器。

* 将已下线的服务器设置为新的服务器的从服务器，当这个就的主服务器重新上线时，它就会成为新的主服务器的从服务器。

