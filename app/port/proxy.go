package port

import (
	handlerresponse "bn_query_order/app/handler_response"
	"context"
)

type IProxy interface {
	QueryOrder(ctx context.Context) (*handlerresponse.QueryOrderResponse, error)
}
