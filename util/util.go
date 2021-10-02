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

type RetCode int

const (
	RetCode_UNKERR = -1000 /**< unknown error */

	//system errors
	RetCode_SCKERR = -105 /**< socket error */
	RetCode_DBERR  = -104 /**< database error */
	RetCode_IOERR  = -103 /**< I/O operation fail */
	RetCode_MEMERR = -102 /**< memory error */
	RetCode_SYSERR = -101 /**< system error */

	//generic error
	RetCode_UNVRSC = -2 /**< unavailable resource */
	RetCode_GENERR = -1 /**< generic error */

	//success, failure [0,1]
	RetCode_OK    = 0 /**< operation ok */
	RetCode_KO    = 1 /**< operation fail */
	RetCode_EXIT  = 2 /**< exit required */
	RetCode_RETRY = 3 /**< request retry */
	RetCode_ABORT = 4 /**< operation aborted */

	//generics
	RetCode_UNSP     = 100 /**< unsupported */
	RetCode_NODATA   = 101 /**< no data */
	RetCode_NOTFOUND = 102 /**< not found */
	RetCode_TIMEOUT  = 103 /**< timeout */

	//contaniers specific
	RetCode_EMPTY = 200 /**< empty */
	RetCode_QFULL = 201 /**< queue full */
	RetCode_OVRSZ = 202 /**< oversize */
	RetCode_BOVFL = 203 /**< buffer overflow */

	//proc. specific
	RetCode_BADARG  = 300 /**< bad argument */
	RetCode_BADIDX  = 301 /**< bad index */
	RetCode_BADSTTS = 302 /**< bad status */
	RetCode_BADCFG  = 303 /**< bad configuration */

	//network specific
	RetCode_DRPPKT  = 400 /**< packet dropped*/
	RetCode_MALFORM = 401 /**< packet malformed */
	RetCode_SCKCLO  = 402 /**< socket closed */
	RetCode_SCKWBLK = 403 /**< socket would block */
	RetCode_PARTPKT = 404 /**< partial packet */
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

/**
 * event type
 */
type EvtType int

const (
	Undef           = iota
	Interrupt       //generic interrupt
	ConnectRequest  //request for TCP connection (peer -> selector)
	IncomingConnect //new incoming TCP connection (selector -> peer)
	SendPacket      //request to send a packet (peer -> selector)
	PacketAvailable //foreign packet available (selector -> peer)
	Disconnect      //connection disconnection event
)

/**
 * An event
 *
 * This struct is used as message between selector/peer.
 *
 * It can model an interrupt (Interrupt).
 * It can be used by peer thread to request selector thread to connect to TCP (ConnectRequest).
 * It can transport an incoming packet from a multicast/unicast socket (PacketAvailable).
 * It can be used by peer thread to request selector to send a packet (SendPacket).
 * It can be used by selector to inform peer thread of a received foreign packet (PacketAvailable).
 * It can be used to manage an event of disconnection (Disconnect).
 */
type Event struct {
	Evt      EvtType
	OptSrcIp string
}
