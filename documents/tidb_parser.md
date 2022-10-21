# tidb_parser

介绍tidb_parser, 了解其使用方法与实现原理.

其实是这两篇文章的复制粘贴整合:
https://cn.pingcap.com/blog/

https://github.com/pingcap/tidb/blob/master/parser/docs/quickstart.md

# 如何使用tidb_parser

https://github.com/pingcap/tidb/blob/master/parser/docs/quickstart.md

This parser is highly compatible with MySQL syntax. You can use it as a library, parse a text SQL into an AST tree, and traverse the AST nodes.

In this example, you will build a project, which can extract all the column names from a text SQL.

## Prerequisites

- [Golang](https://golang.org/dl/) version 1.13 or above. You can follow the instructions in the official [installation page](https://golang.org/doc/install) (check it by `go version`)

## Create a Project

```
mkdir colx && cd colx
go mod init colx && touch main.go
```

## Import Dependencies

First, you need to use `go get` to fetch the dependencies through git hash. The git hashes are available in [release page](https://github.com/pingcap/tidb/releases). Take `v5.3.0` as an example:

```
go get -v github.com/pingcap/tidb/parser@4a1b2e9
```

> **NOTE**
>
> The parser was merged into TiDB repo since v5.3.0. So you can only choose version v5.3.0 or higher in this TiDB repo.
>
> You may want to use advanced API on expressions (a kind of AST node), such as numbers, string literals, booleans, nulls, etc. It is strongly recommended using the `types` package in TiDB repo with the following command:
>
> ```
> go get -v github.com/pingcap/tidb/types/parser_driver@4a1b2e9
> ```
>
> and import it in your golang source code:
>
> ```
> import _ "github.com/pingcap/tidb/types/parser_driver"
> ```

Your directory should contain the following three files:

```
.
├── go.mod
├── go.sum
└── main.go
```

Now, open `main.go` with your favorite editor, and start coding!

## Parse SQL text

To convert a SQL text to an AST tree, you need to:

1. Use the [`parser.New()`](https://pkg.go.dev/github.com/pingcap/tidb/parser?tab=doc#New) function to instantiate a parser, and
2. Invoke the method [`Parse(sql, charset, collation)`](https://pkg.go.dev/github.com/pingcap/tidb/parser?tab=doc#Parser.Parse) on the parser.

```golang
package main

import (
	"fmt"

	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/parser/test_driver"
)

func parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()

	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	return &stmtNodes[0], nil
}

func main() {
	astNode, err := parse("SELECT a, b FROM t")
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return
	}
	fmt.Printf("%v\n", *astNode)
}
```

Test the parser by running the following command:

```
go run main.go
```

If the parser runs properly, you should get a result like this:

```
&{{{{SELECT a, b FROM t}}} {[]} 0xc0000a1980 false 0xc00000e7a0 <nil> 0xc0000a19b0 <nil> <nil> [] <nil> <nil> none [] false false 0 <nil>}
```

> **NOTE**
>
> Here are a few things you might want to know:
>
> - To use a parser, a `parser_driver` is required. It decides how to parse the basic data types in SQL.
>
>   You can use [`github.com/pingcap/tidb/parser/test_driver`](https://pkg.go.dev/github.com/pingcap/tidb/parser/test_driver) as the `parser_driver` for test. Again, if you need advanced features, please use the `parser_driver` in TiDB (run `go get -v github.com/pingcap/tidb/types/parser_driver@4a1b2e9` and import it).
>
> - The instantiated parser object is not goroutine safe. It is better to keep it in a single goroutine.
>
> - The instantiated parser object is not lightweight. It is better to reuse it if possible.
>
> - The 2nd and 3rd arguments of [`parser.Parse()`](https://pkg.go.dev/github.com/pingcap/tidb/parser?tab=doc#Parser.Parse) are charset and collation respectively. If you pass an empty string into it, a default value is chosen.

## Traverse AST Nodes

Now you get the AST tree root of a SQL statement. It is time to extract the column names by traverse.

Parser implements the interface [`ast.Node`](https://pkg.go.dev/github.com/pingcap/tidb/parser/ast?tab=doc#Node) for each kind of AST node, such as SelectStmt, TableName, ColumnName. [`ast.Node`](https://pkg.go.dev/github.com/pingcap/tidb/parser/ast?tab=doc#Node) provides a method `Accept(v Visitor) (node Node, ok bool)` to allow any struct that has implemented [`ast.Visitor`](https://pkg.go.dev/github.com/pingcap/tidb/parser/ast?tab=doc#Visitor) to traverse itself.

[`ast.Visitor`](https://pkg.go.dev/github.com/pingcap/tidb/parser/ast?tab=doc#Visitor) is defined as follows:

```golang
type Visitor interface {
	Enter(n Node) (node Node, skipChildren bool)
	Leave(n Node) (node Node, ok bool)
}
```

Now you can define your own visitor, `colX`(columnExtractor):

```golang
type colX struct{
	colNames []string
}

func (v *colX) Enter(in ast.Node) (ast.Node, bool) {
	if name, ok := in.(*ast.ColumnName); ok {
		v.colNames = append(v.colNames, name.Name.O)
	}
	return in, false
}

func (v *colX) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
```

Finally, wrap `colX` in a simple function:

```golang
func extract(rootNode *ast.StmtNode) []string {
	v := &colX{}
	(*rootNode).Accept(v)
	return v.colNames
}
```

And slightly modify the main function:

```golang
func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: colx 'SQL statement'")
		return
	}
	sql := os.Args[1]
	astNode, err := parse(sql)
	if err != nil {
		fmt.Printf("parse error: %v\n", err.Error())
		return
	}
	fmt.Printf("%v\n", extract(astNode))
}
```

Test your program:

```
go build && ./colx 'select a, b from t'
[a b]
```

You can also try a different SQL statement as an input. For example:

```
$ ./colx 'SELECT a, b FROM t GROUP BY (a, b) HAVING a > c ORDER BY b'
[a b a b a c b]

If necessary, you can deduplicate by yourself.

$ ./colx 'SELECT a, b FROM t/invalid_str'
parse error: line 1 column 19 near "/invalid_str"
```

Enjoy!

## Restore

tidb早些年开展过一项restore计划, 将ast转换回sql语句:

https://github.com/pingcap/tidb/issues/8532

这是restore的全部任务, 现已全部完成.

有了restore, 我们就可以通过修改ast节点对sql语句进行变异.

>  另外, 我们也可以借助这些任务认识tidb的ast结构.

**使用示例**

```golang
if sel, ok := (*rootNode).(*ast.SelectStmt); ok {
		buf := new(bytes.Buffer)
		ctx := format.NewRestoreCtx(format.DefaultRestoreFlags, buf)
		err := sel.Restore(ctx)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(buf.String())
	}
```

可以配合ast节点的修改, 比如:

```golang
if bin, ok := in.(*ast.BinaryOperationExpr); ok {
		if bin.Op == opcode.EQ {
			bin.Op = opcode.GE
		}
	}
```

# tidb_parser的实现原理

## Tidb模块概览

TiDB 的模块非常多，这里做一个整体介绍，大家可以看到每个模块大致是做什么用的，想看相关功能的代码是，可以直接找到对应的模块。

| Package                   | Introduction                                                 |
| ------------------------- | ------------------------------------------------------------ |
| ast                       | 抽象语法树的数据结构定义，例如 `SelectStmt`定义了一条 Select 语句被解析成什么样的数据结构 |
| cmd/benchdb               | 简单的 benchmark 工具，用于性能优化                          |
| cmd/benchfilesort         | 简单的 benchmark 工具，用于性能优化                          |
| cmd/benchkv               | Transactional KV API benchmark 工具，也可以看做 KV 接口的使用样例 |
| cmd/benchraw              | Raw KV API benchmark 工具，也可以看做不带事务的 KV 接口的使用样例 |
| cmd/importer              | 根据表结构以及统计信息伪造数据的工具，用于构造测试数据       |
| config                    | 配置文件相关逻辑                                             |
| context                   | 主要包括 Context 接口，提供一些基本的功能抽象，很多包以及函数都会依赖于这个接口，把这些功能抽象为接口是为了解决包之间的依赖关系 |
| ddl                       | DDL 的执行逻辑                                               |
| distsql                   | 对分布式计算接口的抽象，通过这个包把 Executor 和 TiKV Client 之间的逻辑做隔离 |
| domain                    | domain 可以认为是一个存储空间的抽象，可以在其中创建数据库、创建表，不同的 domain 之间，可以存在相同名称的数据库，有点像 Name Space。一般来说单个 TiDB 实例只会创建一个 Domain 实例，其中会持有 information schema 信息、统计信息等。 |
| executor                  | 执行器相关逻辑，可以认为大部分语句的执行逻辑都在这里，比较杂，后面会专门介绍 |
| expression                | 表达式相关逻辑，包括各种运算符、内建函数                     |
| expression/aggregation    | 聚合表达式相关的逻辑，比如 Sum、Count 等函数                 |
| infoschema                | SQL 元信息管理模块，另外对于 Information Schema 的操作，都会访问这里 |
| kv                        | KV 引擎接口以及一些公用方法，底层的存储引擎需要实现这个包中定义的接口 |
| meta                      | 利用 structure 包提供的功能，管理存储引擎中存储的 SQL 元信息，infoschema/DDL 利用这个模块访问或者修改 SQL 元信息 |
| meta/autoid               | 用于生成全局唯一自增 ID 的模块，除了用于给每个表的自增 ID 之外，还用于生成全局唯一的 Database ID 和 Table ID |
| metrics                   | Metrics 相关信息，所有的模块的 Metrics 信息都在这里          |
| model                     | SQL 元信息数据结构，包括 DBInfo / TableInfo / ColumnInfo / IndexInfo 等 |
| mysql                     | MySQL 相关的常量定义                                         |
| owner                     | TiDB 集群中的一些任务只能由一个实例执行，比如异步 Schema 变更，这个模块用于多个 tidb-server 之间协调产生一个任务执行者。每种任务都会产生自己的执行者。 |
| parser                    | 语法解析模块，主要包括词法解析 (lexer.go) 和语法解析 (parser.y)，这个包对外的主要接口是 Parse()，用于将 SQL 文本解析成 AST |
| parser/goyacc             | 对 GoYacc 的包装                                             |
| parser/opcode             | 关于操作符的一些常量定义                                     |
| perfschema                | Performance Schema 相关的功能，默认不会启用                  |
| plan                      | 查询优化相关的逻辑                                           |
| privilege                 | 用户权限管理接口                                             |
| privilege/privileges      | 用户权限管理功能实现                                         |
| server                    | MySQL 协议以及 Session 管理相关逻辑                          |
| sessionctx/binloginfo     | 向 Binlog 模块输出 Binlog 信息                               |
| sessionctx/stmtctx        | Session 中的语句运行时所需要的信息，比较杂                   |
| sessionctx/variable       | System Variable 相关代码                                     |
| statistics                | 统计信息模块                                                 |
| store                     | 储存引擎相关逻辑，这里是存储引擎和 SQL 层之间的交互逻辑      |
| store/mockoracle          | 模拟 TSO 组件                                                |
| store/mockstore           | 实例化一个 Mock TiKV 的逻辑，主要方法是 NewMockTikvStore，把这部分逻辑从 mocktikv 中抽出来是避免循环依赖 |
| store/mockstore/mocktikv  | 在单机存储引擎上模拟 TiKV 的一些行为，主要作用是本地调试、构造单元测试以及指导 TiKV 开发 Coprocessor 相关逻辑 |
| store/tikv                | TiKV 的 Go 语言 Client                                       |
| store/tikv/gcworker       | TiKV GC 相关逻辑，tidb-server 会根据配置的策略向 TiKV 发送 GC 命令 |
| store/tikv/oracle         | TSO 服务接口                                                 |
| store/tikv/oracle/oracles | TSO 服务的 Client                                            |
| store/tikv/tikvrpc        | TiKV API 的一些常量定义                                      |
| structure                 | 在 Transactional KV API 上定义的一层结构化 API，提供 List/Queue/HashMap 等结构 |
| table                     | 对 SQL 的 Table 的抽象                                       |
| table/tables              | 对 table 包中定义的接口的实现                                |
| tablecodec                | SQL 到 Key-Value 的编解码，每种数据类型的具体编解码方案见 `codec`包 |
| terror                    | TiDB 的 error 封装                                           |
| tidb-server               | 服务的 main 方法                                             |
| types                     | 所有和类型相关的逻辑，包括一些类型的定义、对类型的操作等     |
| types/json                | json 类型相关的逻辑                                          |
| util                      | 一些实用工具，这个目录下面包很多，这里只会介绍几个重要的包   |
| util/admin                | TiDB 的管理语句（ `Admin`语句）用到的一些方法                |
| util/charset              | 字符集相关逻辑                                               |
| util/chunk                | Chunk 是 TiDB 1.1 版本引入的一种数据表示结构。一个 Chunk 中存储了若干行数据，在进行 SQL 计算时，数据是以 Chunk 为单位在各个模块之间流动 |
| util/codec                | 各种数据类型的编解码                                         |
| x-server                  | X-Protocol 实现                                              |

在全部 80 个模块中，下面几个模块是最重要的，希望大家能仔细阅读。

- plan
- expression
- executor
- distsql
- store/tikv
- ddl
- tablecodec
- server
- types
- kv
- tidb

## SQL层架构

![SQL 层架构](tidb_parser.assets/2_a7f34487c3.png)

SQL 层架构

这幅图比上一幅图详细很多，大体描述了 SQL 核心模块，大家可以从左边开始，顺着箭头的方向看。

### Protocol Layer

最左边是 TiDB 的 Protocol Layer，这里是与 Client 交互的接口，目前 TiDB 只支持 MySQL 协议，相关的代码都在 `server` 包中。

这一层的主要功能是管理客户端 connection，解析 MySQL 命令并返回执行结果。具体的实现是按照 MySQL 协议实现，具体的协议可以参考 [MySQL 协议文档 ](https://dev.mysql.com/doc/internals/en/client-server-protocol.html)。这个模块我们认为是当前实现最好的一个 MySQL 协议组件，如果大家的项目中需要用到 MySQL 协议解析、处理的功能，可以参考或引用这个模块。

连接建立的逻辑在 server.go 的 [Run() ](https://github.com/pingcap/tidb/blob/source-code/server/server.go#L236)方法中，主要是下面两行：

```
236:         conn, err := s.listener.Accept()

258:         go s.onConn(conn)
```

单个 Session 处理命令的入口方法是调用 clientConn 类的 [dispatch 方法 ](https://github.com/pingcap/tidb/blob/source-code/server/conn.go#L465)，这里会解析协议并转给不同的处理函数。

#### SQL Layer

大体上讲，一条 SQL 语句需要经过，语法解析-->合法性验证-->制定查询计划-->优化查询计划-->根据计划生成查询器-->执行并返回结果 等一系列流程。这个主干对应于 TiDB 的下列包：

| Package    | 作用                                          |
| ---------- | --------------------------------------------- |
| tidb       | Protocol 层和 SQL 层之间的接口                |
| parser     | 语法解析                                      |
| plan       | 合法性验证 + 制定查询计划 + 优化查询计划      |
| executor   | 执行器生成以及执行                            |
| distsql    | 通过 TiKV Client 向 TiKV 发送以及汇总返回结果 |
| store/tikv | TiKV Client                                   |

#### KV API Layer

TiDB 依赖于底层的存储引擎提供数据的存取功能，但是并不是依赖于特定的存储引擎（比如 TiKV），而是对存储引擎提出一些要求，满足这些要求的引擎都能使用（其中 TiKV 是最合适的一款）。

最基本的要求是『带事务的 Key-Value 引擎，且提供 Go 语言的 Driver』，再高级一点的要求是『支持分布式计算接口』，这样 TiDB 可以把一些计算请求下推到 存储引擎上进行。

这些要求都可以在 `kv`这个包的 [接口 ](https://github.com/pingcap/tidb/blob/source-code/kv/kv.go)中找到，存储引擎需要提供实现了这些接口的 Go 语言 Driver，然后 TiDB 利用这些接口操作底层数据。

对于最基本的要求，可以重点看这几个接口：

- [Transaction ](https://github.com/pingcap/tidb/blob/source-code/kv/kv.go#L121)：事务基本操作
- [Retriever ](https://github.com/pingcap/tidb/blob/source-code/kv/kv.go#L75)：读取数据的接口
- [Mutator ](https://github.com/pingcap/tidb/blob/source-code/kv/kv.go#L91)：修改数据的接口
- [Storage ](https://github.com/pingcap/tidb/blob/source-code/kv/kv.go#L229)：Driver 提供的基本功能
- [Snapshot ](https://github.com/pingcap/tidb/blob/source-code/kv/kv.go#L214)：在数据 Snapshot 上面的操作
- [Iterator ](https://github.com/pingcap/tidb/blob/source-code/kv/kv.go#L255)：`Seek`返回的结果，可以用于遍历数据

有了上面这些接口，可以对数据做各种所需要的操作，完成全部 SQL 功能，但是为了更高效的进行运算，我们还定义了一个高级计算接口，可以关注这三个 Interface/struct :

- [Client ](https://github.com/pingcap/tidb/blob/source-code/kv/kv.go#L150)：向下层发送请求以及获取下层存储引擎的计算能力
- [Request ](https://github.com/pingcap/tidb/blob/source-code/kv/kv.go#L176): 请求的内容
- [Response ](https://github.com/pingcap/tidb/blob/source-code/kv/kv.go#L204): 返回结果的抽象

## SQL的一生

### 概述

上一篇文章讲解了 TiDB 项目的结构以及三个核心部分，本篇文章从 SQL 处理流程出发，介绍哪里是入口，对 SQL 需要做哪些操作，知道一个 SQL 是从哪里进来的，在哪里处理，并从哪里返回。

SQL 有很多种，比如读、写、修改、删除以及管理类的 SQL，每种 SQL 有自己的执行逻辑，不过大体上的流程是类似的，都在一个统一的框架下运转。

### 框架

我们先从整体上看一下，一条语句需要经过哪些方面的工作。如果大家还记得上一篇文章所说的三个核心部分，可以想到首先要经过协议解析和转换，拿到语句内容，然后经过 SQL 核心层逻辑处理，生成查询计划，最后去存储引擎中获取数据，进行计算，返回结果。这个就是一个粗略的处理框架，本篇文章会把这个框架不断细化。

对于第一部分，协议解析和转换，所有的逻辑都在 server 这个包中，主要逻辑分为两块：一是连接的建立和管理，每个连接对应于一个 Session；二是在单个连接上的处理逻辑。第一点本文暂时不涉及，感兴趣的同学可以翻翻代码，看看连接如何建立、如何握手、如何销毁，后面也会有专门的文章讲解。对于 SQL 的执行过程，更重要的是第二点，也就是已经建立了连接，在这个连接上的操作，本文会详细讲解这一点。

对于第二部分，SQL 层的处理是整个 TiDB 最复杂的部分。这部分为什么复杂？原因有三点：

1. SQL 语言本身是一门复杂的语言，语句的种类多、数据类型多、操作符多、语法组合多，这些『多』经过排列组合会变成『很多』『非常多』，所以需要写大量的代码来处理。
2. SQL 是一门表意的语言，只是说『要什么数据』，而不说『如何拿数据』，所以需要一些复杂的逻辑选择『如何拿数据』，也就是选择一个好的查询计划。
3. 底层是一个分布式存储引擎，会面临很多单机存储引擎不会遇到的问题，比如做查询计划的时候要考虑到下层的数据是分片的、网络不通了如何处理等情况，所以需要一些复杂的逻辑处理这些情况，并且需要一个很好的机制将这些处理逻辑封装起来。这些复杂性是看懂源码比较大的障碍，所以本篇文章会尽量排除这些干扰，给大家讲解核心的逻辑是什么。

这一层有几个核心概念，掌握了这几个也就掌握了这一层的框架，请大家关注下面这几个接口：

- [Session](https://github.com/pingcap/tidb/blob/source-code/session.go#L62)
- [RecordSet](https://github.com/pingcap/tidb/blob/source-code/ast/ast.go#L136)
- [Plan](https://github.com/pingcap/tidb/blob/source-code/plan/plan.go#L30)
- [LogicalPlan](https://github.com/pingcap/tidb/blob/source-code/plan/plan.go#L140)
- [PhysicalPlan](https://github.com/pingcap/tidb/blob/source-code/plan/plan.go#L190)
- [Executor](https://github.com/pingcap/tidb/blob/source-code/executor/executor.go#L190)

下面的详细内容中，会讲解这些接口，用这些接口理清楚整个逻辑。

对于第三部分可以认为两块，第一块是 KV 接口层，主要作用是将请求路由到正确的的 KV Server，接收返回消息传给 SQL 层，并在此过程中处理各种异常逻辑；第二块是 KV Server 的具体实现，由于 TiKV 比较复杂，我们可以先看 Mock-TiKV 的实现，这里有所有的 SQL 分布式计算相关的逻辑。 接下来的几节，会对上面的三块详细展开描述。

#### 协议层入口

当和客户端的连接建立好之后，TiDB 中会有一个 Goroutine 监听端口，等待从客户端发来的包，并对发来的包做处理。这段逻辑在 server/conn.go 中，可以认为是 TiDB 的入口，本节介绍一下这段逻辑。 首先看 [clientConn.Run() ](https://github.com/pingcap/tidb/blob/source-code/server/conn.go#L413)，这里会在一个循环中，不断的读取网络包：

```
    445:    data, err := cc.readPacket()
```

然后调用 dispatch() 方法处理收到的请求：

```
    465:        if err = cc.dispatch(data); err != nil {
```

接下来进入 [clientConn.dispatch() ](https://github.com/pingcap/tidb/blob/source-code/server/conn.go#L571)方法：

```
    func (cc *clientConn) dispatch(data []byte) error {
```

这里要处理的包是原始 byte 数组，里面的内容读者可以参考 [MySQL 协议 ](https://dev.mysql.com/doc/internals/en/client-server-protocol.html)，第一个 byte 即为 Command 的类型：

```
        580:     cmd := data[0]
```

然后根据 Command 的类型，调用对应的处理函数，最常用的 Command 是 [COM_QUERY ](https://dev.mysql.com/doc/internals/en/com-query.html#packet-COM_QUERY)，对于大多数 SQL 语句，只要不是用 Prepared 方式，都是 COM_QUERY，本文也只会介绍这个 Command，其他的 Command 请读者对照 MySQL 文档看代码。 对于 Command Query，从客户端发送来的主要是 SQL 文本，处理函数是 [handleQuery() ](https://github.com/pingcap/tidb/blob/source-code/server/conn.go#L849):

```
    func (cc *clientConn) handleQuery(goCtx goctx.Context, sql string) (err error) {
```

这个函数会调用具体的执行逻辑：

```
    850:  rs, err := cc.ctx.Execute(goCtx, sql)
```

这个 Execute 方法的实现在 server/driver_tidb.go 中，

```
    func (tc *TiDBContext) Execute(goCtx goctx.Context, sql string) (rs []ResultSet, err error) {
        rsList, err := tc.session.Execute(goCtx, sql)
```

最重要的就是调用 tc.session.Execute，这个 session.Execute 的实现在 session.go 中，自此会进入 SQL 核心层，详细的实现会在后面的章节中描述。

经过一系列处理，拿到 SQL 语句的结果后会调用 writeResultset 方法把结果写回客户端：

```
        857:        err = cc.writeResultset(goCtx, rs[0], false, false)
```

#### 协议层出口

出口比较简单，就是上面提到的 [writeResultset ](https://github.com/pingcap/tidb/blob/source-code/server/conn.go#L909)方法，按照 MySQL 协议的要求，将结果（包括 Field 列表、每行数据）写回客户端。读者可以参考 MySQL 协议中的 [COM_QUERY Response ](https://dev.mysql.com/doc/internals/en/com-query-response.html)理解这段代码。

接下的几节我们进入核心流程，看看一条文本的 SQL 是如何处理的。我会先介绍所有的流程，然后用一个图把所有的流程串起来。

#### Session

Session 中最重要的函数是 [Execute ](https://github.com/pingcap/tidb/blob/source-code/session.go#L742)，这里会调用下面所述的各种模块，完成语句执行。注意这里在执行的过程中，会考虑 Session 环境变量，比如是否 `AutoCommit`，时区是什么。

#### Lexer & Yacc

这两个组件共同构成了 Parser 模块，调用 Parser，可以将文本解析成结构化数据，也就是抽象语法树 （AST）：

```
    session.go 699:     return s.parser.Parse(sql, charset, collation)
```

在解析过程中，会先用 [lexer ](https://github.com/pingcap/tidb/blob/source-code/parser/lexer.go)不断地将文本转换成 token，交付给 Parser，Parser 是根据 [yacc 语法 ](https://github.com/pingcap/tidb/blob/source-code/parser/parser.y)生成，根据语法不断的决定 Lexer 中发来的 token 序列可以匹配哪条语法规则，最终输出结构化的节点。 例如对于这样一条语句 `SELECT * FROM t WHERE c > 1;`，可以匹配 [SelectStmt 的规则 ](https://github.com/pingcap/tidb/blob/source-code/parser/parser.y#L3936)，被转换成下面这样一个数据结构：

```
    type SelectStmt struct {
        dmlNode
        resultSetNode
    
        // SelectStmtOpts wraps around select hints and switches.
        *SelectStmtOpts
        // Distinct represents whether the select has distinct option.
        Distinct bool
        // From is the from clause of the query.
        From *TableRefsClause
        // Where is the where clause in select statement.
        Where ExprNode
        // Fields is the select expression list.
        Fields *FieldList
        // GroupBy is the group by expression list.
        GroupBy *GroupByClause
        // Having is the having condition.
        Having *HavingClause
        // OrderBy is the ordering expression list.
        OrderBy *OrderByClause
        // Limit is the limit clause.
        Limit *Limit
        // LockTp is the lock type
        LockTp SelectLockType
        // TableHints represents the level Optimizer Hint
        TableHints []*TableOptimizerHint
    }
```

其中，`FROM t`会被解析为 `FROM`字段，`WHERE c > 1`被解析为 `Where`字段，`*`被解析为 `Fields`字段。所有的语句的结构够都被抽象为一个 `ast.StmtNode`，这个接口读者可以自行看注释，了解一下。这里只提一点，大部分 ast 包中的数据结构，都实现了 `ast.Node`接口，这个接口有一个 `Accept`方法，后续对 AST 的处理，主要依赖这个 Accept 方法，以 [Visitor 模式 ](https://en.wikipedia.org/wiki/Visitor_pattern)遍历所有的节点以及对 AST 做结构转换。

#### 制定查询计划以及优化

拿到 AST 之后，就可以做各种验证、变化、优化，这一系列动作的入口在这里：

```
    session.go 805:             stmt, err := compiler.Compile(goCtx, stmtNode)
```

我们进入 [Compile 函数 ](https://github.com/pingcap/tidb/blob/source-code/executor/compiler.go#L37)，可以看到三个重要步骤：

- `plan.Preprocess`: 做一些合法性检查以及名字绑定；
- `plan.Optimize`：制定查询计划，并优化，这个是最核心的步骤之一，后面的文章会重点介绍；
- 构造 `executor.ExecStmt`结构：这个 [ExecStmt ](https://github.com/pingcap/tidb/blob/source-code/executor/adapter.go#L148)结构持有查询计划，是后续执行的基础，非常重要，特别是 Exec 这个方法。

#### 生成执行器

在这个过程中，会将 plan 转换成 executor，执行引擎即可通过 executor 执行之前定下的查询计划，具体的代码见 [ExecStmt.buildExecutor() ](https://github.com/pingcap/tidb/blob/source-code/executor/adapter.go#L318)：

```
    executor/adpter.go 227:  e, err := a.buildExecutor(ctx)
```

生成执行器之后，被 [封装在一个 `recordSet`结构中 ](https://github.com/pingcap/tidb/blob/source-code/executor/adapter.go#L260)：

```
        return &recordSet{
            executor:    e,
            stmt:        a,
            processinfo: pi,
            txnStartTS:  ctx.Txn().StartTS(),
        }, nil
```

这个结构实现了 [`ast.RecordSet`](https://github.com/pingcap/tidb/blob/source-code/ast/ast.go#L136)接口，从字面上大家可以看出，这个接口代表了查询结果集的抽象，我们看一下它的几个方法：

```
    // RecordSet is an abstract result set interface to help get data from Plan.
    type RecordSet interface {
        // Fields gets result fields.
        Fields() []*ResultField
    
        // Next returns the next row, nil row means there is no more to return.
        Next(ctx context.Context) (row types.Row, err error)
    
        // NextChunk reads records into chunk.
        NextChunk(ctx context.Context, chk *chunk.Chunk) error
    
        // NewChunk creates a new chunk with initial capacity.
        NewChunk() *chunk.Chunk
    
        // SupportChunk check if the RecordSet supports Chunk structure.
        SupportChunk() bool
    
        // Close closes the underlying iterator, call Next after Close will
        // restart the iteration.
        Close() error
    }
```

通过注释大家可以看到这个接口的作用，简单来说，可以调用 Fields() 方法获得结果集每一列的类型，调用 Next/NextChunk() 可以获取一行或者一批数据，调用 Close() 可以关闭结果集。

#### 运行执行器

TiDB 的执行引擎是以 Volcano 模型运行，所有的物理 Executor 构成一个树状结构，每一层通过调用下一层的 Next/NextChunk() 方法获取结果。 举个例子，假设语句是 `SELECT c1 FROM t WHERE c2 > 1;`，并且查询计划选择的是全表扫描+过滤，那么执行器树会是下面这样：

![执行器树](tidb_parser.assets/1_c3e07627e9.png)

执行器树

大家可以从图中看到 Executor 之间的调用关系，以及数据的流动方式。那么最上层的 Next 是在哪里调用，也就是整个计算的起始点在哪里，谁来驱动这个流程？ 有两个地方大家需要关注，这两个地方分别处理两类语句。 第一类语句是 Select 这种查询语句，需要对客户端返回结果，这类语句的执行器调用点在 [给客户端返回数据的地方 ](https://github.com/pingcap/tidb/blob/master/server/conn.go#L909)：

```
            row, err = rs.Next(ctx)
```

这里的 `rs`即为一个 `RecordSet`接口，对其不断的调用 `Next()`，拿到更多结果，返回给 MySQL Client。 第二类语句是 Insert 这种不需要返回数据的语句，只需要把语句执行完成即可。这类语句也是通过 `Next`驱动执行，驱动点在 [构造 `recordSet`结构之前 ](https://github.com/pingcap/tidb/blob/source-code/executor/adapter.go#L251)：

```
        // If the executor doesn't return any result to the client, we execute it without delay.
        if e.Schema().Len() == 0 {
            return a.handleNoDelayExecutor(goCtx, e, ctx, pi)
        } else if proj, ok := e.(*ProjectionExec); ok && proj.calculateNoDelay {
            // Currently this is only for the "DO" statement. Take "DO 1, @a=2;" as an example:
            // the Projection has two expressions and two columns in the schema, but we should
            // not return the result of the two expressions.
            return a.handleNoDelayExecutor(goCtx, e, ctx, pi)
        }
```

### 总结

上面描述了整个 SQL 层的执行框架，这里用一幅图来描述整个过程：

![SQL 层执行过程](tidb_parser.assets/2_9c0f3a6934.png)

SQL 层执行过程

通过这篇文章，相信大家已经了解了 TiDB 中语句的执行框架，整个逻辑还是比较简单，框架中具体的模块的详细解释会在后续章节中给出。下一篇文章会用具体的语句为例，帮助大家理解本篇文章。

## SQL Parser的实现

`SQL Parser`的功能是把 SQL 语句按照 SQL 语法规则进行解析，将文本转换成抽象语法树（`AST`），这部分功能需要些背景知识才能比较容易理解，我尝试做下相关知识的介绍，希望能对读懂这部分代码有点帮助。

TiDB 是使用 [goyacc ](https://github.com/cznic/goyacc)根据预定义的 SQL 语法规则文件 [parser.y ](https://github.com/pingcap/tidb/blob/source-code/parser/parser.y)生成 SQL 语法解析器。我们可以在 TiDB 的 [Makefile ](https://github.com/pingcap/tidb/blob/50e98f427e7943396dbe38d23178b9f9dc5398b7/Makefile#L50)文件中看到这个过程，先 build `goyacc`工具，然后使用 `goyacc`根据 `parser.y`生成解析器 `parser.go`：

```makefile
goyacc:
    $(GOBUILD) -o bin/goyacc parser/goyacc/main.go

parser: goyacc
    bin/goyacc -o /dev/null parser/parser.y
    bin/goyacc -o parser/parser.go parser/parser.y 2>&1 ...
```

[goyacc ](https://github.com/cznic/goyacc)是 [yacc ](http://dinosaur.compilertools.net/)的 Golang 版，所以要想看懂语法规则定义文件 [parser.y ](https://github.com/pingcap/tidb/blob/source-code/parser/parser.y)，了解解析器是如何工作的，先要对 [Lex & Yacc ](http://dinosaur.compilertools.net/)有些了解。

### Lex & Yacc 介绍

[Lex & Yacc ](http://dinosaur.compilertools.net/)是用来生成词法分析器和语法分析器的工具，它们的出现简化了编译器的编写。`Lex & Yacc`分别是由贝尔实验室的 [Mike Lesk ](https://en.wikipedia.org/wiki/Mike_Lesk)和 [Stephen C. Johnson ](https://en.wikipedia.org/wiki/Stephen_C._Johnson)在 1975 年发布。对于 Java 程序员来说，更熟悉的是 [ANTLR ](http://www.antlr.org/)，`ANTLR 4`提供了 `Listener`+`Visitor`组合接口， 不需要在语法定义中嵌入`actions`，使应用代码和语法定义解耦。`Spark`的 SQL 解析就是使用了 `ANTLR`。`Lex & Yacc`相对显得有些古老，实现的不是那么优雅，不过我们也不需要非常深入的学习，只要能看懂语法定义文件，了解生成的解析器是如何工作的就够了。我们可以从一个简单的例子开始：

![图例](tidb_parser.assets/2_3a000040e8.png)

上图描述了使用 `Lex & Yacc`构建编译器的流程。`Lex`根据用户定义的 `patterns`生成词法分析器。词法分析器读取源代码，根据 `patterns`将源代码转换成 `tokens`输出。`Yacc`根据用户定义的语法规则生成语法分析器。语法分析器以词法分析器输出的 `tokens`作为输入，根据语法规则创建出语法树。最后对语法树遍历生成输出结果，结果可以是产生机器代码，或者是边遍历 `AST`边解释执行。

从上面的流程可以看出，用户需要分别为 `Lex`提供 `patterns`的定义，为 `Yacc`提供语法规则文件，`Lex & Yacc`根据用户提供的输入文件，生成符合他们需求的词法分析器和语法分析器。这两种配置都是文本文件，并且结构相同：

```
... definitions ...
%%
... rules ...
%%
... subroutines ...
```

文件内容由 `%%`分割成三部分，我们重点关注中间规则定义部分。对于上面的例子，`Lex`的输入文件如下：

```
...
%%
/* 变量 */
[a-z]    {
            yylval = *yytext - 'a';
            return VARIABLE;
         }   
/* 整数 */
[0-9]+   {
            yylval = atoi(yytext);
            return INTEGER;
         }
/* 操作符 */
[-+()=/*\n] { return *yytext; }
/* 跳过空格 */
[ \t]    ;
/* 其他格式报错 */
.        yyerror("invalid character");
%%
...
```

上面只列出了规则定义部分，可以看出该规则使用正则表达式定义了变量、整数和操作符等几种 `token`。例如整数 `token`的定义如下：

```
[0-9]+  {
            yylval = atoi(yytext);
            return INTEGER; 
        }
```

当输入字符串匹配这个正则表达式，大括号内的动作会被执行：将整数值存储在变量 `yylval`中，并返回 `token`类型 `INTEGER`给 `Yacc`。

再来看看 `Yacc`语法规则定义文件：

```
%token INTEGER VARIABLE
%left '+' '-'
%left '*' '/'
...
%%

program:
        program statement '\n' 
        |
        ;

statement:
        expr                    { printf("%d\n", $1); }
        | VARIABLE '=' expr     { sym[$1] = $3; }
        ;
        
expr:
        INTEGER
        | VARIABLE              { $$ = sym[$1]; }
        | expr '+' expr         { $$ = $1 + $3; }
        | expr '-' expr         { $$ = $1 - $3; }
        | expr '*' expr         { $$ = $1 * $3; }
        | expr '/' expr         { $$ = $1 / $3; }
        | '(' expr ')'          { $$ = $2; }
        ;

%%
...
```

第一部分定义了 `token`类型和运算符的结合性。四种运算符都是左结合，同一行的运算符优先级相同，不同行的运算符，后定义的行具有更高的优先级。

语法规则使用了 `BNF`定义。`BNF`可以用来表达上下文无关（*context-free*）语言，大部分的现代编程语言都可以使用 `BNF`表示。上面的规则定义了三个**产生式**。**产生式**冒号左边的项（例如 `statement`）被称为**非终结符**， `INTEGER`和 `VARIABLE`被称为**终结符**,它们是由 `Lex`返回的 `token`。**终结符**只能出现在**产生式**的右侧。可以使用**产生式**定义的语法生成表达式：

```
expr -> expr * expr
     -> expr * INTEGER
     -> expr + expr * INTEGER
     -> expr + INTEGER * INTEGER
     -> INTEGER + INTEGER * INTEGER
```

解析表达式是生成表达式的逆向操作，我们需要归约表达式到一个**非终结符**。`Yacc`生成的语法分析器使用**自底向上**的归约（*shift-reduce*）方式进行语法解析，同时使用堆栈保存中间状态。还是看例子，表达式 `x + y * z`的解析过程：

```
1    . x + y * z
2    x . + y * z
3    expr . + y * z
4    expr + . y * z
5    expr + y . * z
6    expr + expr . * z
7    expr + expr * . z
8    expr + expr * z .
9    expr + expr * expr .
10   expr + expr .
11   expr .
12   statement .
13   program  .
```

点（`.`）表示当前的读取位置，随着 `.`从左向右移动，我们将读取的 `token`压入堆栈，当发现堆栈中的内容匹配了某个**产生式**的右侧，则将匹配的项从堆栈中弹出，将该**产生式**左侧的**非终结符**压入堆栈。这个过程持续进行，直到读取完所有的 `tokens`，并且只有**启始非终结符**（本例为 `program`）保留在堆栈中。

产生式右侧的大括号中定义了该规则关联的动作，例如：

```
expr:  expr '*' expr         { $$ = $1 * $3; }
```

我们将堆栈中匹配该**产生式**右侧的项替换为**产生式**左侧的**非终结符**，本例中我们弹出 `expr '*' expr`，然后把 `expr`压回堆栈。 我们可以使用 `$position`的形式访问堆栈中的项，`$1`引用的是第一项，`$2`引用的是第二项，以此类推。`$$`代表的是归约操作执行后的堆栈顶。本例的动作是将三项从堆栈中弹出，两个表达式相加，结果再压回堆栈顶。

上面例子中语法规则关联的动作，在完成语法解析的同时，也完成了表达式求值。一般我们希望语法解析的结果是一棵抽象语法树（`AST`），可以这么定义语法规则关联的动作：

```
...
%%
...
expr:
    INTEGER             { $$ = con($1); }
    | VARIABLE          { $$ = id($1); }
    | expr '+' expr     { $$ = opr('+', 2, $1, $3); }
    | expr '-' expr     { $$ = opr('-', 2, $1, $3); }
    | expr '*' expr     { $$ = opr('*', 2, $1, $3); } 
    | expr '/' expr     { $$ = opr('/', 2, $1, $3); }
    | '(' expr ')'      { $$ = $2; }
    ; 
%%
nodeType *con(int value) {
    ...
}
nodeType *id(int i) {
    ...
}
nodeType *opr(int oper, int nops, ...) {
    ...
}    
```

上面是一个语法规则定义的片段，我们可以看到，每个规则关联的动作不再是求值，而是调用相应的函数，该函数会返回抽象语法树的节点类型 `nodeType`，然后将这个节点压回堆栈，解析完成时，我们就得到了一颗由 `nodeType`构成的抽象语法树。对这个语法树进行遍历访问，可以生成机器代码，也可以解释执行。

至此，我们大致了解了 `Lex & Yacc`的原理。其实还有非常多的细节，例如如何消除语法的歧义，但我们的目的是读懂 TiDB 的代码，掌握这些概念已经够用了。

### goyacc 简介

[goyacc ](https://github.com/cznic/goyacc)是 golang 版的 `Yacc`。和 `Yacc`的功能一样，`goyacc`根据输入的语法规则文件，生成该语法规则的 go 语言版解析器。`goyacc`生成的解析器 `yyParse`要求词法分析器符合下面的接口：

```
type yyLexer interface {
    Lex(lval *yySymType) int
    Error(e string)
}
```

或者

```
type yyLexerEx interface {
    yyLexer
    // Hook for recording a reduction.
    Reduced(rule, state int, lval *yySymType) (stop bool) // Client should copy *lval.
}
```

TiDB 没有使用类似 `Lex`的工具生成词法分析器，而是纯手工打造，词法分析器对应的代码是 [parser/lexer.go ](https://github.com/pingcap/tidb/blob/source-code/parser/lexer.go)， 它实现了 `goyacc`要求的接口：

```
...
// Scanner implements the yyLexer interface.
type Scanner struct {
    r   reader
    buf bytes.Buffer

    errs         []error
    stmtStartPos int

    // For scanning such kind of comment: /*! MySQL-specific code */ or /*+ optimizer hint */
    specialComment specialCommentScanner

    sqlMode mysql.SQLMode
}
// Lex returns a token and store the token value in v.
// Scanner satisfies yyLexer interface.
// 0 and invalid are special token id this function would return:
// return 0 tells parser that scanner meets EOF,
// return invalid tells parser that scanner meets illegal character.
func (s *Scanner) Lex(v *yySymType) int {
    tok, pos, lit := s.scan()
    v.offset = pos.Offset
    v.ident = lit
    ...
}
// Errors returns the errors during a scan.
func (s *Scanner) Errors() []error {
    return s.errs
}
```

另外 `lexer`使用了 `字典树`技术进行 `token`识别，具体的实现代码在 [parser/misc.go](https://github.com/pingcap/tidb/blob/source-code/parser/misc.go)

### TiDB SQL Parser 的实现

终于到了正题。有了上面的背景知识，对 TiDB 的 `SQL Parser`模块会相对容易理解一些。TiDB 的词法解析使用的 [手写的解析器 ](https://github.com/pingcap/tidb/blob/source-code/parser/lexer.go)（这是出于性能考虑），语法解析采用 `goyacc`。先看 SQL 语法规则文件 [parser.y ](https://github.com/pingcap/tidb/blob/source-code/parser/parser.y)，`goyacc`就是根据这个文件生成SQL语法解析器的。

`parser.y`有 6500 多行，第一次打开可能会被吓到，其实这个文件仍然符合我们上面介绍过的结构：

```
... definitions ...
%%
... rules ...
%%
... subroutines ...
```

`parser.y`第三部分 `subroutines`是空白没有内容的， 所以我们只需要关注第一部分 `definitions`和第二部分 `rules`。

第一部分主要是定义 `token`的类型、优先级、结合性等。注意 `union`这个联合体结构体：

```
%union {
    offset int // offset
    item interface{}
    ident string
    expr ast.ExprNode
    statement ast.StmtNode
}
```

该联合体结构体定义了在语法解析过程中被压入堆栈的**项**的属性和类型。

压入堆栈的**项**可能是 `终结符`，也就是 `token`，它的类型可以是`item`或 `ident`；

这个**项**也可能是 `非终结符`，即产生式的左侧，它的类型可以是 `expr`、 `statement`、 `item`或 `ident`。

`goyacc`根据这个 `union`在解析器里生成对应的 `struct`是：

```
type yySymType struct {
    yys       int
    offset    int // offset
    item      interface{}
    ident     string
    expr      ast.ExprNode
    statement ast.StmtNode
}
```

在语法解析过程中，`非终结符`会被构造成抽象语法树（`AST`）的节点 [ast.ExprNode ](https://github.com/pingcap/tidb/blob/73900c4890dc9708fe4de39021001ca554bc8374/ast/ast.go#L60)或 [ast.StmtNode ](https://github.com/pingcap/tidb/blob/73900c4890dc9708fe4de39021001ca554bc8374/ast/ast.go#L94)。抽象语法树相关的数据结构都定义在 [ast ](https://github.com/pingcap/tidb/tree/source-code/ast)包中，它们大都实现了 [ast.Node ](https://github.com/pingcap/tidb/blob/73900c4890dc9708fe4de39021001ca554bc8374/ast/ast.go#L29)接口：

```
// Node is the basic element of the AST.
// Interfaces embed Node should have 'Node' name suffix.
type Node interface {
    Accept(v Visitor) (node Node, ok bool)
    Text() string
    SetText(text string)
}
```

这个接口有一个 `Accept`方法，接受 `Visitor`参数，后续对 `AST`的处理，主要依赖这个 `Accept`方法，以 `Visitor`模式遍历所有的节点以及对 `AST`做结构转换。

```
// Visitor visits a Node.
type Visitor interface {
    Enter(n Node) (node Node, skipChildren bool)
    Leave(n Node) (node Node, ok bool)
}
```

例如 [plan.preprocess ](https://github.com/pingcap/tidb/blob/source-code/plan/preprocess.go)是对 `AST`做预处理，包括合法性检查以及名字绑定。

`union`后面是对 `token`和 `非终结符`按照类型分别定义：

```
/* 这部分的 token 是 ident 类型 */
%token    <ident>
    ...
    add            "ADD"
    all             "ALL"
    alter            "ALTER"
    analyze            "ANALYZE"
    and            "AND"
    as            "AS"
    asc            "ASC"
    between            "BETWEEN"
    bigIntType        "BIGINT"
    ...

/* 这部分的 token 是 item 类型 */   
%token    <item>
    /*yy:token "1.%d"   */    floatLit        "floating-point literal"
    /*yy:token "1.%d"   */    decLit          "decimal literal"
    /*yy:token "%d"     */    intLit          "integer literal"
    /*yy:token "%x"     */    hexLit          "hexadecimal literal"
    /*yy:token "%b"     */    bitLit          "bit literal"

    andnot        "&^"
    assignmentEq    ":="
    eq        "="
    ge        ">="
    ...

/* 非终结符按照类型分别定义 */
%type    <expr>
    Expression            "expression"
    BoolPri                "boolean primary expression"
    ExprOrDefault            "expression or default"
    PredicateExpr            "Predicate expression factor"
    SetExpr                "Set variable statement value's expression"
    ...

%type    <statement>
    AdminStmt            "Check table statement or show ddl statement"
    AlterTableStmt            "Alter table statement"
    AlterUserStmt            "Alter user statement"
    AnalyzeTableStmt        "Analyze table statement"
    BeginTransactionStmt        "BEGIN TRANSACTION statement"
    BinlogStmt            "Binlog base64 statement"
    ...
    
%type   <item>
    AlterTableOptionListOpt        "alter table option list opt"
    AlterTableSpec            "Alter table specification"
    AlterTableSpecList        "Alter table specification list"
    AnyOrAll            "Any or All for subquery"
    Assignment            "assignment"
    ...

%type    <ident>
    KeyOrIndex        "{KEY|INDEX}"
    ColumnKeywordOpt    "Column keyword or empty"
    PrimaryOpt        "Optional primary keyword"
    NowSym            "CURRENT_TIMESTAMP/LOCALTIME/LOCALTIMESTAMP"
    NowSymFunc        "CURRENT_TIMESTAMP/LOCALTIME/LOCALTIMESTAMP/NOW"
    ...
```

第一部分的最后是对优先级和结合性的定义：

```
...
%precedence sqlCache sqlNoCache
%precedence lowerThanIntervalKeyword
%precedence interval
%precedence lowerThanStringLitToken
%precedence stringLit
...
%right   assignmentEq
%left     pipes or pipesAsOr
%left     xor
%left     andand and
%left     between
...
```

`parser.y`文件的第二部分是 `SQL`语法的产生式和每个规则对应的 `aciton`。SQL语法非常复杂，`parser.y`的大部分内容都是产生式的定义。

`SQL`语法可以参照 MySQL 参考手册的 [SQL Statements ](https://dev.mysql.com/doc/refman/5.7/en/sql-statements.html)部分，例如 [SELECT ](https://dev.mysql.com/doc/refman/5.7/en/select.html)语法的定义如下：

```
SELECT
    [ALL | DISTINCT | DISTINCTROW ]
      [HIGH_PRIORITY]
      [STRAIGHT_JOIN]
      [SQL_SMALL_RESULT] [SQL_BIG_RESULT] [SQL_BUFFER_RESULT]
      [SQL_CACHE | SQL_NO_CACHE] [SQL_CALC_FOUND_ROWS]
    select_expr [, select_expr ...]
    [FROM table_references
      [PARTITION partition_list]
    [WHERE where_condition]
    [GROUP BY {col_name | expr | position}
      [ASC | DESC], ... [WITH ROLLUP]]
    [HAVING where_condition]
    [ORDER BY {col_name | expr | position}
      [ASC | DESC], ...]
    [LIMIT {[offset,] row_count | row_count OFFSET offset}]
    [PROCEDURE procedure_name(argument_list)]
    [INTO OUTFILE 'file_name'
        [CHARACTER SET charset_name]
        export_options
      | INTO DUMPFILE 'file_name'
      | INTO var_name [, var_name]]
    [FOR UPDATE | LOCK IN SHARE MODE]]
```

我们可以在 `parser.y`中找到 `SELECT`语句的产生式：

```
SelectStmt:
    "SELECT" SelectStmtOpts SelectStmtFieldList OrderByOptional SelectStmtLimit SelectLockOpt
    { ... }
|   "SELECT" SelectStmtOpts SelectStmtFieldList FromDual WhereClauseOptional SelectStmtLimit SelectLockOpt
    { ... }  
|   "SELECT" SelectStmtOpts SelectStmtFieldList "FROM"
    TableRefsClause WhereClauseOptional SelectStmtGroup HavingClause OrderByOptional
    SelectStmtLimit SelectLockOpt
    { ... } 
```

产生式 `SelectStmt`和 `SELECT`语法是对应的。

我省略了大括号中的 `action`，这部分代码会构建出 `AST`的 [ast.SelectStmt ](https://github.com/pingcap/tidb/blob/3ac2b34a3491e809a96db358ee2ce8d11a66abb6/ast/dml.go#L451)节点：

```
type SelectStmt struct {
    dmlNode
    resultSetNode

    // SelectStmtOpts wraps around select hints and switches.
    *SelectStmtOpts
    // Distinct represents whether the select has distinct option.
    Distinct bool
    // From is the from clause of the query.
    From *TableRefsClause
    // Where is the where clause in select statement.
    Where ExprNode
    // Fields is the select expression list.
    Fields *FieldList
    // GroupBy is the group by expression list.
    GroupBy *GroupByClause
    // Having is the having condition.
    Having *HavingClause
    // OrderBy is the ordering expression list.
    OrderBy *OrderByClause
    // Limit is the limit clause.
    Limit *Limit
    // LockTp is the lock type
    LockTp SelectLockType
    // TableHints represents the level Optimizer Hint
    TableHints []*TableOptimizerHint
}
```

可以看出，`ast.SelectStmt`结构体内包含的内容和 `SELECT`语法也是一一对应的。

其他的产生式也都是根据对应的 `SQL`语法来编写的。从 `parser.y`的注释看到，这个文件最初是用 [工具 ](https://github.com/cznic/ebnf2y)从 `BNF`转化生成的，从头手写这个规则文件，工作量会非常大。

完成了语法规则文件 `parser.y`的定义，就可以使用 `goyacc`生成语法解析器：

```
bin/goyacc -o parser/parser.go parser/parser.y 2>&1
```

TiDB 对 `lexer`和 `parser.go`进行了封装，对外提供 [parser.yy_parser ](https://github.com/pingcap/tidb/blob/source-code/plan/preprocess.go)进行 SQL 语句的解析：

```
// Parse parses a query string to raw ast.StmtNode.
func (parser *Parser) Parse(sql, charset, collation string) ([]ast.StmtNode, error) {
    ...
}
```

最后，我写了一个简单的例子，使用 TiDB 的 `SQL Parser`进行 SQL 语法解析，构建出 `AST`，然后利用 `visitor`遍历 `AST`：

```golang
package main

import (
    "fmt"
    "github.com/pingcap/parser"
    "github.com/pingcap/parser/ast"
    _ "github.com/pingcap/tidb/types/parser_driver"
)
type visitor struct{}

func (v *visitor) Enter(in ast.Node) (out ast.Node, skipChildren bool) {
    fmt.Printf("%T\n", in)
    return in, false
}

func (v *visitor) Leave(in ast.Node) (out ast.Node, ok bool) {
    return in, true
}

func main() {
    p := parser.New()

    sql := "SELECT /*+ TIDB_SMJ(employees) */ emp_no, first_name, last_name " +
        "FROM employees USE INDEX (last_name) " +
        "where last_name='Aamodt' and gender='F' and birth_date > '1960-01-01'"
    stmtNodes, _, err := p.Parse(sql, "", "")

    if err != nil {
        fmt.Printf("parse error:\n%v\n%s", err, sql)
        return
    }
    for _, stmtNode := range stmtNodes {
        v := visitor{}
        stmtNode.Accept(&v)
    }
}
```

我实现的 `visitor`什么也没干，只是输出了节点的类型。 这段代码的运行结果如下，依次输出遍历过程中遇到的节点类型：

```golang
*ast.SelectStmt
*ast.TableOptimizerHint
*ast.TableRefsClause
*ast.Join
*ast.TableSource
*ast.TableName
*ast.BinaryOperationExpr
*ast.BinaryOperationExpr
*ast.BinaryOperationExpr
*ast.ColumnNameExpr
*ast.ColumnName
*ast.ValueExpr
*ast.BinaryOperationExpr
*ast.ColumnNameExpr
*ast.ColumnName
*ast.ValueExpr
*ast.BinaryOperationExpr
*ast.ColumnNameExpr
*ast.ColumnName
*ast.ValueExpr
*ast.FieldList
*ast.SelectField
*ast.ColumnNameExpr
*ast.ColumnName
*ast.SelectField
*ast.ColumnNameExpr
*ast.ColumnName
*ast.SelectField
*ast.ColumnNameExpr
*ast.ColumnName
```

了解了 TiDB `SQL Parser`的实现，我们就有可能实现 TiDB 当前不支持的语法，例如添加内置函数，也为我们学习查询计划以及优化打下了基础。希望这篇文章对你能有所帮助。

## Select语句概览

   在先前的 [TiDB 源码阅读系列文章（四） ](https://pingcap.com/blog-cn/tidb-source-code-reading-4/)中，我们介绍了 Insert 语句，想必大家已经了解了 TiDB 是如何写入数据，本篇文章介绍一下 Select 语句是如何执行。相比 Insert，Select 语句的执行流程会更复杂，本篇文章会第一次进入优化器、Coprocessor 模块进行介绍。

### 表结构和语句

表结构沿用上篇文章的：

```sql
CREATE TABLE t {
  id   VARCHAR(31),
  name VARCHAR(50),
  age  int,
  key id_idx (id)
};
```

`Select`语句只会讲解最简单的情况：全表扫描+过滤，暂时不考虑索引等复杂情况，更复杂的情况会在后续章节中介绍。语句为：

```sql
SELECT name FROM t WHERE age > 10;
```

### 语句处理流程

相比 Insert 的处理流程，Select 的处理流程中有 3 个明显的不同：

1. 需要经过 Optimize

   Insert 是比较简单语句，在查询计划这块并不能做什么事情（对于 Insert into Select 语句这种，实际上只对 Select 进行优化），而 Select 语句可能会无比复杂，不同的查询计划之间性能天差地别，需要非常仔细的进行优化。

2. 需要和存储引擎中的计算模块交互

   Insert 语句只涉及对 Key-Value 的 Set 操作，Select 语句可能要查询大量的数据，如果通过 KV 接口操作存储引擎，会过于低效，必须要通过计算下推的方式，将计算逻辑发送到存储节点，就近进行处理。

3. 需要对客户端返回结果集数据

   Insert 语句只需要返回是否成功以及插入了多少行即可，而 Select 语句需要返回结果集。

本篇文章会重点说明这些不同的地方，而相同的步骤会尽量化简。

### Parsing

[Select 语句的语法解析规则 ](https://github.com/pingcap/tidb/blob/source-code/parser/parser.y#L3906)，相比 Insert 语句，要复杂很多，大家可以对着 [MySQL 文档 ](https://dev.mysql.com/doc/refman/5.7/en/select.html)看一下具体的解析实现。需要特别注意的是 From 字段，这里可能会非常复杂，其语法定义是递归的。

最终语句被解析成 [ast.SelectStmt ](https://github.com/pingcap/tidb/blob/source-code/ast/dml.go#L451)结构：

```go
type SelectStmt struct {
        dmlNode
        resultSetNode
        // SelectStmtOpts wraps around select hints and switches.
        *SelectStmtOpts
        // Distinct represents whether the select has distinct option.
        Distinct bool
        // From is the from clause of the query.
        From *TableRefsClause
        // Where is the where clause in select statement.
        Where ExprNode
        // Fields is the select expression list.
        Fields *FieldList
        // GroupBy is the group by expression list.
        GroupBy *GroupByClause
        // Having is the having condition.
        Having *HavingClause
        // OrderBy is the ordering expression list.
        OrderBy *OrderByClause
        // Limit is the limit clause.
        Limit *Limit
        // LockTp is the lock type
        LockTp SelectLockType
        // TableHints represents the level Optimizer Hint
        TableHints [](#)*TableOptimizerHint
}
```

对于本文所提到的语句 `SELECT name FROM t WHERE age > 10; `name 会被解析为 Fields，`WHERE age > 10`被解析为 Where 字段，`FROM t`被解析为 From 字段。

### Planning

在 [planBuilder.buildSelect() ](https://github.com/pingcap/tidb/blob/source-code/plan/logical_plan_builder.go#L1452)方法中，我们可以看到 ast.SelectStmt 是如何转换成一个 plan 树，最终的结果是一个 LogicalPlan，每一个语法元素都被转换成一个逻辑查询计划单元，例如 `WHERE c > 10`会被处理为一个 plan.LogicalSelection 的结构：

```go
    if sel.Where != nil {
        p = b.buildSelection(p, sel.Where, nil)
        if b.err != nil {
            return nil
        }
    }  
```

具体的结构如下：

```go
// LogicalSelection represents a where or having predicate.
type LogicalSelection struct {
    baseLogicalPlan

    // Originally the WHERE or ON condition is parsed into a single expression,
    // but after we converted to CNF(Conjunctive normal form), it can be
    // split into a list of AND conditions.
    Conditions []expression.Expression
}
```

其中最重要的就是这个 Conditions 字段，代表了 Where 语句需要计算的表达式，这个表达式求值结果为 True 的时候，表明这一行符合条件。

其他字段的 AST 转 LogicalPlan 读者可以自行研究一下，经过这个 buildSelect() 函数后，AST 变成一个 Plan 的树状结构树，下一步会在这个结构上进行优化。

### Optimizing

让我们回到 [plan.Optimize() 函数 ](https://github.com/pingcap/tidb/blob/source-code/plan/optimizer.go#L61)，Select 语句得到的 Plan 是一个 LogicalPlan，所以 [这里 ](https://github.com/pingcap/tidb/blob/source-code/plan/optimizer.go#L81)可以进入 doOptimize 这个函数，这个函数比较短，其内容如下：

```go
func doOptimize(flag uint64, logic LogicalPlan) (PhysicalPlan, error) {
    logic, err := logicalOptimize(flag, logic)
    if err != nil {
        return nil, errors.Trace(err)
    }
    if !AllowCartesianProduct && existsCartesianProduct(logic) {
        return nil, errors.Trace(ErrCartesianProductUnsupported)
    }
    physical, err := dagPhysicalOptimize(logic)
    if err != nil {
        return nil, errors.Trace(err)
    }
    finalPlan := eliminatePhysicalProjection(physical)
    return finalPlan, nil
}
```

大家可以关注两个步骤：logicalOptimize 和 dagPhysicalOptimize，分别代表逻辑优化和物理优化，这两种优化的基本概念和区别本文不会描述，请大家自行研究（这个是数据库的基础知识）。下面分别介绍一下这两个函数做了什么事情。

#### 逻辑优化

逻辑优化由一系列优化规则组成，对于这些规则会按顺序不断应用到传入的 LogicalPlan Tree 中，见 [logicalOptimize() 函数 ](https://github.com/pingcap/tidb/blob/source-code/plan/optimizer.go#L131)：

```go
func logicalOptimize(flag uint64, logic LogicalPlan) (LogicalPlan, error) {
    var err error
    for i, rule := range optRuleList {
        // The order of flags is same as the order of optRule in the list.
        // We use a bitmask to record which opt rules should be used. If the i-th bit is 1, it means we should
        // apply i-th optimizing rule.
        if flag&(1<<uint(i)) == 0 {
            continue
        }
        logic, err = rule.optimize(logic)
        if err != nil {
            return nil, errors.Trace(err)
        }
    }
    return logic, errors.Trace(err)
}
```

目前 TiDB 已经支持下列优化规则：

```go
var optRuleList = []logicalOptRule{
    &columnPruner{}, 
    &maxMinEliminator{},
    &projectionEliminater{},
    &buildKeySolver{},
    &decorrelateSolver{},
    &ppdSolver{},
    &aggregationOptimizer{},
    &pushDownTopNOptimizer{},
}
```

这些规则并不会考虑数据的分布，直接无脑的操作 Plan 树，因为大多数规则应用之后，一定会得到更好的 Plan（不过上面有一个规则并不一定会更好，读者可以想一下是哪个）。

这里选一个规则介绍一下，其他优化规则请读者自行研究或者是等待后续文章。

columnPruner（列裁剪） 规则，会将不需要的列裁剪掉，考虑这个 SQL: `select c from t;`对于 `from t`这个全表扫描算子（也可能是索引扫描）来说，只需要对外返回 c 这一列的数据即可，这里就是通过列裁剪这个规则实现，整个 Plan 树从树根到叶子节点递归调用这个规则，每层节点只保留上面节点所需要的列即可。

经过逻辑优化，我们可以得到这样一个查询计划：

![logical-select](tidb_parser.assets/1_9f9e92bad2.png)

其中 `FROM t`变成了 DataSource 算子，`WHERE age > 10`变成了 Selection 算子，这里留一个思考题，`SELECT name`中的列选择去哪里了？

#### 物理优化

在物理优化阶段，会考虑数据的分布，决定如何选择物理算子，比如对于 `FROM t WHERE age > 10`这个语句，假设在 age 字段上有索引，需要考虑是通过 TableScan + Filter 的方式快还是通过 IndexScan 的方式比较快，这个选择取决于统计信息，也就是 age > 10 这个条件究竟能过滤掉多少数据。

我们看一下 [dagPhysicalOptimize ](https://github.com/pingcap/tidb/blob/source-code/plan/optimizer.go#L148)这个函数：

```go
func dagPhysicalOptimize(logic LogicalPlan) (PhysicalPlan, error) {
    logic.preparePossibleProperties()
    logic.deriveStats()
    t, err := logic.convert2PhysicalPlan(&requiredProp{taskTp: rootTaskType, expectedCnt: math.MaxFloat64})
    if err != nil {
        return nil, errors.Trace(err)
    }
    p := t.plan()
    p.ResolveIndices()
    return p, nil
}
```

这里的 convert2PhysicalPlan 会递归调用下层节点的 convert2PhysicalPlan 方法，生成物理算子并且估算其代价，然后从中选择代价最小的方案，这两个函数比较重要：

```go
// convert2PhysicalPlan implements LogicalPlan interface.
func (p *baseLogicalPlan) convert2PhysicalPlan(prop *requiredProp) (t task, err error) {
    // Look up the task with this prop in the task map.
    // It's used to reduce double counting.
    t = p.getTask(prop)
    if t != nil {
        return t, nil
    }
    t = invalidTask
    if prop.taskTp != rootTaskType {
        // Currently all plan cannot totally push down.
        p.storeTask(prop, t)
        return t, nil
    }
    for _, pp := range p.self.genPhysPlansByReqProp(prop) {
        t, err = p.getBestTask(t, pp)
        if err != nil {
            return nil, errors.Trace(err)
        }
    }
    p.storeTask(prop, t)
    return t, nil
}

func (p *baseLogicalPlan) getBestTask(bestTask task, pp PhysicalPlan) (task, error) {
    tasks := make([]task, 0, len(p.children))
    for i, child := range p.children {
        childTask, err := child.convert2PhysicalPlan(pp.getChildReqProps(i))
        if err != nil {
            return nil, errors.Trace(err)
        }
        tasks = append(tasks, childTask)
    }
    resultTask := pp.attach2Task(tasks...)
    if resultTask.cost() < bestTask.cost() {
        bestTask = resultTask
    }
    return bestTask, nil
}
```

上面两个方法的返回值都是一个叫 task 的结构，而不是物理计划，这里引入一个概念，叫 **`Task`**，TiDB 的优化器会将 PhysicalPlan 打包成为 Task。Task 的定义在 [task.go ](https://github.com/pingcap/tidb/blob/source-code/plan/task.go)中，我们看一下注释：

```go
// task is a new version of `PhysicalPlanInfo`. It stores cost information for a task.
// A task may be CopTask, RootTask, MPPTask or a ParallelTask.
type task interface {
    count() float64
    addCost(cost float64)
    cost() float64
    copy() task
    plan() PhysicalPlan
    invalid() bool
}
```

在 TiDB 中，Task 的定义是能在单个节点上不依赖于和其他节点进行数据交换即可进行的一系列操作，目前只实现了两种 Task：

- CopTask 是需要下推到存储引擎（TiKV）上进行计算的物理计划，每个收到请求的 TiKV 节点都会做相同的操作
- RootTask 是保留在 TiDB 中进行计算的那部分物理计划

如果了解过 TiDB 的 Explain 结果，那么可以看到每个 Operator 都会标明属于哪种 Task，比如下面这个例子：

![explain](tidb_parser.assets/2_5b4b914b8e.jpg)

整个流程是一个树形动态规划的算法，大家有兴趣可以跟一下相关的代码自行研究或者等待后续的文章。

经过整个优化过程，我们已经得到一个物理查询计划，这个 `SELECT name FROM t WHERE age > 10;`语句能够指定出来的查询计划大概是这样子的：

![simple-select](tidb_parser.assets/3_fa22deaac0.png)

读者可能会比较奇怪，为什么只剩下这样一个物理算子？`WHERR age > 10`哪里去了？实际上 age > 10 这个过滤条件被合并进了 PhysicalTableScan，因为 `age > 10`这个表达式可以下推到 TiKV 上进行计算，所以会把 TableScan 和 Filter 这样两个操作合在一起。哪些表达式会被下推到 TiKV 上的 Coprocessor 模块进行计算呢？对于这个 Query 是在下面 [这个地方 ](https://github.com/pingcap/tidb/blob/source-code/plan/predicate_push_down.go#L72)进行识别：

```go
// PredicatePushDown implements LogicalPlan PredicatePushDown interface.
func (ds *DataSource) PredicatePushDown(predicates []expression.Expression) ([]expression.Expression, LogicalPlan) {
    _, ds.pushedDownConds, predicates = expression.ExpressionsToPB(ds.ctx.GetSessionVars().StmtCtx, predicates, ds.ctx.GetClient())
    return predicates, ds
}
```

在 `expression.ExpressionsToPB`这个方法中，会把能下推 TiKV 上的表达式识别出来（TiKV 还没有实现所有的表达式，特别是内建函数只实现了一部分），放到 DataSource.pushedDownConds 字段中。接下来我们看一下 DataSource 是如何转成 PhysicalTableScan，见 [DataSource.convertToTableScan() ](https://github.com/pingcap/tidb/blob/source-code/plan/physical_plan_builder.go#L523)方法。这个方法会构建出 PhysicalTableScan，并且调用 [addPushDownSelection() ](https://github.com/pingcap/tidb/blob/source-code/plan/physical_plan_builder.go#L610)方法，将一个 PhysicalSelection 加到 PhysicalTableScan 之上，一起放进 copTask 中。

这个查询计划是一个非常简单的计划，不过我们可以用这个计划来说明 TiDB 是如何执行查询操作。

### Executing

一个查询计划如何变成一个可执行的结构以及如何驱动这个结构执行查询已经在前面的两篇文章中做了描述，这里不再敷述，这一节我会重点介绍具体的执行过程以及 TiDB 的分布式执行框架。

#### Coprocessor 框架

Coprocessor 这个概念是从 HBase 中借鉴而来，简单来说是一段注入在存储引擎中的计算逻辑，等待 SQL 层发来的计算请求（序列化后的物理执行计划），处理本地数据并返回计算结果。在 TiDB 中，计算是以 Region 为单位进行，SQL 层会分析出要处理的数据的 Key Range，再将这些 Key Range 根据 PD 中拿到的 Region 信息划分成若干个 Key Range，最后将这些请求发往对应的 Region。

SQL 层会将多个 Region 返回的结果进行汇总，再经过所需的 Operator 处理，生成最终的结果集。

##### DistSQL

请求的分发与汇总会有很多复杂的处理逻辑，比如出错重试、获取路由信息、控制并发度以及结果返回顺序，为了避免这些复杂的逻辑与 SQL 层耦合在一起，TiDB 抽象了一个统一的分布式查询接口，称为 DistSQL API，位于 [distsql ](https://github.com/pingcap/tidb/blob/source-code/distsql/distsql.go)这个包中。

其中最重要的方法是 [SelectDAG ](https://github.com/pingcap/tidb/blob/source-code/distsql/distsql.go#L305)这个函数：

```go
// SelectDAG sends a DAG request, returns SelectResult.
// In kvReq, KeyRanges is required, Concurrency/KeepOrder/Desc/IsolationLevel/Priority are optional.
func SelectDAG(goCtx goctx.Context, ctx context.Context, kvReq *kv.Request, fieldTypes []*types.FieldType) (SelectResult, error) {
    // kvReq 中包含了计算所涉及的数据的 KeyRanges
    // 这里通过 TiKV Client 向 TiKV 集群发送计算请求
    resp := ctx.GetClient().Send(goCtx, kvReq)
    if resp == nil {
        err := errors.New("client returns nil response")
        return nil, errors.Trace(err)
    }

    if kvReq.Streaming {
        return &streamResult{
            resp:       resp,
            rowLen:     len(fieldTypes),
            fieldTypes: fieldTypes,
            ctx:        ctx,
        }, nil
    }
    // 这里将结果进行了封装
    return &selectResult{
        label:      "dag",
        resp:       resp,
        results:    make(chan newResultWithErr, kvReq.Concurrency),
        closed:     make(chan struct{}),
        rowLen:     len(fieldTypes),
        fieldTypes: fieldTypes,
        ctx:        ctx,
    }, nil
}
```

TiKV Client 中的具体逻辑我们暂时跳过，这里只关注 SQL 层拿到了这个 `selectResult`后如何读取数据，下面这个接口是关键。

```go
// SelectResult is an iterator of coprocessor partial results.
type SelectResult interface {
    // NextRaw gets the next raw result.
    NextRaw(goctx.Context) ([]byte, error)
    // NextChunk reads the data into chunk.
    NextChunk(goctx.Context, *chunk.Chunk) error
    // Close closes the iterator.
    Close() error
    // Fetch fetches partial results from client.
    // The caller should call SetFields() before call Fetch().
    Fetch(goctx.Context)
    // ScanKeys gets the total scan row count.
    ScanKeys() int64
```

selectResult 实现了 SelectResult 这个接口，代表了一次查询的所有结果的抽象，计算是以 Region 为单位进行，所以这里全部结果会包含所有涉及到的 Region 的结果。调用 Chunk 方法可以读到一个 Chunk 的数据，通过不断调用 NextChunk 方法，直到 Chunk 的 NumRows 返回 0 就能拿到所有结果。NextChunk 的实现会不断获取每个 Region 返回的 SelectResponse，把结果写入 Chunk。

##### Root Executor

能推送到 TiKV 上的计算请求目前有 TableScan、IndexScan、Selection、TopN、Limit、PartialAggregation 这样几个，其他更复杂的算子，还是需要在单个 tidb-server 上进行处理。所以整个计算是一个多 tikv-server 并行处理 + 单个 tidb-server 进行汇总的模式。

### 总结

Select 语句的处理过程中最复杂的地方有两点，一个是查询优化，一个是如何分布式地执行，这两部分后续都会有文章来更进一步介绍。下一篇文章会脱离具体的 SQL 逻辑，介绍一下如何看懂某一个特定的模块。

