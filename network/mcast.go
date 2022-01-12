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
	"context"
	"fmt"
	"nds/util"
	"net"
	"strings"
	"syscall"

	"golang.org/x/net/ipv4"
)

type MCastHelper struct {
	//config
	Cfg *util.Config

	//logger
	logger util.Logger

	//host network interfaces
	hintfs map[string]bool

	//chosen inet for multicasting
	inet net.Interface

	//incoming multicast connection
	iPktConn  net.PacketConn
	iNPktConn *ipv4.PacketConn
}

func (m *MCastHelper) init() error {
	//logger init
	err := m.logger.Init("mcast.", m.Cfg)
	if err != nil {
		return err
	}

	m.hintfs = make(map[string]bool)

	//we enum net interfaces because we want to recognize foreign packets
	nis, err := net.Interfaces()
	for _, ni := range nis {
		addr, _ := ni.Addrs()
		if len(addr) > 0 {
			addr0 := addr[0].String()
			m.logger.Trace("registering host-intf:%s-%s", ni.Name, addr0)
			m.hintfs[strings.Split(addr0, "/")[0]] = true
			//we choose the first eligible interface different from loopback one
			//this interface will be join with multicast group
			if ni.Name != "lo" {
				m.inet = ni
			}
		}
	}

	return err
}

func (m *MCastHelper) stop() error {
	if m.iPktConn != nil {
		m.iPktConn.Close()
	}
	return nil
}

func (m *MCastHelper) establish_multicast() error {
	m.logger.Trace("establishing multicast: group:%s - port:%d", m.Cfg.MulticastAddress, m.Cfg.MulticastPort)

	config := &net.ListenConfig{Control: mcastIncoRawConnCfg}

	if iPktConn, err := config.ListenPacket(context.Background(), "udp4", fmt.Sprintf("0.0.0.0:%d", m.Cfg.MulticastPort)); err != nil {
		m.logger.Err("ListenPacket:%s", err.Error())
		return err
	} else {
		m.iPktConn = iPktConn
	}

	m.iNPktConn = ipv4.NewPacketConn(m.iPktConn)
	mgroup := net.ParseIP(m.Cfg.MulticastAddress)
	if err := m.iNPktConn.JoinGroup(&m.inet, &net.UDPAddr{IP: mgroup}); err != nil {
		m.logger.Err("JoinGroup:%s", err.Error())
		return err
	}

	return nil
}

func mcastIncoRawConnCfg(network, address string, conn syscall.RawConn) error {
	return conn.Control(func(descriptor uintptr) {
		if err := syscall.SetsockoptInt(int(descriptor), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
			util.DefLog().Err("SetsockoptInt:SO_REUSEADDR", err.Error())
		}
	})
}

func (m *MCastHelper) Run() error {
	if err := m.init(); err != nil {
		return err
	}

	if err := m.establish_multicast(); err != nil {
		return err
	}

	return m.stop()
}
