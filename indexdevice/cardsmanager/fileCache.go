package cardsmanager

import (
	"bytes"
	"os"
	"strings"
	"time"
)

const (
	FILE_NAME = "cache.log"
)

type fileCache interface {
	Write(doc string)

	Load() string

	HasCache() bool
}

type myFileCache struct {
	cacheUrl string
}

func GenFileCache(cacheUrl string) fileCache {
	if pathIsExist(cacheUrl) != true {
		if err := os.MkdirAll(cacheUrl, 0777); err != nil {
			if os.IsPermission(err) {
				return nil
			}
			return nil
		}
	}
	return &myFileCache{
		cacheUrl: cacheUrl,
	}
}

func (cache *myFileCache) Write(data string) {
	var buffer, doc bytes.Buffer
	buffer.WriteString(cache.cacheUrl)
	buffer.WriteString("/")
	buffer.WriteString(FILE_NAME)
	f, err := os.OpenFile(buffer.String(), os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		return
	}
	doc.WriteString("[")
	doc.Write([]byte(time.Now().String()))
	doc.WriteString("]：")
	doc.WriteString(data)
	doc.WriteString("\r\n\r\n")
	f.Write(doc.Bytes())
}

func (cache *myFileCache) Load() string {
	var buffer, doc bytes.Buffer
	buffer.WriteString(cache.cacheUrl)
	buffer.WriteString("/")
	buffer.WriteString(FILE_NAME)
	f, err := os.Open(buffer.String())
	defer f.Close()
	if err != nil {
		return ""
	}
	var a = make([]byte, 5)
	for n, err := f.Read(a); err == nil; n, err = f.Read(a) {
		doc.Write(a[:n])
	}

	cacheStrings := strings.Split(doc.String(), "]：")
	if len(cacheStrings) != 2 {
		return ""
	}
	return cacheStrings[1]
}

func (cache *myFileCache) HasCache() bool {
	var buffer bytes.Buffer
	buffer.WriteString(cache.cacheUrl)
	buffer.WriteString("/")
	buffer.WriteString(FILE_NAME)
	f, err := os.Open(buffer.String())
	defer f.Close()
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func pathIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}
