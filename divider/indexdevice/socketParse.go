package indexdevice

import (
	"bufio"
	cmn "circle/common"
	base "circle/divider/base"
	"fmt"
	"io"
	"net"
)

func (idm *myIndexDeviceManager) handleConn(conn net.Conn) {
	defer conn.Close()
	var reader = bufio.NewReader(conn)
	//	save conn
	idm.connMap[conn.RemoteAddr().String()] = true
	connId := connIdGenertor.GenUint32()
	deviceConn := base.NewDeviceConn(conn, connId)
	idm.connList.PushBack(deviceConn)

	for {
		// conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		readBytes, err := reader.ReadBytes(cmn.DELIMITER)
		if err != nil {
			if err == io.EOF {
				genLogger().Infoln("The connection is closed by another side. (IndexDevice manager)\r\n")
			} else {
				genLogger().Warnln(
					fmt.Sprintf("Read Error: %s (IndexDevice manager)\n", err))
			}
			break
		}
		var readStr = string(readBytes[:len(readBytes)-1])
		if readStr == "p" {
			continue
		}

		genLogger().Infoln(fmt.Sprintf("Received request[%s]: %s (IndexDevice manager)\n",
			conn.RemoteAddr().String(),
			readStr))

		// parse data
		var parseData = cmn.ParseControlMessage(readStr)

		if parseData != nil {
			switch parseData.Accepter {
			case cmn.DEVICE_DIVIDER:

			case cmn.DEVICE_CRAWLER:
				switch parseData.Function {
				case cmn.CRAWLER_QUERYINDEX_RESULT:
					idm.queryIndexResultToCrawler(deviceConn, parseData)
				case cmn.INDEXDEVICE_UPDATEPAGE:
					idm.updatePageToCrawler(deviceConn, parseData)
				}
			case cmn.DEVICE_DATABASE:
				switch parseData.Function {
				case cmn.INDEXDEVICE_SAVEPAGE:
					idm.savePageToDatabase(deviceConn, parseData)
				}
			}
		}
	}
}

func (idm *myIndexDeviceManager) queryIndexResultToCrawler(resultConn base.DeviceConn, parseData *cmn.ControlMessage) {
	//	send data to analyze index device manager
	idm.dataChan <- &cmn.ControlMessage{
		Function:   parseData.Function,
		Doc:        parseData.Doc,
		Sender:     parseData.Sender,
		Accepter:   parseData.Accepter,
		SenderId:   fmt.Sprintf("%d", resultConn.Id()),
		AccepterId: parseData.AccepterId,
	}
}

func (idm *myIndexDeviceManager) savePageToDatabase(resultConn base.DeviceConn, parseData *cmn.ControlMessage) {
	//	send data to analyze index device manager
	idm.dataChan <- &cmn.ControlMessage{
		Function: parseData.Function,
		Doc:      parseData.Doc,
		Sender:   parseData.Sender,
		Accepter: parseData.Accepter,
		SenderId: fmt.Sprintf("%d", resultConn.Id()),
	}
}

func (idm *myIndexDeviceManager) updatePageToCrawler(resultConn base.DeviceConn, parseData *cmn.ControlMessage) {
	//	send data to analyze index device manager
	idm.dataChan <- &cmn.ControlMessage{
		Function: parseData.Function,
		Doc:      parseData.Doc,
		Sender:   parseData.Sender,
		Accepter: parseData.Accepter,
		SenderId: fmt.Sprintf("%d", resultConn.Id()),
	}
}
