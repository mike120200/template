package cmn

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
)

// dailyRotateWriter 实现了 io.Writer 接口
type dailyRotateWriter struct {
	mu       sync.Mutex
	logDir   string
	filename string
	writer   *lumberjack.Logger
	date     string
}

// NewDailyRotateWriter 创建一个新的 dailyRotateWriter
func NewDailyRotateWriter(logDir, filename string) *dailyRotateWriter {
	return &dailyRotateWriter{
		logDir:   logDir,
		filename: filename,
	}
}

// Write 实现了 io.Writer 接口
func (w *dailyRotateWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	today := time.Now().Format("2006-01-02")
	if w.writer == nil || w.date != today {
		// 创建或更新 writer
		dailyDir := filepath.Join(w.logDir, today)
		if err := os.MkdirAll(dailyDir, os.ModePerm); err != nil {
			return 0, fmt.Errorf("无法创建日志目录 %s: %w", dailyDir, err)
		}

		w.writer = &lumberjack.Logger{
			Filename:   filepath.Join(dailyDir, w.filename),
			MaxSize:    5,    // 每个日志文件的最大大小（单位：MB）
			MaxBackups: 10,   // 最大保留的旧日志文件数量
			MaxAge:     30,   // 最长保留天数（单位：天）
			Compress:   true, // 是否启用压缩
		}
		w.date = today
	}

	return w.writer.Write(p)
}

const (
	defaultLogDir = "./logs/"
)

func LoggerInit() error {
	logDir := viper.GetString("log.dir")
	if logDir == "" {
		logDir = defaultLogDir
	}
	// 检查日志目录是否存在
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err := os.MkdirAll(logDir, os.ModePerm)
		if err != nil {
			fmt.Printf("无法创建目录：%v\n", err)
			return err
		}
		fmt.Println("目录已创建：", logDir)
	} else {
		fmt.Println("目录已存在：", logDir)
	}
	cores := GetZapCores(logDir)
	// 构建 logger
	Logger := zap.New(zapcore.NewTee(cores...))
	zap.ReplaceGlobals(Logger)
	Logger.Info("log init success")

	return nil
}

// GetZapCores 创建核心，主要用于日志文件保存
func GetZapCores(logDir string) []zapcore.Core {
	cores := make([]zapcore.Core, 0, 7)
	levelStr := viper.GetString("log.level")
	smallLevel := GetLevel(levelStr)

	// 设置日志编码器
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.TimeKey = "timestamp" // 自定义时间字段名

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	//给每个日志等级创建一个日志核心
	for level := smallLevel; level <= zapcore.FatalLevel; level++ {
		// 使用我们自定义的 dailyRotateWriter
		levelFileName := zapcore.LevelOf(level).CapitalString()
		dailyWriter := NewDailyRotateWriter(logDir, fmt.Sprintf("%s.log", levelFileName))
		fileWriteSyncer := zapcore.AddSync(dailyWriter)
		cores = append(cores, zapcore.NewCore(encoder, fileWriteSyncer, level))
	}
	//控制台输出
	consoleWriteSyncer := zapcore.AddSync(os.Stdout)
	// 控制台使用 Console 编码器
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	cores = append(cores, zapcore.NewCore(consoleEncoder, consoleWriteSyncer, smallLevel))
	return cores
}

// GetLevel 获取最小的日志等级
func GetLevel(s string) zapcore.Level {
	level := strings.ToLower(s)
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.WarnLevel
	case "dpanic":
		return zapcore.DPanicLevel //生产环境中不会奔溃
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

func Logger() *zap.Logger {
	if logger == nil {
		logger = zap.L()
	}
	return logger
}
