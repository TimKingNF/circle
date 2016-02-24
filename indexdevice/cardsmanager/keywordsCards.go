package cardsmanager

import (
	"bytes"
	cmn "circle/common"
	base "circle/indexdevice/base"
	"container/list"
	"fmt"
)

type keywordsCards interface {
	Gen() uint64

	AddPage(
		pageId uint32,
		indexs []cmn.Tag,
		times uint)

	Container() *list.List
	CopyContainer(container *list.List)
	SetContainer(container *list.List)

	RemovePageKeywords(pageId uint32)

	Keyword() string

	Json() string
}

type myKeywordsCards struct {
	container *list.List
	keyword   string
}

func NewKeywordsCards(keyword string) keywordsCards {
	return &myKeywordsCards{
		container: list.New(),
		keyword:   keyword,
	}
}

func (cards *myKeywordsCards) Keyword() string {
	return cards.keyword
}

func (cards *myKeywordsCards) Gen() uint64 {
	var buffer bytes.Buffer
	for e := cards.container.Front(); e != nil; e = e.Next() {
		v := e.Value.(PageCards)
		if v.IsContain() {
			buffer.WriteString("1")
		} else {
			buffer.WriteString("0")
		}
	}
	return cmn.Binary2Decimal(buffer.String())
}

func (cards *myKeywordsCards) AddPage(pageId uint32,
	indexs []cmn.Tag,
	times uint) {
	//	if container length > uint64 length (2 ^ 63)

	pagecards := NewPageCards(
		cards.keyword,
		pageId, indexs, times)
	cards.container.PushBack(pagecards)
}

func (cards *myKeywordsCards) Container() *list.List {
	return cards.container
}

func (cards *myKeywordsCards) CopyContainer(currentList *list.List) {
	var container = list.New()
	for e := currentList.Front(); e != nil; e = e.Next() {
		v := e.Value.(base.MyPage)
		keywords := v.Page.Keywords()
		for k, _ := range keywords {
			v1 := keywords[k]
			if cards.keyword == v1.Keyword {
				container.PushBack(NewPageCards(cards.keyword,
					v.Id,
					v1.Indexs,
					v1.Times))
			}
		}
	}
	cards.container = container
}

func (cards *myKeywordsCards) SetContainer(currentList *list.List) {
	cards.container = currentList
}

func (cards *myKeywordsCards) RemovePageKeywords(pageId uint32) {
	var n *list.Element
	for e := cards.container.Front(); e != nil; e = n {
		v := e.Value.(PageCards)
		n = e.Next()
		if v.PageId() == pageId {
			cards.container.Remove(e)
		}
	}
}

func (cards *myKeywordsCards) Json() string {
	var buffer bytes.Buffer
	for e := cards.container.Front(); e != nil; e = e.Next() {
		v := e.Value.(PageCards)
		buffer.WriteString(v.Json())
		buffer.WriteString(",")
	}
	ret := buffer.String()
	containerJson := ret[:len(ret)-1]
	return fmt.Sprintf(`{"container":[%s],"keyword":"%s"}`, containerJson, cards.keyword)
}
