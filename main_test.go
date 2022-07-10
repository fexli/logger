package logger

import "github.com/fexli/logger"

func main() {
	logger.RootLogger.Debug(logger.WithContent("hello world"))
	logger.GetLogger("AA", true)
}
