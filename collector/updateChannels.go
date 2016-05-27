package collector

import(
  "time"
  "sync"
  "strconv"
  api "github.com/nlopes/slack"
  log "github.com/Sirupsen/logrus"
  "github.com/samarudge/slackmood/slack"
  "github.com/samarudge/slackmood/models"
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
      if c.Name != "slackmood-test"{
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
        var relevantMessages []api.Message
        now := time.Now().UTC()
        for _,m := range h.Messages{
          ts, _ := strconv.ParseFloat(m.Timestamp, 64)
          t := time.Unix(int64(ts), 0)
          if t.After(now.Add(time.Hour*-12)){
            relevantMessages = append(relevantMessages, m)
          }
        }

        models.ParseEmoji(relevantMessages)

        log.WithFields(log.Fields{
          "channel": c.Name,
          "channelId": c.ID,
          "messages": len(h.Messages),
          "relevantMessages": len(relevantMessages),
        }).Debug("Got channel history")
      }

      time.Sleep(time.Second*1)
    }

    time.Sleep(time.Second*5)
  }
}
