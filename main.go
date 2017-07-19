package main

import (
	"context"
	"log"

	elastic "gopkg.in/olivere/elastic.v5"
)

// Env defines our environment variables
type Env struct {
	Host  string
	Index string
	Type  string
}

// Person defines the document we want to index to ElasticSearch
type Person struct {
	Name string
	ID   int
}

func main() {
	data := []Person{}
	ctx := context.Background()

	env := Env{
		Host:  "localhost:9200",
		Index: "index",
	}

	// Index maximum of 5k documents at a time
	LIMIT := 5000

	client, err := elastic.NewClient(
		elastic.SetURL(env.Host),
		elastic.SetSniff(false),
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
				Id("string_id").
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
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}
