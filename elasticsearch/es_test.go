package elasticsearch

import (
	"fmt"
	"testing"
)

func TestNewESClient(t *testing.T) {
	esConfig := &ESConfig{
		Addresses:  []string{"http://127.0.0.1:9200"},
		Username:   "",
		Password:   "",
		SearchSize: 100,
	}
	clinet, err := NewESClient(esConfig, nil, 10)
	if err != nil {
		panic(err)
	}
	query := BuildModuleQuery(BuildBasicQuery(), map[string]interface{}{"title": "yuxingwang"})
	esIndex := "website"
	DocumentType := "blog"
	searchSize := 100
	hits, indexExist, err := clinet.searchAfter(query, esIndex, DocumentType, searchSize, "")
	if err != nil && !indexExist {
		fmt.Printf("es search after failed [%s], "+
			"query: %v, "+"index: %s, type: %s, "+
			"size: %d",
			err, query, esIndex, DocumentType, searchSize)
	} else {
		for _, hit := range hits {
			if hit == nil {
				fmt.Printf("es search after hit has not "+
					"'_source' field, "+"query: %v, index: %s, type: %s, size: %d",
					query, esIndex, DocumentType, searchSize)
			} else {
				fmt.Println(*hit.Source, hit.Id)
			}

		}

	}
}
