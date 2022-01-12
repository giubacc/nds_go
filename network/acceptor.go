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

package network

import (
	"fmt"
	"nds/util"
	"net"
)

type AcceptorStatus int

type Acceptor struct {
	//config
	Cfg *util.Config

	//status
	Status AcceptorStatus

	//listening port
	lport uint

	//listener
	listener net.Listener

	//channel used to serve incoming TCP connections
	EnteringChan chan net.Conn

	//logger
	logger util.Logger
}

func (a *Acceptor) Run() error {
	if err := a.init(); err != nil {
		return err
	}
	a.accept()
	return a.stop()
}

func (a *Acceptor) init() error {
	a.lport = a.Cfg.ListeningPort

	//logger init
	err := a.logger.Init("acpt.", a.Cfg)
	return err
}

func (a *Acceptor) stop() error {
	return a.listener.Close()
}

func (a *Acceptor) accept() {
	for {
		var err error
		if a.listener, err = net.Listen("tcp", fmt.Sprintf("%s%d", ":", a.lport)); err == nil {
			break
		} else {
			a.lport++
			a.logger.Trace("err:%s, try auto-adjusting listening port to:%d ...", err.Error(), a.lport)
		}
	}

	a.logger.Trace("accepting ...")
	for /*@fixme*/ {
		if conn, err := a.listener.Accept(); err != nil {
			a.logger.Err("err:%s, accepting connection ...", err.Error())
		} else {
			a.logger.Trace("connection accepted")
			a.EnteringChan <- conn
		}
	}

	a.logger.Trace("stop accepting")
}
