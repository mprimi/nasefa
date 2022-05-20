package commands

import (
  "flag"
)

const defaultBucketName = "nasefa"

var options struct {
  bucketName          string
}

func RegisterTopLevelFlags() {
  flag.StringVar(&options.bucketName, "bucket", defaultBucketName, "Target a non-default bucket")
}
