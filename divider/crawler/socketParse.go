package crawler

import (
	"bufio"
	cmn "circle/common"
	base "circle/divider/base"
	"fmt"
	"io"
	"net"
	// "time"
)

func (cm *myCrawlerManager) handleConn(conn net.Conn) {
	defer conn.Close()
	var reader = bufio.NewReader(conn)
	//	save conn
	cm.connMap[conn.RemoteAddr().String()] = true
	connId := connIdGenertor.GenUint32()
	deviceConn := base.NewDeviceConn(conn, connId)
	cm.connList.PushBack(deviceConn)

	for {
		// conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		readBytes, err := reader.ReadBytes(cmn.DELIMITER)
		if err != nil {
			if err == io.EOF {
				genLogger().Infoln("The connection is closed by another side. (crawler manager)\r\n")
			} else {
				genLogger().Warnln(
					fmt.Sprintf("Read Error: %s (crawler manager)\n", err))
			}
			break
		}
		var readStr = string(readBytes[:len(readBytes)-1])
		if readStr == "p" {
			continue
		}

		genLogger().Infoln(fmt.Sprintf("Received request[%s]: %s (crawler manager)\n",
			deviceConn.Conn().RemoteAddr().String(),
			readStr))

		// parse data
		var parseData = cmn.ParseControlMessage(readStr)

		if parseData != nil {
			switch parseData.Accepter {
			case cmn.DEVICE_DIVIDER:

			case cmn.DEVICE_INDEXDEVICE:
				switch parseData.Function {
				case cmn.CRAWLER_QUERYINDEX:
					cm.queryIndexToIndexDevice(deviceConn, parseData)
				case cmn.CRAWLER_SAVEPAGE:
					cm.savePageToIndexDevice(deviceConn, parseData)
				case cmn.INDEXDEVICE_UPDATEPAGE:
					cm.updatePageToIndexDevice(deviceConn, parseData)
				}
			case cmn.DEVICE_DATABASE:

			}
		}
	}
}

func (cm *myCrawlerManager) queryIndexToIndexDevice(resultConn base.DeviceConn, parseData *cmn.ControlMessage) {
	//	send data to analyze index device Manager
	cm.dataChan <- &cmn.ControlMessage{
		Function: parseData.Function,
		Doc:      parseData.Doc,
		Sender:   parseData.Sender,
		Accepter: parseData.Accepter,
		SenderId: fmt.Sprintf("%d", resultConn.Id()),
	}
}

func (cm *myCrawlerManager) savePageToIndexDevice(resultConn base.DeviceConn, parseData *cmn.ControlMessage) {
	//	send data to analyze index device Manager
	cm.dataChan <- &cmn.ControlMessage{
		Function: parseData.Function,
		Doc:      parseData.Doc,
		Sender:   parseData.Sender,
		Accepter: parseData.Accepter,
		SenderId: fmt.Sprintf("%d", resultConn.Id()),
	}
}

func (cm *myCrawlerManager) updatePageToIndexDevice(resultConn base.DeviceConn, parseData *cmn.ControlMessage) {
	//	send data to analyze index device Manager
	cm.dataChan <- &cmn.ControlMessage{
		Function:   parseData.Function,
		Doc:        parseData.Doc,
		Sender:     parseData.Sender,
		Accepter:   parseData.Accepter,
		SenderId:   fmt.Sprintf("%d", resultConn.Id()),
		AccepterId: parseData.AccepterId,
	}
}
