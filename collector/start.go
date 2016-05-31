package collector

import(
  "sync"
  log "github.com/Sirupsen/logrus"
  "github.com/yosmudge/slackmood/slack"
)

// Start the emoji collector
// Connects to Slack and starts streaming all public channels
func Start(wg *sync.WaitGroup) bool{
  s, err := slack.Connect()
  if err != nil {
    log.WithFields(log.Fields{
      "error": err,
    }).Error("Could not connect to Slack!")
    return false
  }

  wg.Add(1)
  go collector(s, wg)

  return true
}

func collector(s *slack.Slack, wg *sync.WaitGroup){
  defer wg.Done()
  wg.Add(1)
  go customEmoji(s, wg)
  wg.Add(1)
  go channels(s, wg)

  wg.Add(1)
  go updateChannels(s, wg)
}
