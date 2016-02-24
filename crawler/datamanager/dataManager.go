package datamanager

import (
	cmn "circle/common"
)

type DataManager interface {
	SendDataChan() chan *cmn.ControlMessage
	AcceptDataChan() chan *cmn.ControlMessage
}

type myDataManager struct {
	sendDataChan   chan *cmn.ControlMessage
	acceptDataChan chan *cmn.ControlMessage
}

func GenDataManager() DataManager {
	return &myDataManager{
		sendDataChan:   make(chan *cmn.ControlMessage),
		acceptDataChan: make(chan *cmn.ControlMessage),
	}
}

func (dm *myDataManager) SendDataChan() chan *cmn.ControlMessage {
	return dm.sendDataChan
}

func (dm *myDataManager) AcceptDataChan() chan *cmn.ControlMessage {
	return dm.acceptDataChan
}
