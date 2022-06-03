package commands

import (
  "errors"
  "fmt"
  "time"
  "github.com/google/uuid"
  "github.com/nats-io/nats.go"
)

const (
  kNotificationStreamName = "nasefa_bundles_notification_stream"
  kNotificationSubjectPrefix = "nasefa.bundle"
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

func createBucket(bucket string, ttl time.Duration) (nats.ObjectStore, error)  {
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
    TTL: ttl,
    //TODO: MaxBytes, Storage, Replicas
  }

  objStore, err := js.CreateObjectStore(&objStoreConfig)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Bucket creation error: %s", err))
  }

  return objStore, nil
}

func initNotificationStream() (error) {
  js, err := getJSContext()
  if err != nil {
    return err
  }

  streamConfig := &nats.StreamConfig{
    Name: kNotificationStreamName,
    Subjects: []string{kNotificationSubjectPrefix + ".*"},
    MaxMsgs: 100, //TODO hardcoded
  }
  streamInfo, err := js.AddStream(streamConfig)
  if err != nil {
    return err
  }

  logDebug("Created notification stream: %v", streamInfo)
  return nil
}

func notifyRecipients(bundle *fileBundle, recipients ...string) (error) {
  //TODO creates 2 clients, unnecessarily
  initNotificationStream()

  js, err := getJSContext()
  if err != nil {
    return err
  }

  msg := &nats.Msg{
    Data: []byte("New bundle: " + bundle.name),
    Header: nats.Header{
      // Set header to de-dupe notification on multiple recipients
      // https://docs.nats.io/using-nats/developer/develop_jetstream/model_deep_dive#message-deduplication
      "Nats-Msg-Id": []string{uuid.NewString()},
    },
  }

  for _, recipient := range recipients {
    msg.Subject = kNotificationSubjectPrefix + "." + recipient
    pubAck, err := js.PublishMsg(msg)
    if err != nil {
      return err
    }

    logDebug("Published notification for recipient %s: %v", recipient, pubAck)
  }

  return nil
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
