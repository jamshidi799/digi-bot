package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"log"
	"search/model"
	"strconv"
	"strings"
)

var es *elasticsearch.Client

func InitElasticsearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://my-deployment-3d2ec6.es.us-central1.gcp.cloud.es.io",
		},
		Username: "elastic",
		Password: "1OZeXWG8ZXB4mvVlfPljMQdW",
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

type SearchResponse struct {
	Took int
	Hits struct {
		Total struct {
			Value int
		}
		Hits []struct {
			ID      string        `json:"_id"`
			Product model.Product `json:"_source"`
		}
	}
}

func SearchProduct(ctx context.Context, term string) []model.Product {
	query := fmt.Sprintf(`
{
  "query": {
    "match_bool_prefix" : {
      "Name" : "%s"
    }
  }
}
`, term)

	res, err := es.Search(
		es.Search.WithContext(ctx),
		es.Search.WithIndex("product"),
		es.Search.WithBody(strings.NewReader(query)),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
		es.Search.WithSize(5),
	)

	if err != nil {
		log.Println(err)
		return nil
	}

	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil
		}
		return nil
	}

	var r SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Println(err)
		return nil
	}

	var products []model.Product
	for _, hit := range r.Hits.Hits {
		products = append(products, hit.Product)
	}

	return products
}

//{
//  "query": {
//    "bool": {
//      "should": [
//        {
//          "wildcard": {
//            "name": {
//              "value": "*شارژ*",
//              "boost": 1,
//              "rewrite": "constant_score"
//            }
//          }
//        },
//        {
//          "wildcard": {
//            "name": {
//              "value": "*سی*",
//              "boost": 1,
//              "rewrite": "constant_score"
//            }
//          }
//        }
//      ],
//      "minimum_should_match": 1,
//      "boost": 1
//    }
//  }
//}
