package bnrequest

type QueryOrderRequest struct {
	ID     string           `json:"id"`
	Method string           `json:"method"`
	Params *QueryOrderRaram `json:"params"`
}

func (q *QueryOrderRequest) SetSignature() {
	q.Params.Signature = "signature"
}

type QueryOrderRaram struct {
	APIKey     string `json:"apiKey"`
	RecvWindow int64  `json:"recvWindow"`
	Timestamp  int64  `json:"timestamp"`
	Signature  string `json:"signature"`
}

func (q *QueryOrderRaram) SetSignature(signature string) {
	q.Signature = signature
}
