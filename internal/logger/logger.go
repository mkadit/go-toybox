package logfile

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mkadit/go-toybox/common/config"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var AppLogger Loggers

var (
	// Confiture the ../../../ as need
	ProjectRootPath  = filepath.Join(filepath.Dir(b), "../../")
	_, b, _, _       = runtime.Caller(0)
	FormatTimeforLog = "2006-01-02 15:04:05.000"
	TimeFormat       = "2006-01-02"
	Testing          bool
)

// FormatTimeforLog format for Logging File

type Loggers struct {
	MessageLogger  MultLogger
	EventLogger    MultLogger
	ErrorLogger    MultLogger
	HTTPLogger     MultLogger
	CriticalLogger MultLogger
}

type MultLogger struct {
	ZeroLog zerolog.Logger
	StdLog  *log.Logger
}

type line struct {
	Message  string
	Severity string
}

type logFileWriter struct {
	w        io.Writer
	severity string
}

// Write Write method for standard log
func (e logFileWriter) Write(p []byte) (int, error) {
	l := &line{
		Message:  string(p),
		Severity: e.severity,
	}
	// Format the message with desired elements (timestamp, severity, message)
	formattedMessage := fmt.Sprintf("[%s] %s", l.Severity, l.Message)

	// Write the formatted message to the underlying writer
	n, err := e.w.Write([]byte(formattedMessage))
	if err != nil {
		return n, err
	}
	if n != len(formattedMessage) {
		return n, io.ErrShortWrite
	}
	return len(p), nil
}

// nextMidnight function to midnight
func nextMidnight() time.Duration {
	now := time.Now()
	nextMidnight := now.Truncate(24*time.Hour).AddDate(0, 0, 1)
	return nextMidnight.Sub(now)
}

// createIoWriter Crate IO writer for logs
func createIoWriter(files map[string]*lumberjack.Logger) map[string][]io.Writer {
	writers := make(map[string][]io.Writer)
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	for k, v := range files {
		if k == "debug" {
			continue
		}
		writers[k] = []io.Writer{v, consoleWriter, files["debug"]}
		// if os.Getenv("ENV") != "prod" {
		// 	writers[k] = append(writers[k], files["debug"])
		// }
	}
	writers["critical"] = []io.Writer{files["error"], files["debug"], consoleWriter}
	return writers
}

// CreateLogger Create logger in background
func CreateLogger() {
	err := SetLogger()
	if err != nil {
		LogFatal(err, "SYSTEM: logger cannot be set")
	}
	for {
		nextRotation := nextMidnight()
		time.Sleep(nextRotation)

		// Close and reopen logger to trigger rotation
		err := SetLogger()
		if err != nil {
			LogFatal(err, "SYSTEM: logger cannot be set")
		}

	}
}

// CreateLogWriter Create log writer for rotational logs
func CreateLogWriter(filepath string, maxSize int, maxBackup int, maxAge int, compress bool) *lumberjack.Logger {
	logFolder := ProjectRootPath + "/logs"
	if _, err := os.Stat(logFolder); os.IsNotExist(err) {
		os.Mkdir(logFolder, 0777)
	}

	filename := fmt.Sprintf("%s/%s.%s.log", logFolder, filepath, time.Now().Format(TimeFormat)) // Path to your log file
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}
}

// SetLogger Set logger for app
func SetLogger() (err error) {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.ErrorFieldName = "err"

	logFiles := make(map[string]*lumberjack.Logger)

	logMessageFile := CreateLogWriter("Message/Message", 80, 5, 0, false)
	logEventFile := CreateLogWriter("event/Event", 80, 5, 0, false)
	logErrorFile := CreateLogWriter("error/Error", 80, 5, 0, false)
	logHTTPFile := CreateLogWriter("HTTP/HTTP", 80, 5, 0, false)
	logDebugFile := CreateLogWriter("Debug/Debug", 80, 5, 0, false)
	_ = logDebugFile

	logFiles["debug"] = logDebugFile
	logFiles["event"] = logEventFile
	logFiles["error"] = logErrorFile
	logFiles["ws"] = logMessageFile
	logFiles["http"] = logHTTPFile

	multWriters := createIoWriter(logFiles)
	// log level for std files -> key = file , value = message
	logLevel := map[string]string{
		"info": "INFO", "event": "EVENT", "error": "ERROR",
		"http": "HTTP", "ws": "message", "critical": "CRITICAL",
	}

	// create zerologWriters
	multLoggers := make(map[string]MultLogger)
	for i, writer := range multWriters {
		if i == "debug" {
			continue
		}
		zerologWriter := zerolog.New(zerolog.MultiLevelWriter(writer[1:]...))
		// stdWriters := io.MultiWriter(writer[0:0]...)
		stdLogFile := logFileWriter{writer[0], logLevel[i]}
		stdLogWriter := log.New(stdLogFile, "", 0)
		multLogger := MultLogger{
			ZeroLog: zerologWriter,
			StdLog:  stdLogWriter,
		}
		multLoggers[i] = multLogger
	}
	AppLogger = Loggers{
		MessageLogger:  multLoggers["message"],
		EventLogger:    multLoggers["event"],
		ErrorLogger:    multLoggers["error"],
		CriticalLogger: multLoggers["critical"],
		HTTPLogger:     multLoggers["http"],
	}
	return
}

// LogMsgHTTP Log an HTTP request, response without struct
func LogMsgHTTP(m MessageLog) {
	if Testing {
		return
	}
	timeNow := time.Now().Format(config.FORMAT_TIME_FOR_LOG)
	AppLogger.HTTPLogger.ZeroLog.Info().Timestamp().Str("request-id", m.InternalID).Msgf(m.Msg)
	AppLogger.HTTPLogger.StdLog.Printf("[%s][%s][%s][%s][%d][%s][%s][%s][%s][%s][%s][%s]",
		timeNow, m.SystemName, m.InternalID, m.ReffTrx, m.Step, m.Flow, m.Entity, m.RC, m.TypeTrx, m.Header, m.URL, m.Msg)
}

// LogMsgInterfaceHTTP Log an HTTP request, response with struct
func LogMsgInterfaceHTTP[T any](m MessageLog, object *T) {
	if Testing {
		return
	}
	timeNow := time.Now().Format(config.FORMAT_TIME_FOR_LOG)
	AppLogger.HTTPLogger.ZeroLog.Info().Timestamp().Str("request-id", m.InternalID).Str("type", fmt.Sprintf("%T", object)).Interface("content", object).Msgf(m.Msg)
	AppLogger.HTTPLogger.StdLog.Printf("[%s][%s][%s][%s][%d][%s][%s][%s][%s][%s][%s][%s]",
		timeNow, m.SystemName, m.InternalID, m.ReffTrx, m.Step, m.Flow, m.Entity, m.RC, m.TypeTrx, m.Header, m.URL, m.Msg)
}

// LogMsgHTTP Log an HTTP request, response without struct
func LogMsgEventHTTP[T any](m MessageLog, attr ...any) {
	if Testing {
		return
	}
	timeNow := time.Now().Format(config.FORMAT_TIME_FOR_LOG)
	AppLogger.HTTPLogger.ZeroLog.Info().Timestamp().Str("request-id", m.InternalID).Msgf(m.Msg)
	AppLogger.HTTPLogger.StdLog.Printf("[%s][%s][%s][%s][%d][%s][%s][%s][%s][%s][%s][%s]",
		timeNow, m.SystemName, m.InternalID, m.ReffTrx, m.Step, m.Flow, m.Entity, m.RC, m.TypeTrx, m.Header, m.URL, m.Msg)
}

// LogMsgInterfaceHTTP Log an HTTP request, response with struct
func LogMsgEventInterfaceHTTP[T any](m MessageLog, object *T, attr ...any) {
	if Testing {
		return
	}
	timeNow := time.Now().Format(config.FORMAT_TIME_FOR_LOG)
	AppLogger.HTTPLogger.ZeroLog.Info().Timestamp().Str("request-id", m.InternalID).Str("type", fmt.Sprintf("%T", object)).Interface("content", object).Msgf(m.Msg)
	AppLogger.HTTPLogger.StdLog.Printf("[%s][%s][%s][%s][%d][%s][%s][%s][%s][%s][%s][%s]",
		timeNow, m.SystemName, m.InternalID, m.ReffTrx, m.Step, m.Flow, m.Entity, m.RC, m.TypeTrx, m.Header, m.URL, m.Msg)
}

// LogEvent Log application condition, get new configuration, etc
func LogEvent(message string, attr ...any) {
	if Testing {
		return
	}
	timeNow := time.Now().Format(config.FORMAT_TIME_FOR_LOG)
	msg := fmt.Sprintf(message, attr...)
	AppLogger.EventLogger.ZeroLog.Info().Timestamp().Msgf(message, attr...)
	AppLogger.EventLogger.StdLog.Printf("[%s][][%s]", timeNow, msg)
}

// LogEventInterface Log events along with the desired interface
func LogEventInterface[T any](object *T, message string, attr ...any) {
	if Testing {
		return
	}
	timeNow := time.Now().Format(config.FORMAT_TIME_FOR_LOG)
	msg := fmt.Sprintf(message, attr...)
	AppLogger.EventLogger.ZeroLog.Info().Timestamp().Str("type", fmt.Sprintf("%T", object)).Interface("content", object).Msgf(message, attr...)
	AppLogger.EventLogger.StdLog.Printf("[%s][][%T|%+v|%s]", timeNow, object, object, msg)
}

// LogEventInterfaceHTTP Log an HTTP request events along with the desired interface
func LogEventInterfaceHTTP[T any](uuid string, object *T, message string, attr ...any) {
	if Testing {
		return
	}
	timeNow := time.Now().Format(config.FORMAT_TIME_FOR_LOG)
	msg := fmt.Sprintf(message, attr...)
	AppLogger.EventLogger.ZeroLog.Info().Timestamp().Str("request-id", uuid).Str("type", fmt.Sprintf("%T", object)).Interface("content", object).Msgf(message, attr...)
	AppLogger.EventLogger.StdLog.Printf("[%s][%s][%T|%+v|%s]", timeNow, uuid, object, object, msg)
}

// LogInfoHTTP Used for logging any info need to be placed for HTTP
func LogInfoHTTP(uuid string, message string, attr ...any) {
	if Testing {
		return
	}
	timeNow := time.Now().Format(config.FORMAT_TIME_FOR_LOG)
	msg := fmt.Sprintf(message, attr...)
	AppLogger.HTTPLogger.ZeroLog.Info().Timestamp().Str("request-id", uuid).Msgf(message, attr...)
	AppLogger.HTTPLogger.StdLog.Printf("[%s][][%s]", timeNow, msg)
}

// LogInfoInterfaceHTTP Log an HTTP request, response, and other structs of an id
func LogInfoInterfaceHTTP[T any](uuid string, object *T, message string, attr ...any) {
	if Testing {
		return
	}
	timeNow := time.Now().Format(config.FORMAT_TIME_FOR_LOG)
	msg := fmt.Sprintf(message, attr...)
	AppLogger.HTTPLogger.ZeroLog.Info().Timestamp().Str("request-id", uuid).Str("type", fmt.Sprintf("%T", object)).Interface("content", object).Msgf(message, attr...)
	AppLogger.HTTPLogger.StdLog.Printf("[%s][][%T|%+v|%s]", timeNow, object, object, msg)
}

// LogErrorHTTP Log an error for and id
func LogErrorHTTP(uuid string, err any, message string, attr ...any) {
	if Testing {
		return
	}
	timeNow := time.Now().Format(config.FORMAT_TIME_FOR_LOG)
	msg := fmt.Sprintf(message, attr...)
	AppLogger.ErrorLogger.ZeroLog.Err(fmt.Errorf(message, attr...)).Timestamp().Str("request-id", uuid).Msgf("%s", err)
	AppLogger.ErrorLogger.StdLog.Printf("[%s][%s][%s]", timeNow, uuid, msg)
}

func LogErrorEvent(err any, message string, attr ...any) {
	if Testing {
		return
	}
	timeNow := time.Now().Format(config.FORMAT_TIME_FOR_LOG)
	msg := fmt.Sprintf(message, attr...)
	AppLogger.EventLogger.ZeroLog.Err(fmt.Errorf(message, attr...)).Timestamp().Msgf("%s", err)
	AppLogger.EventLogger.StdLog.Printf("[%s][][%s]", timeNow, msg)
}

// LogCritical Log base application critical errors, like connection DB, TCP
// It doesn't automatically do an os.Exit(1) but does print an error with stack
func LogCritical(err error, message string, attr ...any) {
	if Testing {
		return
	}
	timeNow := time.Now().Format(config.FORMAT_TIME_FOR_LOG)
	msg := fmt.Sprintf(message, attr...)
	AppLogger.CriticalLogger.ZeroLog.Err(fmt.Errorf(message, attr...)).Timestamp().Msgf("%s", err)
	AppLogger.CriticalLogger.StdLog.Printf("[%s][%s][%+v]", timeNow, msg, WithStack(err))
}

// LogFatal Log base application critical errors and do an os.Exit(1)
func LogFatal(err error, message string, attr ...any) {
	if Testing {
		return
	}
	timeNow := time.Now().Format(config.FORMAT_TIME_FOR_LOG)
	msg := fmt.Sprintf(message, attr...)
	AppLogger.CriticalLogger.ZeroLog.Err(fmt.Errorf(message, attr...)).Timestamp().Msgf("%s", err)
	AppLogger.CriticalLogger.StdLog.Fatalf("[%s][%s][%s]", timeNow, msg, err)
}
