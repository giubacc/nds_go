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
	"nds/net"
	"nds/util"
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

	//network selector
	selector net.Selector

	//network acceptor
	acceptor net.Acceptor

	//logger
	logger util.Logger
}

func (p *Peer) Run() util.RetCode {
	var rcode util.RetCode = util.RetCode_OK

	rcode = p.init()
	if rcode != util.RetCode_OK {
		return rcode
	}

	rcode = p.start()
	if rcode != util.RetCode_OK {
		return rcode
	}

	p.processEvents()

	return rcode
}

func (p *Peer) init() util.RetCode {
	var rcode util.RetCode = util.RetCode_OK

	//logger init
	p.logger.Init("peer.", p.Cfg)

	//seconds before this node will auto generate the timestamp
	p.TpInitialSynchWindow = time.Now().Add(time.Second * NodeSynchDuration)

	p.acceptor.Cfg = p.Cfg
	p.selector.Cfg = p.Cfg

	return rcode
}

func (p *Peer) start() util.RetCode {
	var rcode util.RetCode = util.RetCode_OK

	p.logger.Trace("starting acceptor")
	go p.acceptor.Run()
	p.logger.Trace("wait acceptor go accepting")

	p.logger.Trace("starting selector")
	go p.selector.Run()
	p.logger.Trace("wait selector go selecting")

	return rcode
}

func (p *Peer) stop() util.RetCode {
	var rcode util.RetCode = util.RetCode_OK

	p.logger.Stop()
	return rcode
}

func (p *Peer) processEvents() util.RetCode {
	var rcode util.RetCode = util.RetCode_OK

	for !p.ExitRequired {
		p.logger.Trace("processing events...")
		time.Sleep(time.Second * 2)
	}

	return rcode
}

func (p *Peer) processNodeStatus() util.RetCode {
	var rcode util.RetCode = util.RetCode_OK

	return rcode
}

func (p *Peer) foreignEvent(evt *util.Event) bool {
	return true
}

func (p *Peer) processForeignEvent(evt *util.Event) bool {
	return true
}
