package commands

import (
  "flag"
  "github.com/nats-io/nats.go"
)

const defaultBucketName = "nasefa"
const defaultNatsURL = nats.DefaultURL

var options struct {
  bucketName          string
  natsURL             string
}

func RegisterTopLevelFlags() {
  flag.StringVar(&options.bucketName, "bucket",   defaultBucketName,  "Target a non-default bucket")
  flag.StringVar(&options.natsURL,    "natsURL",  defaultNatsURL,     "NATS server URL ")
}
