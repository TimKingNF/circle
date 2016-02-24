/* system condition	*/
package common

var (
	circleCondition = []string{
		"┌────────────┐┌─────────────┐┌─────────────────────────┐",
		"│   Crawler  ││   Database  ││                         │",
		"└──────┬─────┘└──────┬──────┘│     IndexDeviceCtrlM    │",
		"┌──────┴─────┐┌──────┴──────┐│                         │",
		"│CrawlerCtrlM││DatabaseCtrlM││                         │",
		"└──────┬─────┘└──────┬──────┘└─┬─────────────┬─────────┘",
		"┌──────┴─────────────┴─────────┴──┐┌─────────┴─────────┐",
		"│             Divider             ││    WordAnalyzer   │",
		"└────────────────┬────────────────┘└─────────┬─────────┘",
		"┌────────────────┴────────────────┐┌─────────┴─────────┐",
		"│            HttpServer           ││     HttpClient    │",
		"└─────────────────────────────────┘└───────────────────┘"}

	dividerStruct = []string{
		"┌───────────────────────────────────┐┌──────────────┐",
		"│         DividerControlModel       ││ Load Balancer│",
		"└───────────────────────────────────┘└──────────────┘",
		"┌──────────────┐┌───────────────┐┌──────────────────┐",
		"│CrawlerManager││DatabaseManager││IndexDeviceManager│",
		"└──────────────┘└───────────────┘└──────────────────┘",
	}

	crawlerStruct = []string{
		"┌───────────────────────────────────────┐ ┌───────────────────────┐",
		"│                                       │ │         Socket        │",
		"│           CrawlerControlModel         │ └───────────────────────┘",
		"│                                       └─────────────────────────┐",
		"└─────────────────────────────────────────────────────────────────┘",
		"┌─────────────────────────────────────────────────────────────────┐",
		"│                            Scheduler                            │",
		"└─────────────────────────────────────────────────────────────────┘",
		"┌────────────┐┌──────────────┐┌────────────┐┌────────────┐┌───────┐",
		"│RequestCache││DownloaderPool││AnalyzerPool││ItemPipeline││       │",
		"└────────────┘└──────────────┘└────────────┘└────────────┘│       │",
		"┌────────────┐┌──────────────┐┌────────────┐┌────────────┐│       │",
		"│  UrlCache  ││  Downloader  ││  Analyzer  ││ DataManager││Monitor│",
		"└────────────┘└──────────────┘└────────────┘└────────────┘│       │",
		"┌────────────────────────────────────────────────────────┐│       │",
		"│                      middleware                        ││       │",
		"└────────────────────────────────────────────────────────┘└───────┘"}

	crawlerLifeCycle = []string{
		"┌───────────────┐┌───────────────┐┌───────────────┐",
		"│    Running    ││   Wait Time   ││    Sleeping   │",
		"└───────────────┘└───────────────┘└───────────────┘",
		"┌────────────────────────────────┐┌────────────────",
		"│           LifeCycle            ││    ......      ",
		"└────────────────────────────────┘└────────────────"}

	indexDeviceStruct = []string{
		"┌───────────────────────┐ ┌────────────────────┐",
		"│                       │ │        Socket      │",
		"│IndexDeviceControlModel│ └────────────────────┘",
		"│                       └──────────────────────┐",
		"└──────────────────────────────────────────────┘",
		"┌───────┐┌───────────┐┌────────────────────────┐",
		"│       ││           ││       CardsManager     │",
		"│MapPage││MapKeywords│└────────────────────────┘",
		"│       ││           │┌─────────────┐┌─────────┐",
		"│       ││           ││KeywordsCards││pageCards│",
		"└───────┘└───────────┘└─────────────┘└─────────┘",
	}

	databaseStruct = []string{
		"┌────────────────────┐ ┌──────┐",
		"│                    │ │Socket│",
		"│DatabaseControlModel│ └──────┘",
		"│                    └────────┐",
		"└─────────────────────────────┘",
		"┌────────────┐┌───────────────┐",
		"│     API    ││  SQL Diriver  │",
		"└────────────┘└───────────────┘"}
)

func condition(param []string) (ret string) {
	for _, v := range param {
		ret += v + "\n"
	}
	return ret
}

func CrawlerStruct() string {
	return condition(crawlerStruct)
}

func CircleCondition() string {
	return condition(circleCondition)
}

func DividerStruct() string {
	return condition(dividerStruct)
}

func CrawlerLifeCycle() string {
	return condition(crawlerLifeCycle)
}

func IndexDeviceStruct() string {
	return condition(indexDeviceStruct)
}

func DatabaseStruct() string {
	return condition(databaseStruct)
}
