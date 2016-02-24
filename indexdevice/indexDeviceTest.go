package main

import (
	base "circle/indexdevice/base"
	idcm "circle/indexdevice/indexdevicectrlm"
	"fmt"
	"sync"
	// "time"
)

var (
	connected_divider_address = "127.0.0.1:8087"

	indexDeviceArgs base.IndexDeviceArgs = base.NewIndexDeviceArgs(
		connected_divider_address)
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	TestSocket()
}

func TestSocket() {
	idCtrlM := idcm.GenIndexDeviceControlModel()
	err := idCtrlM.Init(indexDeviceArgs)
	if err != nil {
		panic(fmt.Sprintf("index device control model initialize failed, Error: %s", err))
	}

	// idCtrlM.QueryKeyword("梦中旅人")

	//	socket test
	// for i := 0; i < 5; i++ {
	// 	idCtrlM.Send(fmt.Sprintf("send the number[%d]", i))
	// 	time.Sleep(time.Second)
	// }

	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	waitgroup.Wait()
}
