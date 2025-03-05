package utils

import (
	"fmt"
	"log"
	"os"
)

var logInst *log.Logger = nil

func init() {
	logInst = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}
func LogPrintln(data ...any) {
	logInst.Output(2, fmt.Sprintln(data...))
}
func LogPrintf(format string, data ...any) {
	logInst.Output(2, fmt.Sprintf(format, data...))
}
