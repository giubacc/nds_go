package util

import (
	"fmt"
	"io"
	"log"
	"os"
)

type LogLevel int

const (
	Trace = iota
	Info
	Warn
	Err
	Critical
	Off
)

type LogLevelStr string

const (
	TraceStr    = "trace"
	InfoStr     = "info"
	WarnStr     = "warn"
	ErrStr      = "err"
	CriticalStr = "critical"
	OffStr      = "off"
)

var LogLevelStr2LvL = map[string]LogLevel{TraceStr: Trace, InfoStr: Info, WarnStr: Warn, ErrStr: Err, CriticalStr: Critical, OffStr: Off}

type Logger struct {
	Lgr    *log.Logger
	LgrLvl LogLevel
	Class  string
	Out    io.Writer

	//opt logger fd
	l_fd *os.File
}

func (lgr *Logger) Init(class string, cfg Config) RetCode {
	lgr.LgrLvl = LogLevelStr2LvL[cfg.LogLevel]
	lgr.Class = class
	if cfg.LogType == "console" {
		lgr.Out = os.Stdout
	} else {
		var err error
		lgr.l_fd, err = os.OpenFile(class+cfg.LogType, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Println(err.Error())
			return RetCode_IOERR
		}
		lgr.Out = lgr.l_fd
	}

	lgr.Lgr = log.New(lgr.Out, "", log.Lmsgprefix|log.Ltime|log.Lmicroseconds)
	return RetCode_OK
}

func (lgr *Logger) Stop() {
	if lgr.l_fd != nil {
		lgr.l_fd.Close()
	}
}

func (lgr *Logger) Trace(format string, v ...interface{}) {
	if lgr.LgrLvl > Trace {
		return
	}
	lgr.Lgr.SetPrefix("[" + lgr.Class + TraceStr + "] ")
	lgr.Lgr.Printf(format, v...)
}

func (lgr *Logger) Info(format string, v ...interface{}) {
	if lgr.LgrLvl > Info {
		return
	}
	lgr.Lgr.SetPrefix("[" + lgr.Class + InfoStr + "] ")
	lgr.Lgr.Printf(format, v...)
}

func (lgr *Logger) Warn(format string, v ...interface{}) {
	if lgr.LgrLvl > Warn {
		return
	}
	lgr.Lgr.SetPrefix("[" + lgr.Class + WarnStr + "] ")
	lgr.Lgr.Printf(format, v...)
}

func (lgr *Logger) Err(format string, v ...interface{}) {
	if lgr.LgrLvl > Err {
		return
	}
	lgr.Lgr.SetPrefix("[" + lgr.Class + ErrStr + "] ")
	lgr.Lgr.Printf(format, v...)
}

func (lgr *Logger) Critical(format string, v ...interface{}) {
	if lgr.LgrLvl > Critical {
		return
	}
	lgr.Lgr.SetPrefix("[" + lgr.Class + CriticalStr + "] ")
	lgr.Lgr.Printf(format, v...)
}
