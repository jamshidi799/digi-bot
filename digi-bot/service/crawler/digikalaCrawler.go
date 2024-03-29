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

type DigikalaCrawler struct {
}

func (DigikalaCrawler) Crawl(url string) (dto model.ProductDto, err error) {
	regex := regexp.MustCompile(`.*dkp-(\d*).*`)
	id := regex.FindStringSubmatch(url)[1]

	res, err := http.Get(fmt.Sprintf("https://api.digikala.com/v1/product/%s/", id))
	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	var product digikalaProduct
	err = njson.Unmarshal(body, &product)
	if err != nil {
		log.Println(err)
		return
	}

	dto.Sku = product.Sku
	dto.Url = url
	dto.Name = product.Title
	dto.Price = product.Price / 10
	dto.OldPrice = product.OldPrice / 10
	dto.Image = product.Image
	dto.Discount = product.Discount

	dto.Status = 1
	if dto.Price == 0 {
		dto.Status = 0
	}

	return dto, err
}

func (c *DigikalaCrawler) GetDomain() string {
	return "digikala.com"
}

type digikalaProduct struct {
	Sku        string `njson:"data.product.id"`
	Title      string `njson:"data.product.title_fa"`
	Image      string `njson:"data.product.images.main.url.0"`
	Price      int    `njson:"data.product.default_variant.price.selling_price"`
	OldPrice   int    `njson:"data.product.default_variant.price.rrp_price"`
	OrderLimit int    `njson:"data.product.default_variant.price.order_limit"`
	Discount   int    `njson:"data.product.default_variant.price.discount_percent"`
}
