package base

import (
	"net"
)

type DeviceConn interface {
	Conn() net.Conn
	Id() uint32
}

type myConn struct {
	conn net.Conn
	id   uint32
}

func NewDeviceConn(conn net.Conn, id uint32) DeviceConn {
	return &myConn{
		conn: conn,
		id:   id,
	}
}

func (conn *myConn) Conn() net.Conn {
	return conn.conn
}

func (conn *myConn) Id() uint32 {
	return conn.id
}
