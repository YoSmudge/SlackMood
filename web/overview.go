package web

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yosmudge/slackmood/models"
)

type timePeriod struct {
	Name      string
	Period    time.Duration
	Breakdown time.Duration
	Active    bool
}

var timePeriods = []timePeriod{
	timePeriod{"24h", time.Hour * 24, time.Hour, false},
	timePeriod{"7d", time.Hour * 24 * 7, time.Hour * 4, false},
	timePeriod{"31d", time.Hour * 24 * 31, time.Hour * 24, false},
	timePeriod{"90d", time.Hour * 24 * 90, time.Hour * 24 * 7, false},
}

// Overview does something
func Overview(c *gin.Context) {
	periods := timePeriods
	period := timePeriod{}

	var validPeriod bool
	periodName := c.DefaultQuery("period", timePeriods[2].Name)
	for i, p := range periods {
		periods[i].Active = false
		if p.Name == periodName {
			validPeriod = true
			period = p
			periods[i].Active = true
		}
	}

	if !validPeriod {
		c.String(410, "Invalid Period")
		return
	}

	mood := models.GetMood(models.FilterEmoji(time.Now().UTC().Add(period.Period*-1), time.Now().UTC(), models.AllEmoji()))
	graphData := models.GraphMood(period.Period, period.Breakdown)
	graphJSON, _ := json.Marshal(graphData)

	Render(c, "overview.html", gin.H{
		"currentMood":   mood,
		"timePeriods":   timePeriods,
		"moodGraphJson": string(graphJSON),
		"totalEmoji":    len(models.AllEmoji()),
	})
}
