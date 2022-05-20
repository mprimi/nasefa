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
func (*receiveCommand) Synopsis() string { return "Receive a file" }
func (*receiveCommand) Usage() string { return "receive [options] <destination_directory> <file_id> ... \n" }
func (p *receiveCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&p.bucketName, "bucket", defaultBucketName, "Name of the bucket where file is stored")
}

func (p *receiveCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {

  numFiles := len(f.Args()) - 1

  if numFiles < 1 {
    fmt.Printf("⚠️ Usage error: destination directory and at least one file id required\n")
    return subcommands.ExitUsageError
  }

  objStore, err := getObjStore(p.bucketName)
  if err != nil {
    fmt.Printf("❌ %s\n", err)
    return subcommands.ExitFailure
  }

  destinationDirectory := f.Args()[0]
  fileIds := f.Args()[1:]

  for i, fileId := range fileIds {
    fmt.Printf("📥 Receiving file %d/%d: %s\n", i+1, numFiles, fileId)

    destinationFile := path.Join(destinationDirectory, path.Base(fileId))
    err = objStore.GetFile(fileId, destinationFile)
    if err != nil {
      fmt.Printf("❌ Receive failed: %s\n", err)
      return subcommands.ExitFailure
    }
  }

  fmt.Printf("✅ Done\n")
  return subcommands.ExitSuccess
}
