package network

import (
	"fmt"
	"nds/util"
	"net"
)

type AcceptorStatus int

const (
	AcceptorStatus_ZERO     = iota // zero status
	AcceptorStatus_INIT            // acceptor is inited
	AcceptorStatus_ACCEPT          // acceptor is accepting
	AcceptorStatus_STOPPING        // acceptor is stopping
	AcceptorStatus_STOPPED         // acceptor is stopped and can be safely disposed
	AcceptorStatus_ERROR
)

type Acceptor struct {
	//parent
	Cfg util.Config

	//status
	Status AcceptorStatus

	//listening port
	lport uint

	//listener
	listener net.Listener

	//logger
	logger util.Logger
}

func (a *Acceptor) Run() util.RetCode {
	var rcode util.RetCode = util.RetCode_OK

	rcode = a.init()
	if rcode != util.RetCode_OK {
		return rcode
	}

	a.accept()

	rcode = a.stop()
	return rcode
}

func (a *Acceptor) init() util.RetCode {
	var rcode util.RetCode = util.RetCode_OK

	a.lport = a.Cfg.ListeningPort

	//logger init
	a.logger.Init("acpt.", a.Cfg)

	return rcode
}

func (a *Acceptor) stop() util.RetCode {
	var rcode util.RetCode = util.RetCode_OK

	a.listener.Close()
	return rcode
}

func (a *Acceptor) accept() {

	for {
		var err error
		a.listener, err = net.Listen("tcp", fmt.Sprintf("%s%d", ":", a.lport))
		if err == nil {
			break
		} else {
			a.lport++
			a.logger.Trace("err:%s, try auto-adjusting listening port to:%d ...", err.Error(), a.lport)
		}
	}

	a.logger.Trace("accepting ...")
	a.Status = AcceptorStatus_ACCEPT
	for a.Status == AcceptorStatus_ACCEPT {
		conn, err := a.listener.Accept()
		if err != nil {
			a.logger.Err("err:%s, accepting connection ...", err.Error())
		} else {
			a.logger.Trace("connection accepted")
			util.EnteringChan <- conn
		}
	}
	a.logger.Trace("stop accepting")
}
