package commands

import (
  "errors"
  "fmt"
  "github.com/nats-io/nats.go"
)

func getJSContext() (nats.JetStreamContext, error)  {
  nc, err := nats.Connect(options.natsURL)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Connection error: %s", err))
  }

  js, err := nc.JetStream()
  if err != nil {
    return nil, errors.New(fmt.Sprintf("JetStream init error: %s", err))
  }

  return js, nil
}

func createBucket(bucket string) (nats.ObjectStore, error)  {
  nc, err := nats.Connect(options.natsURL)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Connection error: %s", err))
  }

  js, err := nc.JetStream()
  if err != nil {
    return nil, errors.New(fmt.Sprintf("JetStream init error: %s", err))
  }

  err = js.DeleteObjectStore(bucket)
  if err == nats.ErrStreamNotFound {
    // Bucket does not exist, as expected
  } else if err != nil {
    return nil, errors.New(fmt.Sprintf("Bucket wipe error: %s", err))
  }

  objStoreConfig := nats.ObjectStoreConfig{
    Bucket: bucket,
    Description: "nasefa file bundle: " + bucket,
    //TODO: TTL, MaxBytes, Storage, Replicas
  }
  
  objStore, err := js.CreateObjectStore(&objStoreConfig)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Bucket creation error: %s", err))
  }

  return objStore, nil
}

func logDebug(format string, a ...interface{}) (int, error)  {
  return fmt.Printf(" 🐛 " + format + "\n", a...)
}

func logInfo(format string, a ...interface{}) (int, error)  {
  return fmt.Printf(" ℹ️  " + format + "\n", a...)
}

func logWarn(format string, a ...interface{}) (int, error)  {
  return fmt.Printf(" ⚠️ " + format + "\n", a...)
}
