package common

type Keyword struct {
	Keyword string `json:"keyword"`
	Times   uint   `json:"times"`
	Indexs  []Tag  `json:"indexs"`
}

type ParseKeyword struct {
	Keyword string      `json:"keyword"`
	Times   uint        `json:"times"`
	Indexs  interface{} `json:"indexs"`
}
