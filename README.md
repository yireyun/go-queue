# go-queue
前一久看到一篇文章美团高性能队列——Disruptor，时候自己琢磨了一下；经过反复修改，实现了一个相似的无锁队列EsQueue，该无锁队列相对Disruptor，而言少了队列数量属性quantity的CAP操作，因此性能杠杠的，在测试环境：windows10，Core(TM) i5-3320M CPU 2.6G, 8G 内存，go1.7.4，下性能达到1360-1500万之间。现在把代码发布出来，请同行验证一下，代码如下：

```
go1.7.4, Grp:   1, Times:   10000000, miss:     0, use:   552.8887ms,     55ns/op 
go1.7.4, Grp:   2, Times:   20000000, miss:     0, use:   1.4554794s,     72ns/op 
go1.7.4, Grp:   3, Times:   30000000, miss:     0, use:   2.2382081s,     74ns/op 
go1.7.4, Grp:   4, Times:   40000000, miss:     0, use:   2.9799835s,     74ns/op 
go1.7.4, Grp:   5, Times:   50000000, miss:     0, use:   3.7434942s,     74ns/op 
go1.7.4, Grp:   6, Times:   60000000, miss:     0, use:   4.4849934s,     74ns/op 
go1.7.4, Grp:   7, Times:   70000000, miss:     0, use:   5.2675198s,     75ns/op 
go1.7.4, Grp:   8, Times:   80000000, miss:     0, use:   6.0115122s,     75ns/op 
go1.7.4, Grp:   9, Times:   90000000, miss:     0, use:   6.7634953s,     75ns/op 
go1.7.4, Grp:  10, Times:   10000000, miss:     0, use:    765.514ms,     76ns/op 
go1.7.4, Grp:  11, Times:   11000000, miss:     0, use:   827.0602ms,     75ns/op 
go1.7.4, Grp:  12, Times:   12000000, miss:     0, use:   907.1067ms,     75ns/op 
go1.7.4, Grp:  13, Times:   13000000, miss:     0, use:   975.6492ms,     75ns/op 
go1.7.4, Grp:  14, Times:   14000000, miss:     0, use:   1.0634071s,     75ns/op 
go1.7.4, Grp:  15, Times:   15000000, miss:     0, use:    1.136756s,     75ns/op 
go1.7.4, Grp:  16, Times:   16000000, miss:     0, use:   1.2058011s,     75ns/op 
go1.7.4, Grp: Sum, Times:  541000000, miss:     0, use:  40.3788689s,     74ns/op 
```
