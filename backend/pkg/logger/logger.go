package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log *zap.Logger
)

// Config 日志配置
type Config struct {
	Level      string // debug, info, warn, error
	Format     string // json, console
	OutputPath string // stdout, stderr, or file path
}

// Init 初始化日志
func Init(cfg *Config) error {
	if Log != nil {
		return nil // 已初始化
	}
	return initLogger(cfg)
}

func initLogger(cfg *Config) error {
	// 解析日志级别
	level := zapcore.InfoLevel
	if cfg != nil && cfg.Level != "" {
		if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
			level = zapcore.InfoLevel
		}
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 选择编码器
	var encoder zapcore.Encoder
	format := "console"
	if cfg != nil && cfg.Format != "" {
		format = cfg.Format
	}
	if format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 输出目标
	var writeSyncer zapcore.WriteSyncer
	outputPath := "stdout"
	if cfg != nil && cfg.OutputPath != "" {
		outputPath = cfg.OutputPath
	}

	switch outputPath {
	case "stdout":
		writeSyncer = zapcore.AddSync(os.Stdout)
	case "stderr":
		writeSyncer = zapcore.AddSync(os.Stderr)
	default:
		// 确保日志目录存在
		logDir := filepath.Dir(outputPath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
		file, err := os.OpenFile(outputPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		writeSyncer = zapcore.AddSync(file)
	}

	// 创建核心
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 创建日志器
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// InitDefault 使用默认配置初始化
func InitDefault() {
	Init(&Config{
		Level:      "info",
		Format:     "console",
		OutputPath: "stdout",
	})
}

// Debug 输出 debug 级别日志
func Debug(msg string, fields ...zap.Field) {
	if Log == nil {
		InitDefault()
	}
	Log.Debug(msg, fields...)
}

// Info 输出 info 级别日志
func Info(msg string, fields ...zap.Field) {
	if Log == nil {
		InitDefault()
	}
	Log.Info(msg, fields...)
}

// Warn 输出 warn 级别日志
func Warn(msg string, fields ...zap.Field) {
	if Log == nil {
		InitDefault()
	}
	Log.Warn(msg, fields...)
}

// Error 输出 error 级别日志
func Error(msg string, fields ...zap.Field) {
	if Log == nil {
		InitDefault()
	}
	Log.Error(msg, fields...)
}

// Fatal 输出 fatal 级别日志并退出
func Fatal(msg string, fields ...zap.Field) {
	if Log == nil {
		InitDefault()
	}
	Log.Fatal(msg, fields...)
}

// With 创建带有字段的子日志器
func With(fields ...zap.Field) *zap.Logger {
	if Log == nil {
		InitDefault()
	}
	return Log.With(fields...)
}

// Sync 刷新日志缓冲
func Sync() error {
	if Log != nil {
		return Log.Sync()
	}
	return nil
}
