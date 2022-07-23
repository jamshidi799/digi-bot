package crawler

import (
	"digi-bot/model"
	"fmt"
	"github.com/m7shapan/njson"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func (dkCrawler DigikalaCrawler) Crawl(url string) (dto model.ProductDto, err error) {
	regex := regexp.MustCompile(`.*dkp-(\d*).*`)
	id := regex.FindStringSubmatch(url)[1]

	res, err := http.Get(fmt.Sprintf("https://api.digikala.com/v1/product/%s/", id))
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	var product product
	err = njson.Unmarshal(body, &product)
	if err != nil {
		log.Fatalln(err)
	}

	dto.Id = product.ID
	dto.Url = url
	dto.Name = product.Title
	dto.Price = product.Price
	dto.OldPrice = product.OldPrice

	dto.Status = 1
	if dto.Price == 0 {
		dto.Status = 0
	}

	fmt.Printf("%+v\n", dto)
	return
}

type product struct {
	ID         int    `njson:"data.product.id"`
	Title      string `njson:"data.product.title_fa"`
	Image      string `njson:"data.product.images.main.url.0"`
	Price      int    `njson:"data.product.default_variant.price.selling_price"`
	OldPrice   int    `njson:"data.product.default_variant.price.rrp_price"`
	OrderLimit int    `njson:"data.product.default_variant.price.order_limit"`
}

//func (dkCrawler DigikalaCrawler) Crawl2(url string) (model.ProductDto, error) {
//	c := colly.NewCollector(
//		colly.AllowedDomains("digikala.com", "www.digikala.com"),
//	)
//
//	var product model.ProductDto
//	c.OnHTML("script[type=\"application/ld+json\"]", func(e *colly.HTMLElement) {
//		err := json.Unmarshal([]byte(e.Text), &product)
//		if err != nil {
//			log.Println("err")
//		}
//		product.Url = url
//	})
//
//	c.OnHTML(".c-product", func(e *colly.HTMLElement) {
//		price := e.ChildText(".c-product__seller-price-pure.js-price-value")
//		product.Price = utils.FixNumber(price)
//
//		oldPrice := e.ChildText(".c-product__seller-price-prev.js-rrp-price")
//		product.OldPrice = utils.FixNumber(oldPrice)
//
//		product.Status = 1
//		if product.Price == 0 {
//			product.Status = 0
//		}
//
//		//fmt.Printf("%s %s\n", price, oldPrice)
//
//		desc1 := e.ChildText(`div[class="c-product__headline--gallery "]`)
//		product.Desc1 = utils.CleaningString(desc1)
//
//		if desc1 == "" {
//			seconds := e.ChildAttr(`.c-product-gallery__timer.js-counter`, "data-countdownseconds")
//			if seconds != "" {
//				product.Desc1 = utils.CreateAmazingOfferText(seconds)
//				product.Status = 2
//			}
//		}
//
//		desc2 := e.ChildText(".c-product__user-suggestion-line")
//		desc2 = utils.CleaningString(desc2)
//		desc2 = strings.Split(desc2, ".")[0]
//		product.Desc2 = utils.CleaningString(desc2)
//
//		desc3 := e.ChildText(".c-product__engagement-rating")
//		if desc3 == "" {
//			product.Desc3 = ""
//		} else {
//			product.Desc3 = "امتیاز " + utils.CleaningString(desc3)
//		}
//	})
//
//	err := c.Visit(url)
//	if err != nil {
//		println("visit error")
//		return product, errors.New("visit error")
//	}
//
//	return product, nil
//
//}
