package databasectrlm

import (
	cmn "circle/common"
	base "circle/database/base"
	"encoding/json"
)

func (dbCtrlM *myDatabaseControlModel) savePageAnalyze(doc *cmn.ControlMessage) {
	items := parsePageJson(doc.Doc)
	// use apiLayout and go savePage API
	for k, _ := range items {
		dbCtrlM.apiLayout.Insert("page", items[k])
	}
}

func parsePageJson(doc string) []*cmn.Page {
	var pages []*cmn.Page
	var parsePages []base.ParsePage
	err := json.Unmarshal([]byte(doc), &parsePages)
	if err == nil {
		for k, _ := range parsePages {
			var page cmn.Page = cmn.GenPage(parsePages[k].Url)
			page.SetTopic(parsePages[k].Topic)
			page.Snap(parsePages[k].Snap)
			page.SetTitle(parsePages[k].Title)
			page.SetDescription(parsePages[k].Description)

			var keywords []cmn.Keyword
			for _, v := range parsePages[k].Keywords {
				var indexs []cmn.Tag
				for _, v1 := range v.Indexs.([]interface{}) {
					indexs = append(indexs, cmn.GenTag(v1.(map[string]interface{})["tag"].(string)))
				}
				keywords = append(keywords, cmn.Keyword{
					Keyword: v.Keyword,
					Times:   v.Times,
					Indexs:  indexs,
				})
			}

			page.SetKeywords(keywords)
			pages = append(pages, &page)
		}
	}
	return pages
}
