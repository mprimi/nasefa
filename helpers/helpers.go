package helpers

import (
  "log"
  "path"
  "github.com/nats-io/nats.go"
)

const bucket = "nasefa-default"

func getObjStore() (nats.ObjectStore, error)  {
  nc, err := nats.Connect(nats.DefaultURL)
  if err != nil {
    log.Printf("❌ Failed to connect to NATS: %s", err)
    return nil, err
  }

  js, err := nc.JetStream()
  if err != nil {
    log.Printf("❌ Failed to initialize JetStream: %s", err)
    return nil, err
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
      log.Printf("❌ Failed to create bucket %s: %s", bucket, err)
      return nil, err
    }

  } else if err != nil {
    log.Printf("❌ Failed to initialize Object Store bucket %s: %s", bucket, err)
    return nil, err
  }

  return objStore, nil
}

func UploadFile(sourceFilePath string) (error) {
  log.Printf("⏳ Uploading file %s", sourceFilePath)

  objStore, err := getObjStore()
  if err != nil {
    log.Printf("❌ Failed to upload: %s", err)
    return err
  }

  obj, err := objStore.PutFile(sourceFilePath)
  if err != nil {
    log.Printf("❌ Failed to upload: %s", err)
    return err
  }

  log.Printf("File uploaded: %s (%s)", obj.NUID, obj.Digest)

  log.Printf("✅ Done")
  return nil
}

func DownloadFiles(fileId, destinationDirectory string) (error) {
  log.Printf("⏳ Downloading files to %s", destinationDirectory)

  objStore, err := getObjStore()
  if err != nil {
    log.Printf("❌ Failed to download: %s", err)
    return err
  }

  destinationFile := path.Join(destinationDirectory, path.Base(fileId))
  err = objStore.GetFile(fileId, destinationFile)
  if err != nil {
    log.Printf("❌ Failed to download: %s", err)
    return err
  }

  log.Printf("✅ Done")
  return nil
}
