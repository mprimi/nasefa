package commands

import (
  "errors"
  "fmt"
  "io"
  "os"
  "path"
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
  kBundleDescription = "NASEFA files bundle"
)

var (
  kErrBundleNotFound = errors.New("Bundle not found")
  kErrBundleFileNotFound = errors.New("Bundle file not found")
)

func createBucket(bucket string, ttl time.Duration) (nats.ObjectStore, error)  {
  js, err := getJSContext()
  if err != nil {
    return nil, errors.New(fmt.Sprintf("JetStream init error: %s", err))
  }

  err = js.DeleteObjectStore(bucket)
  // TODO revisit later, maybe delete and recreate is not always the proper thing to do
  if err == nats.ErrStreamNotFound {
    // Bucket does not exist, as expected
  } else if err != nil {
    return nil, errors.New(fmt.Sprintf("Bucket wipe error: %s", err))
  }

  objStoreConfig := nats.ObjectStoreConfig{
    Bucket: bucket,
    Description: kBundleDescription,
    TTL: ttl,
    //TODO: MaxBytes, Storage, Replicas
  }

  objStore, err := js.CreateObjectStore(&objStoreConfig)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("Bucket creation error: %s", err))
  }

  return objStore, nil
}

func deleteBucket(bucket string) (error)  {
  js, err := getJSContext()
  if err != nil {
    return errors.New(fmt.Sprintf("JetStream init error: %s", err))
  }

  err = js.DeleteObjectStore(bucket)
  if err == nats.ErrStreamNotFound {
    // Bucket does not exist
    return nil
  } else if err != nil {
    return errors.New(fmt.Sprintf("Bucket delete error: %s", err))
  }

  return nil
}

func newBundle(bundleName string, ttl time.Duration) (*fileBundle, error) {
  if bundleName == "" {
    //TODO use bucket name validation regex
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

func deleteBundle(bundleName string) (error) {
  if bundleName == "" {
    //TODO use bucket name validation regex
    return errors.New("Invalid bucket name: " + bundleName)
  }
  err := deleteBucket(bundleName)
  return err
}

func addFileToBundle(bundle *fileBundle, filePath string) (*bundleFile, error) {

  fileName := path.Base(filePath)

  file, err := os.Open(filePath)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  return _addFileToBundle(bundle, file, fileName)
}

func _addFileToBundle(bundle *fileBundle, reader io.Reader, fileName string) (*bundleFile, error) {

  objMeta := nats.ObjectMeta{
    Name: fileName,
    Description: fmt.Sprintf("Nasefa file %s", fileName),
    Headers: nats.Header{},
  }
  objMeta.Headers.Add(kFilenameHeader, fileName)

  objInfo, err := bundle.objStore.Put(&objMeta, reader)
  if err != nil {
    return nil, err
  }

  bundleFile := &bundleFile{
    bundle: bundle,
    objInfo: objInfo,
    fileName: fileName,
    id: fileName,
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
      bucketName := strings.Replace(streamName, "OBJ_", "", 1)
      logDebug("Found bucket %s", bucketName)
      bucketNames = append(bucketNames, bucketName)
    }
  }

  bundles := []*fileBundle{}
  for _, bucketName := range bucketNames {
    bundle, err := _loadBundle(js, bucketName)
    if err != nil {
      // TODO expired empty bundles show up here
      logWarn("Skipping bucket '%s': %s", bucketName, err)
      continue
    }

    bundles = append(bundles, bundle)
  }

  return bundles, nil
}

func loadBundle(bundleName string) (*fileBundle, error) {

  js, err := getJSContext()
  if err != nil {
    return nil, err
  }

  bundle, err := _loadBundle(js, bundleName)
  if err != nil {
    return nil, err
  }

  return bundle, nil
}

func loadBundleFile(bundleName, fileName string) (*bundleFile, error) {

  js, err := getJSContext()
  if err != nil {
    return nil, err
  }

  bundle, err := _loadBundle(js, bundleName)
  if err != nil {
    return nil, err
  }

  for _, file := range bundle.files {
    if file.fileName == fileName {
      return file, nil
    }
  }

  return nil, kErrBundleFileNotFound
}

func getBundleFileReader(file *bundleFile) (io.Reader, error) {
  return file.bundle.objStore.Get(file.id)
}

func _loadBundle(js nats.JetStreamContext, bundleName string) (*fileBundle, error) {
  objStore, err := js.ObjectStore(bundleName)
  if err == nats.ErrBucketNotFound || err == nats.ErrStreamNotFound {
    return nil, kErrBundleNotFound
  } else if err != nil {
    return nil, err
  }

  objStoreStatus, err := objStore.Status()
  if err != nil {
    return nil, err
  }

  if objStoreStatus.Description() != kBundleDescription {
    // If description doesn't match, this may be a bucket created by some other
    // application, leave it alone.
    return nil, kErrBundleNotFound
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

func downloadBundleFile(file *bundleFile, filePath string) (error)  {
  return file.bundle.objStore.GetFile(file.id, filePath)
}
