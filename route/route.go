package route

import (
    "gopkg.in/gin-gonic/gin.v1"
    "dnt-swr/handler"
)

func RouteAll(listenInfo string) {
    router := gin.Default()
    gin.SetMode(gin.ReleaseMode)
    router.GET("/devops/status", handler.HandleCheckHealth)
    router.HEAD("/devops/status", handler.HandleCheckHealth)
    router.POST("/devops/status", handler.HandleSetHealth)
    router.POST("/service/notifications", handler.HandleDnt)
    router.POST("/service/repoaction", handler.HandleRepoAction)
    router.Run(listenInfo)
}
