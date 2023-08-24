package logger

import (
	"container/list"
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/fexli/logger/logcolor"
	"github.com/fexli/logger/utils"
	"github.com/modern-go/reflect2"
	"github.com/xo/terminfo"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"
	"unsafe"
)

var (
	emptyCurInfo = &CurInfo{
		Function: "unknown",
		FileName: "unknown",
		FilePath: "unknown",
		Line:     0,
	}
	colorableStdout            = logcolor.Colorable(os.Stdout)
	GlobalFileHandler *os.File = nil

	pool                = make(map[string]*Logger)
	RootLogger  *Logger = nil
	LogPrefix           = make(map[LogLevel]*logcolor.LogTextCtx)
	TimeColor           = logcolor.NewColor(logcolor.RGB(127, 255, 237))
	MemCurColor         = logcolor.NewColor(logcolor.RGB(203, 127, 255))
	StackColor          = logcolor.NewColor(logcolor.RGB(168, 209, 135))

	EnableGlobLog = false
	GlobLogFilter = LevelDefault
)

const (
	LevelFatal   LogLevel = 1 << 7
	LevelNotice  LogLevel = 1 << 6
	LevelError   LogLevel = 1 << 5
	LevelWarning LogLevel = 1 << 4
	LevelSystem  LogLevel = 1 << 3
	LevelCommon  LogLevel = 1 << 2
	LevelHelp    LogLevel = 1 << 1
	LevelDebug   LogLevel = 1 << 0

	LevelDefault = ^LogLevel(0)
	LevelShowcur = LevelFatal | LevelError
)

func init() {
	// 内部日志Logger初始化

	RootLogger = GetLogger("sys", false)

	LogPrefix = map[LogLevel]*logcolor.LogTextCtx{
		LevelFatal:   logcolor.ColorString("FATL", logcolor.NewColor(logcolor.TextBlack, logcolor.OpInverse)),
		LevelNotice:  logcolor.LightCyanString("NOTE"),
		LevelError:   logcolor.LightRedString("EROR"),
		LevelWarning: logcolor.LightMagentaString("WARN"),
		LevelSystem:  logcolor.LightGreenString("SYST"),
		LevelCommon:  logcolor.WhiteString("INFO"),
		LevelHelp:    logcolor.LightYellowString("HELP"),
		LevelDebug:   logcolor.LightBlueString("DBUG"),
	}
}

// DisableColor 禁用日志颜色
func DisableColor() {
	colorableStdout.DisableColor()
}

// EnableColor 启用日志颜色
func EnableColor() {
	colorableStdout.EnableColor()
}

// ForceSetColor 强制设置日志颜色
func ForceSetColor(colorMode terminfo.ColorLevel) {
	colorableStdout.ForceSetColor(colorMode)
}

// SetGlobLogFilter 设置全局日志记录等级，默认为LevelDefault，即记录所有等级到日志文件，此项目受到Logger本身logLevel限制
//
// e.g.
//
//	logger.SetGlobLogFilter(logger.LevelFatal | logger.LevelError) // 设置全局日志记录等级为Fatal和Error
func SetGlobLogFilter(filter LogLevel) {
	GlobLogFilter = filter
}

// _pause
func _pause() {
	print("请勿在同一目录下启动多个终端，输入回车以退出...\n")
	b := make([]byte, 1)
	_, _ = os.Stdin.Read(b)
}

// InitGlobLog 初始化全局日志，name为日志文件名，如果name为空，则使用默认文件名，可选logDesc为日志文件描述
//
// e.g.
//
//	logger.InitGlobLog("globlog.log", "awesomeProgram v0.1")
func InitGlobLog(name string, logDesc ...string) {
	if EnableGlobLog {
		return
	}
	EnableGlobLog = true
	if name == "" {
		name = "globlog.log"
	}
	_ = os.Mkdir("logs", 0764)
	if info, err := os.Stat(name); info != nil && err == nil {
		if err = os.Rename(name, path.Join("logs", info.ModTime().Format("2006-01-02-15-04")+"."+utils.RandomStr(2, false, "")+".log")); err != nil {
			_pause()
			os.Exit(0xD01)
		}
	}
	file, e := os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_SYNC, 0764)
	if e != nil {
		_pause()
		os.Exit(0xD00)
	}
	if len(logDesc) > 0 {
		_, err := file.Write([]byte(logDesc[0] + " Running Log [Started At " + time.Now().String() + "]\n"))
		if err != nil {
			_pause()
			os.Exit(0xD02)
		}
	}
	GlobalFileHandler = file

	initErr()
}

func defaultConfig() Config {
	return Config{
		Name: "globlog.log",
		Desc: "",
		OldLogPath: func(info os.FileInfo) string {
			return path.Join("logs", info.ModTime().Format("2006-01-02-15-04")+"."+utils.RandomStr(2, false, "")+".log")
		},
	}
}

// InitGlobLogWithConfig 初始化全局日志，name为日志文件名，如果name为空，则使用默认文件名，可选logDesc为日志文件描述
//
// e.g.
//
//	logger.InitGlobLogWithConfig("globlog.log", "awesomeProgram v0.1")
func InitGlobLogWithConfig(config ...Config) {
	if EnableGlobLog {
		return
	}
	EnableGlobLog = true
	var current Config = defaultConfig()
	if len(config) != 0 {
		if config[0].Name != "" {
			current.Name = config[0].Name
		}
		if config[0].Desc != "" {
			current.Desc = config[0].Desc
		}
		if config[0].OldLogPath != nil {
			current.OldLogPath = config[0].OldLogPath
		}
		if config[0].MaxLogTime != 0 {
			current.MaxLogTime = config[0].MaxLogTime
		}
		if len(config[0].SliceWhen) != 0 {
			current.SliceWhen = config[0].SliceWhen
		}
	}
	_ = os.Mkdir("logs", 0764)
	if info, err := os.Stat(current.Name); info != nil && err == nil {
		if err = os.Rename(current.Name, current.OldLogPath(info)); err != nil {
			_pause()
			os.Exit(0xD01)
		}
	}
	file, e := os.OpenFile(current.Name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_SYNC, 0764)
	if e != nil {
		_pause()
		os.Exit(0xD00)
	}
	if current.Desc != "" {
		_, err := file.Write([]byte(current.Desc + " Running Log [Started At " + time.Now().String() + "]\n"))
		if err != nil {
			_pause()
			os.Exit(0xD02)
		}
	}
	GlobalFileHandler = file
	if current.MaxLogTime != 0 || len(current.SliceWhen) != 0 {
		go sliceGlobByTime(current)
	}
	initErr()
}

func sliceGlobByTime(current Config) {
	for {
		next := findNextByWhen(current)
		time.Sleep(next)
		info, err := os.Stat(current.Name)
		if err != nil {
			continue
		}
		if GlobalFileHandler != nil {
			fr := GlobalFileHandler
			GlobalFileHandler = nil
			fr.Close()
		}
		if err = os.Rename(current.Name, current.OldLogPath(info)); err != nil {
			continue
		}
		file, e := os.OpenFile(current.Name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY|os.O_SYNC, 0764)
		if e != nil {
			continue
		}
		if current.Desc != "" {
			_, err := file.Write([]byte(current.Desc + " Running Log [Started At " + time.Now().String() + "]\n"))
			if err != nil {
				continue
			}
		}
		GlobalFileHandler = file
	}
}

// LogLevel 日志等级
type (
	LogLevel   uint8
	LogPrinter func(info *LoggInfo)
	LogCtx     interface{}
)

// CurInfo 记录日志调用栈信息
type CurInfo struct {
	Function string `json:"function"`
	Line     int    `json:"line"`
	FilePath string `json:"filePath"`
	FileName string `json:"fileName"`
}

// LoggInfo 记录日志信息
type LoggInfo struct {
	Ts     float64              `json:"ts"`
	Cur    *CurInfo             `json:"-"`
	Level  LogLevel             `json:"level"`
	MemCur string               `json:"-"`
	Info   *logcolor.LogTextCtx `json:"info"`
}

// Logger 日志类结构体
type Logger struct {
	Name        string
	logs        []*LoggInfo
	printer     *list.List
	keepPrinter bool
	logLevel    LogLevel
	logShowCur  bool
	DefaultIO   LogPrinter
	latestTs    float64
}

////////////////////////////////////////////////////////////////////////////////
// Functions

func GetCurInfo(backLevel int) *CurInfo {
	if backLevel < 0 {
		return emptyCurInfo
	}
	pcs := make([]uintptr, backLevel+2)
	_ = runtime.Callers(2, pcs) // default skip back
	frames := runtime.CallersFrames(pcs)
	var cnt = 0
	for f, again := frames.Next(); again; f, again = frames.Next() {
		if cnt == backLevel {
			return &CurInfo{
				Function: f.Function,
				Line:     f.Line,
				FilePath: path.Dir(f.File),
				FileName: path.Base(f.File),
			}
		}
		cnt++
	}
	return emptyCurInfo
}

func GetLogger(name string, showCur bool) *Logger {
	get := pool[name]

	if get == nil {
		current := &Logger{
			Name:        name,
			keepPrinter: true,
			logs:        make([]*LoggInfo, 0),
			printer:     list.New(),
			logLevel:    LevelDefault,
			logShowCur:  showCur,
			DefaultIO:   nil,
		}
		current.DefaultIO = current.internalPrinter
		pool[name] = current
		return current
	}
	return get
}

func fillContent(sep string, end string, content ...LogCtx) *logcolor.LogTextCtx {
	s := logcolor.New()
	b := make([]byte, 0)
	ttl := len(content) - 1
	for i, v := range content {
		switch v.(type) {
		case int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
			b = append(b, fmt.Sprintf("%d", v)...)
		case float32, float64:
			b = append(b, fmt.Sprintf("%.3f", v)...)
		case string:
			b = append(b, v.(string)...)
		case []byte:
			b = append(b, v.([]byte)...)
		case bool:
			if v.(bool) {
				b = append(b, 't', 'r', 'u', 'e')
			} else {
				b = append(b, 'f', 'a', 'l', 's', 'e')
			}
		case *logcolor.LogTextCtx:
			if len(b) != 0 {
				s.Then(logcolor.New().WithText(string(b)))
				b = b[:0]
			}
			s.Then(v.(*logcolor.LogTextCtx))
			continue
		default:
			b = append(b, fmt.Sprintf("%+v", v)...)
		}
		if i < ttl {
			b = append(b, sep...)
		}
	}
	b = append(b, end...)
	s.Then(logcolor.New().WithText(string(b)))
	return s
}

func parseOption(opts ...LogComponent) logOptions {
	dopts := defaultOptions()
	for _, opt := range opts {
		opt.apply(&dopts)
	}
	return dopts
}

////////////////////////////////////////////////////////////////////////////////
// CurInfo Functions

func (c *CurInfo) Format() string {
	return "(\"" + c.FileName + "\",in " + c.Function + " line " + strconv.Itoa(c.Line) + ")"
}

////////////////////////////////////////////////////////////////////////////////
// LoggInfo Functions

func (i *LoggInfo) Gt(other LoggInfo) bool {
	return i.Ts > other.Ts
}

////////////////////////////////////////////////////////////////////////////////
// Logger Functions

// ClearLogInfo 清空当前Logger的日志信息
func (l *Logger) ClearLogInfo() *Logger {
	l.logs = make([]*LoggInfo, 0)
	return l
}

// SetLogLevel 设置日志等级
func (l *Logger) SetLogLevel(level LogLevel) *Logger {
	l.logLevel = level
	return l
}

func (l *Logger) SetDebug(flag bool) *Logger {
	if flag {
		l.logLevel |= LevelDebug
	} else {
		l.logLevel &= ^LevelDebug
	}
	return l
}

func (l *Logger) internalPrinter(dump *LoggInfo) {
	ent := logcolor.New()
	//prefix := make([]byte, 0, 16+len(l.Name)+30+len(dump.MemCur))
	ent.Then(
		logcolor.New().WithText(
			"[" + time.Unix(int64(dump.Ts), int64(dump.Ts*1000)%1000*1000000).Format("15:04:05.000") + "]",
		).WithColor(TimeColor),
	)
	if len(dump.MemCur) > 0 {
		ent.Then(
			logcolor.New().WithText(
				"[" + dump.MemCur + "]",
			).WithColor(MemCurColor),
		)
	}
	if l.logShowCur || (dump.Level&LevelShowcur) > 0 {
		ent.Then(
			logcolor.New().WithText(
				dump.Cur.Format(),
			).WithColor(StackColor),
		)
	}

	ent.Then(
		logcolor.New().WithText(
			"<" + l.Name + ">",
		),
	)

	ent.Then(logcolor.New().WithText("["))
	ent.Then(LogPrefix[dump.Level])
	ent.Then(logcolor.New().WithText("]"))

	ent.Then(dump.Info)

	colorableStdout.Println(ent)

	if EnableGlobLog && GlobalFileHandler != nil && (dump.Level&GlobLogFilter != 0) {
		info := ent.GetRawBytes()
		info = append(info, '\n')
		if _, e := GlobalFileHandler.WriteString(*(*string)(unsafe.Pointer(&info))); e != nil {
			_ = GlobalFileHandler.Close()
			GlobalFileHandler = nil
			RootLogger.Error(WithContent("GlobalFileHandler write error:", e.Error()))
		}
	}
}

// AddPrinter 向当前Logger的Printer列表中添加一个Printer
func (l *Logger) AddPrinter(printer LogPrinter) *Logger {
	if printer != nil {
		l.printer.PushBack(printer)
	}
	return l
}

// ClearPrinter 清空当前Printer
func (l *Logger) ClearPrinter() *Logger {
	l.printer.Init()
	return l
}

// RemovePrinter 从当前Logger的Printer列表中移除指定的Printer
func (l *Logger) RemovePrinter(printer LogPrinter) *Logger {
	if printer != nil {
		gPtr := reflect2.PtrOf(printer)
		for e := l.printer.Front(); e != nil; e = e.Next() {
			if reflect2.PtrOf(e.Value.(LogPrinter)) == gPtr {
				l.printer.Remove(e)
				break
			}
		}
	}
	return l
}
func (l *Logger) GetLogs(from float64, levelMask LogLevel, maxCnt int) []*LoggInfo {

	predicate := func(x interface{}) bool {
		log, ok := x.(*LoggInfo)
		return ok && log.Ts > from && (levelMask&log.Level) != 0
	}

	q := linq.From(l.logs).Where(predicate)
	skipCnt := 0
	if maxCnt > 0 {
		skipCnt = q.Count() - maxCnt
		if skipCnt < 0 {
			skipCnt = 0
		}
	}

	var result []*LoggInfo
	if skipCnt == 0 {
		q.ToSlice(&result)
	} else {
		q.Skip(skipCnt).ToSlice(&result)
	}

	return result
}

// GetLatestLog 从当前Logger中获取最后一条记录的信息，如果没有记录，则返回nil
func (l *Logger) GetLatestLog() *LoggInfo {
	if len(l.logs) == 0 {
		return nil
	}
	return l.logs[len(l.logs)-1]
}

func (l *Logger) printerProc(dump *LoggInfo, printer *list.Element) {
	if printer == nil {
		return
	}
	defer func(l *Logger) {
		if r := recover(); r != nil {
			println("Printer Failed To Print:%v", r)
			l.printer.Remove(printer)
		}
	}(l)
	(printer.Value.(LogPrinter))(dump)
}

func (l *Logger) print(dump *LoggInfo) {
	if l.printer != nil {
		if l.printer.Len() > 0 {
			for e := l.printer.Front(); e != nil; e = e.Next() {
				l.printerProc(dump, e)
			}
		}
	}
	if l.keepPrinter {
		if l.DefaultIO != nil {
			defer func(l *Logger) {
				if r := recover(); r != nil {
					println("DefaultIO Failed To Print:%v", r)
					l.keepPrinter = false
				}
			}(l)
			l.DefaultIO(dump)
		}
	}
}

func (l *Logger) _log(info *logcolor.LogTextCtx, level LogLevel, backLevel int, log bool, log2logs bool, cur string) *LoggInfo {
	dump := &LoggInfo{
		Ts:     float64(time.Now().UnixMilli()) / 1000,
		Cur:    GetCurInfo(backLevel + 2),
		Level:  level,
		Info:   info,
		MemCur: cur,
	}
	if log2logs {
		l.logs = append(l.logs, dump)
	}
	l.latestTs = dump.Ts
	if level&l.logLevel > 0 && log {
		l.print(dump)
	}
	return dump
}
func (l *Logger) Log(opts ...LogComponent) {
	dopts := parseOption(opts...)
	l._log(fillContent(dopts.Sep, dopts.End, dopts.Info...), dopts.Level, dopts.BacktraceLevelDelta, dopts.Log, dopts.Log2Logs, dopts.Cur)
}
func (l *Logger) Common(opts ...LogComponent) {
	dopts := parseOption(opts...)
	l._log(fillContent(dopts.Sep, dopts.End, dopts.Info...), LevelCommon, dopts.BacktraceLevelDelta, dopts.Log, dopts.Log2Logs, dopts.Cur)
}

func (l *Logger) Error(opts ...LogComponent) {
	dopts := parseOption(opts...)
	l._log(fillContent(dopts.Sep, dopts.End, dopts.Info...), LevelError, dopts.BacktraceLevelDelta, dopts.Log, dopts.Log2Logs, dopts.Cur)
}

func (l *Logger) Debug(opts ...LogComponent) {
	dopts := parseOption(opts...)
	l._log(fillContent(dopts.Sep, dopts.End, dopts.Info...), LevelDebug, dopts.BacktraceLevelDelta, dopts.Log, dopts.Log2Logs, dopts.Cur)
}

func (l *Logger) Help(opts ...LogComponent) {
	dopts := parseOption(opts...)
	l._log(fillContent(dopts.Sep, dopts.End, dopts.Info...), LevelHelp, dopts.BacktraceLevelDelta, dopts.Log, false, dopts.Cur)
}

func (l *Logger) System(opts ...LogComponent) {
	dopts := parseOption(opts...)
	l._log(fillContent(dopts.Sep, dopts.End, dopts.Info...), LevelSystem, dopts.BacktraceLevelDelta, dopts.Log, dopts.Log2Logs, dopts.Cur)
}

func (l *Logger) Notice(opts ...LogComponent) {
	dopts := parseOption(opts...)
	l._log(fillContent(dopts.Sep, dopts.End, dopts.Info...), LevelNotice, dopts.BacktraceLevelDelta, dopts.Log, dopts.Log2Logs, dopts.Cur)
}

func (l *Logger) Warning(opts ...LogComponent) {
	dopts := parseOption(opts...)
	l._log(fillContent(dopts.Sep, dopts.End, dopts.Info...), LevelWarning, dopts.BacktraceLevelDelta, dopts.Log, dopts.Log2Logs, dopts.Cur)
}

func (l *Logger) Fatal(opts ...LogComponent) {
	dopts := parseOption(opts...)
	l._log(fillContent(dopts.Sep, dopts.End, dopts.Info...), LevelFatal, dopts.BacktraceLevelDelta, true, true, dopts.Cur)
}

////////////////////////////////////////////////////////////////////////////////
