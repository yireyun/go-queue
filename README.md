# go-queue
前一久看到一篇文章美团高性能队列——Disruptor，时候自己琢磨了一下；经过反复修改，实现了一个相似的无锁队列EsQueue，该无锁队列相对Disruptor，而言少了队列数量属性quantity的CAP操作，因此性能杠杠的，在测试环境：windows10，Core(TM) i5-3320M CPU 2.6G, 8G 内存，go1.8.3，下性能达到1460-1600万之间。现在把代码发布出来，请同行验证一下，代码如下：

```
注：请注意本方法已经通过 go test -race, 无警告。
go1.8.3 amd64, Grp:   1, Times:   10000000, miss:     0, use:   573.9109ms,     57ns/op
go1.8.3 amd64, Grp:   2, Times:   20000000, miss:     0, use:   1.1548186s,     57ns/op
go1.8.3 amd64, Grp:   3, Times:   30000000, miss:     0, use:   1.6787567s,     55ns/op
go1.8.3 amd64, Grp:   4, Times:   40000000, miss:     0, use:   2.2651588s,     56ns/op
go1.8.3 amd64, Grp:   5, Times:   50000000, miss:     0, use:   2.8762257s,     57ns/op
go1.8.3 amd64, Grp:   6, Times:   60000000, miss:     0, use:   3.4914045s,     58ns/op
go1.8.3 amd64, Grp:   7, Times:   70000000, miss:     0, use:   4.0040473s,     57ns/op
go1.8.3 amd64, Grp:   8, Times:   80000000, miss:     0, use:   4.5712089s,     57ns/op
go1.8.3 amd64, Grp:   9, Times:   90000000, miss:     0, use:      5.1765s,     57ns/op
go1.8.3 amd64, Grp:  10, Times:   10000000, miss:     0, use:    586.914ms,     58ns/op
go1.8.3 amd64, Grp:  11, Times:   11000000, miss:     0, use:   644.4879ms,     58ns/op
go1.8.3 amd64, Grp:  12, Times:   12000000, miss:     0, use:   694.4974ms,     57ns/op
go1.8.3 amd64, Grp:  13, Times:   13000000, miss:     0, use:   745.5212ms,     57ns/op
go1.8.3 amd64, Grp:  14, Times:   14000000, miss:     0, use:   822.6344ms,     58ns/op
go1.8.3 amd64, Grp:  15, Times:   15000000, miss:     0, use:   868.4927ms,     57ns/op
go1.8.3 amd64, Grp:  16, Times:   16000000, miss:     0, use:   943.6699ms,     58ns/op
go1.8.3 amd64, Grp: Sum, Times:  541000000, miss:     0, use:  31.0982489s,     57ns/op
```
