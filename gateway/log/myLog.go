package log

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Logger struct {
	logLv     uint64
	logDir    string
	logFile   *os.File
	logFileWf *os.File
	sync.Mutex
	buf []byte
}

type UserLogger struct {
	*Logger
	userBuf []byte
}

const (
	TRACE = iota + 1
	DEBUG
	INFO
	WARNING
	FATAL
)

var g_bufPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 0, 500)

	},
}
var mLog = Logger{logLv: 1, logDir: "./log"}

func NewLogger() *Logger {
	return &mLog
}

func NewUserLogger() *UserLogger {
	return &UserLogger{&mLog, g_bufPool.Get().([]byte)}

}

func (this *Logger) GetLogLevel() uint64 {
	return this.logLv

}

func (this *Logger) SetLogLevel(lv uint64) {
	this.logLv = lv

}

func (this *Logger) write(fp *os.File, logtype string, calldepth int, bufInfo string, s string) {

	cNow := time.Now()
	cTime := fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d %06d",
		cNow.Year(), int(cNow.Month()), cNow.Day(), cNow.Hour(), cNow.Minute(), cNow.Second(), cNow.Nanosecond())
	//cTime := time.Now().Format("2006-01-02 15:04:05")
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0

	}

	fileInfo := file + ":" + strconv.Itoa(line)
	midStr := bufInfo + "[" + fileInfo + "]"

	if len(s) == 0 || s[len(s)-1] != '\n' {
		s += "\n"

	}

	ret := fmt.Sprintf("[%s][%s]%s %s", cTime, logtype, midStr, s)
	this.Lock()
	fp.WriteString(ret)
	this.Unlock()
}

func (this *Logger) PrintStdout(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (this *Logger) SetLogPath(dirpath string, name string) {
	this.logDir = dirpath
	err := os.MkdirAll(this.logDir, 0777)
	if err != nil {
		panic(err.Error())

	}

	logFileP := path.Join(dirpath, name)
	dir, err := filepath.Abs(filepath.Dir(logFileP))

	logFile := dir + "/" + name

	logWfFile := logFile + ".wf"

	this.logFile, err = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(this.logFile)
	this.logFileWf, _ = os.OpenFile(logWfFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

}

func (this *UserLogger) AddPairInfo(field string, data string) {
	this.userBuf = append(this.userBuf, fmt.Sprintf("[%s:%s]", field, data)...)
}

func (this *UserLogger) AddSingleInfo(s string) {
	this.userBuf = append(this.userBuf, fmt.Sprintf("[%s]", s)...)

}

func (this *Logger) Debug(format string, v ...interface{}) {
	if this.GetLogLevel() > uint64(DEBUG) {
		return

	}
	this.write(this.logFile, "DEBUG", 2, string(this.buf), fmt.Sprintf(format, v...))

}

func (this *Logger) Trace(format string, v ...interface{}) {
	if this.GetLogLevel() > uint64(TRACE) {
		return

	}
	this.write(this.logFile, "TRACE", 2, string(this.buf), fmt.Sprintf(format, v...))

}

func (this *Logger) Fatal(format string, v ...interface{}) {
	if this.GetLogLevel() > uint64(FATAL) {
		return

	}
	s := fmt.Sprintf(format, v...)
	this.write(this.logFile, "FATAL", 2, string(this.buf), s)
	this.write(this.logFileWf, "FATAL", 2, string(this.buf), s)
	//panic(s)
}

func (this *Logger) Info(format string, v ...interface{}) {
	if this.GetLogLevel() > uint64(INFO) {
		return

	}
	this.write(this.logFile, "INFO", 2, string(this.buf), fmt.Sprintf(format, v...))

}

func (this *Logger) Warn(format string, v ...interface{}) {
	if this.GetLogLevel() > uint64(WARNING) {
		return

	}
	this.write(this.logFile, "WARNING", 2, string(this.buf), fmt.Sprintf(format, v...))

}

func (this *Logger) Flush() {
	this.logFile.Sync()
	this.logFileWf.Sync()

}

func (this *UserLogger) Put() {
	this.userBuf = this.userBuf[:0]
	g_bufPool.Put(this.userBuf)
}
