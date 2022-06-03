package commands

import (
  "errors"
  "fmt"
  "path"
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

func downloadBundle(destinationDirectory, bundleName string) (*fileBundle, error) {

  js, err := getJSContext()
  if err != nil {
    return nil, err // TODO wrap error (here and everywhere)
  }

  bundle, err := _loadBundle(js, bundleName)
  if err != nil {
    return nil, err
  }

  for _, file := range bundle.files {
    destinationPath := path.Join(destinationDirectory, file.fileName)
    err := downloadBundleFile(file, destinationPath)
    if err != nil {
      return nil, err
    }
  }

  return bundle, nil
}
