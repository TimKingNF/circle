package base

import (
	cmn "circle/common"
)

type Item struct {
	Tag       string `json:"tag"`
	Html      string `json:"html"`
	ParentUrl string `json:"parent_url"`
}

type SavePageItem struct {
	Item
	Topic       cmn.PageTopic `json:"topic"`
	Keywords    []cmn.Keyword `json:"keywords"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
}

type ParseSavePageItem struct {
	Item
	Topic          cmn.PageTopic `json:"topic"`
	Keywords       []cmn.Keyword `json:"-"`
	KeywordsString string        `json:"keywords"`
	Title          string        `json:"title"`
	Description    string        `json:"description"`
}
