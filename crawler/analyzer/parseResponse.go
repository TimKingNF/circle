/* crawler scheduler analyzer parseResponse */
package analyzer

import (
	base "circle/crawler/base"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ParseResponse func(resp base.Response) ([]base.Data, []error)

//	获取超链接
func ParseForATag(resp base.Response) ([]base.Data, []error) {
	httpResp := resp.HttpResp()
	if httpResp.StatusCode != 200 {
		err := errors.New(
			fmt.Sprintf("Unsupported status code %d. (httpResponse=%v)", httpResp))
		return nil, []error{err}
	}
	var reqUrl *url.URL = httpResp.Request.URL
	var httpRespBody io.ReadCloser = httpResp.Body
	defer func() {
		if httpRespBody != nil {
			httpRespBody.Close()
		}
	}()
	dataList := make([]base.Data, 0)
	errs := make([]error, 0)
	doc, err := goquery.NewDocumentFromReader(httpRespBody)
	if err != nil {
		errs = append(errs, err)
		return dataList, errs
	}

	//	查找标签并提取属性
	doc.Find("a").Each(func(index int, sel *goquery.Selection) {
		href, exists := sel.Attr("href")
		//	提取超链接 并封装为 base.Request
		if !exists || len(href) == 0 || href == "#" || href == "/" {
			return
		}
		href = strings.TrimSpace(href)
		lowerHref := strings.ToLower(href)
		// 暂不支持对javascript、邮箱的解析
		if len(href) > 0 && !strings.HasPrefix(lowerHref, "javascript") &&
			!strings.HasPrefix(lowerHref, "mailto") {
			aUrl, err := url.Parse(href)
			if err != nil {
				errs = append(errs, err)
				return
			}
			if aUrl.Scheme != "http" && aUrl.Scheme != "https" {
				return
			}
			// 判断是否绝对地址
			if !aUrl.IsAbs() {
				aUrl = reqUrl.ResolveReference(aUrl)
			}
			if !resp.Update() {
				httpReq, err := http.NewRequest("GET", aUrl.String(), nil)
				if err != nil {
					errs = append(errs, err)
				} else {
					req := base.NewRequest(httpReq, resp.PrimaryDomain(), resp.Depth())
					dataList = append(dataList, req)
				}
			}
		}
	})
	return dataList, errs
}

//	获取页面的内容
func ParseForHtmlTag(resp base.Response) ([]base.Data, []error) {
	httpResp := resp.HttpResp()
	if httpResp.StatusCode != 200 {
		err := errors.New(
			fmt.Sprintf("Unsupported status code %d. (httpResponse=%v)", httpResp))
		return nil, []error{err}
	}
	var reqUrl *url.URL = httpResp.Request.URL
	var httpRespBody io.ReadCloser = httpResp.Body
	defer func() {
		if httpRespBody != nil {
			httpRespBody.Close()
		}
	}()
	dataList := make([]base.Data, 0)
	errs := make([]error, 0)
	doc, err := goquery.NewDocumentFromReader(httpRespBody)
	if err != nil {
		errs = append(errs, err)
		return dataList, errs
	}
	//	查找html
	text, _ := doc.Find("html").Html()
	text = strings.TrimSpace(text)
	if len(text) > 0 {
		imap := make(map[string]interface{})
		imap["html"] = text
		imap["tag"] = "html"
		imap["parent_url"] = reqUrl.String()
		if resp.Update() {
			imap["update"] = true
		}
		item := base.Item(imap)
		dataList = append(dataList, &item)
	}
	return dataList, errs
}
