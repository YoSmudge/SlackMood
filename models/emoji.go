package models

import(
  "fmt"
  "time"
  "regexp"
  "strconv"
  "encoding/json"
  "github.com/boltdb/bolt"
  log "github.com/Sirupsen/logrus"
  api "github.com/nlopes/slack"
)

var emojiList EmojiList

type Emoji struct{
  Name      string
  SeenAt    time.Time
  Channel   string
  User      string
}

type EmojiList struct{}

func (e *EmojiList) AddEmoji(emoji string, m api.Message, id string){
  i, _ := strconv.ParseFloat(m.Timestamp, 64)
  seen := time.Unix(int64(i), 0)

  em := &Emoji{}
  em.Name = emoji
  em.SeenAt = seen
  em.Channel = m.Channel
  em.User = m.User

  err := db.Update(func(tx *bolt.Tx) error{
    var err error
    v, err := json.Marshal(em)
    if err != nil {
      return err
    }

    err = tx.Bucket([]byte("emoji")).Put([]byte(id), v)
    return err
  })

  if err != nil {
    log.WithFields(log.Fields{
      "error": err,
      "emoji": emoji,
      "message": m.Text,
      "channel": m.Channel,
      "time": seen,
    }).Warning("Could not save emoji to Bolt")
  }
}

func (e *EmojiList) List() []*Emoji{
  var em []*Emoji
  err := db.View(func(tx *bolt.Tx) error{
    b := tx.Bucket([]byte("emoji"))
    c := b.Cursor()

    for k, v := c.First(); k != nil; k, v = c.Next() {
      e := Emoji{}
      err := json.Unmarshal(v, &e)
      if err != nil {
        return err
      }
      em = append(em, &e)
    }

    return nil
  })

  if err != nil {
    log.WithFields(log.Fields{
      "error": err,
    }).Warning("Could not list emoji by timestamp")
    return []*Emoji{}
  }
  return em
}

// Parse messages returning found emoji and reactions
func ParseEmoji(messages []api.Message){
  r := regexp.MustCompile(`:([a-z0-9_\+\-]+):`)

  for _,m := range messages{
    msgId := fmt.Sprintf("%s-%s-%s", m.Timestamp, m.Channel, m.User)
    for _,r := range m.Reactions{
      emojiList.AddEmoji(r.Name, m, fmt.Sprintf("%s-%s-%s", msgId, m.User, m.Name))
    }

    foundEmoji := r.FindAllStringSubmatch(m.Text, -1)
    for _,em := range foundEmoji{
      emojiList.AddEmoji(em[1], m, msgId)
    }
  }
}
