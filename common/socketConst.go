package common

import (
	"bytes"
	"strings"
)

const (
	SERVER_NETWORK = "tcp"
	DELIMITER      = '\t'

	MESSAGE_HEAD = "CIRCLE"

	DEVICE_CRAWLER     = "crawler"
	DEVICE_DIVIDER     = "divider"
	DEVICE_INDEXDEVICE = "indexDevice"
	DEVICE_DATABASE    = "database"

	MESSAGE_POINT = "|@#$%^|"
	MESSAGE_TO    = "-"
	MESSAGE_IN    = "_"

	CRAWLER_QUERYINDEX        = "queryIndex"
	CRAWLER_QUERYINDEX_RESULT = "queryIndexResult"
	CRAWLER_SAVEPAGE          = "savePage"

	INDEXDEVICE_SAVEPAGE   = "savePage"
	INDEXDEVICE_UPDATEPAGE = "updatePage"
)

//	socketDataTemplate
//	[SYSTEM NAME] | [SENDER-ACCEPTER] | [FUNCTION] | [DATA] \t
//	eg: CIRCLE|crawler-divider|[FUNCTION]|[DATA]\t ...

type ControlMessage struct {
	//	control function
	Function string
	//	data
	Doc string
	//	sender
	Sender string
	//	accepter
	Accepter string
	//	Sender id in divider
	SenderId string
	//	Accepter id in divider
	AccepterId string
}

func (cm *ControlMessage) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(MESSAGE_HEAD)
	buffer.WriteString(MESSAGE_POINT)
	buffer.WriteString(cm.Sender)

	if len(cm.SenderId) > 0 {
		buffer.WriteString(MESSAGE_IN)
		buffer.WriteString(cm.SenderId)
	}
	buffer.WriteString(MESSAGE_TO)
	buffer.WriteString(cm.Accepter)

	if len(cm.AccepterId) > 0 {
		buffer.WriteString(MESSAGE_IN)
		buffer.WriteString(cm.AccepterId)
	}

	buffer.WriteString(MESSAGE_POINT)
	buffer.WriteString(cm.Function)
	buffer.WriteString(MESSAGE_POINT)
	buffer.WriteString(cm.Doc)
	return buffer.String()
}

func ParseControlMessage(data string) *ControlMessage {
	var parseData = strings.Split(data, MESSAGE_POINT)
	if len(parseData) != 4 {
		return nil
	}
	if parseData[0] != MESSAGE_HEAD {
		return nil
	}
	var er = strings.Split(parseData[1], MESSAGE_TO)
	if len(er) != 2 {
		return nil
	}
	var ers0 = strings.Split(er[0], MESSAGE_IN)
	var ers1 = strings.Split(er[1], MESSAGE_IN)
	var senderId, accepterId string
	if len(ers0) == 2 {
		senderId = ers0[1]
	}
	if len(ers1) == 2 {
		accepterId = ers1[1]
	}
	return &ControlMessage{
		Function:   parseData[2],
		Doc:        parseData[3],
		Sender:     ers0[0],
		SenderId:   senderId,
		Accepter:   ers1[0],
		AccepterId: accepterId,
	}
}
