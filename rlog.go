package rlog

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var (
	infomessage      *klogger
	errormessage     *klogger
	debugmessage     *klogger
	warnmessage      *klogger
	logFileDir       string = "logs/"
	logFileNameOld   string
	logFileOld       *os.File
	logLevel         int   = 1
	maxLogFileSize   int64 = 10 * 1024 * 1024 // 日志文件最大size
	stdoutFlag       bool  = true             // 日志是否同时输入到stdout中
	defaultCalldepth int   = 5
)

type klogger struct {
	logger *log.Logger
}

func init() {
	logFileName := time.Now().Format("2006.01.02_15-04-05.log")
	setKLogFile(logFileName)
	SetDefaultCalldepth(3)
	go checkLogFileSize()
}

func newLogger(out io.Writer, prefix string) *klogger {
	logger := new(klogger)
	if nil != out {
		logger.logger = log.New(out, prefix, log.Ldate|log.Ltime)
	}
	return logger
}

func setKLogFile(iFileName string) {
	logFileNameOld = iFileName
	if _, err := os.Stat(logFileDir); os.IsNotExist(err) {
		log.Println(logFileDir)
		err = os.MkdirAll(logFileDir, os.ModePerm)
		if nil != err {
			log.Fatalf("[Error] Create klog dir failed: %s, %v", logFileDir, err)
		}
	}
	logFile, err := os.OpenFile(logFileDir+iFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		log.Fatalf("[Error] Failed to open log logFile: %v, logFile: %v", iFileName, logFile)
	}
	logFileOld.Close()
	logFileOld = logFile
	mw := io.MultiWriter(logFile)
	if stdoutFlag {
		mw = io.MultiWriter(os.Stdout, logFile)
	}
	infomessage = newLogger(mw, fmt.Sprintf("%-8s", "[INFO]"))
	warnmessage = newLogger(mw, fmt.Sprintf("%-8s", "[WARN] "))
	errormessage = newLogger(mw, fmt.Sprintf("%-8s", "[ERROR]"))
	debugmessage = newLogger(mw, fmt.Sprintf("%-8s", "[DEBUG] "))
}

func checkLogFileSize() {
	for {
		time.Sleep(1 * time.Second)
		fInfo, err := os.Stat(logFileDir + logFileNameOld)
		if nil != err {
			fmt.Printf("getLogFileSize Error: %v", err)
			return
		}
		if fInfo.Size() >= maxLogFileSize {
			logFileNameNew := time.Now().Format("2006.01.02_15-04-05.log")
			setKLogFile(logFileNameNew)
		}
	}
}

// 一般信息
func Info(v ...interface{}) {
	if logLevel >= 1 {
		infomessage.logger.Output(defaultCalldepth, fmt.Sprintln(v...))
	}
}

// 警告信息
func Warn(v ...interface{}) {
	if logLevel >= 2 {
		warnmessage.logger.Output(defaultCalldepth, fmt.Sprintln(v...))
	}
}

// 错误信息，会自动退出
func Error(v ...interface{}) {
	errormessage.logger.Output(defaultCalldepth, fmt.Sprintln(v...))
	os.Exit(1)
}

// 详细调试信息
func Debug(v ...interface{}) {
	if logLevel >= 3 {
		debugmessage.logger.Output(defaultCalldepth, fmt.Sprintln(v...))
	}
}

// 是否输出到终端，否则只输入文件
func SetStdOut(iStdoutFlag bool) {
	stdoutFlag = iStdoutFlag
}

// 设置日志显示级别
func SetLogLevel(iLogLevel string) {
	switch {
	case strings.EqualFold("DEBUG", iLogLevel):
		logLevel = 3
	case strings.EqualFold("WARN", iLogLevel):
		logLevel = 2
	case strings.EqualFold("INFO", iLogLevel):
		logLevel = 1
	}
	log.Println(logLevel, iLogLevel)
}

// 设置日志文件存放目录
func SetLogFileDir(iLogFileDir string) {
	logFileDir = iLogFileDir
}

// 以MB为单位设置日志文件最大size
func SetMaxFileSizeMB(size int) {
	maxLogFileSize = int64(size) * 1024 * 1024
}

// 设置默认的打印深度
func SetDefaultCalldepth(calldepth int) {
	defaultCalldepth = calldepth
}
