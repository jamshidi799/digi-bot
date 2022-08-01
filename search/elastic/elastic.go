package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"search/model"
	"strconv"
)

var es *elasticsearch.Client

func InitElasticsearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			//"https://my-deployment-dff137.client.us-central1.gcp.cloud.client.io",
			"https://my-deployment-dff137.es.us-central1.gcp.cloud.es.io",
		},
		Username: "elastic",
		Password: "xDBr9Ij9IfiZvYxaqU9Y0mqC",
	}

	client, _ := elasticsearch.NewClient(cfg)
	log.Println("elasticsearch connected")

	es = client
}

func AddProduct(product model.Product) {
	data, _ := json.Marshal(product)
	req := esapi.IndexRequest{
		Index:      "product",
		DocumentID: strconv.Itoa(product.Id),
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[%s] Error indexing document error: %s", res.Status(), res.String())
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			log.Printf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and indexed document version.
			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
			log.Printf("%v", r)
		}
	}
}
