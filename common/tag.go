package common

import (
	"fmt"
)

type Tag interface {
	Html() string
	String() string
}

type myTag struct {
	Tag string `json:"tag"`
}

func GenTag(tag string) Tag {
	return &myTag{Tag: tag}
}

func (tag *myTag) Html() string {
	return fmt.Sprintf("<%s>", tag.Tag)
}

func (tag *myTag) String() string {
	return tag.Tag
}
