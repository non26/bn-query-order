package proxy

import (
	"bn_query_order/app/bn"
	handlerresponse "bn_query_order/app/handler_response"
	"bn_query_order/app/port"
	"bn_query_order/config"
	"context"
)

type QueryOrderProxy struct {
	config *config.Config
}

func NewQueryOrderProxy(config *config.Config) port.IProxy {
	return &QueryOrderProxy{config: config}
}

func (p *QueryOrderProxy) QueryOrder(ctx context.Context) (*handlerresponse.QueryOrderResponse, error) {
	_bn := bn.NewQueryCurrentOrder()
	_bn.SetApiKey(p.config.BnCredentials.APIKey)
	_bn.SetPrivateKey(p.config.BnCredentials.PrivateKey)
	_bn.SetMethod(p.config.Bn.Method.PositionInformation)
	_bn.SetUrl(p.config.Bn.WsURL)
	_bn.Run()
	// _handlerResponse := _bn.GetResponse()
	return nil, nil

}
