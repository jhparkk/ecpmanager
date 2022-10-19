package prdebug

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

var PrDebug bool = false

func Println(v ...any) {
	if !PrDebug {
		return
	}
	pc, fullFileName, line, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()

	fileName := fullFileName[strings.LastIndex(fullFileName, "/")+1:]
	funcName := callerName[strings.LastIndex(callerName, "/")+1:]
	header := "<" + funcName + ":" + fileName + ":" + strconv.Itoa(line) + ">"
	fmt.Println(header, v)
}

func Printf(format string, v ...any) {
	if !PrDebug {
		return
	}
	pc, fullFileName, line, _ := runtime.Caller(1)
	callerName := runtime.FuncForPC(pc).Name()

	fileName := fullFileName[strings.LastIndex(fullFileName, "/")+1:]
	funcName := callerName[strings.LastIndex(callerName, "/")+1:]

	body := fmt.Sprintf(format, v...)
	header := "<" + funcName + ":" + fileName + ":" + strconv.Itoa(line) + ">"
	fmt.Println(header, body)
}
