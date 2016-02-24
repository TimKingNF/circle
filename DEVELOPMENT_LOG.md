# 开发日志

# 2016-1-11
# 待解决问题
@ 调度器内部 urlmap 添加标志位，在调度器关闭前将 urlmap 同步至索引器
@ 添加 response 在分析器中处理更详细的日志内容

---------------------------------------------------------------------

# 2016-1-12
# 已解决问题
* 修复了启动监控之后 scheduler.Idle() 报错的bug，对scheduler.Start() 中先初始化一系列组件完成之后再修改 scheduler.running(status)
* 修复了监控信息中 urlMap 不显示的问题
* 修复了 scheduler.saveReqToCache 中的 主域比对问题

# 新增加内容
* scheduler 添加id 编号，在 crawlerError 添加了一个新的错误类型 SCHEDULER_ERROR
* 添加了 analyzer 日志

# 待解决问题
@ 下载之后的 response 并不会被分析

---------------------------------------------------------------------

# 2016-1-16
# 已解决问题
* 修正了启动之后 (sched.schedule) 发出了终止信号 (stopSign.Sign) 导致 channelManager 失效，从而引发的 response 不会被分析的bug

# 新增加内容
* 添加了 itempipeline 日志
* 添加了 crawlerControlModel.Init 初始化方法, crawlerControlModel.Accept 接收 url
* 添加了 args.Reset 方法
* 添加了 crawlerControlModel.Test 方法

# 待解决问题
* itempipeline 处理完之后的item 发送给 virtual db

---------------------------------------------------------------------

# 2016-1-19
# 已解决问题
* 完善 crawlerControlModel.Start 和 cralerControlModel.Accept 方法
* 完善 crawlerControlModel.Scheduler 的 restart 流程
* 完善 logging.ConsoleLogger 的读写锁

# 待解决问题
@ crawlerControlModel.Socket 和 crawlerControlModel.Dial 的编写

---------------------------------------------------------------------

# 2016-1-26
# 新增加内容
* 确定了系统整体架构
* 初步确定 divider 架构
* 初步编写完成了 divider.CrawlerManager
* 完成了基础的 crawlerControlModel 与 divider 的通信, 并编写了相应测试代码

# 待解决问题
@ 设计通信协议以便让 divider 分发数据
* itempipeline 的作用在于在同一份文本中以不同方式提取 关键内容

---------------------------------------------------------------------

# 2016-1-27
# 新增加内容
* 初步完成了 packages ( divider.databse 与 divider.indexdevice ) 的编写
* 初步编写 databaseControlModel

# 已解决问题
* 完成了 databseControlModel.Socket 的初步测试
* 初步设计通信协议

# 待增加内容
* 设计完成数据库
@ 完成 database.diriver 与 database.Api 的编写

# 待解决问题
@ 编写从 divider 回送 crawler、database、indexDevice 的测试用例

---------------------------------------------------------------------

# 2016-2-1
# 新增加内容
* 编写了 database.sqldiriver ,并完成测试

---------------------------------------------------------------------

# 2016-2-2
# 修改内容
* 修改 ControlModel.Init 的接受参数类型为 interface{} 使其接受 nil 值传递
* 调整了 logger.go 与 loggerArgs.go 至 common package 处

# 新增加内容
* 编写了 database.API ,并完成测试
* 添加了网页内容的主题相关 pageTopic

# 已解决问题
* 编写从 divider 回送 database, crawler 的测试用例
* 解决当 socket 断开之后重连的问题

---------------------------------------------------------------------

# 2016-2-6
# 新增加内容
* 编写 indexDevice.ControlModel

# 修改内容
* 修改了 ControlModel 有关 socket 的实现方式，改为 ControlModel.Socket
* 修改 logging 以便更好的区分不同时间的日志

# 待解决问题
@ 将 divider 对 crawler、database、index device 的监听 由原来的单例模式改成多例

---------------------------------------------------------------------

# 2016-2-10
# 待解决问题
* 将 crawler 的 requestCache 修改为 广度优先，不同的站点先爬取，尽可能在有限的时间内爬取多的内容
@ 将 divider 接收到的报文进行分发（如将 crawler 的 query urlMap 操作经由 divider 转发至 indexDevice）

# 新增加内容
* 初步确定 index device 架构
* 初步编写 indexDevice.CardsManager、indexDevice.MapPage、indexDevice.MapKeyword、indexDevice.keywordCards、
indexDevice.pageCards
* 编写 common.Tag 与 common.Page
* 添加了 crawler.scheduler.urlCache，定时批量发送数据给 divider，判断请求是否可以下载
* 添加了 crawler.scheduler.DataManager，以便数据在处理过程中可以发送给 divider 处理
* 添加了 crawlerControlModel 对数据收发的处理流程

# 修改内容
@ 修改 crawler.scheduler 部分架构
* 将 crawler.scheduler 中的 urlMap 修改为批量向 index device 发送查询请求

# 已解决问题
* 初步完成 crawlerControlModel 对查询 urlMap 的 发送往 divider 的操作的处理

---------------------------------------------------------------------

# 2016-2-11
# 待解决问题
@ 将 divider 中 crawlerManager、databaseManager、indexDeviceManager 存放的 Conn 改为 list 存放
* 在 queryIndex 操作中，对于不存在的 page 进行 cardsManager 添加 page 处理
@ 对 crawler 等待查询结果的过程添加超时处理
@ crawler 在发送 queryIndex 给 divider 之后，divider 应该保存下该连接以便 index device 回复时能准确找到所发送的 crawler
@ divider 对动态连接上来的 crawler、database、indexDevice 分配 id,在通信传输的过程中带上 id 号
@ crawler.scheduler.saveToReqCache 中的 判断 urlCache.Downloading()，当多个判断请求阻塞时，只返回了一个结果，其他判断请求将会一直阻塞直至超时结束，需要将 crawler.scheduler.urlCache.sign 改为 stopSign

# 新增加内容
* 添加了 divider 对来自 crawler 的数据的转发处理
* 添加了 divider 与 index device 之间的数据收发处理，完善从 crawler 到 divider 到 index device 到 divider 再到 crawler 的一个 queryIndex 的数据传输流程, 并完成测试代码

# 修改内容
* 将部分数据传输结构由原来的 chan string 改为 chan common.ControlMessage

---------------------------------------------------------------------

# 2016-2-13
# 新增加内容
* 添加下载后提取快照内容, 并把快照内容经由 divider 发送给 indexDevice 发送给 divider 发送给 database 的流程

# 修改内容
* 完善 common.ControlMessage 在 divider、crawler、index device 中的传递方式
* 完善数据在转发流程中通过 divider 记录不同 id 号以便区分各个客户端的流程

# 待解决问题
* savaPage 操作中将数据进行负载均衡计算,需要将要添加的 page 发送给哪个 index device, 并 index device 需要详细记录 cardsManager 的一系列变动(方便数据找回，和加载内存), 再由 index device 经由一段空闲时间后 批量发送给 database 执行添加操作
@ 添加 divider.base.LoadBalancer, 并对 divider.crawlerManager、divider.indexDeviceManager、divider.databaseManager 添加 GenLoadBalancing() 进行负载均衡计算，返回id
@ 负载均衡算法计算(轮替法)

---------------------------------------------------------------------

# 2016-2-15
# 备注内容
* 负载均衡算法，假定所有客户端的权重一致，根据使用过程中耗时的长短记录为矩阵，通过矩阵的相乘使的结果不断收敛趋近于客户端的真实权重，在负载均衡算法中，通过衡量客户端权重、最近分配任务来确定 负载均衡器提供的 客户端id

# 新增加内容
* 编写负载均衡算法

---------------------------------------------------------------------

# 2016-2-16
# 新增加内容
* 编写 indexDevice.cardsManager 

# 修改内容
* 修改 indexDevice.cardsManager.keywords 的保存结构为 list
* 重新设计 indexDevice.cardsManager 结构

---------------------------------------------------------------------

# 2016-2-17
# 修改内容
* 修改 crawler.itemPipeline.GenKeywordsFromPage, 增加了对关键词更多的提取规则
* 修改 indexDeviceControlModel.queryIndexAnalyze

# 新增加内容
* 完善 indexDeviceCtrlM.cardsManager, 并添加 filecache
* 添加 indexDeviceCtrlM.pageCache 用于缓存暂未发给 database 的 page 数据
* 添加 indexDeviceCtrlM.analyzePageUpdate 的流程

# 待解决问题
@ 修改 crawler.itemPipeline.GenKeywordsFromPage 的 json 转化规则
* 完成 indexDevice.cardsManager 、indexDevice.updatePage 、indexDevice.savePage 的测试流程

---------------------------------------------------------------------

# 2016-2-18
# 修改内容
* 修正 crawler.itempipeline.processitem.GenKeywordsFromPage
* 修改 crawler.urlCache.SendDataChan 的类型为 cmn.ControlMessage
* 修改 crawler.dataManager.SendDataChan 的类型为 cmn.ControlMessage
* 修改 crawler.scheduler.SendChan 的类型为 cmn.ControlMessage
* crawler.itemPipeline.GenKeywordsFromPage, 增加了对关键词更多的提取规则 (title), 并剔除了无效的 tags
* 修改 indexDevice.parsePageAnalyze 的解析规则

# 新增加内容
* 为 IndexDevice、divider、crawler、database 添加了记录 socket 发送消息的日志

# 待解决问题
@ 编写 indexDevice.cardsManager 完成 filecahce 写入和读取
@ 编写测试 indexDevice.cardsManager.queryKeyword 的代码文件

---------------------------------------------------------------------

# 2016-2-19
# 已解决问题
* 修复 indexDevice.cardsManager.filecache.Write 其中的bug
* 修复 indexDevice.cardsManager.fielcache.Load 其中的bug

# 修改内容
* 修改 indexDevice.cardsManager.pageCards.Page 的类型为 page.Id, 为 page 提供 genId(), 提高数据的复用性

# 待解决问题
* 完善 indexDevice.cardsManager.queryKeyword 的方法，并对其结果进行排序
@ 对 indexDevice.cardsManager.pagecache 发往 database 的存储操作进行测试
@ 对 indexDevice.cardsManager.updatePage 请求 发往 crawler 重新下载并重新发送给 indexDevice 执行 indexDevice.cardsManager.updatePage() 操作
* 对 queryKeyword 操作来说，由 显示服务器发送给 indexDevice 获取到 page 合并之后才开始进行排序
@ 验证 crawler 在重复发送 同样的 page 数据之后，indexDevice 不会重复添加 page 的问题

---------------------------------------------------------------------

# 2016-2-20
# 修改内容
* 修正 database.savePageAnalyze() 不能正确获取 page 的bug

---------------------------------------------------------------------

# 2016-2-22
# 修改内容
* 剔除了 indexDevice.indexDeviceArgs.pageCap 的无效属性
* 添加 divider 在转发 crawler.SavePage 中存在两种可能性，一种是不存在 accepterId 的时候说明该请求由 crawler 自主发出，另一种存在 accepterId 的时候说明该请求由 indexDevice 发出
* 修改了 crawler.updatePage 的传输方式，并添加 divider 的转发
* 修复了 indexDevice 在 loadFileCache 之后 idGenertor 的值会重新计算的问题

# 待解决问题
* crawler.requestCache.GET() 的时候越界的问题

---------------------------------------------------------------------

# 2016-2-23
# 修改内容
* 修改 crawler、database、divider、indexDevice 的 args 包在编译之后 使用 reset 重置 log 日志参数无效的问题
* 修改 crawler.urlcache 在 等待

# 已解决问题
* 修改由于两个 线程竞争抢夺 crawler.requestCache.Get() 造成的数据越界问题 

# 待解决问题
* 编写 工具 启动各个程序

# 添加内容

