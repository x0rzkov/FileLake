package logger

import (
	"HIBL/common/config"
	"HIBL/common/tools/msg"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"sync"
	xormcore "xorm.io/core"
)

var (
	loggerInstance                 *MicroLogger
	sp                             = string(filepath.Separator)
	errWS, warnWS, infoWS, debugWS zapcore.WriteSyncer       // IO输出
	debugConsoleWS                 = zapcore.Lock(os.Stdout) // 控制台标准输出
	errorConsoleWS                 = zapcore.Lock(os.Stderr)
)

type MicroLogger struct {
	*zap.Logger
	atomicLevel zap.AtomicLevel
	serviceName string
	msgResolver *msg.Resolver
	Opts        *config.Zap `json:"opts"`
	zapConfig   zap.Config
	sync.RWMutex
}

// Init 初始化
func Init(serviceName string, options ...zap.Option) {
	//alevel := zap.NewAtomicLevel()
	//如果经初始化 直接跳出
	if loggerInstance != nil {
		return
	}
	//todo 配置中心修改
	//logCfg := zap.NewProductionConfig()
	logCfg := zap.NewDevelopmentConfig()
	//zap.AddCaller(), zap.AddCallerSkip(1)
	logger, err := logCfg.Build(append(options, zap.AddCallerSkip(2))...)
	//logger, err := zap.NewProduction(options...)
	if err != nil {
		panic("logger init failed![" + err.Error() + "]")
	}
	loggerInstance = &MicroLogger{
		Logger:      logger,
		atomicLevel: logCfg.Level,
		serviceName: serviceName,
	}
}
func InitWithConf(zapCfg *config.Zap) {
	loggerInstance = &MicroLogger{
		Opts:        zapCfg,
		serviceName: zapCfg.AppName,
	}
	initLogger()
}

func NewLogger(serviceName string, options ...zap.Option) *MicroLogger {
	logCfg := zap.NewProductionConfig()
	logger, err := logCfg.Build(append(options, zap.AddCallerSkip(1))...)

	if err != nil {
		panic("logger init failed![" + err.Error() + "]")
	}
	return &MicroLogger{
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
func (log *MicroLogger) Debug(msg string, fields ...zap.Field) {
	log.Logger.Debug(msg, append(fields, zap.String("ServiceName", loggerInstance.serviceName))...)
}

func (log *MicroLogger) Info(msg string, fields ...zap.Field) {
	log.Logger.Info(msg, append(fields, zap.String("ServiceName", loggerInstance.serviceName))...)
}

func (log *MicroLogger) Warn(msg string, fields ...zap.Field) {
	log.Logger.Warn(msg, append(fields, zap.String("ServiceName", loggerInstance.serviceName))...)
}

func (log *MicroLogger) Error(msg string, fields ...zap.Field) {
	log.Logger.Error(msg, append(fields, zap.String("ServiceName", loggerInstance.serviceName))...)
}

func (log *MicroLogger) DPanic(msg string, fields ...zap.Field) {
	log.Logger.DPanic(msg, append(fields, zap.String("ServiceName", loggerInstance.serviceName))...)
}

func (log *MicroLogger) Panic(msg string, fields ...zap.Field) {
	log.Logger.Panic(msg, append(fields, zap.String("ServiceName", loggerInstance.serviceName))...)
}

func (log *MicroLogger) Fatal(msg string, fields ...zap.Field) {
	log.Logger.Fatal(msg, append(fields, zap.String("ServiceName", loggerInstance.serviceName))...)
}

func (log *MicroLogger) SetMsgResolver(r *msg.Resolver) {
	log.msgResolver = r
}

func SetMsgResolver(r *msg.Resolver) {
	loggerInstance.SetMsgResolver(r)
}

func (log *MicroLogger) InfoWithCode(invokeID string, code string, v ...interface{}) {
	if log.msgResolver == nil {
		panic("Please call the function [SetMsgResolver(r *msg.Resolver)]")
	}
	msgStruct := log.msgResolver.GetMsg(code)
	_msg := fmt.Sprintf(msgStruct.Lmsg, v...)
	log.Info(_msg, zap.String("InvokeID", invokeID), zap.String("MsgCode", msgStruct.MsgCode))
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (log *MicroLogger) WarnWithCode(invokeID string, code string, v ...interface{}) {
	if log.msgResolver == nil {
		panic("Please call the function [SetMsgResolver(r *msg.Resolver)]")
	}
	msgStruct := log.msgResolver.GetMsg(code)
	_msg := fmt.Sprintf(msgStruct.Lmsg, v...)
	log.Warn(_msg, zap.String("InvokeID", invokeID), zap.String("MsgCode", msgStruct.MsgCode))
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (log *MicroLogger) ErrorWithCode(invokeID string, code string, v ...interface{}) {
	if log.msgResolver == nil {
		panic("Please call the function [SetMsgResolver(r *msg.Resolver)]")
	}
	msgStruct := log.msgResolver.GetMsg(code)
	_msg := fmt.Sprintf(msgStruct.Lmsg, v...)
	log.Error(_msg, zap.String("InvokeID", invokeID), zap.String("MsgCode", msgStruct.MsgCode))
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.

func (log *MicroLogger) DPanicWithCode(invokeID string, code string, v ...interface{}) {
	if log.msgResolver == nil {
		panic("Please call the function [SetMsgResolver(r *msg.Resolver)]")
	}
	msgStruct := log.msgResolver.GetMsg(code)
	_msg := fmt.Sprintf(msgStruct.Lmsg, v...)
	log.DPanic(_msg, zap.String("InvokeID", invokeID), zap.String("MsgCode", msgStruct.MsgCode))
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (log *MicroLogger) PanicWithCode(invokeID string, code string, v ...interface{}) {
	if log.msgResolver == nil {
		panic("Please call the function [SetMsgResolver(r *msg.Resolver)]")
	}
	msgStruct := log.msgResolver.GetMsg(code)
	_msg := fmt.Sprintf(msgStruct.Lmsg, v...)
	log.Panic(_msg, zap.String("InvokeID", invokeID), zap.String("MsgCode", msgStruct.MsgCode))
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (log *MicroLogger) FatalWithCode(invokeID string, code string, v ...interface{}) {
	if log.msgResolver == nil {
		panic("Please call the function [SetMsgResolver(r *msg.Resolver)]")
	}
	msgStruct := log.msgResolver.GetMsg(code)
	_msg := fmt.Sprintf(msgStruct.Lmsg, v...)
	log.Fatal(_msg, zap.String("InvokeID", invokeID), zap.String("MsgCode", msgStruct.MsgCode))
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
	logger  *MicroLogger
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

func InfoWithCode(invokeID string, code string, v ...interface{}) {

	loggerInstance.InfoWithCode(invokeID, code, v...)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func WarnWithCode(invokeID string, code string, v ...interface{}) {
	loggerInstance.WarnWithCode(invokeID, code, v...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func ErrorWithCode(invokeID string, code string, v ...interface{}) {
	loggerInstance.ErrorWithCode(invokeID, code, v...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.

func DPanicWithCode(invokeID string, code string, v ...interface{}) {
	loggerInstance.DPanicWithCode(invokeID, code, v...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func PanicWithCode(invokeID string, code string, v ...interface{}) {
	loggerInstance.PanicWithCode(invokeID, code, v...)
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
	loggerInstance.loadCfg()
	loggerInstance.init()
}

// GetLogger returns logger
func GetLogger() (ret *MicroLogger) {
	return loggerInstance
}

func (log *MicroLogger) init() {

	log.setSyncers()
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

func (log *MicroLogger) loadCfg() {

	//c := config.C()

	//err := c.Path("zap", l.Opts)
	//if err != nil {
	//	panic(err)
	//}

	if log.Opts.Development {
		log.zapConfig = zap.NewDevelopmentConfig()
	} else {
		log.zapConfig = zap.NewProductionConfig()
	}

	// application log output path
	if log.Opts.OutputPaths == nil || len(log.Opts.OutputPaths) == 0 {
		log.zapConfig.OutputPaths = []string{"stdout"}
	}

	//  error of zap-self log
	if log.Opts.ErrorOutputPaths == nil || len(log.Opts.ErrorOutputPaths) == 0 {
		log.zapConfig.OutputPaths = []string{"stderr"}
	}

	// 默认输出到程序运行目录的logs子目录
	if log.Opts.LogFileDir == "" {
		log.Opts.LogFileDir, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
		log.Opts.LogFileDir += sp + "logs" + sp
	}

	if log.Opts.AppName == "" {
		log.Opts.AppName = "app"
	}

	if log.Opts.ErrorFileName == "" {
		log.Opts.ErrorFileName = "error.log"
	}

	if log.Opts.WarnFileName == "" {
		log.Opts.WarnFileName = "warn.log"
	}

	if log.Opts.InfoFileName == "" {
		log.Opts.InfoFileName = "info.log"
	}

	if log.Opts.DebugFileName == "" {
		log.Opts.DebugFileName = "debug.log"
	}

	if log.Opts.MaxSize == 0 {
		log.Opts.MaxSize = 50
	}
	if log.Opts.MaxBackups == 0 {
		log.Opts.MaxBackups = 3
	}
	if log.Opts.MaxAge == 0 {
		log.Opts.MaxAge = 30
	}
}

func (log *MicroLogger) setSyncers() {

	f := func(fN string) zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   log.Opts.LogFileDir + sp + log.Opts.AppName + "-" + fN,
			MaxSize:    log.Opts.MaxSize,
			MaxBackups: log.Opts.MaxBackups,
			MaxAge:     log.Opts.MaxAge,
			Compress:   true,
			LocalTime:  true,
		})
	}

	errWS = f(log.Opts.ErrorFileName)
	warnWS = f(log.Opts.WarnFileName)
	infoWS = f(log.Opts.InfoFileName)
	debugWS = f(log.Opts.DebugFileName)
}

func (log *MicroLogger) cores() zap.Option {

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
