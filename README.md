# go-queue
前一久看到一篇文章美团高性能队列——Disruptor，时候自己琢磨了一下；经过反复修改，实现了一个相似的无锁队列EsQueue，该无锁队列相对Disruptor，而言少了队列数量属性quantity的CAP操作，因此性能杠杠的，在测试环境：windows10，Core(TM) i5-3320M CPU 2.6G, 8G 内存，go1.8.3，下性能达到1460-1600万之间。现在把代码发布出来，请同行验证一下，代码如下：

```go
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

```go
注: 受Meltdown和Spectre处理器漏洞修复的影响，性能下降50%
go1.8.3 amd64, Grp:   1, Times:    1000000, miss:     0, use:    90.8111ms,     90ns/op
go1.8.3 amd64, Grp:   2, Times:    2000000, miss:     0, use:   267.0498ms,    133ns/op
go1.8.3 amd64, Grp:   3, Times:    3000000, miss:     0, use:   325.0141ms,    108ns/op
go1.8.3 amd64, Grp:   4, Times:    4000000, miss:     0, use:   459.9871ms,    114ns/op
go1.8.3 amd64, Grp:   5, Times:    5000000, miss:     0, use:   531.0004ms,    106ns/op
go1.8.3 amd64, Grp:   6, Times:    6000000, miss:     0, use:    675.946ms,    112ns/op
go1.8.3 amd64, Grp:   7, Times:    7000000, miss:     0, use:   742.9081ms,    106ns/op
go1.8.3 amd64, Grp:   8, Times:    8000000, miss:     0, use:     900.09ms,    112ns/op
go1.8.3 amd64, Grp:   9, Times:    9000000, miss:     0, use:   966.0397ms,    107ns/op
go1.8.3 amd64, Grp:  10, Times:    1000000, miss:     0, use:   121.9575ms,    121ns/op
go1.8.3 amd64, Grp:  11, Times:    1100000, miss:     0, use:   123.9134ms,    112ns/op
go1.8.3 amd64, Grp:  12, Times:    1200000, miss:     0, use:   145.0397ms,    120ns/op
go1.8.3 amd64, Grp:  13, Times:    1300000, miss:     0, use:   144.9599ms,    111ns/op
go1.8.3 amd64, Grp:  14, Times:    1400000, miss:     0, use:   167.2686ms,    119ns/op
go1.8.3 amd64, Grp:  15, Times:    1500000, miss:     0, use:   168.7482ms,    112ns/op
go1.8.3 amd64, Grp:  16, Times:    1600000, miss:     0, use:   190.9838ms,    119ns/op
go1.8.3 amd64, Grp:  17, Times:    1700000, miss:     0, use:   186.9993ms,    109ns/op
go1.8.3 amd64, Grp:  18, Times:    1800000, miss:     0, use:   215.0006ms,    119ns/op
go1.8.3 amd64, Grp:  19, Times:    1900000, miss:     0, use:   211.9994ms,    111ns/op
go1.8.3 amd64, Grp:  20, Times:    2000000, miss:     0, use:   239.0021ms,    119ns/op
go1.8.3 amd64, Grp:  21, Times:    2100000, miss:     0, use:   233.9982ms,    111ns/op
go1.8.3 amd64, Grp:  22, Times:    2200000, miss:     0, use:   262.0004ms,    119ns/op
go1.8.3 amd64, Grp:  23, Times:    2300000, miss:     0, use:   245.9992ms,    106ns/op
go1.8.3 amd64, Grp:  24, Times:    2400000, miss:     0, use:   284.9998ms,    118ns/op
go1.8.3 amd64, Grp:  25, Times:    2500000, miss:     0, use:   279.0001ms,    111ns/op
go1.8.3 amd64, Grp:  26, Times:    2600000, miss:     0, use:   308.0004ms,    118ns/op
go1.8.3 amd64, Grp:  27, Times:    2700000, miss:     0, use:   300.0002ms,    111ns/op
go1.8.3 amd64, Grp:  28, Times:    2800000, miss:     0, use:   332.0004ms,    118ns/op
go1.8.3 amd64, Grp:  29, Times:    2900000, miss:     0, use:   323.0204ms,    111ns/op
go1.8.3 amd64, Grp:  30, Times:    3000000, miss:     0, use:   356.9781ms,    118ns/op
go1.8.3 amd64, Grp:  31, Times:    3100000, miss:     0, use:    346.001ms,    111ns/op
go1.8.3 amd64, Grp:  32, Times:    3200000, miss:     0, use:   377.9711ms,    118ns/op
go1.8.3 amd64, Grp: Sum, Times:   93300000, miss:     0, use:  10.5246881s,    112ns/op
```