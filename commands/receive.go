package commands

import (
  "context"
  "flag"
  "fmt"
  "path"
  "github.com/google/subcommands"
)

type receiveCommand struct {
  bucketName    string
}

func ReceiveCommand() (subcommands.Command) {
  return &receiveCommand{}
}

func (*receiveCommand) Name() string     { return "receive" }
func (*receiveCommand) Synopsis() string { return "Receive one or more file bundles" }
func (*receiveCommand) Usage() string { return "receive [options] <destination_directory> <bundle> ...\n" }
func (p *receiveCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&p.bucketName, "bucket", defaultBucketName, "Name of the bucket where file is stored")
}

func (p *receiveCommand) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  numBundles := len(flagSet.Args()) - 1

  if numBundles < 1 {
    fmt.Printf("⚠️ Usage error: destination directory and at least one bundle id required\n")
    return subcommands.ExitUsageError
  }

  js, err := getJSContext()
  if err != nil {
    fmt.Printf("❌ Error connecting: %s\n", err)
    return subcommands.ExitFailure
  }

  destinationDirectory := flagSet.Args()[0]
  bundleNames := flagSet.Args()[1:]

  numFilesReceived := 0
  for _, bundleName := range bundleNames {
    bundle, err := _loadBundle(js, bundleName)
    if err != nil {
      fmt.Printf("❌ Error receiving bundle '%s': %s\n", bundleName, err)
      return subcommands.ExitFailure
    }

    for _, file := range bundle.files {
      destinationPath := path.Join(destinationDirectory, file.fileName)
      err := downloadBundleFile(file, destinationPath)
      if err != nil {
        fmt.Printf("❌ Error downloading file '%s': %s\n", file.fileName, err)
        return subcommands.ExitFailure
      }
      numFilesReceived += 1
    }
  }

  fmt.Printf("✅ Received %d files in %d bundles\n", numFilesReceived, numBundles)
  return subcommands.ExitSuccess
}
