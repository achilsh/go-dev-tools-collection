package utils

import (
	"log"
	"os"
)

var logInst *log.Logger = nil
func init() {
	logInst = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}
func LogPrintln(data ...any) {
	logInst.Println(data...)
}
func LogPrintf(format string, data ...any) {
	logInst.Printf(format, data...)
}
