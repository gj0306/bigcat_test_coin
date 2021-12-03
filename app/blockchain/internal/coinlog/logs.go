package coinlog

import (
	"fmt"
	"github.com/google/wire"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"bigcat_test_coin/app/blockchain/internal/config"
	"os"
	"strings"
	"sync"
)

var (
	// Logger 日志对象
	Logger     *zap.Logger
	loggerOnce sync.Once
)

// NewLogger 获取Logger对象
func NewLogger(conf *config.LogConf) *zap.Logger {
	loggerOnce.Do(func() {
		Logger = logConfig(conf)
	})
	return Logger
}

var ProviderSet = wire.NewSet(NewLogger)

// LogConfig 日志配置
func logConfig(conf *config.LogConf) *zap.Logger {
	if conf.LogPath == "" || conf.LogFileName == "" {
		log.Fatal("配置文件错误: 日志文件夹或者日志文件为空")
	}

	logDir := conf.LogPath
	if !strings.HasSuffix(logDir, "/") {
		logDir = logDir + "/"
	}
	fileName := conf.LogFileName

	hook := lumberjack.Logger{
		Filename:   fmt.Sprintf("./%s%s", logDir, fileName), // 日志文件路径
		MaxSize:    16,                                      // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 10,                                      // 日志文件最多保存多少个备份
		MaxAge:     30,                                      // 文件最多保存多少天
		Compress:   true,                                    // 是否压缩
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "line",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,                      // 小写编码器
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"), // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zap.DebugLevel)

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),                                        // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel,                                                                     // 日志级别
	)

	var _logger *zap.Logger
	if conf.Debug {
		// 开启开发模式，堆栈跟踪
		caller := zap.AddCaller()
		// 开启文件及行号
		development := zap.Development()
		// 设置初始化字段
		// filed := zap.Fields(zap.String("service", "novel-downloader"))
		// 构造日志
		_logger = zap.New(core, caller, development)
	} else {
		_logger = zap.New(core)
	}
	return _logger
}
