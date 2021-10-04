package net

import (
	"fmt"
	"nds/util"
	"net"
	"time"
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
	//configuration
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

	a.Status = AcceptorStatus_ACCEPT

	for a.Status == AcceptorStatus_ACCEPT {
		a.logger.Trace("accepting...")

		_, err := a.listener.Accept()
		if err != nil {
			a.logger.Trace("err:%s, accepting connection ...", err.Error())
		}

		time.Sleep(time.Second * 2)
	}
}
