package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	api "github.com/nlopes/slack"
)

var emojiList EmojiList

// Emoji represents a single use of an emoji
type Emoji struct {
	Name    string
	SeenAt  time.Time
	Channel string
	User    string
}

// EmojiList provides access to the emoji datastore
type EmojiList struct{}

// AddEmoji adds a new Emoji to the emoji datastore
func (e *EmojiList) AddEmoji(emoji string, m api.Message, id string) {
	i, _ := strconv.ParseFloat(m.Timestamp, 64)
	seen := time.Unix(int64(i), 0)

	em := &Emoji{}
	em.Name = emoji
	em.SeenAt = seen
	em.Channel = m.Channel
	em.User = m.User

	err := db.Update(func(tx *bolt.Tx) error {
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
			"error":   err,
			"emoji":   emoji,
			"message": m.Text,
			"channel": m.Channel,
			"time":    seen,
		}).Warning("Could not save emoji to Bolt")
	}
}

// List returns all emojis from the database
func (e *EmojiList) List() []*Emoji {
	var em []*Emoji
	err := db.View(func(tx *bolt.Tx) error {
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

// FilterEmoji filters a slice of Emojis based on their SeenAt values
func FilterEmoji(from time.Time, to time.Time, emoji []*Emoji) []*Emoji {
	var emj []*Emoji
	for _, e := range emoji {
		if e.SeenAt.After(from) && e.SeenAt.Before(to) {
			emj = append(emj, e)
		}
	}

	return emj
}

// AllEmoji returns all Emojis in the datastore
func AllEmoji() []*Emoji {
	return emojiList.List()
}

// ParseEmoji parses messages returning found emoji and reactions
func ParseEmoji(messages []api.Message) {
	r := regexp.MustCompile(`:([a-z0-9_\+\-]+):`)

	for _, m := range messages {
		msgID := fmt.Sprintf("%s-%s-%s", m.Timestamp, m.Channel, m.User)
		for _, r := range m.Reactions {
			emojiList.AddEmoji(r.Name, m, fmt.Sprintf("%s-%s-%s", msgID, m.User, m.Name))
		}

		foundEmoji := r.FindAllStringSubmatch(m.Text, -1)
		for _, em := range foundEmoji {
			emojiList.AddEmoji(em[1], m, msgID)
		}
	}
}
