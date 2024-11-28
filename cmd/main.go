package main

import (
	"bn_query_order/app/handler"
	"bn_query_order/app/proxy"
	"bn_query_order/config"
	"log"

	"github.com/gin-gonic/gin"
)

func registerRoute(router *gin.Engine, _config *config.Config) {
	_proxy := proxy.NewQueryOrderProxy(_config)
	_handler := handler.NewQueryOrderHandler(_proxy)
	router.GET("/get_current_order", _handler.Handler)
}

func main() {
	_config, err := config.ReadConfig()
	if err != nil {
		log.Println("error read config:", err)
		return
	}
	_route := gin.Default()
	registerRoute(_route, _config)
	_route.Run(":8080")
}
