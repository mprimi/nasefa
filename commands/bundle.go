package commands

import (
  "errors"
  "fmt"
  "path"
  "os"
  "strings"
  "time"
  "github.com/nats-io/nats.go"
  "github.com/google/uuid"
)

type fileBundle struct {
  name              string
  objStore          nats.ObjectStore
  objStoreStatus    nats.ObjectStoreStatus
  files             []*bundleFile
}

type bundleFile struct {
  bundle    *fileBundle
  objInfo   *nats.ObjectInfo
  fileName  string
  id        string
}

const (
  kFilenameHeader = "nasefa-filename"
)

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

func newBundle(bundleName string, ttl time.Duration) (*fileBundle, error) {
  if bundleName == "" {
    bundleName = uuid.NewString()
  }

  objStore, err := createBucket(bundleName, ttl)
  if err != nil {
    return nil, err
  }

  objStoreStatus, err := objStore.Status()
  if err != nil {
    return nil, err
  }

  return &fileBundle{
    name: bundleName,
    objStore: objStore,
    objStoreStatus: objStoreStatus,
    files: []*bundleFile{},
  }, nil
}

func addFileToBundle(bundle *fileBundle, filePath, fileId string) (*bundleFile, error) {
  if fileId == "" {
    fileId = uuid.NewString()
  }

  file, err := os.Open(filePath)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  fileName := path.Base(filePath)

  objMeta := nats.ObjectMeta{
    Name: fileId,
    Description: fmt.Sprintf("Nasefa file %s", fileName),
    Headers: nats.Header{},
  }
  objMeta.Headers.Add(kFilenameHeader, fileName)

  objInfo, err := bundle.objStore.Put(&objMeta, file)
  if err != nil {
    return nil, err
  }

  bundleFile := &bundleFile{
    bundle: bundle,
    objInfo: objInfo,
    fileName: fileName,
    id: fileId,
  }

  bundle.files = append(bundle.files, bundleFile)

  return bundleFile, nil
}

func loadBundles() ([]*fileBundle, error) {

  js, err := getJSContext()
  if err != nil {
    return nil, err
  }

  // 🔥 HACK alert: ObjectStore API does not offer a way to list buckets.
  // But we can list streams, and it includes "system" streams.
  // TODO Rather than doing this, you may want to create a special meta-bucket with list of bundles
  bucketNames := []string{}
  for streamName := range js.StreamNames() {
    if strings.HasPrefix(streamName, "OBJ_") {
      bucketName := strings.ReplaceAll(streamName, "OBJ_", "")
      logDebug("Found bucket %s", bucketName)
      bucketNames = append(bucketNames, bucketName)
    }
  }

  bundles := []*fileBundle{}
  for _, bucketName := range bucketNames {
    bundle, err := _loadBundle(js, bucketName)
    if err != nil {
      logWarn("Skipping bucket '%s': %s", bucketName, err)
    }

    bundles = append(bundles, bundle)
  }

  return bundles, nil
}

func _loadBundle(js nats.JetStreamContext, bundleName string) (*fileBundle, error) {
  objStore, err := js.ObjectStore(bundleName)
  if err != nil {
    return nil, err
  }

  objStoreStatus, err := objStore.Status()
  if err != nil {
    return nil, err
  }

  objsInfo, err := objStore.List()
  if err != nil {
    return nil, err
  }

  bundle := &fileBundle{
    name: bundleName,
    objStore: objStore,
    objStoreStatus: objStoreStatus,
    files: []*bundleFile{},
  }

  for _, objInfo := range objsInfo {
    file, err := _loadFile(bundle, objInfo)
    if err != nil {
      logWarn("Skipping file '%s/%s': %s", bundleName, objInfo.Name,  err)
      continue
    }

    bundle.files = append(bundle.files, file)
  }

  return bundle, nil
}

func _loadFile(bundle *fileBundle, objInfo *nats.ObjectInfo) (*bundleFile, error) {
  fileId := objInfo.Name
  fileName := objInfo.Headers.Get(kFilenameHeader)
  if fileName == "" {
    return nil, errors.New("Invalid object metadata")
  }

  bundleFile := &bundleFile{
    bundle: bundle,
    objInfo: objInfo,
    fileName: fileName,
    id: fileId,
  }

  return bundleFile, nil
}

func downloadBundleFile(file *bundleFile, filePath string) (error)  {
  return file.bundle.objStore.GetFile(file.id, filePath)
}
