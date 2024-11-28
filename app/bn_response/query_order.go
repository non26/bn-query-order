package bnresponse

import handlerresponse "bn_query_order/app/handler_response"

type QueryOrderResponse struct {
	ID         string                `json:"id"`
	Status     int                   `json:"status"`
	Result     []QueryOrderResult    `json:"result"`
	RateLimits []QueryOrderRateLimit `json:"rateLimits"`
}

func (q *QueryOrderResponse) ToHandlerResponse() *handlerresponse.QueryOrderResponse {
	if len(q.Result) == 0 {
		return nil
	}

	_handlerResponses := &handlerresponse.QueryOrderResponse{}
	for _, v := range q.Result {
		_queryOrderResponseData := handlerresponse.QueryOrderResponseData{
			Symbol: v.Symbol,
		}
		_handlerResponses.Data = append(_handlerResponses.Data, _queryOrderResponseData)
	}
	return _handlerResponses
}

type QueryOrderResult struct {
	Symbol                 string `json:"symbol"`
	PositionSide           string `json:"positionSide"`
	PositionAmt            string `json:"positionAmt"`
	EntryPrice             string `json:"entryPrice"`
	BreakEvenPrice         string `json:"breakEvenPrice"`
	MarkPrice              string `json:"markPrice"`
	UnrealizedProfit       string `json:"unrealizedProfit"`
	LiquidationPrice       string `json:"liquidationPrice"`
	IsolatedMargin         string `json:"isolatedMargin"`
	Notional               string `json:"notional"`
	MarginAsset            string `json:"marginAsset"`
	IsolatedWallet         string `json:"isolatedWallet"`
	InitialMargin          string `json:"initialMargin"`
	MaintMargin            string `json:"maintMargin"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	Adl                    int    `json:"adl"`
	BidNotional            string `json:"bidNotional"`
	AskNotional            string `json:"askNotional"`
	UpdateTime             int    `json:"updateTime"`
}

type QueryOrderRateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
	Count         int    `json:"count"`
}
