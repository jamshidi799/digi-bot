package crawler

import (
	"digi-bot/crawler/entity"
	"digi-bot/db"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/m7shapan/njson"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

func StartBatchCrawl() {
	categories := getCategories()
	for _, category := range categories {
		crawlCategory(category)
	}
}

func crawlCategory(category string) {
	log.Printf("batchCrawler: start crawl %s", category)

	addProductCount := 0
	totalProduct := 0
	pageNumber := 1
	for {
		data := sendRequest(category, pageNumber)
		batch := data.ToBatchBulkHistory()
		if len(batch) == 0 {
			log.Println("batch len is zero")
			totalProduct = data.FoundItems
			break
		}

		db.DB.Create(&batch)
		addProductCount += len(batch)

		if pageNumber >= data.Pages {
			totalProduct = data.FoundItems
			break
		}

		if pageNumber%20 == 0 {
			log.Printf("batchCrawler: crawling %s is on page %d/%d", category, pageNumber, data.Pages)
		}
		pageNumber++
		time.Sleep(time.Second * 10)

		if pageNumber > 40 {
			totalProduct = data.FoundItems
			break
		}
	}

	log.Printf("batchCrawler: crawling %s done. added product: %d, all product: %d, request count: %d",
		category, addProductCount, totalProduct, pageNumber)
}

func getCategories() []string {
	var categories []string
	c := colly.NewCollector()

	c.OnHTML(".c-navi-new-list__options-container", func(e *colly.HTMLElement) {
		categories = extractCategoryFromResponse(e, categories)
	})

	err := c.Visit("https://www.digikala.com/")
	if err != nil {
		println("visit error")
	}

	return categories
}

func extractCategoryFromResponse(e *colly.HTMLElement, categories []string) []string {
	isVisited := map[string]bool{}

	e.ForEach(".c-navi-new-list__sublist-option.c-navi-new-list__sublist-option--title .c-navi-new__medium-display-title", func(_ int, item *colly.HTMLElement) {
		href := item.Attr("href")
		re := regexp.MustCompile(".*/search/(.*)/.*")

		if re.MatchString(href) {
			match := re.FindStringSubmatch(href)
			category := match[1]
			if !isVisited[category] {
				categories = append(categories, match[1])
			}
			isVisited[category] = true
		}
	})
	return categories
}

func sendRequest(category string, pageNumber int) entity.BatchCrawlResult {
	url := fmt.Sprintf("https://www.digikala.com/ajax/search/%s/?has_selling_stock=1&pageno=%d&sortby=4", category, pageNumber)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	data := parseResponse(resp.Body)
	return data
}

func parseResponse(resp io.ReadCloser) entity.BatchCrawlResult {
	body, err := ioutil.ReadAll(resp)
	if err != nil {
		log.Fatalln(err)
	}

	var data entity.BatchCrawlResult
	err = njson.Unmarshal(body, &data)
	if err != nil {
		log.Println("err")
	}
	return data
}
