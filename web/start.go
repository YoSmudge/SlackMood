package web

import(
  "github.com/gin-gonic/gin"
  "github.com/hoisie/mustache"
  log "github.com/Sirupsen/logrus"
)

var router = gin.New()

func Start(bind string){
  router.Use(gin.Recovery())
  router.GET("/", Overview)
  router.Run(bind)
}

func Render(c *gin.Context, filePath string, obj map[string]interface{}){
  templateData, err := Asset(filePath)

  if err != nil {
    log.WithFields(log.Fields{
      "path": filePath,
    }).Error("Could not find template file")
    c.String(500, "Template not found")
  } else {
    for k, v := range c.Keys {
      obj[k] = v
    }

    mainTemplate, _ := Asset("main.html")
    html := mustache.RenderInLayout(string(templateData), string(mainTemplate), obj)

    if c.Writer.Status() == 200{
      c.Status(200)
    }
    c.Writer.Write([]byte(html))
  }
}
