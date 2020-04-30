package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"sync"
	xormcore "xorm.io/core"
)

var (
	loggerInstance                 *Logger
	sp                             = string(filepath.Separator)
	errWS, warnWS, infoWS, debugWS zapcore.WriteSyncer       // IO输出
	debugConsoleWS                 = zapcore.Lock(os.Stdout) // 控制台标准输出
	errorConsoleWS                 = zapcore.Lock(os.Stderr)
)

type Logger struct {
	*zap.Logger
	atomicLevel zap.AtomicLevel
	serviceName string
	zapConfig   zap.Config
	sync.RWMutex
}

// Init 初始化
func Init(serviceName string, options ...zap.Option) {
	if loggerInstance != nil {
		return
	}

	logCfg := zap.NewDevelopmentConfig()
	logger, err := logCfg.Build(append(options, zap.AddCallerSkip(2))...)
	if err != nil {
		panic("logger init failed![" + err.Error() + "]")
	}
	loggerInstance = &Logger{
		Logger:      logger,
		atomicLevel: logCfg.Level,
		serviceName: serviceName,
	}
}

func NewLogger(serviceName string, options ...zap.Option) *Logger {
	logCfg := zap.NewProductionConfig()
	logger, err := logCfg.Build(append(options, zap.AddCallerSkip(1))...)

	if err != nil {
		panic("logger init failed![" + err.Error() + "]")
	}
	return &Logger{
		Logger:      logger,
		atomicLevel: logCfg.Level,
		serviceName: serviceName,
	}
}

func Sync() {
	if loggerInstance != nil {
		_ = loggerInstance.Sync()
	}
}
func (log *Logger) Debug(msg string, fields ...zap.Field) {
	log.Logger.Debug(msg, append(fields, zap.String("ServiceName", loggerInstance.serviceName))...)
}

func (log *Logger) Info(msg string, fields ...zap.Field) {
	log.Logger.Info(msg, append(fields, zap.String("ServiceName", loggerInstance.serviceName))...)
}

func (log *Logger) Warn(msg string, fields ...zap.Field) {
	log.Logger.Warn(msg, append(fields, zap.String("ServiceName", loggerInstance.serviceName))...)
}

func (log *Logger) Error(msg string, fields ...zap.Field) {
	log.Logger.Error(msg, append(fields, zap.String("ServiceName", loggerInstance.serviceName))...)
}

func (log *Logger) DPanic(msg string, fields ...zap.Field) {
	log.Logger.DPanic(msg, append(fields, zap.String("ServiceName", loggerInstance.serviceName))...)
}

func (log *Logger) Panic(msg string, fields ...zap.Field) {
	log.Logger.Panic(msg, append(fields, zap.String("ServiceName", loggerInstance.serviceName))...)
}

func (log *Logger) Fatal(msg string, fields ...zap.Field) {
	log.Logger.Fatal(msg, append(fields, zap.String("ServiceName", loggerInstance.serviceName))...)
}


// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zap.Field) {
	loggerInstance.Debug(msg, fields...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zap.Field) {
	loggerInstance.Info(msg, fields...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...zap.Field) {
	loggerInstance.Warn(msg, fields...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...zap.Field) {
	loggerInstance.Error(msg, fields...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.

func DPanic(msg string, fields ...zap.Field) {
	loggerInstance.DPanic(msg, fields...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...zap.Field) {
	loggerInstance.Panic(msg, fields...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...zap.Field) {
	loggerInstance.Fatal(msg, fields...)
}

func SetLogLevel(l zapcore.Level) {
	loggerInstance.atomicLevel.SetLevel(l)
}
func GetXormLogger() *XormLogger {
	if loggerInstance == nil {
		panic("loggerInstance is nil!")
	}
	return &XormLogger{
		logger:  loggerInstance,
		showSQL: false,
	}
}

type XormLogger struct {
	logger  *Logger
	showSQL bool
	//level xormcore.LogLevel
}

func (l *XormLogger) Debug(v ...interface{}) {
	//if l.level <= xormcore.LOG_DEBUG {
	l.logger.Debug(fmt.Sprint(v...))
	//}
}
func (l *XormLogger) Debugf(format string, v ...interface{}) {

	l.logger.Debug(fmt.Sprintf(format, v...))
}
func (l *XormLogger) Error(v ...interface{}) {
	l.logger.Error(fmt.Sprint(v...))
}
func (l *XormLogger) Errorf(format string, v ...interface{}) {
	l.logger.Error(fmt.Sprintf(format, v...))
}
func (l *XormLogger) Info(v ...interface{}) {
	l.logger.Info(fmt.Sprint(v...))
}
func (l *XormLogger) Infof(format string, v ...interface{}) {
	//l.logger.Info(fmt.Sprintf(format, v...))
	l.logger.Debug(fmt.Sprintf(format, v...)) //暂时info降级成debug
}
func (l *XormLogger) Warn(v ...interface{}) {
	l.logger.Warn(fmt.Sprint(v...))
}
func (l *XormLogger) Warnf(format string, v ...interface{}) {
	l.logger.Warn(fmt.Sprintf(format, v...))
}

func (l *XormLogger) Level() xormcore.LogLevel {
	switch l.logger.atomicLevel.Level() {
	case zapcore.DebugLevel:
		return xormcore.LOG_DEBUG
	case zapcore.InfoLevel:
		return xormcore.LOG_INFO
	case zapcore.WarnLevel:
		return xormcore.LOG_WARNING
	case zapcore.ErrorLevel:
		return xormcore.LOG_ERR
	case zapcore.DPanicLevel:
		return xormcore.LOG_OFF
	default:
		return xormcore.LOG_UNKNOWN
	}
}
func (l *XormLogger) SetLevel(lv xormcore.LogLevel) {
	// 接口需要，level设置统一由系统完成，此处禁止使用。
}

func (l *XormLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		l.showSQL = true
		return
	}
	l.showSQL = show[0]
}
func (l *XormLogger) IsShowSQL() bool {
	return l.showSQL
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func FatalWithCode(msg string, fields ...zap.Field) {
	loggerInstance.Fatal(msg, fields...)
}

func initLogger() {

	loggerInstance.Lock()
	defer loggerInstance.Unlock()
	loggerInstance.init()
}

// GetLogger returns logger
func GetLogger() (ret *Logger) {
	return loggerInstance
}

func (log *Logger) init() {

	var err error

	log.Logger, err = log.zapConfig.Build(log.cores(), zap.AddCallerSkip(2))
	if err != nil {
		panic(err)
	}

	//defer loggerInstance.Sync()
	//可能有没有拿到log.Logger实例的问题，暂时注释
	//defer func() {
	//	_ = log.Logger.Sync()
	//}()
}



func (log *Logger) cores() zap.Option {

	fileEncoder := zapcore.NewJSONEncoder(log.zapConfig.EncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(log.zapConfig.EncoderConfig)

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.WarnLevel && zapcore.WarnLevel-log.zapConfig.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-log.zapConfig.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-log.zapConfig.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-log.zapConfig.Level.Level() > -1
	})

	cores := []zapcore.Core{
		// region 日志文件

		// error 及以上
		zapcore.NewCore(fileEncoder, errWS, errPriority),

		// warn
		zapcore.NewCore(fileEncoder, warnWS, warnPriority),

		// info
		zapcore.NewCore(fileEncoder, infoWS, infoPriority),

		// debug
		zapcore.NewCore(fileEncoder, debugWS, debugPriority),

		// endregion

		// region 控制台

		// 错误及以上
		zapcore.NewCore(consoleEncoder, errorConsoleWS, errPriority),

		// 警告
		zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority),

		// info
		zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority),

		// debug
		zapcore.NewCore(consoleEncoder, debugConsoleWS, debugPriority),

		// endregion
	}

	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}
