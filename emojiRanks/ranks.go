package ranks

import(
  "os"
  "io"
  "path"
  "bufio"
  "strings"
  "strconv"
  "encoding/csv"
  log "github.com/Sirupsen/logrus"
)

type Rank struct{
  Name      string
  Rank      int64
}

var EmojiRanks []Rank

func init(){
  loadFiles := []string{"custom.csv", "standard.csv"}

  for _,f := range loadFiles{
    fp := path.Join("emojiRanks", f)
    fc, err := os.Open(fp)
    if err != nil {
      log.WithFields(log.Fields{
        "file": fp,
        "error": err,
      }).Error("Could not load emoji rank file")
      os.Exit(1)
    }

    r := csv.NewReader(bufio.NewReader(fc))
    for{
      rc, err := r.Read()
      if err == io.EOF{
        break
      }

      rank := Rank{}
      rank.Name = strings.TrimLeft(strings.TrimRight(rc[0], ":"), ":")
      rank.Rank, _ = strconv.ParseInt(rc[3], 10, 32)
      EmojiRanks = append(EmojiRanks, rank)
    }
  }

  log.WithFields(log.Fields{
    "emoji": len(EmojiRanks),
  }).Info("Loaded emoji rankings")
}
