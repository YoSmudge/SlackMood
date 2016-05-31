#SlackMood

Gauge the average mood of your Slack based on Emoji use.

![](https://s3.amazonaws.com/f.cl.ly/items/0E3W453j2I44451b441x/Screen%20Shot%202016-05-31%20at%2015.01.18.png?v=7d9a7302)

Built for a hack day, I.E. everything here is around 6 hours work. The code is terrible, comments are few and far between, there are no tests, it's not all OOP, there is a lot of weirdness, my usage of BoltDB is not fantastic.

I may get round to doing a bit more work on this at some point, but for now I'm just dumping the code here.

## Usage
Create a config file with your Slack bot token and a path to the BoltDB file

```
slack_token: "abcd"
db_path: "./db.bolt"
```

Install dependencies with Glide and create the static assets

    go-bindata -o web/assets.go -pkg web --prefix "web/html/" web/html/...

Then build and run

    go build
    ./slackmood --config ./config.yml --bind :3044
