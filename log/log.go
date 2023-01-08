/* 带颜色的日志包 */
package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

var (
	// log.Lshortfile用于显示文件名和代码行号
	// [error]为红色
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m", log.LstdFlags|log.Lshortfile)
	// [info ]为蓝色
	infoLog = log.New(os.Stdout, "\033[34m[info ]\033[0m", log.LstdFlags|log.Lshortfile)
	loggers = []*log.Logger{errorLog, infoLog}
	mu      sync.Mutex
)

// 暴露出去的方法
var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

// LogLevels
const (
	InfoLevel = iota
	ErrorLevel
	Disable
)

// 控制日志等级(线程安全)
func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		// 如果日志等级小于设置等级，将日志输出重定向到ioutil.Discard，即不打印该等级日志
		errorLog.SetOutput(ioutil.Discard)
	}

	if InfoLevel < level {
		// 如果日志等级小于设置等级，将日志输出重定向到ioutil.Discard，即不打印该等级日志
		infoLog.SetOutput(ioutil.Discard)
	}
}
