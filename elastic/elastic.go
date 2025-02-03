package elastic

import (
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)

var ES *elasticsearch.Client

func ConnectES() {
	cfg := elasticsearch.Config{
		Addresses: []string{os.Getenv("ElasticSTR")},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal("Error creating Elasticsearch client: ", err)
	}
	ES = client
}