/* crawler scheduler cookie */
package scheduler

import (
	"net/http"
	"net/http/cookiejar"
)

type myPublicSuffixList struct{}

func NewCookiejar() http.CookieJar {
	options := &cookiejar.Options{PublicSuffixList: &myPublicSuffixList{}}
	cj, _ := cookiejar.New(options)
	return cj
}

func (psl *myPublicSuffixList) PublicSuffix(domain string) string {
	suffix, _ := getPrimaryDomain(domain)
	return suffix
}

func (psl *myPublicSuffixList) String() string {
	return "Search Engine(circle) - public suffix list (rev 1.0)"
}
