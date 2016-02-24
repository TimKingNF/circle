package base

import (
	cmn "circle/common"
	"time"
)

type ParsePage struct {
	Topic       cmn.PageTopic      `json:"topic"`
	Url         string             `json:"url"`
	Snap        string             `json:"snap"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Keywords    []cmn.ParseKeyword `json:"keywords"`
	LastUpdate  time.Time          `json:"lastUpdate"`
}
