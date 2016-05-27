package models

import(
  "time"
  "github.com/samarudge/slackmood/emojiRanks"
)

type Mood struct{
  Positive      float32
  Negative      float32
  Neutral       float32
  PositiveCount int32
  NegativeCount int32
  NeutralCount  int32
  TotalCount    int32
}

func percentage(a int32, b int32) float32{
  return float32(a)/float32(b)*100
}

func GetMood(over time.Duration) Mood{
  m := Mood{}

  from := time.Now().UTC().Add(over*-1)
  for _, e := range emojiList.List(){
    if e.SeenAt.After(from){
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
  }

  m.Positive = percentage(m.PositiveCount, m.TotalCount)
  m.Negative = percentage(m.NegativeCount, m.TotalCount)
  m.Neutral = percentage(m.NeutralCount, m.TotalCount)

  return m
}
