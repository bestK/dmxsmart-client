package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/bestk/dmxsmart-client/config"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

var Logger *slog.Logger

func Init() {
	logDir := "logs"

	if config.GlobalConfig.LogDir != "" {
		logDir = config.GlobalConfig.LogDir
	}

	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(fmt.Sprintf("创建日志目录失败: %v", err))
	}

	logFileName := filepath.Join(logDir, "dmxsmart-client", time.Now().Format("2006-01-02")+".log")

	fileHandler := handler.MustRotateFile(logFileName, handler.EveryHour)

	consoleHandler := handler.NewConsoleHandler(slog.AllLevels)

	logger := slog.NewWithHandlers(consoleHandler, fileHandler)
	Logger = logger
}
