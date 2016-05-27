package web

import(
  "time"
  "github.com/gin-gonic/gin"
  "github.com/samarudge/slackmood/models"
)

func Overview(c *gin.Context){
  mood := models.GetMood(time.Hour*12)
  c.JSON(200, mood)
}
