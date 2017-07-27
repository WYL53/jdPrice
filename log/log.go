package log

import (
	"log"
	"os"
	"time"
	"runtime"
	"fmt"
)

var logger *log.Logger
var logFile *os.File

func init() {
	fileName := time.Now().Format("2006-01-02") + ".log"
	logFile,err := os.OpenFile(fileName,os.O_CREATE|os.O_APPEND|os.O_WRONLY,os.ModePerm)
	if err != nil{
		panic(err)
	}
	logger = log.New(logFile,"",log.LstdFlags)
}

func Clear()  {
	logFile.Close()
}


func Println(args ...interface{})  {
	fn,line :=getFileNameAndLine()
	pre := fmt.Sprintf("%s[%d]",fn,line)
	logger.Println(pre,args)
}

func Printf(format string,args ...interface{})  {
	fn,line :=getFileNameAndLine()
	pre := fmt.Sprintf("%s[%d]",fn,line)
	logger.Printf("%s "+format,pre,args)
}

func getFileNameAndLine() (string,int) {
	_,fn,line,_ := runtime.Caller(2)
	return fn,line
}