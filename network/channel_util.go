package network

import "net"

var (
	conn2Session map[net.Conn]*Session = make(map[net.Conn]*Session)
)

func RegisterSession(conn net.Conn, s *Session) {
	conn2Session[conn] = s
}

func GetSession(conn net.Conn) *Session {
	return conn2Session[conn]
}

func UnregisterSession(conn net.Conn) {
	delete(conn2Session, conn)
}
