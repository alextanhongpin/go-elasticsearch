package main

import (
	"context"
	"log"

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

	// Log the initial document count
	count, err := client.Count(env.Index).Do(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Got %d documents\n", count)

	q := elastic.NewTermQuery("name", "john")

	res, err := client.DeleteByQuery().
		Index(env.Index).
		Type(env.Type).
		Query(q).
		Pretty(true).
		Do(ctx)

	if err != nil {
		log.Fatal(err)
	}

	if res == nil {

	}
	// Flush and check count
	_, err = client.Flush().Index(env.Index).Do(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Log the initial document count
	count, err = client.Count(env.Index).Do(ctx)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Got %d documents\n", count)
}
