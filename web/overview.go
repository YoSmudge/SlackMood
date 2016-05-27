package web

import(
  "time"
  "github.com/gin-gonic/gin"
  "github.com/samarudge/slackmood/models"
)

type timePeriod struct{
  Name        string
  Period      time.Duration
  Active      bool
}

var timePeriods = []timePeriod{
  timePeriod{"24h",time.Hour*24,false},
  timePeriod{"7d",time.Hour*24*7,false},
  timePeriod{"31d",time.Hour*24*31,false},
}

func Overview(c *gin.Context){
  periods := timePeriods
  period := timePeriod{}

  var validPeriod bool
  periodName := c.DefaultQuery("period", timePeriods[0].Name)
  for i,p := range periods{
    periods[i].Active = false
    if p.Name == periodName{
      validPeriod = true
      period = p
      periods[i].Active = true
    }
  }

  if !validPeriod{
    c.String(410, "Invalid Period")
    return
  }

  mood := models.GetMood(period.Period)

  Render(c, "overview.html", gin.H{
    "currentMood": mood,
    "timePeriods": timePeriods,
  })
}
