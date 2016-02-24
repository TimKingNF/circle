package common

import (
	"time"
)

type Page interface {
	Url() string

	Topic() PageTopic
	SetTopic(topic PageTopic)

	//	page snapshot
	Snap(doc string)
	//	page snapshot string
	SnapShot() string

	//	keywords
	Keywords() []Keyword
	//	set keywords
	SetKeywords(docs []Keyword)

	//	title
	Title() string
	//	set title
	SetTitle(title string)

	//	description
	Description() string
	//	set description
	SetDescription(description string)

	//	update bool
	IsUpdate() bool
}

type myPage struct {
	Mytopic       PageTopic `json:"topic"`
	Myurl         string    `json:"url"`
	Mysnap        string    `json:"snap"`
	Mytitle       string    `json:"title"`
	Mydescription string    `json:"description"`
	Mykeywords    []Keyword `json:"keywords"`
	MylastUpdate  time.Time `json:"lastUpdate"`
}

func GenPage(url string) Page {
	return &myPage{
		Myurl:        url,
		MylastUpdate: time.Now()}
}

func (page *myPage) Url() string {
	return page.Myurl
}

func (page *myPage) Topic() PageTopic {
	return page.Mytopic
}

func (page *myPage) SetTopic(topic PageTopic) {
	page.Mytopic = topic
	page.MylastUpdate = time.Now()
}

func (page *myPage) Snap(doc string) {
	page.Mysnap = doc
	page.MylastUpdate = time.Now()
}

func (page *myPage) SnapShot() string {
	return page.Mysnap
}

func (page *myPage) Keywords() []Keyword {
	return page.Mykeywords
}

func (page *myPage) SetKeywords(docs []Keyword) {
	page.Mykeywords = docs
	page.MylastUpdate = time.Now()
}

func (page *myPage) Title() string {
	return page.Mytitle
}

func (page *myPage) SetTitle(title string) {
	page.Mytitle = title
}

func (page *myPage) Description() string {
	return page.Mydescription
}

func (page *myPage) SetDescription(description string) {
	page.Mydescription = description
}

func (page *myPage) IsUpdate() bool {
	if time.Now().Sub(page.MylastUpdate) > 24*7*time.Hour {
		return true
	}
	return false

	//	test code
	/*	if time.Now().Sub(page.MylastUpdate) > 5*time.Second {
			return true
		}
		return false*/
}
