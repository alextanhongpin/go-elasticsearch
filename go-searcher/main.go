package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

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

	client, err := elastic.NewClient(
		elastic.SetURL(env.Hosts...),
		elastic.SetSniff(false),
		elastic.SetBasicAuth("elastic", "changeme"),
	)
	if err != nil {
		panic(err)
	}
	defer client.Stop()

	// name is the field we are looking at
	termQuery := elastic.NewTermQuery("name", "john")
	log.Println("Performing search...")
	searchResult, err := client.Search().
		Index(env.Index).
		Query(termQuery).
		Sort("id", true).
		From(0).Size(10).
		Pretty(true).
		Do(ctx)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	var p Person
	for _, item := range searchResult.Each(reflect.TypeOf(p)) {
		if t, ok := item.(Person); ok {
			fmt.Printf("Found user by name=%s and id=%d\n", t.Name, t.ID)
		}
	}

	fmt.Printf("Found a total of %d person\n", searchResult.Hits.TotalHits)

	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var p Person
			err := json.Unmarshal(*hit.Source, &p)
			if err != nil {
				log.Println("Deserialization failed")
			}
			fmt.Printf("User with name %s and id %d", p.Name, p.ID)
		}
	} else {
		fmt.Println("Found not users")
	}
}
