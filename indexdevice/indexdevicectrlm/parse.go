package indexdevicectrlm

import (
	"bytes"
	cmn "circle/common"
	base "circle/indexdevice/base"
	"encoding/json"
	"fmt"
)

type myQueryIndex map[string]string

func (indexDeviceCtrlM *myIndexDeviceControlModel) queryIndexAnalyze(doc *cmn.ControlMessage) {
	go func() {
		var dat myQueryIndex
		if err := json.Unmarshal([]byte(doc.Doc), &dat); err == nil {
			var buffer bytes.Buffer
			//	query the url
			for k, _ := range dat {
				buffer.WriteString("\"")
				buffer.WriteString(k)
				buffer.WriteString("\":")
				if _, ok := indexDeviceCtrlM.mapPage[k]; ok {
					buffer.WriteString("0")
				} else {
					buffer.WriteString("1")
				}
				buffer.WriteString(",")
			}
			ret := fmt.Sprintf("{%s}", buffer.String()[:len(buffer.String())-1])

			//	send query result to divider
			var controlMessage = &cmn.ControlMessage{
				Function:   cmn.CRAWLER_QUERYINDEX_RESULT,
				Doc:        ret,
				Sender:     doc.Accepter,
				Accepter:   doc.Sender,
				AccepterId: doc.SenderId,
			}
			indexDeviceCtrlM.Send(controlMessage.String())
		}
	}()
}

func (indexDeviceCtrlM *myIndexDeviceControlModel) parsePageAnalyze(doc *cmn.ControlMessage) {
	go func() {
		var dat = parseSavePageJson(doc.Doc)
		//	cardsManager add new page
		if dat.Tag != "html" {
			return
		}
		indexDeviceCtrlM.cardsManager.ParsePage(*dat)

		//	cache data and go send to database wait
		indexDeviceCtrlM.pageCache.CachePage(*dat)
	}()
}

func (indexDeviceCtrlM *myIndexDeviceControlModel) updatePageAnalyze(url string) {
	go func() {
		//	update page to divider to crawler
		var controlMessage = &cmn.ControlMessage{
			Function: cmn.INDEXDEVICE_UPDATEPAGE,
			Doc:      url,
			Sender:   cmn.DEVICE_INDEXDEVICE,
			Accepter: cmn.DEVICE_CRAWLER,
		}
		indexDeviceCtrlM.Send(controlMessage.String())
	}()
}

func parseSavePageJson(doc string) *base.SavePageItem {
	var parseSavePageItem base.ParseSavePageItem = base.ParseSavePageItem{}

	err := json.Unmarshal([]byte(doc), &parseSavePageItem)
	if err == nil {
		var parseKeywords []cmn.ParseKeyword
		json.Unmarshal([]byte(parseSavePageItem.KeywordsString), &parseKeywords)
		for k, _ := range parseKeywords {
			var tags []cmn.Tag
			var indexs = parseKeywords[k].Indexs.([]interface{})
			for k1, _ := range indexs {
				tags = append(tags, cmn.GenTag(indexs[k1].(map[string]interface{})["tag"].(string)))
			}
			parseSavePageItem.Keywords = append(parseSavePageItem.Keywords, cmn.Keyword{
				Keyword: parseKeywords[k].Keyword,
				Times:   parseKeywords[k].Times,
				Indexs:  tags,
			})
		}
		item := base.Item{
			Tag:       parseSavePageItem.Tag,
			Html:      parseSavePageItem.Html,
			ParentUrl: parseSavePageItem.ParentUrl,
			Index:     parseSavePageItem.Index,
		}
		return &base.SavePageItem{
			Item:        item,
			Topic:       parseSavePageItem.Topic,
			Keywords:    parseSavePageItem.Keywords,
			Title:       parseSavePageItem.Title,
			Description: parseSavePageItem.Description,
		}
	}
	return nil
}
