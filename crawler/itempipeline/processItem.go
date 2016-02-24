/* crawler scheduler process item */
package itempipeline

import (
	cmn "circle/common"
	base "circle/crawler/base"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"strings"
	"time"
)

// accept a item and process, return a processed item
type ProcessItem func(item base.Item) (result base.Item, err error)

//	gen keywords or head or title or topic or tag form page html
func GenKeywordsFromPage(item base.Item) (result base.Item, err error) {
	if v, ok := item["tag"].(string); ok {
		if v == "html" {
			var pageHtml io.ReadCloser = ioutil.NopCloser(strings.NewReader(item["html"].(string)))
			defer func() {
				if pageHtml != nil {
					pageHtml.Close()
				}
			}()
			doc, err := goquery.NewDocumentFromReader(pageHtml)
			if err != nil {
				return item, err
			}

			var keywords []cmn.Keyword

			//	gen title
			title, _ := doc.Find("title").Html()
			item["title"] = title
			//	title keywords
			if strings.Contains(title, "_") {
				titleKeywords := strings.Split(title, "_")
				for k, _ := range titleKeywords {
					keywords = append(keywords, cmn.Keyword{
						Keyword: titleKeywords[k],
						Indexs: []cmn.Tag{
							cmn.GenTag("head"),
						},
					})
				}
			}

			//	gen tag : <meta name="keywords" content="...">
			keywordContent, _ := doc.Find("meta[name='keywords']").Attr("content")
			if strings.Contains(keywordContent, ",") {
				keywordContents := strings.Split(keywordContent, ",")
				var tmpKeywords []cmn.Keyword
				for k, _ := range keywordContents {
					tmpKeywords = append(tmpKeywords, cmn.Keyword{
						Keyword: keywordContents[k],
						Indexs: []cmn.Tag{
							cmn.GenTag("head"),
						},
					})
				}
				keywords = appendKeywordsByTag(keywords, tmpKeywords)
			}

			//	gen tag : <h1> - <h2>
			genTagKeywordFromPageEach(doc, "h1", &keywords)
			genTagKeywordFromPageEach(doc, "h2", &keywords)
			//	gen tag : <b> <strong>
			genTagKeywordFromPageEach(doc, "b", &keywords)
			genTagKeywordFromPageEach(doc, "strong", &keywords)

			//	gen keywords repeatd times
			bodyText := doc.Find("body").Text()
			for k, _ := range keywords {
				//	search the repeatd times every keyword
				keywords[k].Times = uint(strings.Count(bodyText, keywords[k].Keyword))
			}

			keywordsJson, _ := json.Marshal(keywords)
			item["keywords"] = string(keywordsJson)
			//	gen description
			description, _ := doc.Find("meta[name='description']").Attr("content")
			item["description"] = description
			//	gen topic
			//	...
		}
	}
	time.Sleep(10 * time.Millisecond)
	return item, nil
}

func genTagKeywordFromPageEach(doc *goquery.Document, tag string, keywords *[]cmn.Keyword) {
	var tmpKeywords []cmn.Keyword
	doc.Find(tag).Each(func(index int, sel *goquery.Selection) {
		v := sel.Text()
		tmpKeywords = append(tmpKeywords, cmn.Keyword{
			Keyword: v,
			Indexs: []cmn.Tag{
				cmn.GenTag(tag),
			},
			Times: 1,
		})
	})
	*keywords = appendKeywordsByTag(*keywords, tmpKeywords)
}

func appendKeywordsByTag(keywords, appendData []cmn.Keyword) []cmn.Keyword {
	for _, v := range appendData {
		var sign = true
		for k, _ := range keywords {
			if keywords[k].Keyword == v.Keyword {
				sign = false
				for indexK1, _ := range v.Indexs {
					var tagSign = true
					for indexK2, _ := range keywords[k].Indexs {
						if keywords[k].Indexs[indexK2].String() == v.Indexs[indexK1].String() {
							tagSign = false
							break
						}
					}
					if tagSign {
						keywords[k].Indexs = append(keywords[k].Indexs, v.Indexs[indexK1])
					}
				}

				break
			}
		}
		if sign {
			keywords = append(keywords, v)
		}
	}
	return keywords
}
