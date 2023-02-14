package routerinit

import (
	"github.com/gin-gonic/gin"
)

type Option func(engine *gin.Engine)

var options []Option

func Include(optObjs ...Option) {

	options = append(options, optObjs...)
}

func InitRouters() *gin.Engine {

	router := gin.Default()
	for _, opt := range options {
		opt(router)
	}

	return router
}
