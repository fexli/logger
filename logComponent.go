package logger

import (
	"fmt"
	"github.com/fatih/structs"
)

//如何向func传递默认值

type logOptions struct {
	Info                []LogCtx
	Level               LogLevel
	BacktraceLevelDelta int
	Log                 bool
	Log2Logs            bool
	Sep                 string
	End                 string
	Cur                 string
}

type LogComponent interface {
	apply(*logOptions)
}
type logOption struct {
	f func(*logOptions)
}

func (f *logOption) apply(do *logOptions) {
	f.f(do)
}
func newFuncOption(f func(*logOptions)) *logOption {
	return &logOption{
		f: f,
	}
}
func WithContent(ctx ...LogCtx) LogComponent {
	return newFuncOption(func(o *logOptions) {
		o.Info = append(o.Info, ctx...)
	})
}
func WithBacktraceLevelDelta(level int) LogComponent {
	return newFuncOption(func(o *logOptions) {
		o.BacktraceLevelDelta += level
	})
}
func JmpOutDefer() LogComponent {
	return newFuncOption(func(o *logOptions) {
		o.BacktraceLevelDelta += 4
	})
}
func WithLog(log bool) LogComponent {
	return newFuncOption(func(o *logOptions) {
		o.Log = log
	})
}
func WithSep(sep string) LogComponent {
	return newFuncOption(func(o *logOptions) {
		o.Sep = sep
	})
}
func WithEnd(end string) LogComponent {
	return newFuncOption(func(o *logOptions) {
		o.End = end
	})
}
func WithLevel(level LogLevel) LogComponent {
	return newFuncOption(func(o *logOptions) {
		o.Level = level
	})
}
func WithLog2Logs(do bool) LogComponent {
	return newFuncOption(func(o *logOptions) {
		o.Log2Logs = do
	})
}
func WithCur(cur interface{}) LogComponent {
	return newFuncOption(func(o *logOptions) {
		o.Cur = fmt.Sprintf("%p", cur)
	})
}

//WithKVs 依据 (key1 string, value1 interface, key2 string, value2 interface{}...)参数生成结构化日志.
//
//Example:
//
//WithKVs("Field 'taskId' found nil in response!", "URL", req.url, "code", resp.StatusCode, "raw", string(respBytes))
func WithKVs(keyValues ...interface{}) LogComponent {
	n := len(keyValues)
	if n%2 != 0 {
		keyValues = append(keyValues, "<Empty arg>")
		n += 1
	}

	result := ""
	for i := 0; i < n; i += 2 {
		result += formatKV(keyValues[i], keyValues[i+1])
	}

	return WithContent(result)
}

//WithStruct 打印一个结构体.
//
//NOTE: 只能访问公共字段.
func WithStruct(s interface{}) LogComponent {
	structMap := structs.Map(s)

	result := ""
	for key, value := range structMap {
		result += formatKV(key, value)
	}

	return WithContent(result)
}

func formatKV(key interface{}, value interface{}) string {
	return fmt.Sprintf("\n\t- %-10v= %v", key, value)
}

func defaultOptions() logOptions {
	return logOptions{
		Info:                nil,
		Level:               LevelCommon,
		BacktraceLevelDelta: 0,
		Log:                 true,
		Sep:                 " ",
		End:                 "",
		Cur:                 "",
		Log2Logs:            true,
	}
}
