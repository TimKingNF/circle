package socket

import (
	"bytes"
	cmn "circle/common"
	args "circle/indexdevice/args"
	logging "circle/logging"
	"fmt"
	"net"
	"time"
)

var logger logging.Logger

func genLogger() logging.Logger {
	if logger == nil {
		var loggerArgs cmn.LoggerArgs = args.LoggerArgs

		logger = cmn.NewLogger(cmn.NewLoggerArgs(
			loggerArgs.ConsoleLog(),
			loggerArgs.OutputfileLog(),
			loggerArgs.OutputfilePath(),
			loggerArgs.OutputfilePrefix()+"_socket",
		))
	}
	return logger
}

type Socket interface {
	// socket run
	Run()

	//	get socket remote addr
	RemoteAddr() string

	//	socket between index device and divider
	//	send string to divider after dial divider
	Send(data string)

	//	close socket
	Close()

	//	socket parse data to control model chan
	ParseDataChan() chan *cmn.ControlMessage
}

type mySocket struct {
	conn          *net.Conn
	sendChan      chan string
	acceptChan    chan string
	parseDataChan chan *cmn.ControlMessage
	remoteAddr    string
}

func NewSocket(addr string) Socket {
	conn, err := net.DialTimeout(
		cmn.SERVER_NETWORK,
		addr,
		10*time.Second)
	if err != nil {
		return nil
	}
	return &mySocket{
		conn:       &conn,
		remoteAddr: addr,
	}
}

func (socket *mySocket) dial() {
	go func() {
		var conn net.Conn = *socket.conn
		for {
			data, ok := <-socket.sendChan
			if !ok {
				break
			}
			//	reconnect
			if !socket.ping() {
				for {
					if newConn, err := net.DialTimeout(
						cmn.SERVER_NETWORK,
						socket.remoteAddr,
						10*time.Second); err == nil {
						socket.conn = &newConn
						conn = newConn
						break
					}
					time.Sleep(time.Second)
				}
			}
			//	send data to divider
			go func() {
				var buffer bytes.Buffer
				buffer.WriteString(data)
				buffer.WriteByte(cmn.DELIMITER)
				conn.Write(buffer.Bytes())
			}()
		}
	}()
}

func (socket *mySocket) startAccept() {
	go func() {
		for {
			var conn net.Conn = *socket.conn
			var buffer bytes.Buffer
			var readBytes = make([]byte, 1)
			for {
				_, err := conn.Read(readBytes)
				if err != nil {
					break
				}
				readByte := readBytes[0]
				if readByte == cmn.DELIMITER {
					break
				}
				buffer.WriteByte(readByte)
			}
			if len(buffer.String()) == 0 {
				continue
			}
			socket.acceptChan <- buffer.String()
		}
	}()
}

func (socket *mySocket) parseData() {
	go func() {
		for {
			var conn net.Conn = *socket.conn
			data, ok := <-socket.acceptChan
			if !ok {
				break
			}

			genLogger().Infoln(fmt.Sprintf("Received request[%s]: %s (index device ctrlM)\n",
				conn.RemoteAddr().String(),
				data))

			//	parseData
			var parseData = cmn.ParseControlMessage(data)
			if parseData == nil {
				continue
			}
			if parseData.Accepter != args.DEVICE_NAME {
				continue
			}
			socket.parseDataChan <- parseData
		}
	}()
}

func (socket *mySocket) ping() bool {
	var conn net.Conn = *socket.conn
	var buf bytes.Buffer
	buf.WriteString("p")
	buf.WriteByte(cmn.DELIMITER)
	_, err := conn.Write(buf.Bytes())
	if err != nil {
		return false
	}
	return true
}

func (socket *mySocket) Run() {
	socket.sendChan = make(chan string)
	socket.acceptChan = make(chan string)
	socket.parseDataChan = make(chan *cmn.ControlMessage)

	//	start socket executing
	socket.dial()
	socket.startAccept()
	socket.parseData()
}

func (socket *mySocket) RemoteAddr() string {
	return socket.remoteAddr
}

func (socket *mySocket) Send(data string) {
	if len(data) == 0 {
		return
	}
	genLogger().Infoln(fmt.Sprintf("Send Data[%s]: %s (indexDevice socket)\n",
		(*socket.conn).LocalAddr().String(),
		data))
	socket.sendChan <- data
}

func (socket *mySocket) Close() {
	(*socket.conn).Close()
}

func (socket *mySocket) ParseDataChan() chan *cmn.ControlMessage {
	return socket.parseDataChan
}
