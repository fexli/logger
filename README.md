# logger

### 一个简单、开源、跨平台的日志库，兼容大部分的操作系统。

## 安装

```bash
go get -u github.com/fexli/logger
```

## 开始使用

### 基础输出

如果您仅将日志库用于显示基础日志，那么可以直接使用内置的`logger.RootLogger`作为默认的日志输出器：

```go
package main

import ( // 通过import导入日志库
	"github.com/fexli/logger"
)

func main() {
	// 直接通过默认的日志输出器输出日志
	logger.RootLogger.Debug(logger.WithContent("hello world"))
}

// Output: [00:00:00.000]<sys>hello world
```

> 说明：日志输出同时支持模块调用与链式调用来输出日志内容，如：`logger.RootLogger.NewDebugLog("hello world").Commit()`
> ，请注意若使用链式调用，则需要在最后一个调用后使用`Commit()`方法来提交并输出日志。

### 创建新的日志记录器

通过`logger.GetLogger`
可以获取指定名称的日志记录器，若记录器不存在，则会自动创建，GetLogger接收两个参数：第一个参数作为logger的名称，第二个参数作为该logger是否需要显示日志输出指针（如果不需要，则可以传入false）。

```go
package main

import "github.com/fexli/logger"

func main() {
	ObjLogger := logger.GetLogger("Object", true)
	ObjLogger.Debug(logger.WithContent("hello world"))
}

// Output: [00:00:00.000]("main.go",in main.main line 7)<Object>[DEBUG]hello world
```

...TODO

### 日志记录等级

区别于传统日志等级划分，本日志将记录等级分为：`DEBUG`, `HELP`, `COMMON`, `SYSTEM`, `WARNING`, `ERROR`,`NOTICE`以及`FATAL`：

| 等级      | 显示名称    | 触发函数      | 产生记录* | 强制显示日志位置** |
|---------|---------|-----------|-------|------------|
| DEBUG   | DEBUG   | Debug()   | 是     | 否          |
| HELP    | HELP    | Help()    | 否     | 否          |
| COMMON  | COMMON  | Common()  | 是     | 否          |
| SYSTEM  | SYSTEM  | System()  | 是     | 否          |
| WARNING | WARNING | Warning() | 是     | 否          |
| ERROR   | ERROR   | Error()   | 是     | 否          |
| FATAL   | FATAL   | Fatal()   | 是     | 否          |
| NOTICE  | NOTICE  | Notice()  | 是     | 否          |
| -       | -       | Log()***  | 是     | 否          |

##### *: HELP等级的日志只会输出到控制台，不会被logger记录，因此无法被历史记录(`GetLogs()`)获取。

##### **:当日志记录器设置为关闭显示记录位置时，强制显示日志位置会无视该设置。

##### ***:直接使用Log()需要搭配WithLevel()来设置日志等级。

### 日志记录事务

根据日志记录需求，日志中的每个事务都作为输出的部分：

| 函数                        | 参数               | 描述                                                                     |
|---------------------------|------------------|------------------------------------------------------------------------|
| WithContent()             | 任何参数             | 默认的参数传递函数，用来传递任何输出的内容，接收的参数若为结构体将被fmt进行转换后输出，若为LogTextCtx则会直接输出带颜色的内容。 |
| WithBacktraceLevelDelta() | 回撤等级(int)        | 用于在显示日志位置时显示回撤的等级（如跳出defer显示触发panic的位置、显示被调函数调用者等）                     |
| JmpOutDefer()             | 无                | 等效于`WithBacktraceLevelDelta(4)`                                        |
| WithLog()                 | 是否输出到控制台(bool)   | 用于控制该条内容是否被输出到控制台，无论输出与否都会被记录到历史记录（HELP除外）                             |
| WithSep()                 | 分隔符(string)      | 用于控制由`WithContent`带来的所有参数间的分隔符，默认为空格(` `)                              |
| WithEnd()                 | 结束符(string)      | 用于控制整行日志输出后的结束符，默认为LF(`\n`)                                            |
| WithLevel()               | 日志等级(LogLevel)   | 用于控制日志输出的等级，仅在使用`logger.Log()`时有效                                      |
| WithLog2Logs()            | 是否记录到历史(bool)    | 用于控制日志是否被记录到该记录器的历史记录中                                                 |
| WithCur()                 | 显示指定指针(任意对象指针接口) | 用于显示某个对象的指针，可用于追踪对象的同一性                                                |
| WithStruct()              | 显示指定对象(任意对象指针接口) | 用于显示某个对象的KV值                                                           |

```go
package main

import "github.com/fexli/logger"

func ppp() {
	panic("test panic")
}

type Test struct {
	Name string
	Age  int
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			// 通过JmpOutDefer跳出defer，直接追踪panic的调用位置
			logger.RootLogger.Fatal(logger.WithContent(err), logger.JmpOutDefer())
		}
	}()
	t := Test{
		"test",
		1,
	}
	// 通过WithContent显示任意的内容，同时通过WithSep控制每个字段间的间隔符
	logger.RootLogger.Debug(logger.WithContent("TestT", &t, "Ends"), logger.WithSep("~~~"))
	ppp()
}

// Outputs: 
//[00:00:00.000]<sys>[DEBUG]TestT~~~&{Name:test Age:1}~~~Ends
//[00:00:00.000]("main.go",in main.ppp line 6)<sys>[FATAL]test panic
```

### 获取历史记录

每个记录器都有一个历史记录，可以通过`GetLogs()`获取，GetLogs函数接收三个参数，分别是日志开始时间，日志等级筛选，最大获取日志量。

```go
package main

import (
	"github.com/fexli/logger"
	"encoding/json"
)

func main() {
	for i := 0; i < 10; i++ {
		logger.RootLogger.System(logger.WithContent("Current I=", i))
	}
	logs := logger.RootLogger.GetLogs(0, logger.LevelDefault, 5)
	// 获取从0时刻开始，默认等级（即：不筛选任何日志），最多获取5条日志（倒叙，获取最新的5条）的历史记录
	data, _ := json.Marshal(logs)
	logger.RootLogger.Debug(logger.WithContent(data))
}

// Outputs:
// [00:00:00.000]<sys>[SYSTEM]Current I= 0
// [00:00:00.000]<sys>[SYSTEM]Current I= 1
// [00:00:00.000]<sys>[SYSTEM]Current I= 2
// [00:00:00.000]<sys>[SYSTEM]Current I= 3
// [00:00:00.000]<sys>[SYSTEM]Current I= 4
// [00:00:00.000]<sys>[SYSTEM]Current I= 5
// [00:00:00.000]<sys>[SYSTEM]Current I= 6
// [00:00:00.000]<sys>[SYSTEM]Current I= 7
// [00:00:00.000]<sys>[SYSTEM]Current I= 8
// [00:00:00.000]<sys>[SYSTEM]Current I= 9
// [00:00:00.000]<sys>[DEBUG][{"ts":1660287551.24,"level":8,"info":{"log":"","color":null,"inner":[{"log":"Current I= 5","color":null,"inner":null}]}},{"ts":1660287551.241,"level":8,"info":{"log":"","color":null,"inner":[{"log":"Current I= 6","color":null,"inner":null}]}},{"ts":1660287551.241,"level":8,"info":{"log":"","color":null,"inner":[{"log":"Current I= 7","color":null,"inner":null}]}},{"ts":1660287551.241,"level":8,"info":{"log":"","color":null,"inner":[{"log":"Current I= 8","color":null,"inner":null}]}},{"ts":1660287551.241,"level":8,"info":{"log":"","color":null,"inner":[{"log":"Current I= 9","color":null,"inner":null}]}}]
```

> 如果需要清除所有日志，请使用`ClearLogInfo()`清除

## 进阶使用

### 全局日志等级筛选(GlobLogFilter)

### 自定义输出函数(Printer)

### 文件日志记录(InitGlobLog)

### 色彩日志(LogTextCtx)

### 色彩系统(logcolor)

### 色彩日志适配表格(colorableStdout)

...TODO

## 系统适配重定向表格

| X           | aix              | android | darwin | dragonfly | freebsd | hurd | illumos | ios | js  | linux | nacl | netbsd | openbsd | plan9 | solaris | windows   | zos |
|-------------|------------------|---------|--------|-----------|---------|------|---------|-----|-----|-------|------|--------|---------|-------|---------|-----------|-----|
| 386         |                  |         |        |           | 2       |      |         |     |     | 2     |      | 2      | 2       |       |         | "DLLCALL" |     |     |     |     |     |     |     |     |     |     |     |     |
| amd64       |                  |         | 2      | 2         | 2       |      |         |     |     | 2     |      | 2      | 2       |       | 2@      |           |     |     |     |     |     |     |     |     |     |     |
| amd64p32    |                  |         |        |           |         |      |         |     |     |       |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| arm         |                  |         |        |           | 2       |      |         |     |     | 2 !3  |      | 2      | 2       |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |
| arm64       |                  |         | 2      |           | 2       |      |         |     |     | !3    |      | 2      | 2       |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |
| arm64be     |                  |         |        |           |         |      |         |     |     |       |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| armbe       |                  |         |        |           |         |      |         |     |     |       |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| mips        |                  |         |        |           |         |      |         |     |     | 2     |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| mips64      |                  |         |        |           |         |      |         |     |     | 2     |      |        | 2       |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| mips64le    |                  |         |        |           |         |      |         |     |     | 2     |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| mips64p32   |                  |         |        |           |         |      |         |     |     |       |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| mips64p32le |                  |         |        |           |         |      |         |     |     |       |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| mipsle      |                  |         |        |           |         |      |         |     |     | 2     |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| ppc         | 2@               |         |        |           |         |      |         |     |     |       |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| ppc64       | 2                |         |        |           |         |      |         |     |     | 2     |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| ppc64le     |                  |         |        |           |         |      |         |     |     | 2     |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| riscv       |                  |         |        |           |         |      |         |     |     |       |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| riscv64     |                  |         |        |           |         |      |         |     |     |       |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| s390        |                  |         |        |           |         |      |         |     |     |       |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| s390x       |                  |         |        |           |         |      |         |     |     | 2     |      |        |         |       |         |           | 2@  |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| sparc       |                  |         |        |           |         |      |         |     |     |       |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| sparc64     |                  |         |        |           |         |      |         |     |     |       |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |
| wasm        |                  |         |        |           |         |      |         |     | X   |       |      |        |         |       |         |           |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |     |

- X= ENOSYS功能未实现
- @= golang.org/x/sys 未来可期

## 鸣谢

* [<img alt="skadiD" src="https://avatars.githubusercontent.com/u/50161715?v=4" style="border-radius: 25px;width: 50px">](https://github.com/skadiD)
  测试并纠正arm下linux系统日志重定向问题

* [<img alt="Ooooooutdated" src="https://avatars.githubusercontent.com/u/79091449?v=4" style="border-radius: 25px;width: 50px">](https://github.com/wg138940)
  链式调用日志记录器与记录查找优化
