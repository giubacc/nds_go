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
	"nds/util"
)

type Config struct {
	StartNode        bool
	MulticastAddress string
	MulticastPort    uint
	ListeningPort    uint
	Val              string
	GetVal           bool

	LogType  string
	LogLevel string
}

type Peer struct {

	//the configuration
	Cfg Config

	//the currently timestamp set by this node
	CurrentNodeTS uint32

	//the desired timestamp this node would like to reach.
	//a successful synch with the cluster will transit DesiredClusterTS into CurrentNodeTS.
	DesiredClusterTS uint32

	//the value shared across the cluser
	Data string

	//exit required
	ExitRequired bool
}

func (p *Peer) Run() int {
	return 0
}

func (p *Peer) init() util.RetCode {
	return util.RetCode_OK
}

func (p *Peer) start() util.RetCode {
	return util.RetCode_OK
}

func (p *Peer) stop() util.RetCode {
	return util.RetCode_OK
}

func (p *Peer) processIncomingEvents() util.RetCode {
	return util.RetCode_OK
}

func (p *Peer) processNodeStatus() util.RetCode {
	return util.RetCode_OK
}

func (p *Peer) foreignEvent(evt *util.Event) bool {
	return true
}

func (p *Peer) processForeignEvent(evt *util.Event) bool {
	return true
}
