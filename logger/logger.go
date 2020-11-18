package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 日志库

// * 实现功能
// * 支持多种输出方式stdout/file
// * 支持输出为json 或 plaintext
// * 支持彩色输出
// * 支持log rotate
// *

//Cfg is the struct for log information
type Cfg struct {
	Writers       string `yaml:"writers"`
	Level         string `yaml:"level"`
	File          string `yaml:"file"`
	FormatText    bool   `yaml:"format_text"`
	Color         bool   `yaml:"color"`
	RollingPolicy string `yaml:"rollingPolicy"`
	RotateDate    int    `yaml:"lrotate_date"`
	RotateSize    int    `yaml:"rotate_size"`
	BackupCount   int    `yaml:"backup_count"`
}

// Logger is the global variable
// var Logger *logrus.Logger
var Logger = logrus.New()

// filePath log file path
var filePath string

// definition is having the information about loging
var definition *Cfg = DefaultDefinition()

// constant values for logrotate parameters
const (
	RollingPolicySize = "size"
	LogRotateDate     = 1
	LogRotateSize     = 10
	LogBackupCount    = 7
)

//InitWithConfig 初始化
func InitWithConfig(def *Cfg) (*logrus.Logger, error) {
	definition = def
	return Logger, Init()
}

func init() {
	// Logger = logrus.New()
	Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetOutput(os.Stdout)
	Logger.SetLevel(logrus.DebugLevel)
}

//Init 初始化
func Init() error {
	// Logger = logrus.New()

	//  只输出不低于当前级别是日志数据
	level, err := logrus.ParseLevel(definition.Level)
	Logger.SetLevel(level)

	// 输出日志格式
	var formatter logrus.Formatter
	formatter = &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}
	if definition.FormatText {
		formatter = &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     definition.Color,
		}
	}
	Logger.SetFormatter(formatter)

	// 默认只有 stdout 输出
	writers := []io.Writer{os.Stdout}

	// 如果包含文件输出
	if strings.Index(definition.Writers, "stdout") >= 0 {
		// 查看文件路径是否正确
		if filepath.IsAbs(definition.File) {
			createFile("", definition.File)
			filePath = filepath.Join("", definition.File)
		} else {
			createFile(os.Getenv("CHASSIS_HOME"), definition.File)
			filePath = filepath.Join(os.Getenv("CHASSIS_HOME"), definition.File)
		}

		// 判断文件打开是否正常
		var file io.Writer
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		} else {
			// 文件正常，则使用分割日志的功能
			f.Close()
			file = &lumberjack.Logger{
				Filename:   filePath,
				MaxSize:    100,  // megabytes
				MaxBackups: 30,   //days
				Compress:   true, // disabled by default
			}
			writers = append(writers, file)
		}

	}

	//同时写文件和屏幕
	fileAndStdoutWriter := io.MultiWriter(writers...)
	Logger.SetOutput(fileAndStdoutWriter)

	return err
}

// Trace Trace
func Trace(args ...interface{}) {
	Logger.Trace(args...)
}

// Tracef Tracef
func Tracef(format string, args ...interface{}) {
	Logger.Tracef(format, args...)
}

// Debug Debug
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Debugf Debugf
func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

// Info Info
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Infof Infof
func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

// Warn Warn
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Warnf Warnf
func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

// Error Error
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Errorf Errorf
func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

// Fatal Log with os.exit(1)
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Fatalf Log with os.exit(1)
func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}

// Panic Log with panic
func Panic(args ...interface{}) {
	Logger.Panic(args...)
}

// Panicf Log with panic
func Panicf(format string, args ...interface{}) {
	Logger.Panicf(format, args...)
}

//createLogFile create log file
func createFile(localPath, outputpath string) {
	_, err := os.Stat(strings.Replace(filepath.Dir(filepath.Join(localPath, outputpath)), "\\", "/", -1))
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(strings.Replace(filepath.Dir(filepath.Join(localPath, outputpath)), "\\", "/", -1), os.ModePerm)
	} else if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(strings.Replace(filepath.Join(localPath, outputpath), "\\", "/", -1), os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
}

//DefaultDefinition 预定义
func DefaultDefinition() *Cfg {
	cfg := Cfg{
		Writers:       "stdout,file",
		Level:         "DEBUG",
		File:          "log/chassis.log",
		FormatText:    false,
		Color:         false,
		RollingPolicy: RollingPolicySize,
		RotateDate:    1,
		RotateSize:    10,
		BackupCount:   7,
	}

	return &cfg
}
