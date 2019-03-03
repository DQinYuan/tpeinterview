# 题目描述

题目：设计并实现一个分布式数据库的测试框架。



要求：

- 框架能够按照指定的拓扑关系启动一个 tidb 集群
- 加载并运行 test case，可以自由添加多个 case
- case 在运行期间，可以通过框架提供的功能接口，干掉分布式数据库的某个节点，进行故障模拟
- 以上可以在单机实现，用 VM 或者 Docker 启动集群不限。（这里我选择使用Docker）



(题外话：第一次写Go语言，多多海涵)





# 功能演示

### 使用前的准备：

```
pip install docker-compose
git clone git@github.com:DQinYuan/tpeinterview.git
git clone https://github.com/pingcap/tidb-docker-compose.git
cd tpeinterview
go build
```



### 启动命令行界面：

```
./tpeinterview ../tidb-docker-compose/docker-compose.yml
```



格式是`./tpeinterview  docker-compose配置文件地址`



之后框架会自动使用docker-compose启动tidb集群，并且在mysql客户端测试连接成功后显示出命令行界面：

（输出）

```
tidb starting..., Please wait some seconds
tidb up ok 
test connect to tidb
[mysql] 2019/03/03 20:49:48 packets.go:36: unexpected EOF
[mysql] 2019/03/03 20:49:48 packets.go:36: unexpected EOF
[mysql] 2019/03/03 20:49:48 packets.go:36: unexpected EOF
table create fail, please waiting for tidb start up
test connect to tidb
[mysql] 2019/03/03 20:49:58 packets.go:36: unexpected EOF
[mysql] 2019/03/03 20:49:58 packets.go:36: unexpected EOF
[mysql] 2019/03/03 20:49:58 packets.go:36: unexpected EOF
table create fail, please waiting for tidb start up
test connect to tidb
CREATE TABLE IF NOT EXISTS testtable (TEST_KEY VARCHAR(64) PRIMARY KEY, FIELD0 VARCHAR(100), FIELD1 VARCHAR(100), FIELD2 VARCHAR(100), FIELD3 VARCHAR(100), FIELD4 VARCHAR(100), FIELD5 VARCHAR(100), FIELD6 VARCHAR(100), FIELD7 VARCHAR(100), FIELD8 VARCHAR(100), FIELD9 VARCHAR(100));
»  
```



因为从tidb集群的相关docker容器启动，到tidb真正启动完毕可以使用，中间还有一段时间，所以从上面的输出可以看出程序会先尝试连接tidb，直到连接成功后才显示出交互命令行，在我的电脑上，失败两次后即连接成功，在你的电脑上尝试的次数可能会有所不同。



使用help即可看到交互界面支持的命令：



```
» help
shell command

Usage:
   [command]

Available Commands:
  auth        assert db is readable and data is the same as what is write
  clean       clean test cases loaded
  help        Help about any command
  kill        kill tidb node container by name
  list        list all test cases
  load        load a script file as test case
  nodes       list all starting nodes name
  remove      remove a loaded test case
  run         run all loaded test case
  write       write record_num record into database

Flags:
  -h, --help   help for this command

Use " [command] --help" for more information about a command.
```





### 使用命令行手动测试



```
» write 10
...(输出省略)
» nodes
tidb-docker-compose_tispark-slave0_1
tidb-docker-compose_tispark-master_1
tidb-docker-compose_tidb_1
tidb-docker-compose_tikv0_1
tidb-docker-compose_tikv1_1
tidb-docker-compose_tikv2_1
tidb-docker-compose_grafana_1
tidb-docker-compose_pd1_1
tidb-docker-compose_tidb-vision_1
tidb-docker-compose_pd0_1
tidb-docker-compose_prometheus_1
tidb-docker-compose_pd2_1
tidb-docker-compose_pushgateway_1
» kill tidb-docker-compose_tikv0_1
kill container tidb-docker-compose_tikv0_1 ok
» write 1
INSERT INTO testtable (TEST_KEY, FIELD0, FIELD1, FIELD2, FIELD3, FIELD4, FIELD5, FIELD6, FIELD7, FIELD8, FIELD9) VALUES (? ,? ,? ,? ,? ,? ,? ,? ,? ,? ,?) [10 1010 101010 10101010 1010101010 101010101010 10101010101010 1010101010101010 101010101010101010 10101010101010101010 1010101010101010101010]
» auth
select * from testtable
OK
» exit
stopping tidb..., please wait some seconds
tidb stop ok 
```



含义如下：

- 写入10条数据
- 查看所有的tidb节点
- 删除一个tikv节点模拟故障
- 在写入一条数据
- 查询出所有的数据并且验证是否正确
- 务必使用exit退出交互界面，这样才能保证tidb的相关docker容器正常关闭



### 增加与删除测试用例



也可以选择将这些命令写在一个文件里作为一个测试用例，在本仓库的`testscripts`文件夹下就放着四个这样的示例，可以用`load`命令加载几个这样的用例，并且使用`run`命令一次性将他们全部执行，每个测试用例执行前程序都会将现场还原，所以各个测试用例的是不会相互影响的：



```
» load testscripts/test_*
current test cases:
testscripts/test_3
testscripts/test_4
testscripts/test_1
testscripts/test_2
» list
current test cases:
testscripts/test_1
testscripts/test_2
testscripts/test_3
testscripts/test_4
» run
======start run test case: testscripts/test_1 
...(省略)
======testscripts/test_1 test case success 
======start run test case: testscripts/test_2 
...(省略)
======testscripts/test_2 test case success 
======start run test case: testscripts/test_3 
...(省略)
======testscripts/test_3 test case fail, line num 4 fail, error message: auth fail, query fail 
======start run test case: testscripts/test_4 
...(省略)
======testscripts/test_4 test case success 


 fail cases: [testscripts/test_3] 
```



含义：

- 使用通配符路径加载`testscripts`目录下的全部测试脚本
- `list`命令列出当前已加载的全部测试用例
- `run`命令运行当前已加载的全部测试用例，从输出可以看出，只有`test_3`测试用例未能通过，原因是我在`test_3`中kill掉了两个tikv节点，导致数据库已经无法正常运作了



还可以使用`remove`命令删去加载的测试用例：



```
» list
current test cases:
testscripts/test_4
testscripts/test_1
testscripts/test_2
testscripts/test_3
» remove testscripts/test_3
» list
current test cases:
testscripts/test_1
testscripts/test_2
testscripts/test_4
```



