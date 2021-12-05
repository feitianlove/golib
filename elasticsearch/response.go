package elasticsearch

// Response resp
type Response struct {
	Hits *Hits `json:"hits"`
}

// Hits hits
type Hits struct {
	Hits []*Hit `json:"hits"`
}

// Hit hit
type Hit struct {
	Id              string       `json:"_id"`
	Source          *interface{} `json:"_source"`
	TotalInvokeTime int64        `json:"totalInvokeTime"`
}
