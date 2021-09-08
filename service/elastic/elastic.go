package main

import (
	"bytes"
	"context"
	"digi-bot/db"
	"digi-bot/model"
	"encoding/json"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"strings"
	"sync/atomic"
	"time"
)

var es *elasticsearch.Client

func main() {
	client, err := elasticsearch.NewDefaultClient()
	es = client
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	db.Init()

	var histories []model.BulkHistory
	db.DB.Table("bulk_histories").Find(&histories)

	bulkDocument := make([]interface{}, len(histories))
	for i, v := range histories {
		bulkDocument[i] = v
		break
	}
	//bulkIndexer("history", bulkDocument)
	//addIndex("product", db.GetProductById(4))
}

type StoreConfig struct {
	Client    *elasticsearch.Client
	IndexName string
}

type Store struct {
	es        *elasticsearch.Client
	indexName string
}

func NewStore(c StoreConfig) (*Store, error) {
	indexName := c.IndexName
	if indexName == "" {
		indexName = "xkcd"
	}

	s := Store{es: c.Client, indexName: indexName}
	return &s, nil
}

type Document struct {
	ID        string `json:"id"`
	ImageURL  string `json:"image_url"`
	Published string `json:"published"`

	Title      string `json:"title"`
	Alt        string `json:"alt"`
	Transcript string `json:"transcript"`
	Link       string `json:"link,omitempty"`
	News       string `json:"news,omitempty"`
}

func addIndex(indexName string, document interface{}) {
	res, _ := es.Index(indexName, esutil.NewJSONReader(&document))
	fmt.Println(res)
}

func bulkIndexer(indexName string, documents []interface{}) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  indexName,
		Client: es,
	})
	if err != nil {
		log.Fatalln(err)
	}

	var countSuccessful uint64
	start := time.Now().UTC()

	for _, a := range documents {
		// Prepare the data payload: encode article to JSON
		//
		data, err := json.Marshal(a)
		if err != nil {
			log.Fatalf("Cannot encode document: %s", err)
		}

		// Add an item to the BulkIndexer
		//
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				// Action field configures the operation to perform (index, create, delete, update)
				Action: "index",

				// DocumentID is the (optional) document ID
				//DocumentID: strconv.Itoa(a.ID),

				// Body is an `io.Reader` with the payload
				Body: bytes.NewReader(data),

				// OnSuccess is called for each successful operation
				OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},

				// OnFailure is called for each failed operation
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
	}

	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	biStats := bi.Stats()
	dur := time.Since(start)

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	}
}

func query() {
	var b bytes.Buffer

	res, _ := es.Search(es.Search.WithTrackTotalHits(true))
	b.ReadFrom(res.Body)

	values := gjson.GetManyBytes(b.Bytes(), "hits.total", "took")
	//fmt.Printf(
	//	"[%s] %d hits; took: %dms\n",
	//	res.Status(),
	//	values[0].Int(),
	//	values[1].Int(),
	//)
	fmt.Println(values)
}

func (s *Store) Exists(id string) (bool, error) {
	res, err := s.es.Exists(s.indexName, id)
	if err != nil {
		return false, err
	}
	switch res.StatusCode {
	case 200:
		return true, nil
	case 404:
		return false, nil
	default:
		return false, fmt.Errorf("[%s]", res.Status())
	}
}

const searchAll = `
	"query" : { "match_all" : {} },
	"size" : 25,
	"sort" : { "published" : "desc", "_doc" : "asc" }`

const searchMatch = `
	"query" : {
		"multi_match" : {
			"query" : %q,
			"fields" : ["title^100", "alt^10", "transcript"],
			"operator" : "and"
		}
	},
	"highlight" : {
		"fields" : {
			"title" : { "number_of_fragments" : 0 },
			"alt" : { "number_of_fragments" : 0 },
			"transcript" : { "number_of_fragments" : 5, "fragment_size" : 25 }
		}
	},
	"size" : 25,
	"sort" : [ { "_score" : "desc" }, { "_doc" : "asc" } ]`

//func (s *Store) Search(query string, after ...string) (*SearchResults, error) {
//	var results SearchResults
//
//	res, err := s.es.Search(
//		s.es.Search.WithIndex(s.indexName),
//		s.es.Search.WithBody(s.buildQuery(query, after...)),
//	)
//	if err != nil {
//		return &results, err
//	}
//	defer res.Body.Close()
//
//	if res.IsError() {
//		var e map[string]interface{}
//		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
//			return &results, err
//		}
//		return &results, fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
//	}
//
//	type envelopeResponse struct {
//		Took int
//		Hits struct {
//			Total struct {
//				Value int
//			}
//			Hits []struct {
//				ID         string          `json:"_id"`
//				Source     json.RawMessage `json:"_source"`
//				Highlights json.RawMessage `json:"highlight"`
//				Sort       []interface{}   `json:"sort"`
//			}
//		}
//	}
//
//	var r envelopeResponse
//	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
//		return &results, err
//	}
//
//	results.Total = r.Hits.Total.Value
//
//	if len(r.Hits.Hits) < 1 {
//		results.Hits = []*Hit{}
//		return &results, nil
//	}
//
//	for _, hit := range r.Hits.Hits {
//		var h Hit
//		h.ID = hit.ID
//		h.Sort = hit.Sort
//		h.URL = strings.Join([]string{baseURL, h.ID, ""}, "/")
//
//		if err := json.Unmarshal(hit.Source, &h); err != nil {
//			return &results, err
//		}
//
//		if len(hit.Highlights) > 0 {
//			if err := json.Unmarshal(hit.Highlights, &h.Highlights); err != nil {
//				return &results, err
//			}
//		}
//
//		results.Hits = append(results.Hits, &h)
//	}
//
//	return &results, nil
//}

func (s *Store) buildQuery(query string, after ...string) io.Reader {
	var b strings.Builder

	b.WriteString("{\n")

	if query == "" {
		b.WriteString(searchAll)
	} else {
		b.WriteString(fmt.Sprintf(searchMatch, query))
	}

	if len(after) > 0 && after[0] != "" && after[0] != "null" {
		b.WriteString(",\n")
		b.WriteString(fmt.Sprintf(`	"search_after": %s`, after))
	}

	b.WriteString("\n}")

	// fmt.Printf("%s\n", b.String())
	return strings.NewReader(b.String())
}
