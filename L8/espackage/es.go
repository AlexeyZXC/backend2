package espackage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"l8es/domain"
	"log"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"

	"github.com/google/uuid"
)

type ES struct {
	ES *elasticsearch.Client
}

func New() (*ES, error) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return &ES{ES: es}, nil
}

func (es *ES) Add(product *domain.Product) error {
	jProduct, err := json.Marshal(product)
	if err != nil {
		return fmt.Errorf("error marshaling product: %w", err)
	}

	// request := esapi.IndexRequest{Index: "stsc", DocumentID: uid, Body: strings.NewReader(string(jsonString))}
	request := esapi.IndexRequest{Index: "products", DocumentID: uuid.NewString(), Body: bytes.NewReader(jProduct)}
	resp, err := request.Do(context.Background(), es.ES)
	if err != nil {
		return fmt.Errorf("error IndexRequest product: %w", err)
	}
	defer resp.Body.Close()
	log.Println("resp: ", resp)
	return nil
}

func (es *ES) Search(key, value, querytype string) []map[string]any {
	// key := ReadText(reader, "Enter key")
	// value := ReadText(reader, "Enter value")
	// var buffer bytes.Buffer
	query := map[string]any{
		"query": map[string]any{
			querytype: map[string]any{
				key: value,
			},
		},
	}
	// json.NewEncoder(&buffer).Encode(query)
	buf, _ := json.Marshal(query)
	response, _ := es.ES.Search(es.ES.Search.WithIndex("products"),
		// es.ES.Search.WithBody(&buffer))
		es.ES.Search.WithBody(bytes.NewReader(buf)))

	defer response.Body.Close()
	var result map[string]any
	ret := []map[string]any{}

	json.NewDecoder(response.Body).Decode(&result)
	for _, hit := range result["hits"].(map[string]any)["hits"].([]any) {
		craft :=
			hit.(map[string]any)["_source"].(map[string]any)
		// Print(craft)
		ret = append(ret, craft)
	}

	return ret
}
