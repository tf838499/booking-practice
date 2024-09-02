package router

import (
	"context"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetGeneralMiddlewares add general-purpose middlewares
func SetGeneralMiddlewares(ctx context.Context, ginRouter *gin.Engine) {
	// set swagger
	if mode := gin.Mode(); mode == gin.DebugMode {
		ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	ginRouter.Use(gin.Recovery())
	// ginRouter.Use(CORSMiddleware())
	ginRouter.Use(requestid.New())
	// ginRouter.Use(LoggerMiddleware(ctx))

}
