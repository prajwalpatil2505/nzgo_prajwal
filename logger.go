package nzgo

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Fatal *log.Logger
)

/*Valid log levels : DEBUG, INFO, FATAL, OFF*/
type PDALogger struct {
	logLevel string
}

var elog PDALogger

/* Create logger handler with some predefined prefix setting,
 * this will be overwritten in actual logging.
 * Mostly this setting is not used
 */
func Init() {

	Debug = log.New(ioutil.Discard,
		"DEBUG: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)

	Info = log.New(ioutil.Discard,
		"INFO: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)

	Fatal = log.New(ioutil.Discard,
		"FATAL: ",
		log.Ldate|log.Lmicroseconds|log.Lshortfile)
}

/* Initialize logger and set output to file */
func (elog PDALogger) initialize() {
	/* Set loglevel here, invalid loglevel will discard all log messages */
	elog.logLevel = "DEBUG" //This is default log level

	/* Overwrite log level mentioned in conf, if its blank use default case */
	if configuration.LogLevel != "" {
		elog.logLevel = strings.ToUpper(configuration.LogLevel)
	}

	if elog.logLevel == "OFF" {
		Init() //It will initialize and discard all stream output. Log file won't be created
		return
	}
	fname := fmt.Sprintf("nzgolang_nz%d.log", os.Getpid())

	/* Open file with permissions USER:read and write; GROUP&OTHERS:read */
	fh, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		errorf("Error opening logger file")
	}
	logBanner(fh)

	Init()
	switch elog.logLevel {
	// Sequence of log level case matters. Should not be changed
	case "DEBUG":
		Debug.SetOutput(fh)
		fallthrough
	case "INFO":
		Info.SetOutput(fh)
		fallthrough
	case "FATAL":
		Fatal.SetOutput(fh)

		// case default : //It will do nothing to discard the log output but log file with banner will be generated
	}

}

/* Used to write banner for logger. ToDo:Add more info related to server */
func logBanner(fh io.Writer) {
	fmt.Fprintln(fh, "---------------- IBM PDA Log -----------------")
	fmt.Fprintln(fh, "----------------------------------------------")

}

/* Prefix string created with specific format, will be used in debug and fatal logging */
func prefixString() string {
	prefixStr := fmt.Sprintf("%s [%d] ", time.Now().UTC().Format("2006-01-02 15:04:05 EST"), os.Getpid())
	return prefixStr
}

/* Wrappers to print functions to change pefix format */
func (elog PDALogger) Debugf(fname string, s string, args ...interface{}) {
	prefixStr := prefixString() + "[DEBUG] " + fname + " "
	Debug.SetFlags(0)
	Debug.SetPrefix(prefixStr)

	Debug.Printf(s, args...)
}

/* Used for adding debug log without format */
func (elog PDALogger) Debugln(args ...interface{}) {
	prefixStr := prefixString() + "[DEBUG] "
	Debug.SetFlags(0)
	Debug.SetPrefix(prefixStr)

	Debug.Println(args...)
}

/* Info logger adds messages for client */
func (elog PDALogger) infof(s string, args ...interface{}) {
	prefixStr := prefixString() + "[INFO] : "
	Info.SetFlags(0)
	Info.SetPrefix(prefixStr)

	Info.Printf(s, args...)
}

func (elog PDALogger) Infoln(args ...interface{}) {
	prefixStr := prefixString() + "[INFO] : "
	Info.SetFlags(0)
	Info.SetPrefix(prefixStr)

	Info.Println(args...)
}

/* Fatal logs error and panic, forcing application to exit with message on stdout */
func (elog PDALogger) Fatalf(fname string, s string, args ...interface{}) {
	prefixStr := prefixString() + "[FATAL] " + fname + " "
	Fatal.SetFlags(0)
	Fatal.SetPrefix(prefixStr)

	Fatal.Panic(fmt.Sprintf(s, args...))
}

func (elog PDALogger) Fatalln(args ...interface{}) {
	prefixStr := prefixString() + "[FATAL] "
	Fatal.SetFlags(0)
	Fatal.SetPrefix(prefixStr)

	Fatal.Panic(args...)
}

/* Function name is determined from caller stack at runtime
 * Returns function name with line number
 */
func funName(depthList ...int) string {
	var depth int
	if depthList == nil {
		depth = 1
	} else {
		depth = depthList[0]
	}
	function, _, line, _ := runtime.Caller(depth)
	return fmt.Sprintf("%s %d :", runtime.FuncForPC(function).Name(), line)
}

/* Returns only <filename>.<function_name> and <line_number> */
func chopPath(orig string) string {
	ind := strings.LastIndex(orig, "/")
	if ind == -1 {
		return orig
	} else {
		return orig[ind+1:]
	}
}