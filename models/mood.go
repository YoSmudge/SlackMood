package models

import(
  "time"
  "github.com/yosmudge/slackmood/emojiRanks"
)

type Mood struct{
  Positive      float32
  Negative      float32
  Neutral       float32
  PositiveCount int32
  NegativeCount int32
  NeutralCount  int32
  TotalCount    int32
  Time          time.Time
  TimeString    string
}

func percentage(a int32, b int32) float32{
  if b == 0{
    return 0
  } else {
    return float32(a)/float32(b)*100
  }
}

func GetMood(emoji []*Emoji) Mood{
  m := Mood{}

  for _, e := range emoji{
    for _,r := range ranks.EmojiRanks{
      if r.Name == e.Name{
        switch r.Rank {
        case 1:
          m.PositiveCount += 1
        case 0:
          m.NeutralCount += 1
        case -1:
          m.NegativeCount += 1
        }
        m.TotalCount += 1
        break
      }
    }
  }

  m.Positive = percentage(m.PositiveCount, m.TotalCount)
  m.Negative = percentage(m.NegativeCount, m.TotalCount)
  m.Neutral = percentage(m.NeutralCount, m.TotalCount)

  return m
}

func GraphMood(over time.Duration, interval time.Duration) []Mood{
  var points []Mood

  now := time.Now().UTC()
  dataPointCount := int(over.Seconds()/interval.Seconds())
  endTime := time.Unix(int64(interval.Seconds())*int64(now.Unix()/int64(interval.Seconds())), 0)
  periodEmoji := FilterEmoji(endTime.Add(over*-1), endTime, AllEmoji())
  for i:=0;i<dataPointCount;i++{
    offset := int(interval.Seconds())*(dataPointCount-i)
    startTime := endTime.Add(time.Second*time.Duration(offset)*-1)

    m := GetMood(FilterEmoji(startTime, startTime.Add(interval), periodEmoji))
    m.Time = startTime
    m.TimeString = startTime.Format("Jan _2")
    points = append(points, m)
  }

  return points
}
