package main

import (
	"context"
	"log"

	"github.com/alextanhongpin/go-elasticsearch/golang/elasticsearch"

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
	Name string
	ID   int
}

func main() {
	data := []Person{}
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

// func endpoints() []string {
// 	var endpoints []string
// 	config := consul.DefaultConfig()
// 	// config.Address = "localhost:8500"
// 	c, err := consul.NewClient(config)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	addrs, meta, err := c.Health().Service("global-elastichsearch-check", "search", true, nil)

// 	if len(addrs) == 0 && err == nil {
// 		log.Println("Service not found")
// 		return endpoints
// 	}
// 	if err != nil {
// 		log.Println(err)
// 		return endpoints
// 	}
// 	log.Println(meta)
// 	for _, v := range addrs {
// 		log.Printf("%#v", v.Service)
// 		log.Printf("%#v", v.Service.Address)
// 		log.Printf("%#v", v.Service.Port)
// 		endpoints = append(endpoints, "http://"+v.Service.Address+":"+strconv.Itoa(v.Service.Port))
// 	}
// 	return endpoints
// }
