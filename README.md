# logger

### 一个简单、开源、跨平台的日志库，兼容大部分的操作系统。

## 使用

如果您仅将日志库用于显示基础日志，那么可以直接使用：

```go
import (
"github.com/fexli/logger"
)

func main() {
logger.RootLogger.Debug(logger.WithContent("hello world"))
}
```

编译并运行上述程序，会输出：

```
```bash
[00:00:00.000]<sys>hello world
```

也可以通过链式调用来输出日志内容，如：

```go
logger.RootLogger.NewDebugLog("hello world").Commit()
```

> 通过`logger.GetLogger`可以获取指定名称的日志记录器，若记录器不存在，则会自动创建。

...TODO

## 鸣谢
* [<img alt="skadiD" src="https://avatars.githubusercontent.com/u/50161715?v=4" style="border-radius: 25px;width: 50px">](https://github.com/skadiD) 测试并纠正arm下linux系统日志重定向问题

*  [<img alt="Ooooooutdated" src="https://avatars.githubusercontent.com/u/79091449?v=4" style="border-radius: 25px;width: 50px">](https://github.com/wg138940) 链式调用日志记录器与记录查找优化
