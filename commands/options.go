package commands

import (
  "flag"
  "github.com/nats-io/nats.go"
)

const (
  defaultNatsURL = nats.DefaultURL
  defaultDebug = false
  defaultQuiet = false
  defaultNoEmojis = false
)

var options struct {
  natsURL             string
  debug               bool
  quiet               bool
  noEmojis            bool
}

func RegisterTopLevelFlags() {
  flag.StringVar(&options.natsURL, "natsURL", defaultNatsURL, "NATS server URL ")
  flag.BoolVar(&options.debug, "debug", defaultDebug, "Print debug statements")
  flag.BoolVar(&options.quiet, "quiet", defaultQuiet, "Quiet, only print fatal errors")
  flag.BoolVar(&options.noEmojis, "noEmoji", defaultNoEmojis, "Disable emojis in console messages")
}
