package cardsmanager

import (
	cmn "circle/common"
	"encoding/json"
)

type myKeyword struct {
	Keyword string
	times   uint
}

type PageCards interface {
	Keyword() string

	Times() uint

	PageId() uint32

	IsContain() bool

	Indexs() []cmn.Tag

	Json() string
}

type myPageCards struct {
	Mykeyword string `json:"keyword"`

	//	keyword repeatd times
	Mytimes uint `json:"times"`

	//	page id
	Mypageid uint32 `json:"pageid"`

	//	keyword index localtion, like <head> <title> <body> ...
	Myindexs []cmn.Tag `json:"indexs"`
}

func NewPageCards(
	keyword string,
	pageid uint32,
	indexs []cmn.Tag,
	times uint,
) PageCards {
	return &myPageCards{
		Mykeyword: keyword,
		Mypageid:  pageid,
		Myindexs:  indexs,
		Mytimes:   times,
	}
}

func (card *myPageCards) Keyword() string {
	return card.Mykeyword
}

func (card *myPageCards) Times() uint {
	return card.Mytimes
}

func (card *myPageCards) PageId() uint32 {
	return card.Mypageid
}

func (card *myPageCards) IsContain() bool {
	if card.Mytimes > 0 {
		return true
	}
	return false
}

func (card *myPageCards) Indexs() []cmn.Tag {
	return card.Myindexs
}

func (card *myPageCards) Json() string {
	ret, _ := json.Marshal(card)
	return string(ret)
}
