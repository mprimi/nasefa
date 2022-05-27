package commands

import (
  "context"
  "flag"
  "fmt"
  "path"
  "github.com/google/subcommands"
  "github.com/nats-io/nats.go"
)

type receiveCommand struct {
  bucketName    string
}

func ReceiveCommand() (subcommands.Command) {
  return &receiveCommand{}
}

func (*receiveCommand) Name() string     { return "receive" }
func (*receiveCommand) Synopsis() string { return "Receive one or multiple files" }
func (*receiveCommand) Usage() string { return "receive [options] <destination_directory> <file_id> ...\n" }
func (p *receiveCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&p.bucketName, "bucket", defaultBucketName, "Name of the bucket where file is stored")
}

func (p *receiveCommand) Execute(_ context.Context, flagSet *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  numFiles := len(flagSet.Args()) - 1

  if numFiles < 1 {
    fmt.Printf("⚠️ Usage error: destination directory and at least one file id required\n")
    return subcommands.ExitUsageError
  }

  objStore, err := getObjStore(p.bucketName)
  if err != nil {
    fmt.Printf("❌ %s\n", err)
    return subcommands.ExitFailure
  }

  destinationDirectory := flagSet.Args()[0]
  fileIds := flagSet.Args()[1:]

  for i, fileId := range fileIds {

    objInfo, err := objStore.GetInfo(fileId)
    if err == nats.ErrObjectNotFound {
      fmt.Printf("❌ No such object '%s'\n", fileId)
      return subcommands.ExitFailure
    } else if err != nil {
      fmt.Printf("❌ Receive lookup error '%s': %s\n", fileId, err)
      return subcommands.ExitFailure
    }

    filename := getFilename(objInfo)
    if filename == "" {
      fmt.Printf("❌ Receive error, '%s' is not a file\n", fileId)
      return subcommands.ExitFailure
    }

    destinationFile := path.Join(destinationDirectory, filename)

    fmt.Printf("📥 Receiving file %d/%d: %s/%s => %s\n", i+1, numFiles, p.bucketName, fileId, destinationFile)

    err = objStore.GetFile(fileId, destinationFile)
    if err != nil {
      fmt.Printf("❌ Receive failed: %s\n", err)
      return subcommands.ExitFailure
    }
  }

  fmt.Printf("✅ Done\n")
  return subcommands.ExitSuccess
}
