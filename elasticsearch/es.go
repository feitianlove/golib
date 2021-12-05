package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/feitianlove/golib/common/security"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type ESConfig struct {
	Addresses  []string `json:"addresses"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	SearchSize int      `json:"search_size"`
}

type ESClient struct {
	es7           *elasticsearch7.Client
	searchSize    int
	searchTimeout int
}

func NewESClient(cfg *ESConfig, security *security.SecurityConfig, timeout int) (*ESClient, error) {
	transport := &http.Transport{
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 200,
		MaxConnsPerHost:     200,
		IdleConnTimeout:     90 * time.Second,
	}
	if security != nil && len(security.SSLCAPath) != 0 {
		tlsConfig, err := security.ToTLSConfig()
		if err != nil {
			return nil, fmt.Errorf("to tls config failed %s", err)
		}
		transport.TLSClientConfig = tlsConfig
	}

	es7, err := elasticsearch7.NewClient(elasticsearch7.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		Transport: transport,
	})
	if err != nil {
		return nil, fmt.Errorf("new esclient v7 client failed %s , addresses: %v", err, cfg.Addresses)
	}
	return &ESClient{
		es7:           es7,
		searchSize:    cfg.SearchSize,
		searchTimeout: timeout,
	}, nil
}

func BuildBasicQuery() map[string]interface{} {
	query := map[string]interface{}{
		"constant_score": map[string]interface{}{
			"filter": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": []map[string]interface{}{
						// { // wildcard
						// 	// "wildcard": map[string]interface{}{
						// 	// 	"actionName": action,
						// 	// },
						// },
					}, // end of must
				}, // end of bool
			}, // end of filter
		}, // end of score
	}
	return query
}

func BuildModuleQuery(query, data map[string]interface{}) map[string]interface{} {
	must := query["constant_score"].(map[string]interface{})["filter"].(map[string]interface{})["bool"].(map[string]interface{})["must"].([]map[string]interface{})
	for key, value := range data {
		must = append(must, map[string]interface{}{ // term
			"term": map[string]interface{}{
				key: value,
			},
		})
	}
	query["constant_score"].(map[string]interface{})["filter"].(map[string]interface{})["bool"].(map[string]interface{})["must"] = must
	return query
}

func (c *ESClient) searchAfter(query map[string]interface{}, index,
	documentType string, size int,
	descCondition string) ([]*Hit, bool, error) {
	//startTime := time.Now()
	dsl := map[string]interface{}{
		"query": query,
	}
	// build the request body.
	var dslBuf bytes.Buffer
	err := json.NewEncoder(&dslBuf).Encode(dsl)
	if err != nil {
		return nil, false, fmt.Errorf("json encode dsl failed err %s , dsl: %v", err, dsl)
	}
	dslStr := dslBuf.String()
	// sort = id:asc
	debugPattern := fmt.Sprintf("Get %s/%s/_search?sort=%s&size=%d",
		index, documentType, descCondition, size)
	searchCtx, searchCancel := context.WithTimeout(context.Background(),
		time.Duration(c.searchTimeout)*time.Second)
	defer searchCancel()
	// first search
	resp, err := c.es7.Search(
		c.es7.Search.WithContext(searchCtx),
		c.es7.Search.WithIndex(index),
		c.es7.Search.WithDocumentType(documentType),
		c.es7.Search.WithBody(&dslBuf),
		c.es7.Search.WithSize(size),
		c.es7.Search.WithSort(descCondition),
	)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, false, errors.Errorf("es search failed %s, pattern: %s, dsl: %s",
			err, debugPattern, dslStr)
	}
	if resp.IsError() {
		if resp.Status() == "404" {
			return nil, true, errors.Errorf("es search failed, status-code: %s, pattern: %s, dsl: %s",
				resp.Status(), debugPattern, dslStr)

		} else {
			return nil, false, errors.Errorf("es search failed, status-code: %s, pattern: %s, dsl: %s",
				resp.Status(), debugPattern, dslStr)
		}
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, false, errors.Errorf("read response body failed %s d, dsl: %s",
			err, dslStr)
	}
	response := &Response{}
	err = json.Unmarshal(respBody, response)
	if err != nil {
		return nil, false, errors.Errorf("json unmarshal response body failed %s, pattern: %s, dsl: %s",
			err, debugPattern, dslStr)
	}
	if response.Hits == nil {
		return nil, false, errors.Errorf("response has not 'hits' field, pattern: %s, dsl: %s",
			debugPattern, dslStr)
	}
	return response.Hits.Hits, false, nil
}
