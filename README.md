# go-queue
前一久看到一篇文章美团高性能队列——Disruptor，时候自己琢磨了一下；经过反复修改，实现了一个相似的无锁队列EsQueue，该无锁队列相对Disruptor，而言少了队列数量属性quantity的CAP操作，因此性能杠杠的，在测试环境：windows10，Core(TM) i5-3320M CPU 2.6G, 8G 内存，go1.8.3，下性能达到1460-1600万之间。现在把代码发布出来，请同行验证一下，代码如下：

```go
注：请注意本方法已经通过 go test -race, 无警告。     
go1.8.3 amd64, Grp:   1, Times:   10000000, use:  573.9109ms, 57ns/op
go1.8.3 amd64, Grp:   2, Times:   20000000, use:  1.1548186s, 57ns/op
go1.8.3 amd64, Grp:   3, Times:   30000000, use:  1.6787567s, 55ns/op
go1.8.3 amd64, Grp:   4, Times:   40000000, use:  2.2651588s, 56ns/op
go1.8.3 amd64, Grp:   5, Times:   50000000, use:  2.8762257s, 57ns/op
go1.8.3 amd64, Grp:   6, Times:   60000000, use:  3.4914045s, 58ns/op
go1.8.3 amd64, Grp:   7, Times:   70000000, use:  4.0040473s, 57ns/op
go1.8.3 amd64, Grp:   8, Times:   80000000, use:  4.5712089s, 57ns/op
go1.8.3 amd64, Grp:   9, Times:   90000000, use:     5.1765s, 57ns/op
go1.8.3 amd64, Grp:  10, Times:   10000000, use:   586.914ms, 58ns/op
go1.8.3 amd64, Grp:  11, Times:   11000000, use:  644.4879ms, 58ns/op
go1.8.3 amd64, Grp:  12, Times:   12000000, use:  694.4974ms, 57ns/op
go1.8.3 amd64, Grp:  13, Times:   13000000, use:  745.5212ms, 57ns/op
go1.8.3 amd64, Grp:  14, Times:   14000000, use:  822.6344ms, 58ns/op
go1.8.3 amd64, Grp:  15, Times:   15000000, use:  868.4927ms, 57ns/op
go1.8.3 amd64, Grp:  16, Times:   16000000, use:  943.6699ms, 58ns/op
go1.8.3 amd64, Grp: Sum, Times:  541000000, use: 31.0982489s, 57ns/op
```

```go
注: 受Meltdown和Spectre处理器漏洞修复的影响，性能下降50%, 对应 pprof2
go1.8.3 amd64, Grp:   1, Times:    1000000, use:   90.8111ms, 90ns/op
go1.8.3 amd64, Grp:   2, Times:    2000000, use:  267.0498ms, 133ns/op
go1.8.3 amd64, Grp:   3, Times:    3000000, use:  325.0141ms, 108ns/op
go1.8.3 amd64, Grp:   4, Times:    4000000, use:  459.9871ms, 114ns/op
go1.8.3 amd64, Grp:   5, Times:    5000000, use:  531.0004ms, 106ns/op
go1.8.3 amd64, Grp:   6, Times:    6000000, use:   675.946ms, 112ns/op
go1.8.3 amd64, Grp:   7, Times:    7000000, use:  742.9081ms, 106ns/op
go1.8.3 amd64, Grp:   8, Times:    8000000, use:    900.09ms, 112ns/op
go1.8.3 amd64, Grp:   9, Times:    9000000, use:  966.0397ms, 107ns/op
go1.8.3 amd64, Grp:  10, Times:    1000000, use:  121.9575ms, 121ns/op
go1.8.3 amd64, Grp:  11, Times:    1100000, use:  123.9134ms, 112ns/op
go1.8.3 amd64, Grp:  12, Times:    1200000, use:  145.0397ms, 120ns/op
go1.8.3 amd64, Grp:  13, Times:    1300000, use:  144.9599ms, 111ns/op
go1.8.3 amd64, Grp:  14, Times:    1400000, use:  167.2686ms, 119ns/op
go1.8.3 amd64, Grp:  15, Times:    1500000, use:  168.7482ms, 112ns/op
go1.8.3 amd64, Grp:  16, Times:    1600000, use:  190.9838ms, 119ns/op
go1.8.3 amd64, Grp: Sum, Times:   93300000, use: 10.5246881s, 112ns/op
```

新增批量操作 Puts() Gets()后性能进一步提升

```go
使用Puts, Gets 进行测试，块尺寸32, 性能可以更快, 对应 pprof4
go1.8.3 amd64, Grp:   1, Times:    1500000, use:   908.254ms, 605ns/32op  18/op
go1.8.3 amd64, Grp:   2, Times:    3000000, use:  1.4947772s, 498ns/32op  15/op
go1.8.3 amd64, Grp:   3, Times:    3000000, use:  1.4490451s, 483ns/32op  15/op
go1.8.3 amd64, Grp:   4, Times:    4000000, use:  2.1125661s, 528ns/32op  16/op
go1.8.3 amd64, Grp:   5, Times:    5000000, use:  2.3802556s, 476ns/32op  14/op
go1.8.3 amd64, Grp:   6, Times:    3000000, use:  1.5050799s, 501ns/32op  15/op
go1.8.3 amd64, Grp:   7, Times:    3500000, use:  1.6807146s, 480ns/32op  15/op
go1.8.3 amd64, Grp:   8, Times:    4000000, use:  1.8279384s, 456ns/32op  14/op
go1.8.3 amd64, Grp:   9, Times:    3600000, use:  1.6087893s, 446ns/32op  13/op
go1.8.3 amd64, Grp:  10, Times:    4000000, use:  1.8343257s, 458ns/32op  14/op
go1.8.3 amd64, Grp:  11, Times:    4400000, use:  1.9333989s, 439ns/32op  13/op
go1.8.3 amd64, Grp:  12, Times:    3600000, use:  1.5931753s, 442ns/32op  13/op
go1.8.3 amd64, Grp:  13, Times:    3900000, use:  1.7232328s, 441ns/32op  13/op
go1.8.3 amd64, Grp:  14, Times:    4200000, use:  1.8263283s, 434ns/32op  13/op
go1.8.3 amd64, Grp:  15, Times:    3000000, use:  1.2999684s, 433ns/32op  13/op
go1.8.3 amd64, Grp:  16, Times:    3200000, use:  1.3860389s, 433ns/32op  13/op
go1.8.3 amd64, Grp: Sum, Times:   56900000, use: 26.5638885s, 466ns/32op  14/op

```

```go
使用Put, Gets 进行测试，块尺寸32, 性能提升不多, 对应 pprof5
----块尺寸-32----
go1.10.3 amd64, Grp:   1, Times:   48000000, use:  2.0869798s, 43ns/op
go1.10.3 amd64, Grp:   2, Times:   96000000, use:   4.253081s, 44ns/op
go1.10.3 amd64, Grp:   3, Times:   76800000, use:   3.448918s, 44ns/op
go1.10.3 amd64, Grp:   4, Times:  102400000, use:  4.5190001s, 44ns/op
go1.10.3 amd64, Grp:   5, Times:  128000000, use:  5.7560188s, 44ns/op
go1.10.3 amd64, Grp:   6, Times:   96000000, use:  4.0059824s, 41ns/op
go1.10.3 amd64, Grp:   7, Times:  112000000, use:  5.1509986s, 45ns/op
go1.10.3 amd64, Grp:   8, Times:  128000000, use:  5.5909991s, 43ns/op
go1.10.3 amd64, Grp:   9, Times:  115200000, use:  5.2770015s, 45ns/op
go1.10.3 amd64, Grp:  10, Times:  128000000, use:  5.5629135s, 43ns/op
go1.10.3 amd64, Grp:  11, Times:  140800000, use:  6.4389843s, 45ns/op
go1.10.3 amd64, Grp:  12, Times:  115200000, use:  5.1078895s, 44ns/op
go1.10.3 amd64, Grp:  13, Times:  124800000, use:  5.7341286s, 45ns/op
go1.10.3 amd64, Grp:  14, Times:  134400000, use:  5.8718558s, 43ns/op
go1.10.3 amd64, Grp:  15, Times:   96000000, use:  4.4031446s, 45ns/op
go1.10.3 amd64, Grp:  16, Times:  102400000, use:  4.4728723s, 43ns/op
go1.10.3 amd64, Grp: Sum, Times: 1744000000, use: 1m17.68076s, 44ns/op

```