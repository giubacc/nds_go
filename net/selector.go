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

package net

import (
	"nds/util"
	"time"
)

type SelectorStatus int

const (
	SelectorStatus_ZERO     = iota // zero status
	SelectorStatus_INIT            // selector is inited
	SelectorStatus_READY           // selector has completed all init steps and it is ready to select
	SelectorStatus_SELECT          // selector is currently selecting over sockets
	SelectorStatus_STOPPING        // selector is stopping
	SelectorStatus_STOPPED         // selector is stopped and can be safely disposed
	SelectorStatus_ERROR
)

type Selector struct {
	//configuration
	Cfg util.Config

	//status
	Status SelectorStatus

	//logger
	logger util.Logger
}

func (s *Selector) Run() util.RetCode {
	var rcode util.RetCode = util.RetCode_OK

	rcode = s.init()
	if rcode != util.RetCode_OK {
		return rcode
	}

	s.doSelect()

	rcode = s.stop()
	return rcode
}

func (s *Selector) init() util.RetCode {
	var rcode util.RetCode = util.RetCode_OK

	//logger init
	s.logger.Init("sel.", s.Cfg)

	s.Status = SelectorStatus_INIT
	return rcode
}

func (s *Selector) stop() util.RetCode {
	var rcode util.RetCode = util.RetCode_OK
	s.Status = SelectorStatus_STOPPING

	s.Status = SelectorStatus_STOPPED
	return rcode
}

func (s *Selector) doSelect() util.RetCode {
	var rcode util.RetCode = util.RetCode_OK
	s.Status = SelectorStatus_SELECT

	for s.Status == SelectorStatus_SELECT {
		s.logger.Trace("selecting...")
		time.Sleep(time.Second * 2)
	}

	return rcode
}
