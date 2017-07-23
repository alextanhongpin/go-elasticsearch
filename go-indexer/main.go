package main

import (
	"context"
	"log"
	"strconv"

	"github.com/alextanhongpin/go-elasticsearch/go-indexer/elasticsearch"

	elastic "gopkg.in/olivere/elastic.v5"
)

// Env defines our environment variables
type Env struct {
	Hosts []string
	Index string
	Type  string
}

// Person defines the document we want to index to ElasticSearch
type Person struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

func main() {
	data := []Person{
		Person{"john", 1},
		Person{"doe", 2},
		Person{"john", 3},
	}
	ctx := context.Background()
	esClient, err := elasticsearch.MakeClient()
	if err != nil {
		panic(err)
	}

	hosts, err := esClient.Service("global-elastichsearch-check", "search")
	if err != nil {
		panic(err)
	}
	log.Println("connecting to host", hosts)
	env := Env{
		Hosts: hosts,
		Index: "index",
		Type:  "Person",
	}

	// Index maximum of 5k documents at a time
	LIMIT := 5000

	client, err := elastic.NewClient(
		elastic.SetURL(env.Hosts...),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "changeme"),
	)
	if err != nil {
		panic(err)
	}
	defer client.Stop()

	dataLen := len(data)
	if dataLen == 0 {
		log.Println("No data to index")
		return
	}
	for i := 0; i < dataLen; i += LIMIT {
		batch := data[i:min(i+LIMIT, dataLen)]
		log.Printf("At batch %d\n", i/LIMIT)

		bulkRequest := client.Bulk()
		for _, item := range batch {
			req := elastic.NewBulkIndexRequest().
				Index(env.Index).
				Type(env.Type).
				Id(strconv.Itoa(item.ID)).
				Doc(item)
			bulkRequest = bulkRequest.Add(req)
		}

		log.Println("Performing bulk index...")
		bulkResponse, err := bulkRequest.Do(ctx)

		if err != nil {
			log.Println(err)
		}
		if bulkResponse != nil {
			log.Println("Done indexing for index:", i/LIMIT)
		}
	}
	log.Printf("Done indexing %d item", dataLen)
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
