package logger

type LogAction func(...LogComponent)

type LogBuilder struct {
	components []LogComponent
	onCommit   LogAction
}

/*******	 	LogBuilder Builders 		*******/

func newBuilder(onCommit LogAction) *LogBuilder {
	return &LogBuilder{
		components: []LogComponent{},
		onCommit:   onCommit,
	}
}

func (l *Logger) NewLog(ctx ...LogCtx) *LogBuilder {
	return newBuilder(l.Log).WithContent(ctx...)
}

func (l *Logger) NewCommonLog(ctx ...LogCtx) *LogBuilder {
	return newBuilder(l.Common).WithContent(ctx...)
}

func (l *Logger) NewErrorLog(ctx ...LogCtx) *LogBuilder {
	return newBuilder(l.Error).WithContent(ctx...)
}

func (l *Logger) NewDebugLog(ctx ...LogCtx) *LogBuilder {
	return newBuilder(l.Debug).WithContent(ctx...)
}

func (l *Logger) NewHelpLog(ctx ...LogCtx) *LogBuilder {
	return newBuilder(l.Help).WithContent(ctx...)
}

func (l *Logger) NewSystemLog(ctx ...LogCtx) *LogBuilder {
	return newBuilder(l.System).WithContent(ctx...)
}

func (l *Logger) NewNoticeLog(ctx ...LogCtx) *LogBuilder {
	return newBuilder(l.Notice).WithContent(ctx...)
}

func (l *Logger) NewWarningLog(ctx ...LogCtx) *LogBuilder {
	return newBuilder(l.Warning).WithContent(ctx...)
}

func (l *Logger) NewFatalLog(ctx ...LogCtx) *LogBuilder {
	return newBuilder(l.Fatal).WithContent(ctx...)
}

/*******	 LogBuilder Chain Extension Methods 	*******/

func (builder *LogBuilder) WithComponent(comp LogComponent) *LogBuilder {
	builder.components = append(builder.components, comp)
	return builder
}

func (builder *LogBuilder) WithContent(ctx ...LogCtx) *LogBuilder {
	return builder.WithComponent(WithContent(ctx...))
}

func (builder *LogBuilder) WithBacktraceLevelDelta(level int) *LogBuilder {
	return builder.WithComponent(WithBacktraceLevelDelta(level))
}

func (builder *LogBuilder) WithLog(log bool) *LogBuilder {
	return builder.WithComponent(WithLog(log))
}

func (builder *LogBuilder) WithSep(sep string) *LogBuilder {
	return builder.WithComponent(WithSep(sep))
}

func (builder *LogBuilder) WithEnd(end string) *LogBuilder {
	return builder.WithComponent(WithEnd(end))
}

func (builder *LogBuilder) WithLevel(level LogLevel) *LogBuilder {
	return builder.WithComponent(WithLevel(level))
}

func (builder *LogBuilder) WithLog2Logs(enabled bool) *LogBuilder {
	return builder.WithComponent(WithLog2Logs(enabled))
}

func (builder *LogBuilder) WithKVs(keyValues ...interface{}) *LogBuilder {
	return builder.WithComponent(WithKVs(keyValues...))
}

func (builder *LogBuilder) WithStruct(s interface{}) *LogBuilder {
	return builder.WithComponent(WithStruct(s))
}

func (builder *LogBuilder) Commit() {
	//防止回溯到Commit函数
	builder.WithBacktraceLevelDelta(1).onCommit(builder.components...)
}
