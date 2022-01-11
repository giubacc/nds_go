/* Original Work Copyright (c) 2021 Giuseppe Baccini - giuseppe.baccini@live.com

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

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

func (lgr *Logger) Init(class string, cfg *Config) error {
	lgr.LgrLvl = LogLevelStr2LvL[cfg.LogLevel]
	lgr.Class = class
	if cfg.LogType == "console" {
		lgr.Out = os.Stdout
	} else {
		var err error
		if lgr.l_fd, err = os.OpenFile(class+cfg.LogType, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
			fmt.Println(err.Error())
			return &NDSError{RetCode_IOERR}
		}
		lgr.Out = lgr.l_fd
	}

	lgr.Lgr = log.New(lgr.Out, "", log.Lmsgprefix|log.Ltime|log.Lmicroseconds)
	return nil
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
