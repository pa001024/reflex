MoeCron
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
2. type command `go get github.com/pa001024/MoeCron`
3. cd into `main` folder and type command `go build`
4. run `./MoeCron` to start daemon

#### (Another way)Makefile

1. install [Go](http://golang.org/doc/install)
2. type command `go get github.com/pa001024/MoeCron`
3. cd into `MoeCron` folder and type command `make`
4. run `make test`

信息
-----------

1. 实现多向微博接口
2. 使用所配置的规则获取需要更新的feed信息 (比如 RSS/atom/API/MQ等)
3. 将获取到的更新信息过滤后发送到目标微博接口 (可设置限额)


架构
-----------

使用源(source)->过滤器(filter)->目标(target)的三层架构

### 工作过程

首先由source负责抓取目标页面 将数据传给target 如果target设置了filter 则将数据交给filter过滤

数据过滤完成后 执行target指定的action 比如发送微博 发送帖子 等

扩展性良好 无需修改主文件 只需按需增加source/filter/target并修改相应配置文件即可

