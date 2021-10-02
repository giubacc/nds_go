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

package main

import (
	"flag"
	"nds/peer"
)

var pr peer.Peer

func main() {
	flag.BoolVar(&pr.Cfg.StartNode, "n", false, "spawn a new node")
	flag.StringVar(&pr.Cfg.MulticastAddress, "j", "232.232.200.82", "join the cluster at specified multicast group")
	flag.UintVar(&pr.Cfg.MulticastPort, "jp", 8745, "join the cluster at specified multicast group")
	flag.UintVar(&pr.Cfg.ListeningPort, "p", 31582, "listen on the specified port")

	flag.StringVar(&pr.Cfg.LogType, "l", "console", "specify logging type [console (default), file name]")
	flag.StringVar(&pr.Cfg.LogLevel, "v", "info", "specify logging verbosity [off, trace, info (default), warn, err]")

	flag.StringVar(&pr.Cfg.Val, "set", "", "set the value shared across the cluster")
	flag.BoolVar(&pr.Cfg.GetVal, "get", false, "get the value shared across the cluster")

	flag.Parse()
	pr.Run()
}
