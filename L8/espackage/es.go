package espackage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"l8es/domain"
	"log"
	"strings"

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

// query := `{ "query": { "match_all": {} } }`
// client.Search(
//     client.Search.WithIndex("my_index"),
//     client.Search.WithBody(strings.NewReader(query)),
// )

func (es *ES) Search(key, value, querytype string) []map[string]any {
	// key := ReadText(reader, "Enter key")
	// value := ReadText(reader, "Enter value")
	// var buffer bytes.Buffer

	// query := map[string]any{
	// 	"query": map[string]any{
	// 		querytype: map[string]any{
	// 			key: value,
	// 		},
	// 	},
	// }

	// query := `{ "query": { "match_all": {} } }`

	// querytype ; prefix or match
	query := fmt.Sprintf(`{ "query": { "%s": { "%s": "%s" } } }`, querytype, key, value)

	// json.NewEncoder(&buffer).Encode(query)
	// buf, _ := json.Marshal(query)
	response, _ := es.ES.Search(es.ES.Search.WithIndex("products"),
		// es.ES.Search.WithBody(&buffer))
		// es.ES.Search.WithBody(bytes.NewReader(buf)))
		es.ES.Search.WithBody(strings.NewReader(query)))

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

// https://www.elastic.co/docs/reference/query-languages/query-dsl/query-dsl-combined-fields-query

// query := fmt.Sprintf(`{ "query": { "combined_fields": { "query": "%s %s", "fields": ["name", "description"] } } }`, name, desc)
// name=name1 desc=desc3
//   {
//     "_score": 0.9808291,
//       "name": "name1",
//       "description": "desc1"
//   },
//   {
//     "_score": 0.9808291,
//       "name": "name3",
//       "description": "desc3"
//   }

// query := fmt.Sprintf(`{ "query": { "combined_fields": { "query": "%s %s", "fields": ["name", "description^2"] } } }`, name, desc)
// name=name1 desc=desc3
//   {
//     "_score": 1.3486402,
//       "name": "name3",
//       "description": "desc3"
//   }
//   {
//     "_score": 0.9808291,
//       "name": "name1",
//       "description": "desc1"
//   },

func (es *ES) SearchAllFields(name, price, desc string) []map[string]any {
	// query := fmt.Sprintf(`{ "query": { "combined_fields": { "query": "%s %s", "fields": ["name", "description"] } } }`, name, desc)

	query := fmt.Sprintf(`{ "query": { "combined_fields": { "query": "%s %s", "fields": ["name", "description^2"] } } }`, name, desc)

	response, _ := es.ES.Search(es.ES.Search.WithIndex("products"),
		es.ES.Search.WithBody(strings.NewReader(query)))

	defer response.Body.Close()
	var result map[string]any
	ret := []map[string]any{}

	bodyStr, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading response body: %s", err)
		return nil
	}

	var prettyJson bytes.Buffer
	err = json.Indent(&prettyJson, []byte(bodyStr), "", "  ")
	if err != nil {
		log.Printf("Error indenting response body: %s", err)
	}
	log.Printf("ES SearchAllFields: prettyJson: %s\n", prettyJson.String())

	// json.NewDecoder(response.Body).Decode(&result)
	json.NewDecoder(strings.NewReader(string(bodyStr))).Decode(&result)
	for _, hit := range result["hits"].(map[string]any)["hits"].([]any) {
		craft :=
			hit.(map[string]any)["_source"].(map[string]any)
		// Print(craft)
		ret = append(ret, craft)
	}

	log.Printf("ES SearchAllFields: result: %+v\n", result)

	return ret
}

// query := `{ "query": { "match_all": {} } }`
func (es *ES) SearchGetAll() []map[string]any {
	query := `{ "query": { "match_all": {} } }`

	response, _ := es.ES.Search(es.ES.Search.WithIndex("products"),
		es.ES.Search.WithBody(strings.NewReader(query)))

	defer response.Body.Close()
	var result map[string]any
	ret := []map[string]any{}

	json.NewDecoder(response.Body).Decode(&result)
	for _, hit := range result["hits"].(map[string]any)["hits"].([]any) {
		craft :=
			hit.(map[string]any)["_source"].(map[string]any)
		ret = append(ret, craft)
	}

	return ret
}
