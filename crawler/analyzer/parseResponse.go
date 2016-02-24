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

/*用于生成解析http提取相应标签属性的方法的闭包函数*/
func GenTagAttrParseResponse(
	tag, // 标签
	attr string, // 属性
) ParseResponse {

	return func(resp base.Response) ([]base.Data, []error) {
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
		doc.Find(tag).Each(func(index int, sel *goquery.Selection) {
			tagAttr, exists := sel.Attr(attr)
			//	提取超链接 并封装为 base.Request
			if tag == "a" && attr == "href" {
				if !exists || len(tagAttr) == 0 || tagAttr == "#" || tagAttr == "/" {
					return
				}
				tagAttr = strings.TrimSpace(tagAttr)
				lowerTagAttr := strings.ToLower(tagAttr)
				// 暂不支持对javascript、邮箱的解析
				if len(tagAttr) > 0 && !strings.HasPrefix(lowerTagAttr, "javascript") {
					aUrl, err := url.Parse(tagAttr)
					if err != nil {
						errs = append(errs, err)
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
			} else {
				text := strings.TrimSpace(sel.Text())
				if len(text) > 0 && len(tagAttr) > 0 {
					imap := make(map[string]interface{})
					imap[fmt.Sprintf("%s.text", tag)] = text
					imap[fmt.Sprintf("%s.index", tag)] = index
					imap[fmt.Sprintf("%s.%s", tag, attr)] = tagAttr
					imap["tag"] = tag
					imap["parent_url"] = reqUrl.String()
					if resp.Update() {
						imap["update"] = true
					}
					item := base.Item(imap)
					dataList = append(dataList, &item)
				}
			}
		})
		return dataList, errs
	}
}

/*用于生存解析http提取相应标签html内容的方法的闭包函数*/
func GenTagHtmlParseResponse(
	tag string, // 标签
) ParseResponse {

	return func(resp base.Response) ([]base.Data, []error) {
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
		doc.Find(tag).Each(func(index int, sel *goquery.Selection) {
			text, _ := sel.Html()
			text = strings.TrimSpace(text)
			if len(text) > 0 {
				imap := make(map[string]interface{})
				imap["html"] = text
				imap["index"] = index
				imap["tag"] = tag
				imap["parent_url"] = reqUrl.String()
				if resp.Update() {
					imap["update"] = true
				}
				item := base.Item(imap)
				dataList = append(dataList, &item)
			}
		})
		return dataList, errs
	}
}
