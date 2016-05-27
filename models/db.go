package models

import(
  "github.com/boltdb/bolt"
  "github.com/samarudge/slackmood/config"
)

var db *bolt.DB

func OpenDB() error{
  var err error
  db, err = bolt.Open(config.Config.Db, 0600, nil)
  if err != nil {
    return err
  }

  // IGNORE ALL THE ERRORS \o/
  db.Update(func(tx *bolt.Tx) error{
    tx.CreateBucketIfNotExists([]byte("emoji"))
    return nil
  })

  return nil
}
