package collector

import(
  "time"
  "sync"
  api "github.com/nlopes/slack"
  log "github.com/Sirupsen/logrus"
  "github.com/yosmudge/slackmood/slack"
  "github.com/yosmudge/slackmood/models"
)

// Fetch each channel in sequence and get the messages
func updateChannels(s *slack.Slack, wg *sync.WaitGroup){
  defer wg.Done()
  time.Sleep(time.Second*1)

  for{
    channels := channelList
    log.WithFields(log.Fields{
      "channels": len(channels),
    }).Debug("Fetching channel history")

    for _,c := range channels{
      if c.IsArchived{
        continue
      }
      hp := api.NewHistoryParameters()
      hp.Count = 1000
      h, err := s.Api.GetChannelHistory(c.ID, hp)

      if err != nil {
        log.WithFields(log.Fields{
          "error": err,
          "channelId": c.ID,
          "channel": c,
        }).Warning("Could not fetch channel history")
      } else {
        models.ParseEmoji(h.Messages)

        log.WithFields(log.Fields{
          "channel": c.Name,
          "channelId": c.ID,
          "messages": len(h.Messages),
        }).Debug("Got channel history")
      }

      //time.Sleep(time.Second*1)
    }

    time.Sleep(time.Second*5)
  }
}
