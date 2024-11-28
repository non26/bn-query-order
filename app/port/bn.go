package port

import handlerresponse "bn_query_order/app/handler_response"

type IBN interface {
	Run()
	SetApiKey(apiKey string)
	SetPrivateKey(privateKey string)
	SetMethod(method string)
	SetUrl(url string)
	SetSymbol(symbol string)
	GetResponse() *handlerresponse.QueryOrderResponse
}
