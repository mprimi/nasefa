package commands

import (
  "flag"
  "github.com/nats-io/nats.go"
)

const defaultNatsURL = nats.DefaultURL

var options struct {
  natsURL             string
}

func RegisterTopLevelFlags() {
  flag.StringVar(&options.natsURL,    "natsURL",  defaultNatsURL,     "NATS server URL ")
}
