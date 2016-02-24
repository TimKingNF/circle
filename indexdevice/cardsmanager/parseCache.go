package cardsmanager

import (
	cmn "circle/common"
	"time"
)

type parseCardsManager struct {
	Container map[string]parseKeywordsCards `json:"container"`
	PageList  []parseMyPage                 `json:"pageList"`
}

type parseKeywordsCards struct {
	Container []parsePageCards `json:"container"`
	Keyword   string           `json:"keyword"`
}

type parsePageCards struct {
	Keyword  string     `json:"keyword"`
	Times    uint       `json:"times"`
	Pageid   uint32     `json:"pageid"`
	Myindexs []parseTag `json:"indexs"`
}

type parseTag struct {
	Tag string `json:"tag"`
}

type parseMyPage struct {
	Page parsePage `json:"page"`
	Id   uint32    `json:"id"`
}

type parsePage struct {
	Topic       cmn.PageTopic      `json:"topic"`
	Url         string             `json:"url"`
	Snap        string             `json:"snap"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Keywords    []cmn.ParseKeyword `json:"keywords"`
	LastUpdate  time.Time          `json:"lastUpdate"`
}
