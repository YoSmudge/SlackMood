package collector

import (
	log "github.com/Sirupsen/logrus"
	api "github.com/nlopes/slack"
	"github.com/yosmudge/slackmood/slack"
)

var channelList []api.Channel

// Get a list of all the channels
func channels(s *slack.Slack) {
	log.Debug("Fetching channels")

	chn, err := s.Api.GetChannels(false)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warning("Could not fetch channels from Slack")
	} else {
		var allChannels []api.Channel
		for _, channel := range chn {
			allChannels = append(allChannels, channel)
		}

		channelList = allChannels
	}

	log.WithFields(log.Fields{
		"channels": len(channelList),
	}).Debug("Updated channel list")
}
