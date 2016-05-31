package collector

import(
  "time"
  "sync"
  api "github.com/nlopes/slack"
  log "github.com/Sirupsen/logrus"
  "github.com/yosmudge/slackmood/slack"
)

var channelList []api.Channel

// Get a list of all the channels
func channels(s *slack.Slack, wg *sync.WaitGroup){
  defer wg.Done()

  for{
    log.Debug("Fetching channels")

    chn, err := s.Api.GetChannels(false)
    if err != nil {
      log.WithFields(log.Fields{
        "error": err,
      }).Warning("Could not fetch channels from Slack")
    } else {
      var allChannels []api.Channel
      for _,channel := range chn{
        allChannels = append(allChannels, channel)
      }

      channelList = allChannels
    }

    log.WithFields(log.Fields{
      "channels": len(channelList),
    }).Debug("Updated channel list")

    time.Sleep(time.Minute*5)
  }
}
