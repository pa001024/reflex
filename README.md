MoeWorker
=========

Info
-----------

1. implementation of some microblog api
2. use some rules(in config.json) to get feed info which need to update (like RSS/atom/API/MQ)
3. filtering to post to target (with limit)

Structure
-----------

Build
-----------

1. install [Go](http://golang.org/doc/install)
2. type command `go get github.com/pa001024/MoeWorker`
3. cd into `main` folder and type command `go build`
4. run `./MoeWorker` to start daemon

#### (Another way)Makefile

1. install [Go](http://golang.org/doc/install)
2. type command `go get github.com/pa001024/MoeWorker`
3. cd into `MoeWorker` folder and type command `make`
4. run `make test`

Using lib
-----------

1. [freetype-go](https://code.google.com/p/freetype-go/)

信息
-----------

1. 实现多向微博接口
2. 使用所配置的规则获取需要更新的feed信息 (比如 RSS/atom/API/MQ等)
3. 将获取到的更新信息过滤后发送到目标微博接口 (可设置限额)


结构
-----------

工作组由 源(source) -> 过滤器(filter) -> 目标(target) 的三部分组成

### 工作过程

1. 首先由 __源__ 从指定源(如http)抓取目标页面
2. __源__的处理器将抓取到的数据进行简单标准化 封装成 __标准结构__:
    FeedInfo(id,title,content,author,date,picurl)
3. __源__将标准化后的结构传给 __目标__  如果 __目标__ 设置了 __过滤器__ , 则将数据交给 __过滤器__ 过滤
4. __过滤器__组依次按指定的规则和参数对 __标准结构__ 的各项数据进行修改后 将最终的数据传递给 __目标__
5. __目标__按指定的规则和参数发送 __标准结构__

类结构由 接口(Interface) -> 容器(Container) -> 实现(Implementation) 三部分组成

### 类结构

1. __接口__: 提供给外部的通用API
2. __容器__: 存放实现所需的数据 实现持久化 可配置
3. __实现__: 真正的实例 继承接口和容器 实现接口表述的方法

使用类库
-----------

1. [freetype-go](https://code.google.com/p/freetype-go/)