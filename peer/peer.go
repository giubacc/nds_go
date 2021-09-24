package peer

type Cfg struct {
	StartNode        bool
	MulticastAddress string
	MulticastPort    uint
	ListeningPort    uint
	Val              string
	GetVal           bool

	LogType  string
	LogLevel string
}
