package commands

import (
  "github.com/nats-io/nats.go"
)

// Singleton NATS and JetStream clients that command can use.
// As of today there should be no contention, this is not thread safe,
// and that's ok!
var singleton struct {
  nc      *nats.Conn
  js      nats.JetStreamContext
  subs    []*nats.Subscription
}

func getJSContext() (nats.JetStreamContext, error)  {

  if singleton.nc == nil {
    log.debug("Connecting to %s", options.natsURL)
    nc, err := nats.Connect(options.natsURL)
    if err != nil {
      return nil, err
    }
    singleton.nc = nc
  }

  if singleton.js == nil {
    js, err := singleton.nc.JetStream()
    if err != nil {
      return nil, err
    }
    singleton.js = js
  }

  return singleton.js, nil
}

func trackSubscriptions(subs ...*nats.Subscription)  {
  if singleton.subs == nil {
    singleton.subs = []*nats.Subscription{}
  }
  for _, sub := range subs {
    singleton.subs = append(singleton.subs, sub)
  }
}

func ClientCleanup()  {
  if singleton.subs != nil {
    for _, sub := range singleton.subs {
      sub.Unsubscribe()
    }
    singleton.subs = nil
  }
  if singleton.js != nil {
    singleton.js = nil
  }
  if singleton.nc != nil {
    singleton.nc.Close()
    singleton.nc = nil
  }
}
