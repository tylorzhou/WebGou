package baaplogger

import (
	"fmt"
	"path/filepath"
	"time"
)

// RFC5424 log message levels.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

var levelPrefix = [LevelDebug + 1]string{"[Emergency] ", "[Alert] ", "[Critical] ", "[Error]", "[Warning] ", "[Notice] ", "[Informational] ", "[Debug] "}
var (
	//Log for baap service
	Log *Baaplogger
	//Ginlog for gin service
	Ginlog *Logger
)

func init() {
	initlogger()
}

//Baaplogger log for baap service
type Baaplogger struct {
	Level int
	Log   *Logger
}

func (bl *Baaplogger) writeMsg(logLevel int, msg string, v ...interface{}) (int, error) {

	when := time.Now()

	if len(v) > 0 {
		msg = fmt.Sprintf(msg, v...)
	}
	msg = when.Format("2006-01-02 15:04:05.000") + " " + levelPrefix[logLevel] + msg + "\n"

	return bl.Log.Write([]byte(msg))
}

// Emergency Log EMERGENCY level message.
func (bl *Baaplogger) Emergency(format string, v ...interface{}) {
	if LevelEmergency > bl.Level {
		return
	}
	bl.writeMsg(LevelEmergency, format, v...)
}

// Alert Log ALERT level message.
func (bl *Baaplogger) Alert(format string, v ...interface{}) {
	if LevelAlert > bl.Level {
		return
	}
	bl.writeMsg(LevelAlert, format, v...)
}

// Critical Log CRITICAL level message.
func (bl *Baaplogger) Critical(format string, v ...interface{}) {
	if LevelCritical > bl.Level {
		return
	}
	bl.writeMsg(LevelCritical, format, v...)
}

// Error Log ERROR level message.
func (bl *Baaplogger) Error(format string, v ...interface{}) {
	if LevelError > bl.Level {
		return
	}
	bl.writeMsg(LevelError, format, v...)
}

// Warning Log WARNING level message.
func (bl *Baaplogger) Warning(format string, v ...interface{}) {
	if LevelWarning > bl.Level {
		return
	}
	bl.writeMsg(LevelWarning, format, v...)
}

// Notice Log NOTICE level message.
func (bl *Baaplogger) Notice(format string, v ...interface{}) {
	if LevelNotice > bl.Level {
		return
	}
	bl.writeMsg(LevelNotice, format, v...)
}

// Informational Log INFORMATIONAL level message.
func (bl *Baaplogger) Informational(format string, v ...interface{}) {
	if LevelInformational > bl.Level {
		return
	}
	bl.writeMsg(LevelInformational, format, v...)
}

// Debug Log DEBUG level message.
func (bl *Baaplogger) Debug(format string, v ...interface{}) {
	if LevelDebug > bl.Level {
		return
	}
	bl.writeMsg(LevelDebug, format, v...)
}

func initlogger() {
	dir, err := filepath.Abs("log")
	if err != nil {
		panic(err)
	}

	// this for baap API log
	Log = &Baaplogger{
		Level: LevelDebug,
		Log: &Logger{
			Filename:   filepath.Join(dir, "baap.log"),
			MaxSize:    500, // megabytes
			MaxBackups: 6,
			MaxAge:     28, // days
		},
	}

	// this for log inforamtion in gin
	Ginlog = &Logger{
		Filename:   filepath.Join(dir, "gin.log"),
		MaxSize:    500, // megabytes
		MaxBackups: 6,
		MaxAge:     28, // days
	}

}
