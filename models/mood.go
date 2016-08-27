package models

import (
	"fmt"
	"time"

	"github.com/yosmudge/slackmood/emojiRanks"
)

// Mood represents the team's mood for a given duration
type Mood struct {
	Positive        float32
	PositiveDisplay string
	Negative        float32
	NegativeDisplay string
	Neutral         float32
	NeutralDisplay  string
	PositiveCount   int32
	NegativeCount   int32
	NeutralCount    int32
	TotalCount      int32
	Time            time.Time
	TimeString      string
}

func percentage(a int32, b int32) float32 {
	if b == 0 {
		return 0
	}

	return float32(a) / float32(b) * 100
}

// GetMood returns the mood based on a slice of Emojis
func GetMood(emoji []*Emoji) Mood {
	m := Mood{}

	for _, e := range emoji {
		for _, r := range ranks.EmojiRanks {
			if r.Name == e.Name {
				switch r.Rank {
				case 1:
					m.PositiveCount++
				case 0:
					m.NeutralCount++
				case -1:
					m.NegativeCount++
				}
				m.TotalCount++
				break
			}
		}
	}

	displayFormat := "%0.1f"
	m.Positive = percentage(m.PositiveCount, m.TotalCount)
	m.PositiveDisplay = fmt.Sprintf(displayFormat, m.Positive)
	m.Negative = percentage(m.NegativeCount, m.TotalCount)
	m.NegativeDisplay = fmt.Sprintf(displayFormat, m.Negative)
	m.Neutral = percentage(m.NeutralCount, m.TotalCount)
	m.NeutralDisplay = fmt.Sprintf(displayFormat, m.Neutral)

	return m
}

// GraphMood creates a structure suitable for charting
func GraphMood(over time.Duration, interval time.Duration) []Mood {
	var points []Mood

	now := time.Now().UTC()
	dataPointCount := int(over.Seconds() / interval.Seconds())
	endTime := time.Unix(int64(interval.Seconds())*int64(now.Unix()/int64(interval.Seconds())), 0)
	periodEmoji := FilterEmoji(endTime.Add(over*-1), endTime, AllEmoji())
	for i := 0; i < dataPointCount; i++ {
		offset := int(interval.Seconds()) * (dataPointCount - i)
		startTime := endTime.Add(time.Second * time.Duration(offset) * -1)

		m := GetMood(FilterEmoji(startTime, startTime.Add(interval), periodEmoji))
		m.Time = startTime
		m.TimeString = startTime.Format("Jan _2")
		points = append(points, m)
	}

	return points
}
