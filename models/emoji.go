package models

import(
  "time"
  "regexp"
  "strconv"
  api "github.com/nlopes/slack"
)

var emojiList EmojiList

type Emoji struct{
  Name      string
  SeenAt    time.Time
}

type EmojiList struct{
  emoji     []*Emoji
}

func (e *EmojiList) AddEmoji(emoji string, seen time.Time){
  em := &Emoji{}
  em.Name = emoji
  em.SeenAt = seen
  e.emoji = append(e.emoji, em)
}

func (e *EmojiList) List() []*Emoji{
  return e.emoji
}

// Parse messages returning found emoji and reactions
func ParseEmoji(messages []api.Message){
  r := regexp.MustCompile(`:([a-z0-9_\+\-]+):`)

  for _,m := range messages{
    ts, _ := strconv.ParseFloat(m.Timestamp, 64)
    t := time.Unix(int64(ts), 0)

    for _,r := range m.Reactions{
      emojiList.AddEmoji(r.Name, t)
    }

    foundEmoji := r.FindAllStringSubmatch(m.Text, -1)
    for _,em := range foundEmoji{
      emojiList.AddEmoji(em[1], t)
    }
  }
}
