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

package peer

import (
	"encoding/binary"
	"nds/network"
	"nds/util"
	"net"
	"time"
)

const NodeSynchDuration = 2

type Peer struct {
	//configuration
	Cfg util.Config

	//the time point at which this node will generate itself the timestamp;
	//this will happen if no other node will respond to initial alive sent by this node.
	//not daemon nodes ("pure" setter or getter nodes) will shutdown at this time point.
	TpInitialSynchWindow time.Time

	//the currently timestamp set by this node
	CurrentNodeTS uint32

	//the desired timestamp this node would like to reach.
	//a successful synch with the cluster will transit DesiredClusterTS into CurrentNodeTS.
	DesiredClusterTS uint32

	//the value shared across the cluster
	Data string

	//exit required
	ExitRequired bool

	//network acceptor
	acceptor network.Acceptor

	//multicast manager
	mcastHelper network.MCastHelper

	//channel used to serve incoming TCP connections
	EnteringChan chan net.Conn

	//channels used to send/receive alive messages (UDP multicast)
	AliveChanIncoming chan []byte
	AliveChanOutgoing chan []byte

	//logger
	logger util.Logger
}

func (p *Peer) genTS() {
	p.DesiredClusterTS = uint32(time.Now().Unix())
	p.CurrentNodeTS = p.DesiredClusterTS
}

func (p *Peer) Run() error {

	if err := p.init(); err != nil {
		return err
	}

	if err := p.start(); err != nil {
		return err
	}

	if p.Cfg.Val != "" {
		p.Data = p.Cfg.Val
		p.genTS()
	}

	return p.processEvents()
}

func (p *Peer) init() error {
	//logger init
	if err := p.logger.Init("peer.", &p.Cfg); err != nil {
		return err
	}

	p.EnteringChan = make(chan net.Conn)
	p.AliveChanIncoming = make(chan []byte)
	p.AliveChanOutgoing = make(chan []byte)

	//seconds before this node will auto generate the timestamp
	p.TpInitialSynchWindow = time.Now().Add(time.Second * NodeSynchDuration)

	p.acceptor.Cfg = &p.Cfg
	p.acceptor.EnteringChan = p.EnteringChan

	p.mcastHelper.Cfg = &p.Cfg
	p.mcastHelper.AliveChanIncoming = p.AliveChanIncoming
	p.mcastHelper.AliveChanOutgoing = p.AliveChanOutgoing

	return nil
}

func (p *Peer) start() error {

	p.logger.Trace("starting acceptor ...")
	go p.acceptor.Run()
	//@fixme wait

	p.logger.Trace("starting multicast ...")
	go p.mcastHelper.Run()
	//@fixme wait

	return nil
}

func (p *Peer) stop() error {
	p.logger.Stop()
	return nil
}

func (p *Peer) processEvents() error {
	p.logger.Trace("start processing events ...")

	interrupter := time.NewTicker(time.Second * 2)
	defer interrupter.Stop()

out:
	for {
		select {
		case <-interrupter.C:
			if err := p.processNodeStatus(); err != nil && err.Code == util.RetCode_EXIT {
				break out
			}
		case conn := <-p.EnteringChan:
			p.sendDataMessage(conn)
		case buff := <-p.AliveChanIncoming:
			msgUB := 4 + binary.LittleEndian.Uint32(buff[0:])
			msgStr := string(buff[4:msgUB])
			p.logger.Trace("evt:%s", msgStr)
			msg := util.AliveMsg{}
			msg.UnmarshalJSON(buff[4:msgUB])
		}
	}

	p.logger.Trace("end process events")
	return nil
}

func (p *Peer) processNodeStatus() *util.NDSError {
	now := time.Now()

	if p.ExitRequired {
		return &util.NDSError{Code: util.RetCode_EXIT}
	}

	//"pure" setter or getter nodes must shutdown.
	if !p.Cfg.StartNode && now.After(p.TpInitialSynchWindow) {
		return &util.NDSError{Code: util.RetCode_EXIT}
	}

	//if no other node has still responded to initial alive, the node generates itself the timestamp
	if p.CurrentNodeTS == 0 && p.DesiredClusterTS == 0 && now.After(p.TpInitialSynchWindow) {
		p.genTS()
		p.logger.Trace("auto generated timestamp: %d", p.CurrentNodeTS)
		p.sendAliveMessage()
	}

	return nil
}

func (p *Peer) buildAliveMessage() ([]byte, error) {
	msg := util.AliveMsg{Lp: uint16(p.acceptor.ListenPort), Pt: util.MsgPktTypeAlive, Si: p.acceptor.Listener.Addr().String(), Ts: uint64(p.CurrentNodeTS)}
	return msg.MarshalJSON()
}

func (p *Peer) sendAliveMessage() error {
	if msg, err := p.buildAliveMessage(); err != nil {
		p.logger.Err("building alive msg:%s", err.Error())
		return err
	} else {
		p.AliveChanOutgoing <- msg
	}
	return nil
}

func (p *Peer) buildDataMessage() ([]byte, error) {
	msg := util.DataMsg{Dv: p.Data, Pt: util.MsgPktTypeData, Ts: uint64(p.CurrentNodeTS)}
	return msg.MarshalJSON()
}

func (p *Peer) sendDataMessage(conn net.Conn) error {

	if msg, err := p.buildDataMessage(); err != nil {
		p.logger.Err("building data msg:%s", err.Error())
		return err
	} else {
		sent, err := conn.Write(msg)
		if err != nil {
			p.logger.Err("sending data msg:%s", err.Error())
		}
		p.logger.Trace("sent %d bytes to: %s", sent, conn.RemoteAddr().String())
	}

	return nil
}
