package commands

import (
  "github.com/google/uuid"
  "github.com/nats-io/nats.go"
)

const (
  kNotificationStreamName = "nasefa_bundles_notification_stream"
  kNotificationSubjectPrefix = "nasefa.bundle"
  kBundleNameHeader = "nasefa_bundle_name"
)

func initNotificationStream(js nats.JetStreamContext) (error) {
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

  js, err := getJSContext()
  if err != nil {
    return err
  }
  initNotificationStream(js)

  headers := nats.Header{}
  headers.Add(kBundleNameHeader, bundle.name)
  // De-duplicates message published to multiple recipients.
  // Consumers use jetstream so they see a single message per bundle.
  // https://docs.nats.io/using-nats/developer/develop_jetstream/model_deep_dive#message-deduplication
  headers.Add(nats.MsgIdHdr, uuid.NewString())
  msg := &nats.Msg{
    Header: headers,
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

func watchBundles(recipientTags ...string) (<-chan string, error) {
  js, err := getJSContext()
  if err != nil {
    return nil, err
  }
  initNotificationStream(js)

  bundlesCh := make(chan string)
  subs := []*nats.Subscription{}
  doCleanup := true

  // Unsubscribe and close channel,
  // unless everything completed successfully.
  defer func() {
    if doCleanup {
      for _, sub := range subs {
        defer sub.Unsubscribe()
      }
      close(bundlesCh)
    }
  }()

  handleBundleNotification := func (msg *nats.Msg) {
    bundleName := msg.Header.Get(kBundleNameHeader)
    if bundleName == "" {
      logWarn("Invalid notification lacks bundle name")
      return
    }
    logDebug("Bundle notification: %s (subject: %s)", bundleName, msg.Subject)
    bundlesCh <- bundleName
  }

  subOpts := []nats.SubOpt{
    // Skip bundles published earlier than now
    nats.DeliverNew(),
  }

  for _, recipientTag := range recipientTags {
    subject := kNotificationSubjectPrefix + "." + recipientTag
    sub, err := js.Subscribe(subject, handleBundleNotification, subOpts...)
    if err != nil {
      return nil, err
    }
    subs = append(subs, sub)
    logDebug("Subscribed to bundle notifications subject: %s", subject)
  }

  doCleanup = false
  trackSubscriptions(subs...)
  return bundlesCh, nil
}
