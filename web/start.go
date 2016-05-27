package web

import(
  "github.com/gin-gonic/gin"
)

var router = gin.New()

func Start(bind string){
  router.Use(gin.Recovery())
  router.GET("/", Overview)
  router.Run(bind)
}
