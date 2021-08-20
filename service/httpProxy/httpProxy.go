package httpProxy

import (
	"digi-bot/service/httpProxy/entity"
	"fmt"
	"github.com/m7shapan/njson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func Init() {
	for {
		proxies := getFreshHttpProxies()
		log.Println(proxies)
		for _, proxy := range proxies.List {
			httpProxy := fmt.Sprintf("http://%s:%s", proxy.Ip, proxy.Port)
			log.Println(httpProxy)

			proxyUrl, err := url.Parse(httpProxy)
			if err != nil {
				log.Fatalln(err)
			}
			client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

			resp, err := client.Get("http://185.110.189.137/")
			if err != nil {
				log.Println(err)
			} else {
				log.Println(resp.Body)
				resp.Body.Close()
			}
		}
	}
}

func getFreshHttpProxies() entity.Proxies {

	resp, err := http.Get("https://proxylist.geonode.com/api/proxy-list?limit=50&page=1&sort_by=lastChecked&sort_type=desc")
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var proxies entity.Proxies
	err = njson.Unmarshal(body, &proxies)
	if err != nil {
		log.Println("err")
	}
	return proxies
}
