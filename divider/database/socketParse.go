package database

import (
	"bufio"
	"bytes"
	cmn "circle/common"
	base "circle/divider/base"
	"fmt"
	"io"
	"net"
	// "time"
)

func (dbm *myDatabaseManager) handleConn(conn net.Conn) {
	defer conn.Close()
	var reader = bufio.NewReader(conn)
	//	save conn
	dbm.connMap[conn.RemoteAddr().String()] = true
	connId := connIdGenertor.GenUint32()
	deviceConn := base.NewDeviceConn(conn, connId)
	dbm.connList.PushBack(deviceConn)

	for {
		// conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		readBytes, err := reader.ReadBytes(cmn.DELIMITER)
		if err != nil {
			if err == io.EOF {
				genLogger().Infoln("The connection is closed by another side. (database manager)\r\n")
			} else {
				genLogger().Warnln(
					fmt.Sprintf("Read Error: %s (database manager)\r", err))
			}
			break
		}
		var readStr = string(readBytes[:len(readBytes)-1])
		if readStr == "p" {
			continue
		}

		genLogger().Infoln(fmt.Sprintf("Received request[%s]: %s (database manager)\n",
			conn.RemoteAddr().String(),
			readStr))

		// parse data
		dbm.connMap[conn.RemoteAddr().String()] = true

		//	send data to database
		var buffer bytes.Buffer
		buffer.WriteString(readStr)
		buffer.WriteByte(cmn.DELIMITER)
		conn.Write(buffer.Bytes())

		//	send data to channel
		dbm.dataChan <- readStr
	}
}
