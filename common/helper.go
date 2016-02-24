/* common helper method */
package common

import (
	"container/list"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math"
	"net/url"
)

func Str2Md5(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func IsUrl(data string) bool {
	u, _ := url.Parse(data)
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if len(u.Host) <= 0 {
		return false
	}
	return true
}

func GetPrimaryDomain(data string) (string, error) {
	u, _ := url.Parse(data)
	if len(u.Host) > 0 {
		return u.Scheme + "://" + u.Host, nil
	}
	return "", errors.New("Unrecognized host!")
}

func Binary2Decimal(doc string) uint64 {
	stack := list.New()
	var sum uint64
	var stnum, conum float64 = 0, 2

	for _, c := range doc {
		// 入栈 type rune
		stack.PushBack(c)
	}

	// 出栈
	for e := stack.Back(); e != nil; e = e.Prev() {
		// rune是int32的别名
		v := e.Value.(int32)
		sum += uint64(v-48) * uint64(math.Pow(conum, stnum))
		stnum++
	}
	return sum
}
