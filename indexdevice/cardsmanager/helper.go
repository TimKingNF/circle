package cardsmanager

import (
	cmn "circle/common"
	base "circle/indexdevice/base"
	"strings"
)

func (cm *myCardsManager) genTrueIndexsFromQueryResultBinary(binaryString string) (indexs []int64) {
	var currentIndex = 0
	for {
		if strings.Contains(binaryString, "1") {
			index := strings.Index(binaryString, "1")
			indexs = append(indexs, int64(currentIndex+index))
			binaryString = binaryString[index+1:]
			currentIndex += index + 1
			continue
		} else {
			break
		}
	}
	return
}

func (cm *myCardsManager) genPageFromPageListIndex(indexs []int64) (pages []cmn.Page) {
	for _, v := range indexs {
		var index int64 = 0
		for e := cm.pageList.Front(); e != nil; e = e.Next() {
			if v == index {
				pages = append(pages, e.Value.(base.MyPage).Page)
				break
			}
			index++
		}
	}
	return
}
