package main

import(
  "os"
  "sync"
  log "github.com/Sirupsen/logrus"
  "github.com/voxelbrain/goptions"

  "github.com/samarudge/slackmood/web"
  "github.com/samarudge/slackmood/models"
  "github.com/samarudge/slackmood/config"
  "github.com/samarudge/slackmood/collector"
)

type options struct {
  Verbose   bool            `goptions:"-v, --verbose, description='Log verbosely'"`
  Help      goptions.Help   `goptions:"-h, --help, description='Show help'"`
  Config    string          `goptions:"-c, --config, description='Config Yaml file to use'"`
  Bind      string          `goptions:"-b, --bind, description='Port/Address to bind on, can also be specified with WEB_BIND environment variable'"`

  goptions.Verbs
}

func main() {
  parsedOptions := options{}

  parsedOptions.Config = "./config.yml"
  parsedOptions.Bind = os.Getenv("WEB_BIND")

  goptions.ParseAndFail(&parsedOptions)

  if parsedOptions.Verbose{
    log.SetLevel(log.DebugLevel)
  } else {
    log.SetLevel(log.InfoLevel)
  }

  log.SetFormatter(&log.TextFormatter{FullTimestamp:true})

  log.Debug("Logging verbosely!")

  err := config.LoadConfig(parsedOptions.Config)
  if err != nil{
    log.WithFields(log.Fields{
      "configFile": parsedOptions.Config,
      "error": err,
    }).Error("Could not load config file")
    os.Exit(1)
  }

  err = models.OpenDB()
  if err != nil {
    log.WithFields(log.Fields{
      "error": err,
    }).Error("Could not open database")
    return
  }

  var wg sync.WaitGroup

  if !collector.Start(&wg){
    os.Exit(1)
  }
  web.Start(parsedOptions.Bind)

  wg.Wait()
}
