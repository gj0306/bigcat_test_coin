package coinlog

import (
	"fmt"
	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func NewZapLog(logDir,logName string,kv map[string]interface{})(logger log.Logger){
	hook := lumberjack.Logger{
		Filename:   fmt.Sprintf("./%s/%s.log", logDir, logName), // 日志文件路径
		MaxSize:    64,                                            // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 10,                                            // 日志文件最多保存多少个备份
		MaxAge:     7,                                             // 文件最多保存多少天
		Compress:   true,                                          // 是否压缩
	}
	errHook := lumberjack.Logger{
		Filename:   fmt.Sprintf("./%s/%s.error.log", logDir, logName), // 日志文件路径
		MaxSize:    64,                                           // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 10,                                           // 日志文件最多保存多少个备份
		MaxAge:     7,                                            // 文件最多保存多少天
		Compress:   true,                                         // 是否压缩
	}
	// 日志编码配置
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
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	// 指定info和error的日志写入文件
	infoMultiWriter := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	errMultiWriter := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&errHook))
	// 定义日志级别
	infoLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zapcore.WarnLevel
	})
	errLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.WarnLevel
	})
	// 组合所有配置
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, infoMultiWriter, infoLevel),
		zapcore.NewCore(encoder, errMultiWriter, errLevel),
	)
	var _logger *zap.Logger
	_logger = zap.New(core)
	kl := make([]interface{},0,len(kv)*2)
	for k,v := range kv{
		kl = append(kl,k )
		kl = append(kl,v )
	}
	l := kzap.NewLogger(_logger)
	return log.With(l,kl...)
}