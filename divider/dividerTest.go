package main

import (
	dividercm "circle/divider/dividerctrlm"
	"fmt"
	"sync"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	dividerCtrlM := dividercm.NewDivider()
	dividerCtrlM.Init()

	var waitgroup sync.WaitGroup
	waitgroup.Add(1)
	waitgroup.Wait()
}
