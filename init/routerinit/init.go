package routerinit

import (
	"TM-Spike/init/configinit"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Option func(engine *gin.Engine)

var options []Option

func Include(optObjs ...Option) {

	options = append(options, optObjs...)
}

func InitRouters() *gin.Engine {

	router := gin.Default()
	router.Use(otelgin.Middleware(configinit.ServiceName))
	for _, opt := range options {
		opt(router)
	}

	return router
}
