package cardsmanager

import (
	"bytes"
	cmn "circle/common"
	base "circle/indexdevice/base"
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

/*
             Page      Page   Page   Page   Page   Page   Page
keywords   PageCards                                            KeywordsCards
keywords
keywords
		  PageKeywords
*/

var pageIdGenertor cmn.IdGenertor = cmn.NewIdGenertor()

type CardsManager interface {
	//	query keyword
	QueryKeyword(keyword string) []cmn.Page

	//	parse page (contain add and update)
	ParsePage(item base.SavePageItem)

	PageList() *list.List

	//	gen MyPage id form pagelist
	GenPageId(url string) (uint32, error)
}

type myCardsManager struct {
	//	map [keyword] keywordsCards
	container map[string]keywordsCards

	//	list Page (page)
	pageList *list.List

	fileCache fileCache

	mapKeyword []string

	mutex sync.Mutex
}

func (cm *myCardsManager) activateFileCache() {
	go func() {
		//	if has cache and load cache
		if cm.fileCache.HasCache() {
			var cacheString = cm.fileCache.Load()
			if len(cacheString) > 0 {
				//	parse map[string]keywordsCards and pageList
				cm.loadCache(cacheString)
			}
		}

		for {
			//	test time
			// time.Sleep(5 * time.Second)

			time.Sleep(30 * time.Minute)
			if len(cm.container) > 0 {
				//	gen json from CardsManager.container and CardsManager.pageList
				doc := cm.genJsonFromCardsManager()
				cm.fileCache.Write(doc)
			}
		}
	}()
}

func GenCardsManager(fileCacheUrl string) CardsManager {
	cardsManager := &myCardsManager{
		container:  make(map[string]keywordsCards),
		pageList:   list.New(),
		fileCache:  GenFileCache(fileCacheUrl),
		mapKeyword: make([]string, 0),
	}
	cardsManager.activateFileCache()
	return cardsManager
}

func (cm *myCardsManager) QueryKeyword(keyword string) []cmn.Page {
	for {
		if len(cm.container) > 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	var retMap = make([]uint64, 0)
	//	精确查找
	if keywordsCards, ok := cm.container[keyword]; ok {
		retMap = append(retMap, keywordsCards.Gen())
	}
	//	模糊查询
	for k, _ := range cm.container {
		if strings.Contains(k, keyword) && k != keyword {
			retMap = append(retMap, cm.container[k].Gen())
		}
	}
	if len(retMap) > 0 {
		//	综合结果
		var ret uint64 = retMap[0]
		for k, v := range retMap {
			if k != 0 {
				ret = ret & v
			}
		}
		retBinary := fmt.Sprintf("%b", ret)
		pageIndexs := cm.genTrueIndexsFromQueryResultBinary(retBinary)
		pages := cm.genPageFromPageListIndex(pageIndexs)
		return pages
	}
	return nil
}

func (cm *myCardsManager) ParsePage(item base.SavePageItem) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	//	pageList container page, go update page
	//	else go add page
	for e := cm.pageList.Front(); e != nil; e = e.Next() {
		v := e.Value.(base.MyPage)
		if v.Page.Url() == item.ParentUrl {
			cm.updatePage(item)
			return
		}
	}
	cm.addPage(item)
}

func (cm *myCardsManager) updatePage(item base.SavePageItem) {
	var n *list.Element
	for keyword, _ := range cm.container {
		myKeywordsCards := cm.container[keyword]
		pageid, err := cm.GenPageId(item.ParentUrl)
		if err == nil {
			myKeywordsCards.RemovePageKeywords(pageid)
		}
	}
	for e := cm.pageList.Front(); e != nil; e = n {
		v := e.Value.(base.MyPage)
		n = e.Next()
		if v.Page.Url() == item.ParentUrl {
			cm.pageList.Remove(e)
		}
	}
	cm.addPage(item)
}

func (cm *myCardsManager) addPage(item base.SavePageItem) {
	//	keyword control
	for k, _ := range item.Keywords {
		v := item.Keywords[k]
		if _, ok := cm.container[v.Keyword]; !ok {
			cm.container[v.Keyword] = NewKeywordsCards(v.Keyword)
			cm.mapKeyword = append(cm.mapKeyword, v.Keyword)
		}
	}
	for k, _ := range item.Keywords {
		if cm.pageList.Len() > 0 && cm.container[item.Keywords[k].Keyword].Container().Len() == 0 {
			cm.container[item.Keywords[k].Keyword].CopyContainer(cm.pageList)
		}
	}

	if len(cm.container) > 0 {
		//	add page for all keyword
		page := cmn.GenPage(item.ParentUrl)
		page.SetKeywords(item.Keywords)
		page.SetTitle(item.Title)
		page.SetDescription(item.Description)
		page.Snap(item.Html)
		page.SetTopic(item.Topic)
		pageId := pageIdGenertor.GenUint32()

		for keyword, _ := range cm.container {
			myKeywordsCards := cm.container[keyword]
			for k, _ := range item.Keywords {
				v := item.Keywords[k]
				if v.Keyword == keyword {
					myKeywordsCards.AddPage(
						pageId,
						v.Indexs,
						v.Times)
				}
			}
		}

		myPage := base.MyPage{
			Page: page,
			Id:   pageId,
		}
		cm.pageList.PushBack(myPage)
	}
}

func (cm *myCardsManager) PageList() *list.List {
	return cm.pageList
}

func (cm *myCardsManager) GenPageId(url string) (uint32, error) {
	for e := cm.pageList.Front(); e != nil; e = e.Next() {
		v := e.Value.(base.MyPage)
		if url == v.Page.Url() {
			return v.Id, nil
		}
	}
	return 0, errors.New("No search url!")
}

func (cm *myCardsManager) genJsonFromCardsManager() string {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	var buffer bytes.Buffer
	buffer.WriteString(`{"container":`)

	var containerBuf bytes.Buffer
	for k, _ := range cm.container {
		containerBuf.WriteString(`"`)
		containerBuf.WriteString(k)
		containerBuf.WriteString(`":`)
		containerBuf.WriteString(cm.container[k].Json())
		containerBuf.WriteString(",")
	}
	currentJson := containerBuf.String()
	containerJson := fmt.Sprintf(`{%s}`, currentJson[:len(currentJson)-1])

	buffer.WriteString(containerJson)
	buffer.WriteString(`,"pagelist":[`)

	var pageListBuf bytes.Buffer
	for e := cm.pageList.Front(); e != nil; e = e.Next() {
		v := e.Value.(base.MyPage)
		pageListBuf.WriteString(`{"page":`)
		pageBytes, _ := json.Marshal(v.Page)
		pageListBuf.Write(pageBytes)
		pageListBuf.WriteString(`,"id":`)
		pageListBuf.WriteString(fmt.Sprintf("%d", v.Id))
		pageListBuf.WriteString(`},`)
	}
	pageListRet := pageListBuf.String()
	pageListJson := pageListRet[:len(pageListRet)-1]

	buffer.WriteString(pageListJson)
	buffer.WriteString(`]}`)
	return buffer.String()
}

func (cm *myCardsManager) loadCache(cacheString string) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	var cache parseCardsManager
	if err := json.Unmarshal([]byte(cacheString), &cache); err == nil {
		//	set cardsManager container
		for keyword, parseKeywordsCards := range cache.Container {
			//	gen keywordsCards container list
			var keywordsCardsContainer = list.New()
			for _, parsePageCards := range parseKeywordsCards.Container {
				var indexs []cmn.Tag
				for k, _ := range parsePageCards.Myindexs {
					indexs = append(indexs, cmn.GenTag(parsePageCards.Myindexs[k].Tag))
				}
				keywordsCardsContainer.PushBack(NewPageCards(
					keyword,
					parsePageCards.Pageid,
					indexs,
					parsePageCards.Times))
			}
			keywordsCards := NewKeywordsCards(keyword)
			keywordsCards.SetContainer(keywordsCardsContainer)
			cm.container[keyword] = keywordsCards
		}

		var maxGenId uint32 = 0

		//	set cardsManager pagelist
		for k, _ := range cache.PageList {
			//	gen page
			var page = cmn.GenPage(cache.PageList[k].Page.Url)
			page.SetTopic(cache.PageList[k].Page.Topic)
			page.Snap(cache.PageList[k].Page.Snap)
			page.SetTitle(cache.PageList[k].Page.Title)
			page.SetDescription(cache.PageList[k].Page.Description)

			var keywords = make([]cmn.Keyword, 0)
			for _, v := range cache.PageList[k].Page.Keywords {
				var indexs []cmn.Tag
				for _, indexsTags := range v.Indexs.([]interface{}) {
					for _, indexsTag := range indexsTags.(map[string]interface{}) {
						indexs = append(indexs, cmn.GenTag(indexsTag.(string)))
					}
				}
				keywords = append(keywords, cmn.Keyword{
					Keyword: v.Keyword,
					Times:   v.Times,
					Indexs:  indexs,
				})
			}
			page.SetKeywords(keywords)

			var mypage base.MyPage = base.MyPage{
				Id:   cache.PageList[k].Id,
				Page: page,
			}
			if maxGenId < mypage.Id {
				maxGenId = mypage.Id
			}
			cm.pageList.PushBack(mypage)
		}

		//	set pageIdGenertor start
		pageIdGenertor.Start(maxGenId + 1)

		//	set cardsManager mapKeyword
		for k, _ := range cm.container {
			cm.mapKeyword = append(cm.mapKeyword, k)
		}
	}
}
