package commands

import (
  "errors"
  "fmt"
  "github.com/nats-io/nats.go"
)

func getObjStore(bucket string) (nats.ObjectStore, error)  {
  nc, err := nats.Connect(options.natsURL)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Connection error: %s", err))
  }

  js, err := nc.JetStream()
  if err != nil {
    return nil, errors.New(fmt.Sprintf("JetStream init error: %s", err))
  }

  objStore, err := js.ObjectStore(bucket)
  if err == nats.ErrStreamNotFound {
    // TODO: for now, create the bucket, should be an option later

    objStoreConfig := nats.ObjectStoreConfig{
      Bucket: bucket,
      Description: "Default bucket for nasefa file uploads",
      //TODO: TTL, MaxBytes, Storage, Replicas
    }
    objStore, err = js.CreateObjectStore(&objStoreConfig)
    if err != nil {
      return nil, errors.New(fmt.Sprintf("Bucket creation error: %s", err))

    }

  } else if err != nil {
    return nil, errors.New(fmt.Sprintf("Bucket lookup error: %s", err))
  }

  return objStore, nil
}
